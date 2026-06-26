package user

import (
	"errors"
	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Service interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) Service {
	return &service{repo: repo, jwtService: jwtService}
}

func (s *service) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// default role
	role := req.Role
	if role == "" {
		role = "driver"
	}

	user := &User{
		Name:  req.Name,
		Email: req.Email,
		Role:  role,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func (s *service) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := user.CheckPassword(req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
