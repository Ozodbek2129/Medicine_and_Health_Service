package mongoDb

import (
	"context"
	"encoding/json"
	"fmt"
	pb "health/genproto/health_analytics"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// AddWearableData yangi kiyiladigan qurilma ma'lumotlarini qo'shish uchun
func (h *Health) ConsumeWearableDataQueue() {
	// RabbitMQ queue’dan xabarlarni olish
	messages, err := h.RabbitMQChannel.Consume(
		"wearable_data_queue", // Queue nomi
		"",                    // Consumer tag
		true,                  // Auto-ack
		false,                 // Exclusive
		false,                 // No-local
		false,                 // No-wait
		nil,                   // Arguments
	)
	if err != nil {
		log.Fatal("failed to register a consumer: %w", err)
		return
	}

	for msg := range messages {
		var message struct {
			Id string `json:"id"`
			*pb.AddWearableDataRequest
		}
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			h.Logger.Error("Failed to unmarshal message", "error", err)
			return
		}

		wearableData := bson.M{
			"id":                message.Id,
			"userid":            message.UserId,
			"devicetype":        message.DeviceType,
			"datatype":          message.DataType,
			"datavalue":         message.DataValue,
			"recordedtimestamp": message.RecordedTimestamp,
			"createdat":         time.Now().Format(time.RFC3339),
			"updatedat":         time.Now().Format(time.RFC3339),
			"deletedat":         "0",
		}
		// fmt.Println(wearableData)

		_, err = h.Db.Collection("wearable_data").InsertOne(context.Background(), wearableData)
		if err != nil {
			h.Logger.Error("Failed to insert wearable data into MongoDB", "error", err)
			return
		}
	}
}

func (h *Health) ConsumeHealthRecommendationsQueue() {
	// RabbitMQ queue’dan xabarlarni olish
	messages, err := h.RabbitMQChannel.Consume(
		"health_recommendations_queue", // Queue nomi
		"",                             // Consumer tag
		true,                           // Auto-ack
		false,                          // Exclusive
		false,                          // No-local
		false,                          // No-wait
		nil,                            // Arguments
	)
	if err != nil {
		h.Logger.Error("Failed to register a consumer", "error", err)
		return
	}

	for msg := range messages {
		// Xabarni JSON formatidan chiqarish
		var message struct {
			UserId             string `json:"user_id"`
			RecommendationType string `json:"recommendation_type"`
			Description        string `json:"description"`
			Priority           int    `json:"priority"`
		}

		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			h.Logger.Error("Failed to unmarshal message", "error", err)
			continue
		}

		// MongoDB kolleksiyasiga ulanadi
		coll := h.Db.Collection("health")

		// UUID va sana yaratish
		id := uuid.NewString()
		date := time.Now().Format("2006/01/02")

		// MongoDB ga yangi tavsiya kiritish
		_, err = coll.InsertOne(context.Background(), bson.M{
			"id":                  id,
			"user_id":             message.UserId,
			"recommendation_type": message.RecommendationType,
			"description":         message.Description,
			"priority":            message.Priority,
			"created_at":          date,
			"updated_at":          date,
			"deleted_at":          "0",
		})

		fmt.Println(message.UserId)
		fmt.Println(message.Description)
		fmt.Println(message.Priority)
		fmt.Println(message.RecommendationType)

		if err != nil {
			h.Logger.Error("Failed to insert health recommendation into MongoDB", "error", err)
			continue
		}

		// Redis uchun ma'lumotlarni tayyorlash va saqlash
		redisKey := message.UserId
		redisValue, err := json.Marshal(bson.M{
			"id":                  id,
			"user_id":             message.UserId,
			"recommendation_type": message.RecommendationType,
			"description":         message.Description,
			"priority":            message.Priority,
			"created_at":          date,
			"updated_at":          date,
		})
		if err != nil {
			h.Logger.Error("Failed to marshal recommendation for Redis", "error", err)
			continue
		}

		// Redisga yozish
		err = h.Redis.Set(context.Background(), redisKey, redisValue, 0).Err()
		if err != nil {
			h.Logger.Error("Failed to write recommendation to Redis", "error", err)
			continue
		}
	}
}

// coll := h.Db.Collection("health")

// id := uuid.NewString()
// date := time.Now().Format("2006/01/02")

// _, err := coll.InsertOne(ctx, bson.M{
// 	"id":                  id,
// 	"user_id":             req.UserId,
// 	"recommendation_type": req.RecommendationType,
// 	"description":         req.Description,
// 	"priority":            req.Priority,
// 	"created_at":          date,
// 	"updated_at":          date,
// 	"deleted_at":          "0",
// })

// if err != nil {
// 	return nil, err
// }

// recod := pb.HealthRecommendation{
// 	Id:                 id,
// 	UserId:             req.UserId,
// 	RecommendationType: req.RecommendationType,
// 	Description:        req.Description,
// 	Priority:           req.Priority,
// 	CreatedAt:          date,
// 	UpdatedAt:          date,
// }

// // Redisga yozish
// redisKey := req.UserId
// redisValue, err := json.Marshal(recod)
// if err != nil {
// 	return nil, err
// }

// err = h.Redis.Set(ctx, redisKey, redisValue, 0).Err()
// if err != nil {
// 	return nil, err
// }

// return &pb.GenerateHealthRecommendationsResponse{
// 	Recommendations: &recod,
// }, nil
// }
