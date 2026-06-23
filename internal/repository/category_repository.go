package repository

import (
	"context"
	"time"

	"github.com/luxe/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type categoryRepository struct {
	col *mongo.Collection
}

func NewCategoryRepository(db *Database) domain.CategoryRepository {
	return &categoryRepository{col: db.Collection("categories")}
}

func (r *categoryRepository) Create(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	cat.ID = primitive.NewObjectID()
	cat.Slug = slugify(cat.Name)
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, cat)
	return cat, err
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var cat domain.Category
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&cat)
	return &cat, err
}

func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var cat domain.Category
	err := r.col.FindOne(ctx, bson.M{"slug": slug}).Decode(&cat)
	return &cat, err
}

func (r *categoryRepository) Update(ctx context.Context, id string, update map[string]interface{}) (*domain.Category, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	update["updated_at"] = time.Now()
	if name, ok := update["name"].(string); ok {
		update["slug"] = slugify(name)
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated domain.Category
	err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, opts).Decode(&updated)
	return &updated, err
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *categoryRepository) List(ctx context.Context, activeOnly bool) ([]*domain.Category, error) {
	filter := bson.M{}
	if activeOnly {
		filter["is_active"] = true
	}
	cursor, err := r.col.Find(ctx, filter, options.Find().SetSort(bson.M{"sort_order": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var cats []*domain.Category
	return cats, cursor.All(ctx, &cats)
}

func (r *categoryRepository) GetTree(ctx context.Context) ([]*domain.Category, error) {
	all, err := r.List(ctx, true)
	if err != nil {
		return nil, err
	}

	// Build parent->children map
	childMap := make(map[string][]*domain.Category)
	for _, c := range all {
		if c.ParentID != nil {
			pid := c.ParentID.Hex()
			childMap[pid] = append(childMap[pid], c)
		}
	}

	// Root categories
	var roots []*domain.Category
	for _, c := range all {
		if c.ParentID == nil {
			c.Children = childMap[c.ID.Hex()]
			roots = append(roots, c)
		}
	}
	return roots, nil
}
