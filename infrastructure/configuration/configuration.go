package configuration

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	MongoDB  MongoDBConfig
	Postgres PostgresConfig
	Supabase SupabaseConfig
}

type MongoDBConfig struct {
	URL string
}

type PostgresConfig struct {
	DSN string
}

type SupabaseConfig struct {
	URL string
	Key string
}

type CacheConfig struct {
	Redis RedisConfig
}

type RedisConfig struct {
	URL string
}

// LoadEnv loads environment variables from .env file and returns configuration.
// Falls back to defaults and logs warnings for missing values. No panic.
func LoadEnv() Config {
	log.Println("🔄 Loading environment configuration...")

	// Try to load .env file (optional - may not exist in production)
	err := godotenv.Load()
	if err != nil {
		log.Printf("⚠️ No .env file found (this is OK in production): %v", err)
		log.Println("📝 Will use environment variables from system")
	} else {
		log.Println("✅ Loaded configuration from .env file")
	}

	// Load SERVER_PORT with default
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
		log.Printf("⚠️ SERVER_PORT not set, using default: %s", port)
	} else {
		log.Printf("✅ SERVER_PORT: %s", port)
	}

	// Load MONGODB_URL (preferred), fallback to legacy MONGODB_ATLAS_URL
	mongoURI := os.Getenv("MONGODB_URL")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGODB_ATLAS_URL")
		if mongoURI != "" {
			log.Println("⚠️ MONGODB_ATLAS_URL is deprecated, please use MONGODB_URL")
		}
	}
	if mongoURI == "" {
		log.Println("⚠️ MONGODB_URL not set - MongoDB will not be available")
	} else {
		log.Println("✅ MONGODB_URL configured")
	}

	// Load POSTGRES_DSN (preferred), fallback to legacy SUPABASE_URL when it is a DSN
	postgresDSN := os.Getenv("POSTGRES_DSN")
	legacySupabaseURL := os.Getenv("SUPABASE_URL")
	if postgresDSN == "" && legacySupabaseURL != "" && !strings.HasPrefix(legacySupabaseURL, "http") {
		postgresDSN = legacySupabaseURL
		log.Println("⚠️ SUPABASE_URL as DSN is deprecated, please use POSTGRES_DSN")
	}
	if postgresDSN == "" {
		log.Println("⚠️ POSTGRES_DSN not set - PostgreSQL will not be available")
	} else {
		log.Println("✅ POSTGRES_DSN configured")
	}

	// Load SUPABASE API URL (optional, for hosted Supabase client)
	supabaseURL := os.Getenv("SUPABASE_API_URL")
	if supabaseURL == "" && (strings.HasPrefix(legacySupabaseURL, "http://") || strings.HasPrefix(legacySupabaseURL, "https://")) {
		supabaseURL = legacySupabaseURL
		log.Println("⚠️ SUPABASE_URL for API client is deprecated, please use SUPABASE_API_URL")
	}
	if supabaseURL == "" {
		log.Println("⚠️ SUPABASE_API_URL not set - Supabase API client will not be available")
	} else {
		log.Println("✅ SUPABASE_API_URL configured")
	}

	// Load SUPABASE_KEY
	supabaseKey := os.Getenv("SUPABASE_KEY")
	if supabaseKey == "" {
		log.Println("⚠️ SUPABASE_KEY not set - Supabase will not be available")
	} else {
		log.Println("✅ SUPABASE_KEY configured")
	}

	redisURL := strings.TrimSpace(os.Getenv("REDIS_URL"))
	if redisURL == "" {
		log.Println("⚠️ REDIS_URL not set - cache will run in no-cache mode")
	} else {
		log.Println("✅ REDIS_URL configured")
	}

	log.Println("✅ Configuration loaded successfully")

	return Config{
		Server: ServerConfig{
			Port: port,
		},
		Database: DatabaseConfig{
			MongoDB: MongoDBConfig{
				URL: mongoURI,
			},
			Postgres: PostgresConfig{
				DSN: postgresDSN,
			},
			Supabase: SupabaseConfig{
				URL: supabaseURL,
				Key: supabaseKey,
			},
		},
		Cache: CacheConfig{
			Redis: RedisConfig{
				URL: redisURL,
			},
		},
	}
}

// Validate checks if the configuration is valid and returns an error if not.
// This allows graceful handling of configuration issues.
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	// At least one database must be configured
	hasDB := false
	if c.Database.MongoDB.URL != "" {
		hasDB = true
	}
	if c.Database.Postgres.DSN != "" {
		hasDB = true
	}
	if c.Database.Supabase.URL != "" && c.Database.Supabase.Key != "" {
		hasDB = true
	}

	if !hasDB {
		return fmt.Errorf("at least one database (MongoDB or PostgreSQL/Supabase) must be configured")
	}

	return nil
}
