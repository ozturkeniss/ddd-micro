package application

import (
	"time"

	"github.com/ddd-micro/internal/user/application/command"
	"github.com/ddd-micro/internal/user/application/query"
	"github.com/ddd-micro/internal/user/domain"
	"github.com/google/wire"
)

// ProviderSet is the Wire provider set for application layer
var ProviderSet = wire.NewSet(
	// Command handlers
	command.NewCreateUserHandler,
	command.NewUpdateUserHandler,
	command.NewUpdateUserByAdminHandler,
	command.NewDeleteUserHandler,
	command.NewChangePasswordHandler,
	command.NewAssignRoleHandler,
	ProvideLoginHandler,

	// Query handlers
	query.NewGetUserByIDHandler,
	query.NewGetUserByEmailHandler,
	query.NewListUsersHandler,

	// Service
	ProvideUserServiceCQRS,
	ProvideJWTHelper,
	ProvideTokenDuration,
)

// ProvideLoginHandler provides login command handler
func ProvideLoginHandler(repo domain.UserRepository, jwtSecret string, tokenDuration time.Duration) *command.LoginHandler {
	return command.NewLoginHandler(repo, jwtSecret, tokenDuration)
}

// ProvideUserServiceCQRS provides CQRS-based user service
func ProvideUserServiceCQRS(
	createUserHandler *command.CreateUserHandler,
	updateUserHandler *command.UpdateUserHandler,
	updateUserByAdminHandler *command.UpdateUserByAdminHandler,
	deleteUserHandler *command.DeleteUserHandler,
	changePasswordHandler *command.ChangePasswordHandler,
	assignRoleHandler *command.AssignRoleHandler,
	loginHandler *command.LoginHandler,
	getUserByIDHandler *query.GetUserByIDHandler,
	getUserByEmailHandler *query.GetUserByEmailHandler,
	listUsersHandler *query.ListUsersHandler,
	jwtHelper *JWTHelper,
) *UserServiceCQRS {
	return &UserServiceCQRS{
		createUserHandler:        createUserHandler,
		updateUserHandler:        updateUserHandler,
		updateUserByAdminHandler: updateUserByAdminHandler,
		deleteUserHandler:        deleteUserHandler,
		changePasswordHandler:    changePasswordHandler,
		assignRoleHandler:        assignRoleHandler,
		loginHandler:             loginHandler,
		getUserByIDHandler:       getUserByIDHandler,
		getUserByEmailHandler:    getUserByEmailHandler,
		listUsersHandler:         listUsersHandler,
		jwtHelper:                jwtHelper,
	}
}

// ProvideJWTHelper provides JWT helper
func ProvideJWTHelper(jwtSecret string, tokenDuration time.Duration) *JWTHelper {
	return NewJWTHelper(jwtSecret, tokenDuration)
}

// ProvideTokenDuration provides token duration (24 hours)
func ProvideTokenDuration() time.Duration {
	return 24 * time.Hour
}
