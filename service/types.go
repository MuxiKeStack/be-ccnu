package service

import (
	"context"
	"github.com/MuxiKeStack/be-ccnu/domain"
	"time"
)

type CCNUService interface {
	Login(ctx context.Context, studentId string, password string) (bool, error)
	GetSelfCourseList(ctx context.Context, studentId, password, year, term string) ([]domain.Course, error)
	GetSelfGradeList(ctx context.Context, studentId, password, year, term string) ([]domain.Grade, error)
}

type ccnuService struct {
	timeout time.Duration
}

func NewCCNUService() CCNUService {
	return &ccnuService{timeout: time.Second * 5}
}
