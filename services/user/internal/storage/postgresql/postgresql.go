package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/guluzadehh/go_eshop/services/user/internal/domain/models"
	"github.com/guluzadehh/go_eshop/services/user/internal/storage"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.postgresql.UserByEmail"

	var user models.User

	const query = `
		SELECT users.id, users.email, users.password, users.created_at, users.updated_at, users.is_active
		FROM users
		WHERE users.email = $1;
	`

	if err := s.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.UserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) CreateUser(ctx context.Context, email string, password string) (*models.User, error) {
	const op = "storage.postgresql.CreateUser"

	const query = `
		INSERT INTO users(email, password)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at, is_active
	`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var lastInsertedId int64
	var createdAt, updatedAt time.Time
	var isActive bool

	err = stmt.QueryRowContext(ctx, email, password).Scan(&lastInsertedId, &createdAt, &updatedAt, &isActive)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code.Name() == "unique_violation" {
			return nil, storage.UserExists
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.User{
		Id:        lastInsertedId,
		Email:     email,
		Password:  []byte(password),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		IsActive:  isActive,
	}, nil
}

func (s *Storage) UserById(ctx context.Context, id int) (*models.User, error) {
	const op = "storage.postgresql.UserById"

	var user models.User

	const query = `
		SELECT users.id, users.email, users.password, users.created_at, users.updated_at, users.is_active
		FROM users
		WHERE users.id = $1;
	`

	if err := s.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.IsActive); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.UserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) ProfileById(ctx context.Context, userId int64) (*models.Profile, error) {
	const op = "storage.postgresql.ProfileById"

	var profile models.Profile

	const query = `
		SELECT 
			users.id, users.email, 
			users.password, 
			users.created_at, users.updated_at, 
			users.is_active,
			profiles.first_name,
			profiles.last_name,
			profiles.phone_number,
			profiles.profile_pic
		FROM users
		LEFT JOIN profiles ON profiles.id = users.id
		WHERE users.id = $1;
	`

	if err := s.db.QueryRowContext(ctx, query, userId).Scan(
		&profile.Id, &profile.Email, &profile.Password, &profile.CreatedAt,
		&profile.UpdatedAt, &profile.IsActive, &profile.FirstName, &profile.LastName,
		&profile.Phone, &profile.Picture,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.ProfileNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &profile, nil
}

func (s *Storage) DeleteProfile(ctx context.Context, userId int64) error {
	const op = "storage.postgresql.DeleteProfile"

	const queryUserDelete = `DELETE FROM users WHERE users.id = $1;`
	const queryProfileDelete = `DELETE FROM profiles WHERE profiles.user_id = $1;`

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	stmt, err := s.db.PrepareContext(ctx, queryProfileDelete)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = s.db.PrepareContext(ctx, queryUserDelete)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	affects, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if affects == 0 {
		return fmt.Errorf("%s: %w", op, storage.UserNotFound)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveProfile(
	ctx context.Context,
	userId int64,
	firstName string,
	lastName string,
	phone string,
) (*models.Profile, error) {
	const op = "storage.postgresql.SaveProfile"

	const query = `
		INSERT INTO profiles(user_id, first_name, last_name, phone_number)
		VALUES ($1, $2, $3, $4);
	`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.ExecContext(ctx, userId, firstName, lastName, phone); err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
			if postgresErr.Code == "23503" { // Foreign key violation error code
				return nil, fmt.Errorf("%s: %w", op, storage.UserNotFound)
			}
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	profile, err := s.ProfileById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return profile, nil
}
