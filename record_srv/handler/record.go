package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"record_srv/global"
	"record_srv/model"
	"record_srv/proto"
)

type RecordServer struct {
	proto.UnimplementedRecordServer
}

// Paginate 将数据进行分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// CreateRecord 创建一个初始的记录
func (s *RecordServer) CreateRecord(ctx context.Context, req *proto.CreateRecordRequest) (*proto.RecordInfo, error) {
	record := model.RecordModel{
		UID:        req.UID,
		QID:        req.QID,
		Lang:       req.Lang,
		Status:     req.Status,
		ErrCode:    req.ErrCode,
		ErrMsg:     "",
		TimeLimit:  1000,
		MemLimit:   3000,
		SubmitCode: req.SubmitCode,
	}
	result := global.DB.Create(&record)
	if result.RowsAffected == 0 {
		zap.S().Info("创建新记录失败")
		return nil, status.Errorf(codes.Internal, "创建失败")
	}
	return &proto.RecordInfo{
		UID:       record.UID,
		QID:       record.QID,
		Lang:      record.Lang,
		Status:    record.Status,
		ErrCode:   record.ErrCode,
		ErrMsg:    record.ErrMsg,
		TimeLimit: record.TimeLimit,
		MemLimit:  record.MemLimit,
		ID:        record.ID,
	}, nil
}

// GetAllRecordByUID 查找指定uid的所有记录
func (s *RecordServer) GetAllRecordByUID(ctx context.Context, req *proto.UIDRequest) (*proto.RecordInfoList, error) {
	var records []model.RecordModel
	result := global.DB.Scopes(Paginate(int(req.PNum), int(req.PSize))).Where("u_id=?", req.Uid).Find(&records)
	if result.RowsAffected == 0 {
		zap.S().Info("该用户没有记录")
		return nil, status.Errorf(codes.NotFound, "该用户没有记录")
	}
	var rsp proto.RecordInfoList
	rsp.Total = int32(result.RowsAffected)
	for _, v := range records {
		recordInfo := proto.RecordInfo{
			UID:        v.UID,
			QID:        v.QID,
			Lang:       v.Lang,
			Status:     v.Status,
			ErrCode:    v.ErrCode,
			ErrMsg:     v.ErrMsg,
			TimeLimit:  v.TimeLimit,
			MemLimit:   v.MemLimit,
			ID:         v.ID,
			SubmitCode: v.SubmitCode,
		}
		rsp.Data = append(rsp.Data, &recordInfo)
	}
	return &rsp, nil
}

// GetRecordByID 通过记录id查询记录
func (s *RecordServer) GetRecordByID(ctx context.Context, req *proto.IDRequest) (*proto.RecordInfo, error) {
	var record model.RecordModel
	result := global.DB.First(&record, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有这个记录")
	}
	return &proto.RecordInfo{
		UID:        record.UID,
		QID:        record.QID,
		Lang:       record.Lang,
		Status:     record.Status,
		ErrCode:    record.ErrCode,
		ErrMsg:     record.ErrMsg,
		TimeLimit:  record.TimeLimit,
		MemLimit:   record.MemLimit,
		ID:         record.ID,
		SubmitCode: record.SubmitCode,
	}, nil
}

// UpdateRecord 更新记录状态信息
func (s *RecordServer) UpdateRecord(ctx context.Context, req *proto.UpdateRecordRequest) (*proto.RecordInfo, error) {
	var record model.RecordModel
	if result := global.DB.First(&record, req.ID); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有对应的记录")
	}
	record.Status = req.Status
	record.ErrCode = req.ErrCode
	record.ErrMsg = req.ErrMsg
	if result := global.DB.Save(&record); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "更新记录出错")
	}
	return &proto.RecordInfo{
		UID:        record.UID,
		QID:        record.QID,
		Lang:       record.Lang,
		Status:     record.Status,
		ErrCode:    record.ErrCode,
		ErrMsg:     record.ErrMsg,
		TimeLimit:  record.TimeLimit,
		MemLimit:   record.MemLimit,
		ID:         record.ID,
		SubmitCode: record.SubmitCode,
	}, nil
}
