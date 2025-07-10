package users

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lokeam/qko-beta/internal/appcontext"
	"github.com/lokeam/qko-beta/internal/interfaces"
	"github.com/lokeam/qko-beta/internal/models"
	"github.com/lokeam/qko-beta/internal/postgres"
)


type UserDbAdapter struct {
	client *postgres.PostgresClient
	db     *sqlx.DB
	logger interfaces.Logger
}

func NewUserDbAdapter(appContext *appcontext.AppContext) (*UserDbAdapter, error) {
	appContext.Logger.Debug("Creating UserDbAdapter", map[string]any{"appContext": appContext})

	client, err := postgres.NewPostgresClient(appContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres client: %w", err)
	}

	db, err := sqlx.Connect("pgx", appContext.Config.Postgres.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqlx connection: %w", err)
	}

	return &UserDbAdapter{
		client: client,
		db:     db,
		logger: appContext.Logger,
	}, nil
}

func (ua *UserDbAdapter) GetSingleUser(ctx context.Context, userID string) (models.User, error) {
	ua.logger.Debug("Getting single user", map[string]any{"userID": userID})

	var user models.User
	err := ua.db.GetContext(ctx, &user, GetUserQuery, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found: %w", err)
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	ua.logger.Debug("User retrieved successfully", map[string]any{"user": user})
	return user, nil
}

func (ua *UserDbAdapter) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	ua.logger.Debug("Creating user", map[string]any{"user": user})

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := ua.db.GetContext(
		ctx,
		&user,
		CreateUserQuery,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	ua.logger.Debug("User created successfully", map[string]any{
		"user_id":  user.ID,
		"email":    user.Email,
	})

	return user, nil
}

func (ua *UserDbAdapter) UpdateUserProfile(
	ctx context.Context,
	userID string,
	firstName string,
	lastName string,
) (models.User, error) {
	ua.logger.Debug("UpdateUserProfile", map[string]any{
		"user_id": userID,
		"first_name": firstName,
		"last_name": lastName,
	})

	var user models.User
	err := ua.db.GetContext(
		ctx,
		&user,
		UpdateUserProfileQuery,
		firstName,
		lastName,
		time.Now(),
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found: %w", err)
		}
		return models.User{}, fmt.Errorf("failed to update user profile: %w")
	}

	ua.logger.Info("User profile updated", map[string]any{
		"user_id": user.ID,
		"first_name": firstName,
		"last_name": lastName,
	})

	return user, nil
}

func (ua *UserDbAdapter) HasCompleteProfile(ctx context.Context, userID string) (bool, error) {
	ua.logger.Debug("HasCompleteProfile called", map[string]any{"user_id": userID})

	var firstName, lastName string
	err := ua.db.QueryRowContext(
		ctx,
		HasCompleteProfileQuery,
		userID,
	).Scan(&firstName, &lastName)
	if err != nil {
		if err == sql.ErrNoRows {
				return false, fmt.Errorf("user not found: %w", err)
		}
		return false, fmt.Errorf("failed to check user profile: %w", err)
	}

	hasCompleteProfile := firstName != "" && lastName != ""
	ua.logger.Debug("HasCompleteProfile result", map[string]any{
		"user_id": userID,
		"has_complete_profile": hasCompleteProfile,
	})

	return hasCompleteProfile, nil
}

func (ua *UserDbAdapter) GetSingleUserByEmail(ctx context.Context, email string) (models.User, error) {
	ua.logger.Debug("GetSingleUserByEmail called", map[string]any{"email": email})

	var user models.User
	err := ua.db.GetContext(
		ctx,
		&user,
		GetSingleUserByEmailQuery,
		email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found: %w", err)
		}
		return models.User{}, fmt.Errorf("error getting user by email: %w", err)
	}

	ua.logger.Debug("User retrieved successfully", map[string]any{
		"user_id": user.ID,
		"email": user.Email,
	})

	return user, nil
}
