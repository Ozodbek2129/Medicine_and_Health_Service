package mongoDb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	logger "health/pkg"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	pb "health/genproto/health_analytics"

	"github.com/streadway/amqp"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Health struct {
	Logger          *slog.Logger
	Db              *mongo.Database
	Redis           *redis.Client
	RabbitMQChannel *amqp.Channel
}

func NewHealth(mdb *mongo.Database, rdb *redis.Client, amqpChannel *amqp.Channel) *Health {
	return &Health{
		Logger:          logger.NewLogger(),
		Db:              mdb,
		Redis:           rdb,
		RabbitMQChannel: amqpChannel,
	}
}

func (h *Health) AddMedicalRecord(ctx context.Context, req *pb.AddMedicalRecordRequest) (*pb.AddMedicalRecordResponse, error) {
	// Yangi UUID yaratish
	id := uuid.New().String()
	createdAt := time.Now().Format(time.RFC3339)
	updatedAt := time.Now().Format(time.RFC3339)

	recordBson := bson.M{
		"id":          id,
		"user_id":     req.UserId,
		"record_type": req.RecordType,
		"record_date": req.RecordDate,
		"description": req.Description,
		"doctor_id":   req.DoctorId,
		"attachments": req.Attachments,
		"created_at":  createdAt,
		"updated_at":  updatedAt,
		"deleted_at":  "0",
	}

	_, err := h.Db.Collection("medical_records").InsertOne(ctx, recordBson)
	if err != nil {
		h.Logger.Error("Failed to add medical record", "error", err)
		return nil, err
	}

	medicalRecord := &pb.MedicalRecord{
		Id:          id,
		UserId:      req.UserId,
		RecordType:  req.RecordType,
		RecordDate:  req.RecordDate,
		Description: req.Description,
		DoctorId:    req.DoctorId,
		Attachments: req.Attachments,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return &pb.AddMedicalRecordResponse{MedicalRecord: medicalRecord}, nil
}

func (h *Health) GetMedicalRecord(ctx context.Context, req *pb.GetMedicalRecordRequest) (*pb.GetMedicalRecordResponse, error) {
	var record pb.MedicalRecord
	err := h.Db.Collection("medical_records").FindOne(ctx, bson.M{"$and": []bson.M{{"id": req.Id}, {"deleted_at": "0"}}}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			h.Logger.Warn("Medical record not found", "record_id", req.Id)
			return nil, errors.New("tibbiy yozuv topilmadi")
		}
		h.Logger.Error("Failed to get medical record", "error", err)
		return nil, err
	}

	return &pb.GetMedicalRecordResponse{MedicalRecord: &record}, nil
}

func (h *Health) UpdateMedicalRecord(ctx context.Context, req *pb.UpdateMedicalRecordRequest) (*pb.UpdateMedicalRecordResponse, error) {
	updatedAt := time.Now().Format(time.RFC3339)

	update := bson.M{
		"$set": bson.M{
			"record_type": req.RecordType,
			"record_date": req.RecordDate,
			"description": req.Description,
			"doctor_id":   req.DoctorId,
			"attachments": req.Attachments,
			"updated_at":  updatedAt,
		},
	}

	result, err := h.Db.Collection("medical_records").UpdateOne(ctx, bson.M{"$and": []bson.M{{"id": req.Id}, {"deleted_at": "0"}}}, update)
	if err != nil {
		h.Logger.Error("Failed to update medical record", "error", err)
		return nil, err
	}

	if result.MatchedCount == 0 {
		h.Logger.Warn("Medical record not found for update", "record_id", req.Id)
		return &pb.UpdateMedicalRecordResponse{Success: false}, errors.New("tibbiy yozuv topilmadi")
	}

	return &pb.UpdateMedicalRecordResponse{Success: true}, nil
}

func (h *Health) DeleteMedicalRecord(ctx context.Context, req *pb.DeleteMedicalRecordRequest) (*pb.DeleteMedicalRecordResponse, error) {
	currentTime := time.Now().Format(time.RFC3339)

	result, err := h.Db.Collection("medical_records").UpdateOne(
		ctx,
		bson.M{"id": req.Id},
		bson.M{"$set": bson.M{"deleted_at": currentTime}},
	)
	if err != nil {
		h.Logger.Error("Failed to delete medical record", "error", err)
		return nil, err
	}

	if result.MatchedCount == 0 {
		h.Logger.Warn("Medical record not found for deletion", "record_id", req.Id)
		return &pb.DeleteMedicalRecordResponse{Success: false}, errors.New("tibbiy yozuv topilmadi")
	}

	return &pb.DeleteMedicalRecordResponse{Success: true}, nil
}

