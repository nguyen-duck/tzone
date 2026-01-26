package database

import (
	"fmt"
	"log"

	"github.com/supabase-community/supabase-go"
)

// ConnectSupabase initializes a Supabase client using the provided URL and service/anon key.
// Returns nil client with an error when either URL or key is missing.
// This function does not panic - all errors are returned to the caller for graceful handling.
func ConnectSupabase(url, key string) (*supabase.Client, error) {
	log.Println("🔄 Attempting to connect to Supabase...")

	// Validate URL
	if url == "" {
		log.Println("❌ Supabase connection failed: URL is empty")
		return nil, fmt.Errorf("supabase URL cannot be empty")
	}

	// Validate Key
	if key == "" {
		log.Printf("❌ Supabase connection failed: API key is empty (URL: %s)", url)
		return nil, fmt.Errorf("supabase API key cannot be empty")
	}

	log.Printf("🔑 Initializing Supabase client for URL: %s", url)

	// Attempt to create client
	client, err := supabase.NewClient(url, key, nil)
	if err != nil {
		log.Printf("❌ Failed to create Supabase client: %v", err)
		return nil, fmt.Errorf("failed to initialize supabase client: %w", err)
	}

	// Validate client was created
	if client == nil {
		log.Println("❌ Supabase client is nil after initialization")
		return nil, fmt.Errorf("supabase client initialization returned nil")
	}

	log.Println("✅ Successfully connected to Supabase instance")
	return client, nil
}
