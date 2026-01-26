package server

import (
	"fmt"
	"log"

	"github.com/LuuDinhTheTai/tzone/infrastructure/configuration"
	"github.com/gin-gonic/gin"
	supabase "github.com/supabase-community/supabase-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

type Server struct {
	r              *gin.Engine
	cfg            configuration.Config
	db             *gorm.DB
	mongoClient    *mongo.Client
	supabaseClient *supabase.Client
}

// NewServer creates a new server instance with the provided dependencies.
// Accepts nil values for optional database clients (MongoDB, Supabase).
func NewServer(
	r *gin.Engine,
	cfg configuration.Config,
	db *gorm.DB,
	mongoClient interface{},
	supabaseClient interface{},
) *Server {
	log.Println("🔧 Creating new server instance...")

	// Type assert the database clients
	var mongoClientTyped *mongo.Client
	var supaClientTyped *supabase.Client

	if mongoClient != nil {
		if mc, ok := mongoClient.(*mongo.Client); ok {
			mongoClientTyped = mc
			log.Println("✅ MongoDB client attached to server")
		}
	} else {
		log.Println("⚠️ Server created without MongoDB client")
	}

	if supabaseClient != nil {
		if sc, ok := supabaseClient.(*supabase.Client); ok {
			supaClientTyped = sc
			log.Println("✅ Supabase client attached to server")
		}
	} else {
		log.Println("⚠️ Server created without Supabase client")
	}

	return &Server{
		r:              r,
		cfg:            cfg,
		db:             db,
		mongoClient:    mongoClientTyped,
		supabaseClient: supaClientTyped,
	}
}

// Run starts the HTTP server. No panic - all errors are returned.
func (s *Server) Run() error {
	log.Println("🚀 Starting HTTP server...")

	// Re-initialize Gin engine with default middleware
	s.r = gin.Default()

	// Map all handlers
	log.Println("🔧 Mapping HTTP handlers...")
	if err := s.MapHandlers(); err != nil {
		log.Printf("❌ Failed to map handlers: %v", err)
		return fmt.Errorf("failed to map handlers: %w", err)
	}
	log.Println("✅ Handlers mapped successfully")

	// Validate port configuration
	if s.cfg.Server.Port == "" {
		log.Println("❌ Server port is not configured")
		return fmt.Errorf("server port cannot be empty")
	}

	// Start the server
	addr := ":" + s.cfg.Server.Port
	log.Printf("🌐 Server listening on %s", addr)
	log.Println("✅ HTTP server ready to accept connections")

	if err := s.r.Run(addr); err != nil {
		log.Printf("❌ Server failed to start: %v", err)
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

// HasMongoDB returns true if MongoDB client is available
func (s *Server) HasMongoDB() bool {
	return s.mongoClient != nil
}

// HasSupabase returns true if Supabase client is available
func (s *Server) HasSupabase() bool {
	return s.supabaseClient != nil
}
