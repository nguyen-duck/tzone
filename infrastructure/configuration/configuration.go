package configuration

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	MongoDbAtlas MongoDBConfig
	Supabase     SupabaseConfig
}

type MongoDBConfig struct {
	URL string
}

type SupabaseConfig struct {
	URL string
	Key string
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

	// Load MONGODB_ATLAS_URL
	mongoURI := os.Getenv("MONGODB_ATLAS_URL")
	if mongoURI == "" {
		log.Println("⚠️ MONGODB_ATLAS_URL not set - MongoDB will not be available")
	} else {
		log.Println("✅ MONGODB_ATLAS_URL configured")
	}

	// Load SUPABASE_URL
	supabaseURL := os.Getenv("SUPABASE_URL")
	if supabaseURL == "" {
		log.Println("⚠️ SUPABASE_URL not set - Supabase will not be available")
	} else {
		log.Println("✅ SUPABASE_URL configured")
	}

	// Load SUPABASE_KEY
	supabaseKey := os.Getenv("SUPABASE_KEY")
	if supabaseKey == "" {
		log.Println("⚠️ SUPABASE_KEY not set - Supabase will not be available")
	} else {
		log.Println("✅ SUPABASE_KEY configured")
	}

	log.Println("✅ Configuration loaded successfully")

	return Config{
		Server: ServerConfig{
			Port: port,
		},
		Database: DatabaseConfig{
			MongoDbAtlas: MongoDBConfig{
				URL: mongoURI,
			},
			Supabase: SupabaseConfig{
				URL: supabaseURL,
				Key: supabaseKey,
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
	if c.Database.MongoDbAtlas.URL != "" {
		hasDB = true
	}
	if c.Database.Supabase.URL != "" && c.Database.Supabase.Key != "" {
		hasDB = true
	}

	if !hasDB {
		return fmt.Errorf("at least one database (MongoDB or Supabase) must be configured")
	}

	return nil
}
