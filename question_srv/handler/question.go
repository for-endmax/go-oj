package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"question_srv/global"
	"question_srv/model"
	"question_srv/proto"
)

type QuestionServer struct {
	proto.UnimplementedQuestionServer
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

// GetQuestionList 查看题目列表
func (s *QuestionServer) GetQuestionList(ctx context.Context, req *proto.PageInfoRequest) (*proto.QuestionListInfoResponse, error) {
	// 获取分页信息
	page, pageSize := int(req.PNum), int(req.PSize)
	// 查询信息
	var questions []model.Question
	var total int32
	result := global.DB.Scopes(Paginate(page, pageSize)).Find(&questions)
	total = int32(result.RowsAffected)

	var rsp proto.QuestionListInfoResponse
	rsp.Total = total
	for _, v := range questions {
		rsp.Data = append(rsp.Data, &proto.QuestionInfoResponse{
			Id:      v.ID,
			Seq:     v.Seq,
			Name:    v.Name,
			Content: v.Content,
		})
	}
	return &rsp, nil
}

// GetQuestionInfo 通过id获取题目信息
func (s *QuestionServer) GetQuestionInfo(ctx context.Context, req *proto.GetQuestionInfoRequest) (*proto.QuestionInfoResponse, error) {
	question := model.Question{}
	result := global.DB.First(&question, req.Id)
	if result.RowsAffected == 0 {
		zap.S().Infof("题目不存在  id=%d", req.Id)
		return nil, status.Errorf(codes.NotFound, "题目不存在")
	}

	return &proto.QuestionInfoResponse{
		Id:      question.ID,
		Seq:     question.Seq,
		Name:    question.Name,
		Content: question.Content,
	}, nil
}

// AddQuestion 增加题目
func (s *QuestionServer) AddQuestion(ctx context.Context, req *proto.AddQuestionRequest) (*proto.QuestionInfoResponse, error) {
	question := model.Question{Seq: req.Seq}
	result := global.DB.Where(&question).First(&model.Question{})
	if result.RowsAffected != 0 {
		zap.S().Infof("题目已存在  seq=%d", req.Seq)
		return nil, status.Errorf(codes.AlreadyExists, "题目已存在")
	}
	result = global.DB.Create(&model.Question{
		Seq:     req.Seq,
		Name:    req.Name,
		Content: req.Content,
	})
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建题目出错")
	}
	global.DB.Where(&model.Question{Seq: req.Seq}).First(&question)
	return &proto.QuestionInfoResponse{
		Id:      question.ID,
		Seq:     question.Seq,
		Name:    question.Name,
		Content: question.Content,
	}, nil
}

// DelQuestion 删除题目
func (s *QuestionServer) DelQuestion(ctx context.Context, req *proto.DelQuestionRequest) (*proto.DelQuestionResponse, error) {
	//查找题目是否存在
	result := global.DB.First(&model.Question{}, req.Id)

	if result.RowsAffected != 1 {
		return &proto.DelQuestionResponse{Success: false}, status.Errorf(codes.Internal, "题目不存在")
	}

	// 删除题目
	result = global.DB.Delete(&model.Question{}, req.Id)
	if result.Error != nil || result.RowsAffected == 0 {
		return &proto.DelQuestionResponse{Success: false}, status.Errorf(codes.Internal, "删除题目时出错")
	}

	return &proto.DelQuestionResponse{Success: true}, nil
}

// UpdateQuestion 修改题目
func (s *QuestionServer) UpdateQuestion(ctx context.Context, req *proto.UpdateQuestionRequest) (*proto.UpdateQuestionResponse, error) {
	//查找题目是否存在
	result := global.DB.First(&model.Question{}, req.Id)

	if result.RowsAffected != 1 {
		return &proto.UpdateQuestionResponse{Success: false}, status.Errorf(codes.Internal, "题目不存在")
	}

	// 修改题目
	result = global.DB.Save(&model.Question{
		Seq:     req.Seq,
		Name:    req.Name,
		Content: req.Content,
	})
	if result.Error != nil || result.RowsAffected == 0 {
		return &proto.UpdateQuestionResponse{Success: false}, status.Errorf(codes.Internal, "删除题目时出错")
	}

	return &proto.UpdateQuestionResponse{Success: true}, nil
}

//
//rpc GetTestInfo(GetTestRequest) returns (TestInfoResponse); //获取测试信息
//rpc AddTest(AddTestRequest) returns (TestInfoResponse);//增加测试信息
//rpc DelTest(DelTestRequest) returns (DelTestResponse);// 删除测试信息
//rpc UpdateTest(UpdateTestRequest) returns (UpdateTestResponse);//修改测试信息
