package tests

import (
	"context"
	"testing"

	v1 "github.com/Nikita-Filonov/api-go-autotests-server/gen/go/v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestGRPCCreateUser(t *testing.T) {
	request := &v1.CreateUserRequest{
		Email:      gofakeit.Email(),
		Password:   gofakeit.Password(true, true, true, true, false, 12),
		LastName:   gofakeit.LastName(),
		FirstName:  gofakeit.FirstName(),
		MiddleName: gofakeit.FirstName(),
	}

	connection, err := grpc.NewClient(
		"localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer connection.Close()

	usersClient := v1.NewUsersServiceClient(connection)

	response, err := usersClient.CreateUser(
		context.Background(),
		request,
	)

	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.GetUser())
	require.NotEmpty(t, response.GetUser().GetId())

	assert.Equal(t, request.GetEmail(), response.GetUser().GetEmail())
	assert.Equal(t, request.GetLastName(), response.GetUser().GetLastName())
	assert.Equal(t, request.GetFirstName(), response.GetUser().GetFirstName())
	assert.Equal(t, request.GetMiddleName(), response.GetUser().GetMiddleName())

	t.Logf(
		"created user with ID %s",
		response.GetUser().GetId(),
	)
}

func TestGRPCUpdateUser(t *testing.T) {
	userRequest := &v1.CreateUserRequest{
		Email:      gofakeit.Email(),
		Password:   gofakeit.Password(true, true, true, true, false, 12),
		LastName:   gofakeit.LastName(),
		FirstName:  gofakeit.FirstName(),
		MiddleName: gofakeit.FirstName(),
	}

	connection, err := grpc.NewClient(
		"localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer connection.Close()

	usersClient := v1.NewUsersServiceClient(connection)

	createdUser, err := usersClient.CreateUser(
		context.Background(),
		userRequest,
	)

	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotNil(t, createdUser.GetUser())
	require.NotEmpty(t, createdUser.GetUser().GetId())

	authenticationClient :=
		v1.NewAuthenticationServiceClient(connection)

	authRequest := &v1.LoginRequest{
		Email:    userRequest.GetEmail(),
		Password: userRequest.GetPassword(),
	}

	loginResponse, err := authenticationClient.Login(
		context.Background(),
		authRequest,
	)

	require.NoError(t, err)
	require.NotNil(t, loginResponse)
	require.NotNil(t, loginResponse.GetToken())
	require.NotEmpty(
		t,
		loginResponse.GetToken().GetAccessToken(),
	)

	email := gofakeit.Email()
	lastName := gofakeit.LastName()
	firstName := gofakeit.FirstName()
	middleName := gofakeit.FirstName()

	request := &v1.UpdateUserRequest{
		Id:         createdUser.GetUser().GetId(),
		Email:      &email,
		LastName:   &lastName,
		FirstName:  &firstName,
		MiddleName: &middleName,
	}

	authorizationMetadata := metadata.Pairs(
		"authorization",
		"Bearer "+loginResponse.GetToken().GetAccessToken(),
	)

	contextWithToken := metadata.NewOutgoingContext(
		context.Background(),
		authorizationMetadata,
	)

	response, err := usersClient.UpdateUser(
		contextWithToken,
		request,
	)

	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.GetUser())

	assert.Equal(t, createdUser.GetUser().GetId(), response.GetUser().GetId())
	assert.Equal(t, request.GetEmail(), response.GetUser().GetEmail())
	assert.Equal(t, request.GetLastName(), response.GetUser().GetLastName())
	assert.Equal(t, request.GetFirstName(), response.GetUser().GetFirstName())
	assert.Equal(t, request.GetMiddleName(), response.GetUser().GetMiddleName())

	t.Logf("updated user with ID %s", response.GetUser().GetId())
}

func TestGRPCGetUserMe(t *testing.T) {
	userRequest := &v1.CreateUserRequest{
		Email:      gofakeit.Email(),
		Password:   gofakeit.Password(true, true, true, true, false, 12),
		LastName:   gofakeit.LastName(),
		FirstName:  gofakeit.FirstName(),
		MiddleName: gofakeit.FirstName(),
	}

	connection, err := grpc.NewClient(
		"localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer connection.Close()

	usersClient := v1.NewUsersServiceClient(connection)
	authenticationClient := v1.NewAuthenticationServiceClient(connection)

	createdUser, err := usersClient.CreateUser(
		context.Background(),
		userRequest,
	)

	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotNil(t, createdUser.GetUser())
	require.NotEmpty(t, createdUser.GetUser().GetId())

	authRequest := &v1.LoginRequest{
		Email:    userRequest.GetEmail(),
		Password: userRequest.GetPassword(),
	}

	loginResponse, err := authenticationClient.Login(
		context.Background(),
		authRequest,
	)

	require.NoError(t, err)
	require.NotNil(t, loginResponse)
	require.NotNil(t, loginResponse.GetToken())
	require.NotEmpty(t, loginResponse.GetToken().GetAccessToken())

	authorizationMetadata := metadata.Pairs(
		"authorization",
		"Bearer "+loginResponse.GetToken().GetAccessToken(),
	)

	contextWithToken := metadata.NewOutgoingContext(
		context.Background(),
		authorizationMetadata,
	)

	response, err := usersClient.GetMe(
		contextWithToken,
		&v1.Empty{},
	)

	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.GetUser())

	assert.Equal(t, createdUser.GetUser().GetId(), response.GetUser().GetId())
	assert.Equal(t, createdUser.GetUser().GetEmail(), response.GetUser().GetEmail())
	assert.Equal(t, createdUser.GetUser().GetLastName(), response.GetUser().GetLastName())
	assert.Equal(t, createdUser.GetUser().GetFirstName(), response.GetUser().GetFirstName())
	assert.Equal(t, createdUser.GetUser().GetMiddleName(), response.GetUser().GetMiddleName())
}
