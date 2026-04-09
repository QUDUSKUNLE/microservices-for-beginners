package service

import (
	"authservice/internal/model"
	"authservice/internal/repository"
	"authservice/internal/telemetry"
	"authservice/pkg/hasher"
	"authservice/pkg/token"
	"context"
	"errors"
)

type AuthService struct {
	repo *repository.UserRepo
}

func NewAuthService(r *repository.UserRepo) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Register(ctx context.Context, email, password, address string) error {
	ctx, span := telemetry.Tracer().Start(ctx, "service.register")
	defer span.End()

	h, err := hasher.Hash(password)
	if err != nil {
		return err
	}
	_, err = s.repo.Create(ctx, &model.User{
		Email:        email,
		PasswordHash: h,
		Address:      address,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {

	ctx, span := telemetry.Tracer().Start(ctx, "service.login")
	defer span.End()

	user, err := s.repo.Get(ctx,email)
	if err != nil {
		return "", errors.New("cant find user" + err.Error())
	}
	if err := hasher.Compare(user.PasswordHash, password); err != nil {
		return "", errors.New("invalid credentials" + err.Error())
	}
	return token.Generate(email)
}
