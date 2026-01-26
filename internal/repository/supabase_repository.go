package repository

import (
	"fmt"
	"log"

	supabase "github.com/supabase-community/supabase-go"
)

type SupabaseRepository struct {
	client *supabase.Client
}

// NewSupabaseRepository creates a new Supabase repository instance.
// Accepts nil client for graceful handling when Supabase is not available.
func NewSupabaseRepository(client *supabase.Client) *SupabaseRepository {
	if client == nil {
		log.Println("⚠️ SupabaseRepository created with nil client")
	} else {
		log.Println("✅ SupabaseRepository created with valid client")
	}

	return &SupabaseRepository{
		client: client,
	}
}

// IsAvailable returns true if the Supabase client is available
func (r *SupabaseRepository) IsAvailable() bool {
	return r.client != nil
}

// QueryTable queries a Supabase table and returns the results.
// Returns error if client is not available or query fails. No panic.
func (r *SupabaseRepository) QueryTable(tableName string) (interface{}, error) {
	if r.client == nil {
		log.Printf("⚠️ Cannot query table '%s': Supabase client not available", tableName)
		return nil, fmt.Errorf("supabase client is not available")
	}

	if tableName == "" {
		log.Println("❌ Cannot query: table name is empty")
		return nil, fmt.Errorf("table name cannot be empty")
	}

	log.Printf("🔄 Querying Supabase table: %s", tableName)

	// Example of using Supabase Postgrest API
	// Uncomment and adapt based on your schema:
	/*
		result, err := r.client.Postgrest.From(tableName).Select("*", "exact", false).Execute()
		if err != nil {
			log.Printf("❌ Query failed for table '%s': %v", tableName, err)
			return nil, fmt.Errorf("failed to query table %s: %w", tableName, err)
		}

		log.Printf("✅ Successfully queried table '%s'", tableName)
		return result, nil
	*/

	log.Printf("⚠️ QueryTable is not yet implemented for table '%s'", tableName)
	return nil, fmt.Errorf("query method not yet implemented")
}

// Insert inserts a record into a Supabase table. No panic on error.
func (r *SupabaseRepository) Insert(tableName string, data interface{}) error {
	if r.client == nil {
		log.Printf("⚠️ Cannot insert into table '%s': Supabase client not available", tableName)
		return fmt.Errorf("supabase client is not available")
	}

	if tableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}

	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}

	log.Printf("🔄 Inserting into Supabase table: %s", tableName)

	// Example implementation:
	/*
		_, err := r.client.Postgrest.From(tableName).Insert(data).Execute()
		if err != nil {
			log.Printf("❌ Insert failed for table '%s': %v", tableName, err)
			return fmt.Errorf("failed to insert into table %s: %w", tableName, err)
		}

		log.Printf("✅ Successfully inserted into table '%s'", tableName)
		return nil
	*/

	log.Printf("⚠️ Insert is not yet implemented for table '%s'", tableName)
	return fmt.Errorf("insert method not yet implemented")
}

// Update updates records in a Supabase table. No panic on error.
func (r *SupabaseRepository) Update(tableName string, filter map[string]interface{}, data interface{}) error {
	if r.client == nil {
		log.Printf("⚠️ Cannot update table '%s': Supabase client not available", tableName)
		return fmt.Errorf("supabase client is not available")
	}

	if tableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}

	log.Printf("🔄 Updating Supabase table: %s", tableName)
	log.Printf("⚠️ Update is not yet implemented for table '%s'", tableName)
	return fmt.Errorf("update method not yet implemented")
}

// Delete deletes records from a Supabase table. No panic on error.
func (r *SupabaseRepository) Delete(tableName string, filter map[string]interface{}) error {
	if r.client == nil {
		log.Printf("⚠️ Cannot delete from table '%s': Supabase client not available", tableName)
		return fmt.Errorf("supabase client is not available")
	}

	if tableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}

	log.Printf("🔄 Deleting from Supabase table: %s", tableName)
	log.Printf("⚠️ Delete is not yet implemented for table '%s'", tableName)
	return fmt.Errorf("delete method not yet implemented")
}
