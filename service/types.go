package service

import (
	"context"
	"time"
)

type CCNUService interface {
	Login(ctx context.Context, studentId string, password string) (bool, error)
}

type ccnuService struct {
	timeout time.Duration
}

func NewCCNUService() CCNUService {
	return &ccnuService{timeout: time.Second * 5}
}
