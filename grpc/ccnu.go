package grpc

import (
	"context"
	"github.com/MuxiKeStack/be-api/gen/proto/ccnu"
	"github.com/MuxiKeStack/be-ccnu/service"
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
	//TODO implement me
	panic("implement me")
}