func (h *Health) ListMedicalRecords(ctx context.Context, req *pb.ListMedicalRecordsRequest) (*pb.ListMedicalRecordsResponse, error) {
	cursor, err := h.Db.Collection("medical_records").Find(ctx, bson.M{"$and": []bson.M{{"user_id": req.UserId}, {"deleted_at": "0"}}}, options.Find())
	if err != nil {
		h.Logger.Error("Failed to list medical records", "error", err)
		return nil, err
	}
	defer func() {
		if cursor != nil {
			cursor.Close(ctx)
		}
	}()

	var records []*pb.MedicalRecord
	for cursor.Next(ctx) {
		var record pb.MedicalRecord
		if err := cursor.Decode(&record); err != nil {
			h.Logger.Warn("Failed to decode medical record", "error", err)
			continue
		}
		records = append(records, &record)
	}

	if err := cursor.Err(); err != nil {
		h.Logger.Error("Cursor error while listing medical records", "error", err)
		return nil, err
	}

	return &pb.ListMedicalRecordsResponse{MedicalRecords: records}, nil
}

// AddLifestyleData yangi turmush tarzi ma'lumotlarini qo'shadi
func (h *Health) AddLifestyleData(ctx context.Context, req *pb.AddLifestyleDataRequest) (*pb.AddLifestyleDataResponse, error) {

	vaqt := time.Now().Format(time.RFC3339)
	id := uuid.NewString()

	lifestyleData := bson.M{
		"id":           id,
		"userid":       req.UserId,
		"datatype":     req.DataType,
		"datavalue":    req.DataValue,
		"recordeddate": req.RecordedDate,
		"createdat":    vaqt,
		"updatedat":    vaqt,
		"deletedat":    "0",
	}

	_, err := h.Db.Collection("lifestyle_data").InsertOne(ctx, lifestyleData)
	if err != nil {
		h.Logger.Error("Failed to add lifestyle data", "error", err)
		return nil, err
	}

	return &pb.AddLifestyleDataResponse{LifestyleData: &pb.LifestyleData{
		Id:           id,
		UserId:       req.UserId,
		DataType:     req.DataType,
		DataValue:    req.DataValue,
		RecordedDate: req.RecordedDate,
		CreatedAt:    vaqt,
		UpdatedAt:    vaqt,
	}}, nil
}

// GetAllLifestyleData metodi
func (h *Health) GetAllLifestyleData(ctx context.Context, req *pb.GetAllLifestyleDataRequest) (*pb.GetAllLifestyleDataResponse, error) {
	collection := h.Db.Collection("lifestyle_data")

	skip := (req.Page - 1) * req.Limit

	findOptions := options.Find()
	findOptions.SetLimit(req.GetLimit())
	findOptions.SetSkip(skip)

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		h.Logger.Error("Error finding lifestyle data", "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var lifestyleDataList []*pb.LifestyleData
	for cursor.Next(ctx) {
		var data pb.LifestyleData
		if err := cursor.Decode(&data); err != nil {
			h.Logger.Error("Error decoding lifestyle data", "error", err)
			return nil, err
		}
		lifestyleDataList = append(lifestyleDataList, &data)
	}

	if err := cursor.Err(); err != nil {
		h.Logger.Error("Cursor error", "error", err)
		return nil, err
	}

	return &pb.GetAllLifestyleDataResponse{
		Lifestyledata: lifestyleDataList,
	}, nil
}

// GetLifestyleData turmush tarzi ma'lumotlarini olish uchun
func (h *Health) GetLifestyleData(ctx context.Context, req *pb.GetLifestyleDataRequest) (*pb.GetLifestyleDataResponse, error) {
	var lifestyleData pb.LifestyleData

	err := h.Db.Collection("lifestyle_data").FindOne(ctx, bson.M{"$and": []bson.M{{"id": req.Id}, {"deletedat": "0"}}}).Decode(&lifestyleData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			h.Logger.Warn("Lifestyle data not found", "id", req.Id)
			return nil, errors.New("turmush tarzi ma'lumotlari topilmadi")
		}
		h.Logger.Error("Failed to get lifestyle data", "error", err)
		return nil, err
	}

	return &pb.GetLifestyleDataResponse{LifestyleData: &lifestyleData}, nil
}

