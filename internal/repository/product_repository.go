package repository

import (
	"context"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/luxe/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productRepository struct {
	col *mongo.Collection
}

// NewProductRepository creates a new MongoDB-backed product repository
func NewProductRepository(db *Database) domain.ProductRepository {
	return &productRepository{col: db.Collection("products")}
}

func slugify(s string) string {
	s = strings.ToLower(s)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func isUpper(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return false
		}
	}
	return true
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	product.ID = primitive.NewObjectID()
	product.Slug = slugify(product.Name)
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	if product.Status == "" {
		product.Status = domain.ProductStatusActive
	}
	if product.Thumbnail == "" && len(product.Images) > 0 {
		product.Thumbnail = product.Images[0]
	}

	_, err := r.col.InsertOne(ctx, product)
	return product, err
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var product domain.Product
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&product)
	return &product, err
}

func (r *productRepository) FindBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	var product domain.Product
	err := r.col.FindOne(ctx, bson.M{"slug": slug}).Decode(&product)
	return &product, err
}

func (r *productRepository) Update(ctx context.Context, id string, update map[string]interface{}) (*domain.Product, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	update["updated_at"] = time.Now()
	if name, ok := update["name"].(string); ok {
		update["slug"] = slugify(name)
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated domain.Product
	err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, opts).Decode(&updated)
	return &updated, err
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *productRepository) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	query := bson.M{}

	if filter.Status != "" {
		query["status"] = filter.Status
	} else {
		query["status"] = domain.ProductStatusActive
	}

	if filter.CategoryID != "" {
		if cid, err := primitive.ObjectIDFromHex(filter.CategoryID); err == nil {
			query["category_id"] = cid
		}
	}

	if filter.MinPrice > 0 || filter.MaxPrice > 0 {
		priceQuery := bson.M{}
		if filter.MinPrice > 0 {
			priceQuery["$gte"] = filter.MinPrice
		}
		if filter.MaxPrice > 0 {
			priceQuery["$lte"] = filter.MaxPrice
		}
		query["price"] = priceQuery
	}

	if filter.IsFeatured != nil {
		query["is_featured"] = *filter.IsFeatured
	}
	if filter.IsBestSeller != nil {
		query["is_best_seller"] = *filter.IsBestSeller
	}
	if filter.IsNewArrival != nil {
		query["is_new_arrival"] = *filter.IsNewArrival
	}

	if filter.Search != "" {
		query["$text"] = bson.M{"$search": filter.Search}
	}

	if len(filter.IDs) > 0 {
		oids := []primitive.ObjectID{}
		for _, id := range filter.IDs {
			if oid, err := primitive.ObjectIDFromHex(id); err == nil {
				oids = append(oids, oid)
			}
		}
		query["_id"] = bson.M{"$in": oids}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$all": filter.Tags}
	}

	total, err := r.col.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Sort
	sort := bson.D{}
	switch filter.SortBy {
	case "price_asc":
		sort = bson.D{{Key: "price", Value: 1}}
	case "price_desc":
		sort = bson.D{{Key: "price", Value: -1}}
	case "rating":
		sort = bson.D{{Key: "average_rating", Value: -1}}
	case "popular":
		sort = bson.D{{Key: "sold_count", Value: -1}}
	default: // newest
		sort = bson.D{{Key: "created_at", Value: -1}}
	}

	page := filter.Page
	limit := filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(sort)

	cursor, err := r.col.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id string, quantity int) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"stock": quantity, "updated_at": time.Now()},
	})
	return err
}

func (r *productRepository) UpdateRating(ctx context.Context, id string, avgRating float64, count int) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{
			"average_rating": avgRating,
			"review_count":   count,
			"updated_at":     time.Now(),
		},
	})
	return err
}

func (r *productRepository) IncrementSoldCount(ctx context.Context, id string, qty int) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$inc": bson.M{"sold_count": qty, "stock": -qty},
	})
	return err
}

func (r *productRepository) GetFeatured(ctx context.Context, limit int) ([]*domain.Product, error) {
	return r.findWithFilter(ctx, bson.M{"is_featured": true, "status": domain.ProductStatusActive}, limit)
}

func (r *productRepository) GetBestSellers(ctx context.Context, limit int) ([]*domain.Product, error) {
	return r.findWithFilter(ctx, bson.M{"is_best_seller": true, "status": domain.ProductStatusActive}, limit)
}

func (r *productRepository) GetNewArrivals(ctx context.Context, limit int) ([]*domain.Product, error) {
	return r.findWithFilter(ctx, bson.M{"is_new_arrival": true, "status": domain.ProductStatusActive}, limit)
}

func (r *productRepository) GetRelated(ctx context.Context, categoryID, excludeID string, limit int) ([]*domain.Product, error) {
	cid, _ := primitive.ObjectIDFromHex(categoryID)
	eid, _ := primitive.ObjectIDFromHex(excludeID)
	return r.findWithFilter(ctx, bson.M{
		"category_id": cid,
		"_id":         bson.M{"$ne": eid},
		"status":      domain.ProductStatusActive,
	}, limit)
}

func (r *productRepository) findWithFilter(ctx context.Context, filter bson.M, limit int) ([]*domain.Product, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}
