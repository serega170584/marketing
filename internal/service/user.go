package service

import (
	"context"
	"errors"
	"fmt"
	"marketing/internal/dto"
	"marketing/internal/repository"
)

type repo interface {
	GetUserByID(
		ctx context.Context,
		ID uint64,
	) (
		*dto.User,
		error,
	)
}

type UserService struct {
	repo repo
}

func NewUserService(repo repo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(
	ctx context.Context,
	id uint64,
) (
	*dto.User,
	error,
) {
	u, err := s.repo.GetUserByID(ctx, id)

	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"userService: failed get user by id %w",
			err,
		)
	}

	return u, nil
}
