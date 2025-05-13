package user

import (
	"context"
	"github.com/ALexfonSchneider/food-delivery-user-service/gen/grpc/go/user"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/application/auth"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (u Server) RegisterUser(ctx context.Context, request *user.RegisterUserRequest) (*user.RegisterUserResponse, error) {
	var lastName *string
	if request.LastName != nil {
		lastName = &request.LastName.Value
	}

	registerUserRequest := auth.RegisterUserRequest{
		Email:     request.GetEmail(),
		FirstName: request.GetFirstName(),
		LastName:  lastName,
		Password:  request.Password,
		Phone:     request.GetPhone(),
		CreatedAt: request.GetCreatedAt().AsTime(),
	}
	usr, err := u.auth.RegisterUser(ctx, &registerUserRequest)
	if err != nil {
		return nil, err
	}

	return &user.RegisterUserResponse{Id: usr.Id}, nil
}

func (u Server) LoginUser(ctx context.Context, request *user.LoginRequest) (*user.LoginResponse, error) {
	tokens, err := u.auth.LoginUser(ctx, &auth.LoginUserRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (u Server) GetProfile(ctx context.Context, request *user.GetProfileRequest) (*user.GetProfileResponse, error) {
	usr, err := u.user.GetUserById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, domain.UserNotFoundError
	}

	var LastName wrapperspb.StringValue
	if usr.LastName != nil {
		LastName = wrapperspb.StringValue{Value: *usr.LastName}
	}

	return &user.GetProfileResponse{
		Id:        usr.Id,
		FirstName: usr.FirstName,
		LastName:  &LastName,
		Email:     usr.Email,
		Phone:     usr.Phone,
		CreatedAt: &timestamppb.Timestamp{Nanos: int32(usr.CreatedAt.Nanosecond()), Seconds: int64(usr.CreatedAt.Second())},
	}, nil
}
