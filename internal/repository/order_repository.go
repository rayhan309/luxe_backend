package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/luxe/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderRepository struct {
	col *mongo.Collection
}

func NewOrderRepository(db *Database) domain.OrderRepository {
	return &orderRepository{col: db.Collection("orders")}
}

func generateOrderNumber() string {
	return fmt.Sprintf("LUXE-%d", time.Now().UnixNano()/1000000)
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	order.ID = primitive.NewObjectID()
	order.OrderNumber = generateOrderNumber()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	if order.Status == "" {
		order.Status = domain.OrderStatusPending
	}
	if order.PaymentStatus == "" {
		order.PaymentStatus = domain.PaymentPending
	}
	order.StatusHistory = []domain.StatusHistory{{
		Status:    order.Status,
		ChangedAt: time.Now(),
	}}
	_, err := r.col.InsertOne(ctx, order)
	return order, err
}

func (r *orderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var order domain.Order
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&order)
	return &order, err
}

func (r *orderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	var order domain.Order
	err := r.col.FindOne(ctx, bson.M{"order_number": orderNumber}).Decode(&order)
	return &order, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id string, status string, note string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	history := domain.StatusHistory{Status: status, Note: note, ChangedAt: time.Now()}
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set":  bson.M{"status": status, "updated_at": time.Now()},
		"$push": bson.M{"status_history": history},
	})
	return err
}

func (r *orderRepository) UpdatePaymentStatus(ctx context.Context, id string, status string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"payment_status": status, "updated_at": time.Now()},
	})
	return err
}

func (r *orderRepository) SetTrackingNumber(ctx context.Context, id string, tracking string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"tracking_number": tracking, "updated_at": time.Now()},
	})
	return err
}

func (r *orderRepository) ListByUser(ctx context.Context, userID string, page, limit int) ([]*domain.Order, int64, error) {
	uid, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"user_id": uid}
	total, _ := r.col.CountDocuments(ctx, filter)
	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var orders []*domain.Order
	return orders, total, cursor.All(ctx, &orders)
}

func (r *orderRepository) List(ctx context.Context, status string, page, limit int, search string) ([]*domain.Order, int64, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	if search != "" {
		filter["$or"] = bson.A{
			bson.M{"order_number": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	total, _ := r.col.CountDocuments(ctx, filter)
	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	var orders []*domain.Order
	return orders, total, cursor.All(ctx, &orders)
}

func (r *orderRepository) GetStats(ctx context.Context) (*domain.OrderStats, error) {
	stats := &domain.OrderStats{}

	total, _ := r.col.CountDocuments(ctx, bson.M{})
	stats.TotalOrders = total

	pending, _ := r.col.CountDocuments(ctx, bson.M{"status": domain.OrderStatusPending})
	stats.PendingOrders = pending

	// Today
	today := time.Now().Truncate(24 * time.Hour)
	todayOrders, _ := r.col.CountDocuments(ctx, bson.M{"created_at": bson.M{"$gte": today}})
	stats.TodayOrders = todayOrders

	todayRevenuePipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"created_at":     bson.M{"$gte": today},
			"payment_status": domain.PaymentPaid,
		}}},
		{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$total"}}}},
	}
	if cursor, err := r.col.Aggregate(ctx, todayRevenuePipeline); err == nil {
		defer cursor.Close(ctx)
		var result []struct{ Total float64 `bson:"total"` }
		if cursor.All(ctx, &result) == nil && len(result) > 0 {
			stats.TodayRevenue = result[0].Total
		}
	}

	// Revenue aggregation
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"payment_status": domain.PaymentPaid}}},
		{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$total"}}}},
	}
	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)
		var result []struct{ Total float64 `bson:"total"` }
		if cursor.All(ctx, &result) == nil && len(result) > 0 {
			stats.TotalRevenue = result[0].Total
		}
	}

	if total > 0 {
		stats.AverageOrder = stats.TotalRevenue / float64(total)
	}

	return stats, nil
}

func (r *orderRepository) GetRevenueByPeriod(ctx context.Context, days int) ([]domain.RevenuePoint, error) {
	since := time.Now().AddDate(0, 0, -days)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"created_at":     bson.M{"$gte": since},
			"payment_status": domain.PaymentPaid,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":     bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$created_at"}},
			"revenue": bson.M{"$sum": "$total"},
			"orders":  bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var raw []struct {
		Date    string  `bson:"_id"`
		Revenue float64 `bson:"revenue"`
		Orders  int64   `bson:"orders"`
	}
	if err := cursor.All(ctx, &raw); err != nil {
		return nil, err
	}

	points := make([]domain.RevenuePoint, len(raw))
	for i, r := range raw {
		points[i] = domain.RevenuePoint{Date: r.Date, Revenue: r.Revenue, Orders: r.Orders}
	}
	return points, nil
}
