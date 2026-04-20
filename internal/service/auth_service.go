package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/LuuDinhTheTai/tzone/internal/model"
	"github.com/LuuDinhTheTai/tzone/internal/repository"
	"github.com/LuuDinhTheTai/tzone/util/email"
	"github.com/LuuDinhTheTai/tzone/util/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

const (
	otpPurposeRegister       = "register"
	otpPurposeForgotPassword = "forgot_password"
	otpPurposeChangePassword = "change_password"
	otpTTL                   = 5 * time.Minute
	otpMaxAttempts           = 5
)

// auth
type AuthService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.RefreshTokenRepository
	otpStore  *otpMemoryStore
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.RefreshTokenRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		otpStore:  newOTPMemoryStore(),
	}
}

func normalizeEmail(emailAddr string) string {
	return strings.TrimSpace(strings.ToLower(emailAddr))
}

func generateOTPCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (s *AuthService) issueOTP(emailAddr string, purpose string) error {
	code, err := generateOTPCode()
	if err != nil {
		return errors.New("failed to generate verification code")
	}

	now := time.Now()
	s.otpStore.set(emailAddr, purpose, code, now.Add(otpTTL))

	if err := email.SendOTP(emailAddr, code, purpose); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) verifyOTP(emailAddr string, purpose string, otp string) error {
	now := time.Now()
	entry, ok := s.otpStore.get(emailAddr, purpose, now)
	if !ok {
		return errors.New("invalid or expired verification code")
	}

	if entry.Attempts >= otpMaxAttempts {
		s.otpStore.delete(emailAddr, purpose)
		return errors.New("verification code has been locked, request a new code")
	}

	if strings.TrimSpace(otp) != entry.Code {
		s.otpStore.incrementAttempts(emailAddr, purpose)
		return errors.New("invalid or expired verification code")
	}

	s.otpStore.delete(emailAddr, purpose)
	return nil
}

func (s *AuthService) SendRegisterOTP(emailAddr string) error {
	emailAddr = normalizeEmail(emailAddr)

	existing, _ := s.userRepo.FindByEmail(emailAddr)
	if existing != nil {
		return errors.New("email already exists")
	}

	return s.issueOTP(emailAddr, otpPurposeRegister)
}

func (s *AuthService) SendResetPasswordOTP(emailAddr string) error {
	emailAddr = normalizeEmail(emailAddr)

	existing, _ := s.userRepo.FindByEmail(emailAddr)
	if existing == nil {
		return errors.New("email is not registered")
	}

	return s.issueOTP(emailAddr, otpPurposeForgotPassword)
}

func (s *AuthService) SendChangePasswordOTP(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	return s.issueOTP(normalizeEmail(user.Email), otpPurposeChangePassword)
}

