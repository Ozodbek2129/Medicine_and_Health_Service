package mongodb

// import (
// 	"context"
// 	"fmt"
// 	pb "health-service/genproto/health"
// 	"time"

// 	"github.com/google/uuid"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type HealthRepository interface {
// 	GenerateHealthRecommendations(ctx context.Context, req *pb.GenerateHealthRecommendationsReq) (*pb.GenerateHealthRecommendationsRes, error)
// 	GetRealtimeHealthMonitoring(ctx context.Context, req *pb.GetRealtimeHealthMonitoringReq) (*pb.GetRealtimeHealthMonitoringRes, error)
// 	GetDailyHealthSummary(ctx context.Context, req *pb.GetDailyHealthSummaryReq) (*pb.GetDailyHealthSummaryRes, error)
// 	GetWeeklyHealthSummary(ctx context.Context, req *pb.GetWeeklyHealthSummaryReq) (*pb.GetWeeklyHealthSummaryRes, error)
// }

// type healthRepositoryImpl struct {
// 	coll *mongo.Database
// }

// func NewHealthRepository(db *mongo.Database) HealthRepository {
// 	return &healthRepositoryImpl{coll: db}
// }

// func (r *healthRepositoryImpl) GenerateHealthRecommendations(ctx context.Context, req *pb.GenerateHealthRecommendationsReq) (*pb.GenerateHealthRecommendationsRes, error) {
// 	coll := r.coll.Collection("health")

// 	_, err := coll.InsertOne(ctx, bson.M{
// 		"_id":                 uuid.NewString(),
// 		"user_id":             req.UserId,
// 		"recommendation_type": req.RecommendationType,
// 		"description":         req.Description,
// 		"priority":            req.Priority,
// 		"created_at":          time.Now().Format("2006/01/02"),
// 		"updated_at":          time.Now().Format("2006/01/02"),
// 		"deleted_at":          0,
// 	})
// 	if err != nil {
// 		return &pb.GenerateHealthRecommendationsRes{
// 			Message: false,
// 		}, err
// 	}

// 	return &pb.GenerateHealthRecommendationsRes{
// 		Message: true,
// 	}, nil
// }

// func (r *healthRepositoryImpl) GetRealtimeHealthMonitoring(ctx context.Context, req *pb.GetRealtimeHealthMonitoringReq) (*pb.GetRealtimeHealthMonitoringRes, error) {
// 	var user pb.GetRealtimeHealthMonitoringRes
// 	coll := r.coll.Collection("health")

// 	err := coll.FindOne(ctx, bson.M{"$and": []bson.M{{"user_id": req.UserId}, {"deleted_at": 0}, {"created_at": time.Now().Format("2006/01/02")}}}).Decode(&user)
// 	if err != nil {
// 		return nil, fmt.Errorf("realtime health monitoring not found")
// 	}

// 	return &user, nil
// }

// func (r *healthRepositoryImpl) GetDailyHealthSummary(ctx context.Context, req *pb.GetDailyHealthSummaryReq) (*pb.GetDailyHealthSummaryRes, error) {
// 	var summary pb.GetDailyHealthSummaryRes
// 	coll := r.coll.Collection("health")

// 	err := coll.FindOne(ctx, bson.M{"$and": []bson.M{{"user_id": req.UserId}, {"deleted_at": 0}, {"created_at": req.Date}}}).Decode(&summary)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &summary, nil
// }

// func (r *healthRepositoryImpl) GetWeeklyHealthSummary(ctx context.Context, req *pb.GetWeeklyHealthSummaryReq) (*pb.GetWeeklyHealthSummaryRes, error) {
// 	var summary pb.GetWeeklyHealthSummaryRes
// 	coll := r.coll.Collection("health")

// 	cursor, err := coll.Find(ctx, bson.M{
// 		"$and": []bson.M{
// 			{"user_id": req.UserId},
// 			{"deleted_at": 0},
// 			{"created_at": bson.M{
// 				"$gte": req.StartDate,
// 				"$lte": req.EndDate,
// 			}},
// 		},
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer cursor.Close(ctx)

// 	for cursor.Next(ctx) {
// 		var doc pb.HealthRecommendation
// 		if err := cursor.Decode(&doc); err != nil {
// 			return nil, fmt.Errorf("error decoding document: %v", err)
// 		}
// 		summary.Recommendations = append(summary.Recommendations, &doc)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, fmt.Errorf("cursor error: %v", err)
// 	}

// 	return &summary, nil
// }
