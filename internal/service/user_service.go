// Package service implements the business logic for escape-room-ops.
package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/model"
	"github.com/lp/escape-room-ops/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository is the data access contract for user rows.
type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	UpdateProfile(ctx context.Context, id uint, nickname string) error
	Count(ctx context.Context) (int64, error)
}

// UserService coordinates player registration, login and profile flows.
type UserService struct {
	users     UserRepository
	jwtSecret string
	jwtExpire int
	logger    *slog.Logger
}

// NewUserService wires the user service dependencies.
func NewUserService(users UserRepository, jwtSecret string, jwtExpire int, logger *slog.Logger) *UserService {
	return &UserService{users: users, jwtSecret: jwtSecret, jwtExpire: jwtExpire, logger: logger}
}

// Register creates a new player account.
func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*model.User, error) {
	if _, err := s.users.FindByPhone(ctx, req.Phone); err == nil {
		return nil, util.NewAppError(409, constants.CodeConflict, constants.MsgPhoneAlreadyUsed, nil)
	} else if !errors.Is(err, util.ErrNotFound) {
		return nil, util.WrapAppError(fmt.Errorf("user register lookup: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("user register hash: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	user := &model.User{
		Phone: req.Phone, PasswordHash: string(hash), Nickname: req.Nickname,
		Role: constants.UserRolePlayer,
	}
	if err := s.users.Create(ctx, user); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogUserRegisterFailed, req.Phone, err))
		return nil, util.WrapAppError(fmt.Errorf("user register create: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogUserRegisterSuccess, req.Phone, user.ID))
	return user, nil
}

// Login verifies credentials and returns a JWT plus profile.
func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.users.FindByPhone(ctx, req.Phone)
	if err != nil {
		s.logger.Info(fmt.Sprintf(constants.LogUserLoginFailed, req.Phone, "user not found"))
		return nil, util.NewAppError(401, constants.CodeUnauthorized, constants.MsgPhoneOrPassword, nil)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		s.logger.Info(fmt.Sprintf(constants.LogUserLoginFailed, req.Phone, "password mismatch"))
		return nil, util.NewAppError(401, constants.CodeUnauthorized, constants.MsgPhoneOrPassword, nil)
	}
	token, err := util.GenerateToken(s.jwtSecret, user.ID, user.Phone, user.Role, s.jwtExpire)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("user login token: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogUserLoginSuccess, req.Phone, user.ID, user.Role))
	return &dto.LoginResponse{Token: token, User: ToUserView(user)}, nil
}

// GetProfile returns the current user profile.
func (s *UserService) GetProfile(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("user[id=%d] get profile: %w", userID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	return user, nil
}

// UpdateProfile updates the nickname.
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*model.User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("user[id=%d] update find: %w", userID, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	nickname := user.Nickname
	if req.Nickname != "" {
		nickname = req.Nickname
	}
	if err := s.users.UpdateProfile(ctx, userID, nickname); err != nil {
		return nil, util.WrapAppError(fmt.Errorf("user[id=%d] update profile: %w", userID, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogUserProfileUpdateSuccess, userID))
	user.Nickname = nickname
	return user, nil
}

// ToUserView converts a user model into the public view.
func ToUserView(u *model.User) *dto.UserView {
	if u == nil {
		return nil
	}
	return &dto.UserView{ID: u.ID, Phone: u.Phone, Nickname: u.Nickname, Role: u.Role}
}
