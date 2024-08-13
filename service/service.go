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

func (s *HealthService) AddLifestyleData(ctx context.Context,req *pb.AddLifestyleDataRequest)(*pb.AddLifestyleDataResponse,error){
	resp,err:=s.health.AddLifestyleData(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("AddLifestyleData service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) GetAllLifestyleData(ctx context.Context,req *pb.GetAllLifestyleDataRequest)(*pb.GetAllLifestyleDataResponse,error){
	resp,err:=s.health.GetAllLifestyleData(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("GetAllLifestyleData service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) GetLifestyleData(ctx context.Context,req *pb.GetLifestyleDataRequest)(*pb.GetLifestyleDataResponse,error){
	resp,err:=s.health.GetLifestyleData(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("GetLifestyleData service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) UpdateLifestyleData(ctx context.Context,req *pb.UpdateLifestyleDataRequest)(*pb.UpdateLifestyleDataResponse,error){
	resp,err:=s.health.UpdateLifestyleData(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("UpdateLifestyleData service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) DeleteLifestyleData(ctx context.Context,req *pb.DeleteLifestyleDataRequest)(*pb.DeleteLifestyleDataResponse,error){
	resp,err:=s.health.DeleteLifestyleData(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("DeleteLifestyleData service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

// func (s *HealthService) AddWearableData(ctx context.Context,req *pb.AddWearableDataRequest)(*pb.AddWearableDataResponse,error){
// 	resp,err:=s.health.AddWearableData(ctx,req)
// 	if err!=nil{
// 		s.log.Error(fmt.Sprintf("AddWearableData service da xatolik: %v",err))
// 		return nil,err
// 	}
// 	return resp,nil
// }

func (s *HealthService) GetAllWearableData(ctx context.Context,req *pb.GetAllWearableDataRequest)(*pb.GetAllWearableDataResponse,error){
	resp,err:=s.health.GetAllWearableData(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("GetAllWearableData service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) GetWearableData(ctx context.Context,req *pb.GetWearableDataRequest)(*pb.GetWearableDataResponse,error){
	resp,err:=s.health.GetWearableData(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("GetWearableData service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) UpdateWearableData(ctx context.Context,req *pb.UpdateWearableDataRequest)(*pb.UpdateWearableDataResponse,error){
	resp,err:=s.health.UpdateWearableData(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("UpdateWearableData service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) DeleteWearableData(ctx context.Context,req *pb.DeleteWearableDataRequest)(*pb.DeleteWearableDataResponse,error){
	resp,err:=s.health.DeleteWearableData(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("DeleteWearableData service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

// func (s *HealthService) GenerateHealthRecommendations(ctx context.Context,req *pb.GenerateHealthRecommendationsRequest)(*pb.GenerateHealthRecommendationsResponse,error){
// 	resp,err:=s.health.GenerateHealthRecommendations(ctx,req)
// 	if err!=nil{
// 		s.log.Error(fmt.Sprintf("GenerateHealthRecommendations service da xatolik: %v",err))
// 		return nil,err
// 	}
// 	return resp,nil
// }

func (s *HealthService) GenerateHealthRecommendationsId(ctx context.Context,req *pb.GenerateHealthRecommendationsIdRequest)(*pb.GenerateHealthRecommendationsIdResponse,error){
	resp,err:=s.health.GenerateHealthRecommendationsId(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("GenerateHealthRecommendationsId service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) GetRealtimeHealthMonitoring(ctx context.Context,req *pb.GetRealtimeHealthMonitoringRequest)(*pb.GetRealtimeHealthMonitoringResponse,error){
	resp,err:=s.health.GetRealtimeHealthMonitoring(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("GetRealtimeHealthMonitoring service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) GetDailyHealthSummary(ctx context.Context,req *pb.GetDailyHealthSummaryRequest)(*pb.GetDailyHealthSummaryResponse,error){
	resp,err:=s.health.GetDailyHealthSummary(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("GetDailyHealthSummary service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (s *HealthService) GetWeeklyHealthSummary(ctx context.Context,req *pb.GetWeeklyHealthSummaryRequest)(*pb.GetWeeklyHealthSummaryResponse,error){
	resp,err:=s.health.GetWeeklyHealthSummary(ctx,req)
	if err!=nil{
		s.log.Error(fmt.Sprintf("GetWeeklyHealthSummary service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}