// UpdateLifestyleData turmush tarzi ma'lumotlarini yangilash uchun
func (h *Health) UpdateLifestyleData(ctx context.Context, req *pb.UpdateLifestyleDataRequest) (*pb.UpdateLifestyleDataResponse, error) {
	update := bson.M{
		"$set": bson.M{
			"userid":       req.UserId,
			"datatype":     req.DataType,
			"datavalue":    req.DataValue,
			"recordeddate": req.RecordedDate,
			"updatedat":    time.Now().Format(time.RFC3339),
		},
	}

	result, err := h.Db.Collection("lifestyle_data").UpdateOne(ctx, bson.M{"$and": []bson.M{{"id": req.Id}, {"deletedat": "0"}}}, update)
	if err != nil {
		h.Logger.Error("Failed to update lifestyle data", "error", err)
		return nil, err
	}

	if result.MatchedCount == 0 {
		h.Logger.Warn("Lifestyle data not found for update", "id", req.Id)
		return &pb.UpdateLifestyleDataResponse{Success: false}, errors.New("turmush tarzi ma'lumotlari topilmadi")
	}

	return &pb.UpdateLifestyleDataResponse{Success: true}, nil
}

// DeleteLifestyleData turmush tarzi ma'lumotlarini o'chirish uchun
func (h *Health) DeleteLifestyleData(ctx context.Context, req *pb.DeleteLifestyleDataRequest) (*pb.DeleteLifestyleDataResponse, error) {
	// Hozirgi vaqtni olish
	currentTime := time.Now().Format(time.RFC3339)

	// Turmush tarzi ma'lumotlarini yangilash
	result, err := h.Db.Collection("lifestyle_data").UpdateOne(
		ctx,
		bson.M{"id": req.Id},
		bson.M{"$set": bson.M{"deletedat": currentTime}},
	)
	if err != nil {
		h.Logger.Error("Failed to delete lifestyle data", "error", err)
		return nil, err
	}

	// Agar hech qanday yozuv yangilanmagan bo'lsa
	if result.MatchedCount == 0 {
		h.Logger.Warn("Lifestyle data not found for deletion", "id", req.Id)
		return &pb.DeleteLifestyleDataResponse{Success: false}, errors.New("turmush tarzi ma'lumotlari topilmadi")
	}

	return &pb.DeleteLifestyleDataResponse{Success: true}, nil
}

// AddWearableData yangi kiyiladigan qurilma ma'lumotlarini qo'shish uchun
// func (h *Health) ConsumeWearableDataQueue() {
// 	// RabbitMQ queue’dan xabarlarni olish
// 	messages, err := h.RabbitMQChannel.Consume(
// 		"wearable_data_queue", // Queue nomi
// 		"",                    // Consumer tag
// 		true,                  // Auto-ack
// 		false,                 // Exclusive
// 		false,                 // No-local
// 		false,                 // No-wait
// 		nil,                   // Arguments
// 	)
// 	if err != nil {
// 		log.Fatal("failed to register a consumer: %w", err)
// 		return
// 	}

// 	for msg := range messages {
// 		var message struct {
// 			Id string `json:"id"`
// 			*pb.AddWearableDataRequest
// 		}
// 		err := json.Unmarshal(msg.Body, &message)
// 		if err != nil {
// 			h.Logger.Error("Failed to unmarshal message", "error", err)
// 			return
// 		}

// 		wearableData := bson.M{
// 			"id":                message.Id,
// 			"userid":            message.UserId,
// 			"devicetype":        message.DeviceType,
// 			"datatype":          message.DataType,
// 			"datavalue":         message.DataValue,
// 			"recordedtimestamp": message.RecordedTimestamp,
// 			"createdat":         time.Now().Format(time.RFC3339),
// 			"updatedat":         time.Now().Format(time.RFC3339),
// 			"deletedat":         "0",
// 		}
// 		fmt.Println(wearableData)

// 		_, err = h.Db.Collection("wearable_data").InsertOne(context.Background(), wearableData)
// 		if err != nil {
// 			h.Logger.Error("Failed to insert wearable data into MongoDB", "error", err)
// 			return
// 		}
// 	}
// }

