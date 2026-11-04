package usecase

import (
	"context"
	"errors"
	"fmt"
	"marketing/internal/dto"
	"marketing/internal/service"
)

type userService interface {
	GetUserByID(
		ctx context.Context,
		ID uint64,
	) (
		*dto.User,
		error,
	)
}

type UserUseCase struct {
	UserService userService
}

func NewUserUseCase(
	userService userService,
) *UserUseCase {
	return &UserUseCase{
		UserService: userService,
	}
}

func (uuc *UserUseCase) GetUserByID(
	ctx context.Context,
	id uint64,
) (
	*dto.User,
	error,
) {
	user, err := uuc.UserService.GetUserByID(ctx, id)

	if errors.Is(err, service.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("userUseCase: failed to get user by id: %w", err)
	}

	return user, nil
}
