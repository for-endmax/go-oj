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

//////////////////////////////////////////////////////////////////////////////////////

// GetTestInfo 获取某个题目的所有测试信息
func (s *QuestionServer) GetTestInfo(ctx context.Context, req *proto.GetTestRequest) (*proto.TestInfoListResponse, error) {
	var tests []model.Test
	result := global.DB.Where(&model.Test{QID: req.QId}).Find(&tests)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "测试信息不存在")
	}
	var rsp proto.TestInfoListResponse
	rsp.Total = int32(result.RowsAffected)
	for _, v := range tests {
		rsp.Data = append(rsp.Data, &proto.TestInfoResponse{
			Id:        v.ID,
			QId:       v.QID,
			TimeLimit: v.TimeLimit,
			MemLimit:  v.MemLimit,
			Input:     v.Input,
			Output:    v.Output,
		})
	}
	return &rsp, nil
}

// AddTest 增加测试信息
func (s *QuestionServer) AddTest(ctx context.Context, req *proto.AddTestRequest) (*proto.TestInfoResponse, error) {
	result := global.DB.Where("q_id = ?", req.QId).First(&model.Test{})
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "对应的题目已有测试信息")
	}

	test := model.Test{
		QID:       req.QId,
		TimeLimit: req.TimeLimit,
		MemLimit:  req.MemLimit,
		Input:     req.Input,
		Output:    req.Output,
	}
	result = global.DB.Create(&test)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "添加失败")
	}
	return &proto.TestInfoResponse{
		Id:        test.ID,
		QId:       test.QID,
		TimeLimit: test.TimeLimit,
		MemLimit:  test.MemLimit,
		Input:     test.Input,
		Output:    test.Output,
	}, nil
}

// DelTest 删除测试信息
func (s *QuestionServer) DelTest(ctx context.Context, req *proto.DelTestRequest) (*proto.DelTestResponse, error) {
	result := global.DB.First(&model.Test{}, req.Id)
	if result.RowsAffected == 0 {
		return &proto.DelTestResponse{Success: false}, status.Errorf(codes.NotFound, "测试信息不存在")
	}
	result = global.DB.Delete(&model.Test{}, req.Id)
	if result.RowsAffected == 0 {
		return &proto.DelTestResponse{Success: false}, status.Errorf(codes.Internal, "删除测试信息出错")
	}
	return &proto.DelTestResponse{Success: true}, nil
}

// UpdateTest 修改测试信息
func (s *QuestionServer) UpdateTest(ctx context.Context, req *proto.UpdateTestRequest) (*proto.UpdateTestResponse, error) {
	var test model.Test
	result := global.DB.First(&test, req.Id)
	if result.RowsAffected == 0 {
		return &proto.UpdateTestResponse{Success: false}, status.Errorf(codes.NotFound, "测试信息不存在")
	}

	test.QID = req.QId
	test.TimeLimit = req.TimeLimit
	test.MemLimit = req.MemLimit
	test.Input = req.Input
	test.Output = req.Output
	result = global.DB.Updates(&test)
	if result.RowsAffected == 0 {
		return &proto.UpdateTestResponse{Success: false}, status.Errorf(codes.Internal, "测试信息更新失败")
	}
	return &proto.UpdateTestResponse{
		Success: true,
	}, nil
}
