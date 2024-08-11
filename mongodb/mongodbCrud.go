package mongoDb

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid" // UUID paketini import qilish
	logger "health/pkg"
	"log/slog"

	pb "health/genproto/health_analytics"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Health struktura MongoDB bilan ishlash uchun
type Health struct {
	Logger *slog.Logger
	Db     *mongo.Database
}

// NewHealth yangi Health strukturasini yaratish
func NewHealth(mdb *mongo.Database) *Health {
	return &Health{
		Logger: logger.NewLogger(),
		Db:     mdb,
	}
}

// AddMedicalRecord tibbiy yozuvni qo'shish funksiyasi
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

	// Response yaratish
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

// GetMedicalRecord tibbiy yozuvni olish funksiyasi
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

// UpdateMedicalRecord tibbiy yozuvni yangilash funksiyasi
func (h *Health) UpdateMedicalRecord(ctx context.Context, req *pb.UpdateMedicalRecordRequest) (*pb.UpdateMedicalRecordResponse, error) {
	req.MedicalRecord.UpdatedAt = time.Now().Format(time.RFC3339)

	update := bson.M{
		"$set": bson.M{
			"record_type": req.MedicalRecord.RecordType,
			"record_date": req.MedicalRecord.RecordDate,
			"description": req.MedicalRecord.Description,
			"doctor_id":   req.MedicalRecord.DoctorId,
			"attachments": req.MedicalRecord.Attachments,
			"updated_at":  req.MedicalRecord.UpdatedAt,
		},
	}

	result, err := h.Db.Collection("medical_records").UpdateOne(ctx, bson.M{"id": req.MedicalRecord.Id}, update)
	if err != nil {
		h.Logger.Error("Failed to update medical record", "error", err)
		return nil, err
	}

	if result.MatchedCount == 0 {
		h.Logger.Warn("Medical record not found for update", "record_id", req.MedicalRecord.Id)
		return &pb.UpdateMedicalRecordResponse{Success: false}, errors.New("tibbiy yozuv topilmadi")
	}

	return &pb.UpdateMedicalRecordResponse{Success: true}, nil
}

// DeleteMedicalRecord tibbiy yozuvni o'chirish funksiyasi
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

// ListMedicalRecords bir foydalanuvchiga tegishli barcha tibbiy yozuvlarni olish funksiyasi
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
