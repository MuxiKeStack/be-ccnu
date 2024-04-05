package grpc

import (
	"context"
	ccnuv1 "github.com/MuxiKeStack/be-api/gen/proto/ccnu/v1"
	"github.com/MuxiKeStack/be-ccnu/domain"
	"github.com/MuxiKeStack/be-ccnu/service"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

type CCNUServiceServer struct {
	ccnuv1.UnimplementedCCNUServiceServer
	ccnu service.CCNUService
}

func NewCCNUServiceServer(ccnu service.CCNUService) *CCNUServiceServer {
	return &CCNUServiceServer{ccnu: ccnu}
}

func (s *CCNUServiceServer) Register(server grpc.ServiceRegistrar) {
	ccnuv1.RegisterCCNUServiceServer(server, s)
}

func (s *CCNUServiceServer) Login(ctx context.Context, request *ccnuv1.LoginRequest) (*ccnuv1.LoginResponse, error) {
	success, err := s.ccnu.Login(ctx, request.GetStudentId(), request.GetPassword())
	return &ccnuv1.LoginResponse{Success: success}, err
}

func (s *CCNUServiceServer) CourseList(ctx context.Context, request *ccnuv1.CourseListRequest) (*ccnuv1.CourseListResponse, error) {
	courses, err := s.ccnu.GetSelfCourseList(ctx, request.GetStudentId(), request.GetPassword(),
		request.GetYear(), request.GetTerm())
	return &ccnuv1.CourseListResponse{
		Courses: slice.Map(courses, func(idx int, src domain.Course) *ccnuv1.Course {
			return &ccnuv1.Course{
				CourseId: src.CourseId,
				Name:     src.Name,
				Teacher:  src.Teacher,
				Year:     src.Year,
				Term:     src.Term,
			}
		}),
	}, err
}
