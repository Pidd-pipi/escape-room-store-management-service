package service

import (
	"context"
	"errors"
	"testing"

	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/util"
)

func TestRegisterNewPhoneShouldSucceed(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newTestUserService(repo)
	u, err := svc.Register(context.Background(), &dto.RegisterRequest{Phone: "13800003333", Password: "123456", Nickname: "新玩家"})
	if err != nil {
		t.Fatalf("new phone register should succeed, got err: %v", err)
	}
	if u == nil || u.Phone != "13800003333" {
		t.Fatalf("unexpected user: %+v", u)
	}
}

func TestLoginUnknownPhoneUnauthorized(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newTestUserService(repo)
	_, err := svc.Login(context.Background(), &dto.LoginRequest{Phone: "13800009999", Password: "123456"})
	if err == nil {
		t.Fatal("unknown phone login should fail")
	}
	var appErr *util.AppError
	if !errors.As(err, &appErr) || appErr.Code != constants.CodeUnauthorized {
		t.Fatalf("expected CodeUnauthorized, got %v", err)
	}
}
