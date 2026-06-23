package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database holds the MongoDB client and database reference
type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// Connect establishes a connection to MongoDB
func Connect(uri, dbName string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(uri).
		SetServerSelectionTimeout(10 * time.Second).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify the connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(dbName)
	database := &Database{Client: client, DB: db}

	// Create indexes
	if err := database.createIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return database, nil
}

// Disconnect gracefully closes the MongoDB connection
func (d *Database) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = d.Client.Disconnect(ctx)
}

// Collection returns a named collection
func (d *Database) Collection(name string) *mongo.Collection {
	return d.DB.Collection(name)
}

// createIndexes sets up all necessary MongoDB indexes
func (d *Database) createIndexes(ctx context.Context) error {
	// Users: unique email
	_, err := d.DB.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "role", Value: 1}}},
	})
	if err != nil {
		return err
	}

	// Products: unique slug, text search, category
	_, err = d.DB.Collection("products").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "category_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "is_featured", Value: 1}}},
		{Keys: bson.D{{Key: "is_best_seller", Value: 1}}},
		{Keys: bson.D{{Key: "price", Value: 1}}},
		{Keys: bson.D{{Key: "average_rating", Value: -1}}},
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "description", Value: "text"},
				{Key: "tags", Value: "text"},
			},
		},
	})
	if err != nil {
		return err
	}

	// Categories: unique slug
	_, err = d.DB.Collection("categories").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "parent_id", Value: 1}}},
	})
	if err != nil {
		return err
	}

	// Orders
	_, err = d.DB.Collection("orders").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "order_number", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	if err != nil {
		return err
	}

	// Coupons
	_, err = d.DB.Collection("coupons").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "code", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return err
	}

	// Reviews
	_, err = d.DB.Collection("reviews").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "product_id", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "is_approved", Value: 1}}},
	})

	return err
}
