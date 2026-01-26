package database

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Close gracefully disconnects from MongoDB and cancels the context.
// No panic - errors are logged instead.
func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	if client == nil {
		log.Println("⚠️ MongoDB client is nil, skipping disconnect")
		return
	}

	log.Println("🔄 Disconnecting from MongoDB...")

	if err := client.Disconnect(ctx); err != nil {
		log.Printf("❌ Failed to disconnect from MongoDB: %v", err)
		slog.Error("MongoDB disconnect error", "error", err)
	} else {
		log.Println("✅ Successfully disconnected from MongoDB")
	}
}

// Connect establishes a connection to MongoDB with the provided URI.
// Returns client, context, cancel function, and error. No panic on failure.
func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	log.Println("🔄 Attempting to connect to MongoDB...")

	if uri == "" {
		log.Println("❌ MongoDB connection failed: URI is empty")
		return nil, nil, nil, fmt.Errorf("mongodb URI cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	log.Printf("🔑 MongoDB URI configured (timeout: 30s)")

	opts := options.Client().ApplyURI(uri)
	opts.SetConnectTimeout(30 * time.Second)

	client, err := mongo.Connect(opts)
	if err != nil {
		cancel() // Clean up context on error
		log.Printf("❌ Failed to connect to MongoDB: %v", err)
		return nil, nil, nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if client == nil {
		cancel() // Clean up context on error
		log.Println("❌ MongoDB client is nil after connection")
		return nil, nil, nil, fmt.Errorf("mongodb client initialization returned nil")
	}

	log.Println("✅ MongoDB connection established")
	return client, ctx, cancel, nil
}

// Ping verifies the MongoDB connection by sending a ping command.
// Returns error if ping fails. No panic.
func Ping(client *mongo.Client, ctx context.Context) error {
	if client == nil {
		log.Println("❌ Cannot ping MongoDB: client is nil")
		return fmt.Errorf("mongodb client is nil")
	}

	log.Println("🔄 Pinging MongoDB server...")

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Printf("❌ MongoDB ping failed: %v", err)
		slog.Error("MongoDB ping failed", "error", err)
		return fmt.Errorf("mongodb ping failed: %w", err)
	}

	log.Println("✅ MongoDB ping successful")
	slog.Info("MongoDB database connected!")
	return nil
}
