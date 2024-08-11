package mongoDb

import (
	"context"
	"errors"
	"fmt"
	"time"

	logger "health/pkg"
	"log/slog"

	"github.com/google/uuid"

	pb "health/genproto/health_analytics"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Health struct {
	Logger *slog.Logger
	Db     *mongo.Database
}

func NewHealth(mdb *mongo.Database) *Health {
	return &Health{
		Logger: logger.NewLogger(),
		Db:     mdb,
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
	err := h.Db.Collection("medical_records").FindOne(ctx, bson.M{"id": req.Id}).Decode(&record)
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

	result, err := h.Db.Collection("medical_records").UpdateOne(ctx, bson.M{"id": req.Id}, update)
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
	result, err := h.Db.Collection("medical_records").DeleteOne(ctx, bson.M{"id": req.Id})
	if err != nil {
		h.Logger.Error("Failed to delete medical record", "error", err)
		return nil, err
	}

	if result.DeletedCount == 0 {
		h.Logger.Warn("Medical record not found for deletion", "record_id", req.Id)
		return &pb.DeleteMedicalRecordResponse{Success: false}, errors.New("tibbiy yozuv topilmadi")
	}

	return &pb.DeleteMedicalRecordResponse{Success: true}, nil
}

func (h *Health) ListMedicalRecords(ctx context.Context, req *pb.ListMedicalRecordsRequest) (*pb.ListMedicalRecordsResponse, error) {
	cursor, err := h.Db.Collection("medical_records").Find(ctx, bson.M{"user_id": req.UserId}, options.Find())
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
	lifestyleData := &pb.LifestyleData{
		Id:           uuid.NewString(),
		UserId:       req.UserId,
		DataType:     req.DataType,
		DataValue:    req.DataValue,
		RecordedDate: req.RecordedDate,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	_, err := h.Db.Collection("lifestyle_data").InsertOne(ctx, lifestyleData)
	if err != nil {
		h.Logger.Error("Failed to add lifestyle data", "error", err)
		return nil, err
	}

	return &pb.AddLifestyleDataResponse{LifestyleData: lifestyleData}, nil
}

// GetLifestyleData turmush tarzi ma'lumotlarini olish uchun
func (h *Health) GetLifestyleData(ctx context.Context, req *pb.GetLifestyleDataRequest) (*pb.GetLifestyleDataResponse, error) {
	var lifestyleData pb.LifestyleData

	err := h.Db.Collection("lifestyle_data").FindOne(ctx, bson.M{"id": req.Id}).Decode(&lifestyleData)
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

	result, err := h.Db.Collection("lifestyle_data").UpdateOne(ctx, bson.M{"id": req.Id}, update)
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
	result, err := h.Db.Collection("lifestyle_data").DeleteOne(ctx, bson.M{"id": req.Id})
	if err != nil {
		h.Logger.Error("Failed to delete lifestyle data", "error", err)
		return nil, err
	}

	if result.DeletedCount == 0 {
		h.Logger.Warn("Lifestyle data not found for deletion", "id", req.Id)
		return &pb.DeleteLifestyleDataResponse{Success: false}, errors.New("turmush tarzi ma'lumotlari topilmadi")
	}

	return &pb.DeleteLifestyleDataResponse{Success: true}, nil
}

// AddWearableData yangi kiyiladigan qurilma ma'lumotlarini qo'shish uchun
func (h *Health) AddWearableData(ctx context.Context, req *pb.AddWearableDataRequest) (*pb.AddWearableDataResponse, error) {
	wearableData := &pb.WearableData{
		Id:                uuid.NewString(),
		UserId:            req.UserId,
		DeviceType:        req.DeviceType,
		DataType:          req.DataType,
		DataValue:         req.DataValue,
		RecordedTimestamp: req.RecordedTimestamp,
		CreatedAt:         time.Now().Format(time.RFC3339),
		UpdatedAt:         time.Now().Format(time.RFC3339),
	}

	_, err := h.Db.Collection("wearable_data").InsertOne(ctx, wearableData)
	if err != nil {
		h.Logger.Error("Failed to add wearable data", "error", err)
		return nil, err
	}

	return &pb.AddWearableDataResponse{WearableData: wearableData}, nil
}

// GetWearableData kiyiladigan qurilma ma'lumotlarini olish uchun
func (h *Health) GetWearableData(ctx context.Context, req *pb.GetWearableDataRequest) (*pb.GetWearableDataResponse, error) {
	var wearableData pb.WearableData

	err := h.Db.Collection("wearable_data").FindOne(ctx, bson.M{"id": req.Id}).Decode(&wearableData)
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

	result, err := h.Db.Collection("wearable_data").UpdateOne(ctx, bson.M{"id": req.Id}, update)
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
	result, err := h.Db.Collection("wearable_data").DeleteOne(ctx, bson.M{"id": req.Id})
	if err != nil {
		h.Logger.Error("Failed to delete wearable data", "error", err)
		return nil, err
	}

	if result.DeletedCount == 0 {
		h.Logger.Warn("Wearable data not found for deletion", "id", req.Id)
		return &pb.DeleteWearableDataResponse{Success: false}, errors.New("kiyiladigan qurilma ma'lumotlari topilmadi")
	}

	return &pb.DeleteWearableDataResponse{Success: true}, nil
}

func (h *Health) GenerateHealthRecommendations(ctx context.Context, req *pb.GenerateHealthRecommendationsRequest) (*pb.GenerateHealthRecommendationsResponse, error) {
	coll := h.Db.Collection("health")

	id := uuid.NewString()
	date := time.Now().Format("2006/01/02")

	_, err := coll.InsertOne(ctx, bson.M{
		"id":                  id,
		"user_id":             req.UserId,
		"recommendation_type": req.RecommendationType,
		"description":         req.Description,
		"priority":            req.Priority,
		"created_at":          date,
		"updated_at":          date,
		"deleted_at":          0,
	})

	if err != nil {
		return nil, err
	}

	recod := pb.HealthRecommendation{
		Id:                 id,
		UserId:             req.UserId,
		RecommendationType: req.RecommendationType,
		Description:        req.Description,
		Priority:           req.Priority,
		CreatedAt:          date,
		UpdatedAt:          date,
	}

	return &pb.GenerateHealthRecommendationsResponse{
		Recommendations: &recod,
	}, nil
}

func (h *Health) GetRealtimeHealthMonitoring(ctx context.Context, req *pb.GetRealtimeHealthMonitoringRequest) (*pb.GetRealtimeHealthMonitoringResponse, error) {
	var user pb.GetRealtimeHealthMonitoringResponse
	coll := h.Db.Collection("health")

	err := coll.FindOne(ctx, bson.M{"$and": []bson.M{{"user_id": req.UserId}, {"deleted_at": 0}, {"created_at": time.Now().Format("2006/01/02")}}}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("realtime health monitoring not found")
	}

	return &user, nil
}

func (h *Health) GetDailyHealthSummary(ctx context.Context, req *pb.GetDailyHealthSummaryRequest) (*pb.GetDailyHealthSummaryResponse, error) {
	var summary pb.GetDailyHealthSummaryResponse
	coll := h.Db.Collection("health")

	err := coll.FindOne(ctx, bson.M{"$and": []bson.M{{"user_id": req.UserId}, {"deleted_at": 0}, {"created_at": req.Date}}}).Decode(&summary)
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
			{"deleted_at": 0},
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
