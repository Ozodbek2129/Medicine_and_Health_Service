package service

import (
	"context"
	"fmt"
	pb "health/genproto/health_analytics"
	mongoDb "health/mongodb"
	logger "health/pkg"
	"log/slog"
)

type HealthService struct {
	pb.UnimplementedHealthAnalyticsServiceServer
	health *mongoDb.Health
	log    *slog.Logger
}

func NewHealthService(health *mongoDb.Health) *HealthService {
	return &HealthService{
		health: health,
		log:    logger.NewLogger(),
	}
}

func (s *HealthService) AddMedicalRecord(ctx context.Context,req *pb.AddMedicalRecordRequest)(*pb.AddMedicalRecordResponse,error){
	resp,err:=s.health.AddMedicalRecord(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("AddMedicalRecord serviceda xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) GetMedicalRecord(ctx context.Context,req *pb.GetMedicalRecordRequest)(*pb.GetMedicalRecordResponse,error){
	resp,err:=s.health.GetMedicalRecord(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("GetMedicalRecord service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) UpdateMedicalRecord(ctx context.Context,req *pb.UpdateMedicalRecordRequest)(*pb.UpdateMedicalRecordResponse,error){
	resp,err:=s.health.UpdateMedicalRecord(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("UpdateMedicalRecord service xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) DeleteMedicalRecord(ctx context.Context,req *pb.DeleteMedicalRecordRequest)(*pb.DeleteMedicalRecordResponse,error){
	resp,err:=s.health.DeleteMedicalRecord(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("DeleteMedicalRecord service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) ListMedicalRecords(ctx context.Context,req *pb.ListMedicalRecordsRequest)(*pb.ListMedicalRecordsResponse,error){
	resp,err:=s.health.ListMedicalRecords(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("ListMedicalRecords service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}