// GetAllWearableData metodi
func (h *Health) GetAllWearableData(ctx context.Context, req *pb.GetAllWearableDataRequest) (*pb.GetAllWearableDataResponse, error) {
	collection := h.Db.Collection("wearable_data")

	skip := (req.GetPage() - 1) * req.GetLimit()

	findOptions := options.Find()
	findOptions.SetLimit(req.GetLimit())
	findOptions.SetSkip(skip)

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		h.Logger.Error("Error finding wearable data", "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var wearableDataList []*pb.WearableData
	for cursor.Next(ctx) {
		var data pb.WearableData
		if err := cursor.Decode(&data); err != nil {
			h.Logger.Error("Error decoding wearable data", "error", err)
			return nil, err
		}
		wearableDataList = append(wearableDataList, &data)
	}

	if err := cursor.Err(); err != nil {
		h.Logger.Error("Cursor error", "error", err)
		return nil, err
	}

	return &pb.GetAllWearableDataResponse{
		Wearabledata: wearableDataList,
	}, nil
}

// GetWearableData kiyiladigan qurilma ma'lumotlarini olish uchun
func (h *Health) GetWearableData(ctx context.Context, req *pb.GetWearableDataRequest) (*pb.GetWearableDataResponse, error) {
	var wearableData pb.WearableData

	err := h.Db.Collection("wearable_data").FindOne(ctx, bson.M{"$and": []bson.M{{"id": req.Id}, {"deletedat": "0"}}}).Decode(&wearableData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			h.Logger.Warn("Wearable data not found", "id", req.Id)
			return nil, errors.New("kiyiladigan qurilma ma'lumotlari topilmadi")
		}
		h.Logger.Error("Failed to get wearable data", "error", err)
		return nil, err
	}

	return &pb.GetWearableDataResponse{WearableData: &wearableData}, nil
}

// UpdateWearableData kiyiladigan qurilma ma'lumotlarini yangilash uchun
func (h *Health) UpdateWearableData(ctx context.Context, req *pb.UpdateWearableDataRequest) (*pb.UpdateWearableDataResponse, error) {
	update := bson.M{
		"$set": bson.M{
			"userid":            req.UserId,
			"devicetype":        req.DeviceType,
			"datatype":          req.DataType,
			"datavalue":         req.DataValue,
			"recordedtimestamp": req.RecordedTimestamp,
			"updatedat":         time.Now().Format(time.RFC3339),
		},
	}

	result, err := h.Db.Collection("wearable_data").UpdateOne(ctx, bson.M{"$and": []bson.M{{"id": req.Id}, {"deletedat": "0"}}}, update)
	if err != nil {
		h.Logger.Error("Failed to update wearable data", "error", err)
		return nil, err
	}

	if result.MatchedCount == 0 {
		h.Logger.Warn("Wearable data not found for update", "id", req.Id)
		return &pb.UpdateWearableDataResponse{Success: false}, errors.New("kiyiladigan qurilma ma'lumotlari topilmadi")
	}

	return &pb.UpdateWearableDataResponse{Success: true}, nil
}

// DeleteWearableData kiyiladigan qurilma ma'lumotlarini o'chirish uchun
func (h *Health) DeleteWearableData(ctx context.Context, req *pb.DeleteWearableDataRequest) (*pb.DeleteWearableDataResponse, error) {
	// Hozirgi vaqtni olish
	currentTime := time.Now().Unix()

	// Kiyiladigan qurilma ma'lumotlarini yangilash
	result, err := h.Db.Collection("wearable_data").UpdateOne(
		ctx,
		bson.M{"id": req.Id},
		bson.M{"$set": bson.M{"deletedat": currentTime}},
	)
	if err != nil {
		h.Logger.Error("Failed to delete wearable data", "error", err)
		return nil, err
	}

	// Agar hech qanday yozuv yangilanmagan bo'lsa
	if result.MatchedCount == 0 {
		h.Logger.Warn("Wearable data not found for deletion", "id", req.Id)
		return &pb.DeleteWearableDataResponse{Success: false}, errors.New("kiyiladigan qurilma ma'lumotlari topilmadi")
	}

	return &pb.DeleteWearableDataResponse{Success: true}, nil
}

// func (h *Health) GenerateHealthRecommendations(ctx context.Context, req *pb.GenerateHealthRecommendationsRequest) (*pb.GenerateHealthRecommendationsResponse, error) {
// 	coll := h.Db.Collection("health")

// 	id := uuid.NewString()
// 	date := time.Now().Format("2006/01/02")