// register
func (s *AuthService) Register(emailAddr string, password string, otp string) error {
	emailAddr = normalizeEmail(emailAddr)

	if err := s.verifyOTP(emailAddr, otpPurposeRegister, otp); err != nil {
		return err
	}

	// check email tồn tại
	existing, _ := s.userRepo.FindByEmail(emailAddr)
	if existing != nil {
		return errors.New("email already exists")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user := model.User{
		ID:           uuid.New(),
		Email:        emailAddr,
		PasswordHash: stringPtr(string(hash)),
	}

	// Gắn role mặc định là User cho tài khoản mới đăng ký
	return s.userRepo.Create(&user, model.RoleUser)
}

func (s *AuthService) ResetPassword(emailAddr string, otp string, newPassword string) error {
	emailAddr = normalizeEmail(emailAddr)

	if err := s.verifyOTP(emailAddr, otpPurposeForgotPassword, otp); err != nil {
		return err
	}

	user, err := s.userRepo.FindByEmail(emailAddr)
	if err != nil || user == nil {
		return errors.New("email is not registered")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	if err := s.userRepo.UpdatePasswordHash(user.ID.String(), string(newHash)); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (s *AuthService) ChangePassword(userID string, oldPassword string, newPassword string, otp string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if !hasPassword(user) {
		return errors.New("password is not set, please setup password first")
	}

	if err := s.verifyOTP(normalizeEmail(user.Email), otpPurposeChangePassword, otp); err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return errors.New("old password is incorrect")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	if err := s.userRepo.UpdatePasswordHash(userID, string(newHash)); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

// login
func (s *AuthService) Login(email string, password string) (string, string, *model.User, string, error) {
	email = normalizeEmail(email)

	user, roleName, err := s.userRepo.FindByEmailWithRole(email)

	if err != nil {
		return "", "", nil, "", errors.New("invalid email or password")
	}
	if !hasPassword(user) {
		return "", "", nil, "", errors.New("password is not set for this account")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(*user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", "", nil, "", errors.New("invalid email or password")
	}

	accessToken, refreshToken, err := s.issueTokenPair(user.ID)
	if err != nil {
		return "", "", nil, "", err
	}

	return accessToken, refreshToken, user, roleName, nil
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, idToken string) (string, string, *model.User, string, error) {
	clientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	if clientID == "" {
		return "", "", nil, "", errors.New("google sign-in is not configured")
	}

	payload, err := idtoken.Validate(ctx, idToken, clientID)
	if err != nil {
		return "", "", nil, "", errors.New("invalid google token")
	}

	emailValue, _ := payload.Claims["email"].(string)
	emailValue = normalizeEmail(emailValue)
	if emailValue == "" {
		return "", "", nil, "", errors.New("google account email is missing")
	}

	emailVerified, _ := payload.Claims["email_verified"].(bool)
	if !emailVerified {
		return "", "", nil, "", errors.New("google account email is not verified")
	}

	googleSub := strings.TrimSpace(payload.Subject)
	if googleSub == "" {
		return "", "", nil, "", errors.New("google account subject is missing")
	}

	user, roleName, err := s.userRepo.FindByEmailWithRole(emailValue)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, "", errors.New("failed to query user")
		}

		newUser := model.User{
			ID:        uuid.New(),
			Email:     emailValue,
			GoogleSub: stringPtr(googleSub),
		}
		if createErr := s.userRepo.Create(&newUser, model.RoleUser); createErr != nil {
			return "", "", nil, "", errors.New("failed to create user")
		}

		user, roleName, err = s.userRepo.FindByEmailWithRole(emailValue)
		if err != nil {
			return "", "", nil, "", errors.New("failed to load user")
		}
	} else {
		if user.GoogleSub == nil {
			if updateErr := s.userRepo.UpdateGoogleSub(user.ID.String(), googleSub); updateErr != nil {
				return "", "", nil, "", errors.New("failed to link google account")
			}
			user.GoogleSub = stringPtr(googleSub)
		} else if strings.TrimSpace(*user.GoogleSub) != googleSub {
			return "", "", nil, "", errors.New("google account does not match this email")
		}
	}

	accessToken, refreshToken, err := s.issueTokenPair(user.ID)
	if err != nil {
		return "", "", nil, "", err
	}

	return accessToken, refreshToken, user, roleName, nil
}

func (s *AuthService) SetupPassword(userID string, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if hasPassword(user) {
		return errors.New("password is already set")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	if err := s.userRepo.UpdatePasswordHash(userID, string(newHash)); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (s *AuthService) issueTokenPair(userID uuid.UUID) (string, string, error) {
	jti := uuid.New()
	accessToken, refreshToken, err := jwt.GenerateTokenPair(userID, jti)
	if err != nil {
		return "", "", errors.New("failed to generate tokens")
	}

	rtRecord := &model.RefreshToken{
		ID:        jti,
		UserID:    userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.tokenRepo.Create(rtRecord); err != nil {
		return "", "", errors.New("failed to save session")
	}

	return accessToken, refreshToken, nil
}

func stringPtr(value string) *string {
	v := value
	return &v
}

func hasPassword(user *model.User) bool {
	return user != nil && user.PasswordHash != nil && strings.TrimSpace(*user.PasswordHash) != ""
}

// RefreshToken handles generating a new token pair from a valid refresh token
func (s *AuthService) RefreshToken(tokenString string) (string, string, uuid.UUID, error) {
	userID, jti, err := jwt.ValidateRefreshToken(tokenString)
	if err != nil {
		return "", "", uuid.Nil, errors.New("invalid or expired refresh token")
	}

	// Check if this JTI exists in the database
	_, err = s.tokenRepo.FindByID(jti)
	if err != nil {
		// ALARM: The token is structurally valid but NOT in DB!
		// This likely means it was already used (Token Reuse) or forged.
		// Security action: Revoke ALL active sessions for this user.
		_ = s.tokenRepo.DeleteAllByUserID(userID)
		return "", "", uuid.Nil, errors.New("security breach detected: token reuse. All sessions revoked")
	}

	// Consume the old Refresh Token (Rotation)
	_ = s.tokenRepo.DeleteByID(jti)

	// Issue a new token pair
	newJTI := uuid.New()
	newAccessToken, newRefreshToken, err := jwt.GenerateTokenPair(userID, newJTI)
	if err != nil {
		return "", "", uuid.Nil, errors.New("failed to generate new tokens")
	}

	// Save new Refresh Token in DB
	rtRecord := &model.RefreshToken{
		ID:        newJTI,
		UserID:    userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.tokenRepo.Create(rtRecord); err != nil {
		return "", "", uuid.Nil, errors.New("failed to save new session")
	}

	return newAccessToken, newRefreshToken, userID, nil
}

// Logout consumes a refresh token to end the session
func (s *AuthService) Logout(tokenString string) error {
	_, jti, err := jwt.ValidateRefreshToken(tokenString)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	// Delete from DB regardless, preventing further use
	return s.tokenRepo.DeleteByID(jti)
}