// 	_, err := coll.InsertOne(ctx, bson.M{
// 		"id":                  id,
// 		"user_id":             req.UserId,
// 		"recommendation_type": req.RecommendationType,
// 		"description":         req.Description,
// 		"priority":            req.Priority,
// 		"created_at":          date,
// 		"updated_at":          date,
// 		"deleted_at":          "0",
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	recod := pb.HealthRecommendation{
// 		Id:                 id,
// 		UserId:             req.UserId,
// 		RecommendationType: req.RecommendationType,
// 		Description:        req.Description,
// 		Priority:           req.Priority,
// 		CreatedAt:          date,
// 		UpdatedAt:          date,
// 	}

// 	// Redisga yozish
// 	redisKey := req.UserId
// 	redisValue, err := json.Marshal(recod)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = h.Redis.Set(ctx, redisKey, redisValue, 0).Err()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &pb.GenerateHealthRecommendationsResponse{
// 		Recommendations: &recod,
// 	}, nil
// }

func (h *Health) GenerateHealthRecommendationsId(ctx context.Context, req *pb.GenerateHealthRecommendationsIdRequest) (*pb.GenerateHealthRecommendationsIdResponse, error) {
	collection := h.Db.Collection("health")

	filter := bson.M{"id": req.Id, "deleted_at": "0"}

	var recommendation pb.HealthRecommendation

	err := collection.FindOne(ctx, filter).Decode(&recommendation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			h.Logger.Error("Bunday id yoki deleted_at '0' bo'lgan hujjat topilmadi")
			return nil, fmt.Errorf(codes.NotFound.String(), "Bunday id yoki deleted_at '0' bo'lgan hujjat topilmadi")
		}
		h.Logger.Error("MongoDB dan hujjat olishda xatolik", "error", err)
		return nil, fmt.Errorf(codes.Internal.String(), "Hujjatni olishda xatolik: %v", err)
	}

	response := &pb.GenerateHealthRecommendationsIdResponse{
		Recommendations: &recommendation,
	}

	return response, nil
}

func (h *Health) GetRealtimeHealthMonitoring(ctx context.Context, req *pb.GetRealtimeHealthMonitoringRequest) (*pb.GetRealtimeHealthMonitoringResponse, error) {
	var user pb.HealthRecommendation

	redisKey := req.UserId

	val, err := h.Redis.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		// Redisda ma'lumot topilmadi
		return nil, fmt.Errorf("realtime health monitoring not found")
	} else if err != nil {
		// Redis bilan bog'liq xatolik
		return nil, err
	}

	// Redisdan olingan JSON ma'lumotni decoding qilish
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}

	// Bugungi sana
	date := time.Now().Format("2006/01/02")

	// Redisdan olingan ma'lumotdagi `CreatedAt` maydoni bilan bugungi sanani solishtirish
	if user.CreatedAt != date {
		return nil, fmt.Errorf("no health data found for today's date")
	}

	resp := pb.GetRealtimeHealthMonitoringResponse{
		RecommendationType: user.RecommendationType,
		Description:        user.Description,
		Priority:           user.Priority,
	}

	return &resp, nil
}

func (h *Health) GetDailyHealthSummary(ctx context.Context, req *pb.GetDailyHealthSummaryRequest) (*pb.GetDailyHealthSummaryResponse, error) {
	var summary pb.GetDailyHealthSummaryResponse
	coll := h.Db.Collection("health")

	err := coll.FindOne(ctx, bson.M{"$and": []bson.M{{"user_id": req.UserId}, {"deleted_at": "0"}, {"created_at": req.Date}}}).Decode(&summary)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func (h *Health) GetWeeklyHealthSummary(ctx context.Context, req *pb.GetWeeklyHealthSummaryRequest) (*pb.GetWeeklyHealthSummaryResponse, error) {
	var summary pb.GetWeeklyHealthSummaryResponse
	coll := h.Db.Collection("health")

	startDateStr := req.StartDate

	startDate, err := time.Parse("2006/01/02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing start date: %v", err)
	}

	weekAgo := startDate.AddDate(0, 0, -7)
	weekAgoStr := weekAgo.Format("2006/01/02")

	cursor, err := coll.Find(ctx, bson.M{
		"$and": []bson.M{
			{"user_id": req.UserId},
			{"deleted_at": "0"},
			{"created_at": bson.M{
				"$gte": weekAgoStr,
				"$lte": startDateStr,
			}},
		},
	})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc pb.HealthRecommendation
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("error decoding document: %v", err)
		}
		summary.Health = append(summary.Health, &doc)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return &summary, nil
}
