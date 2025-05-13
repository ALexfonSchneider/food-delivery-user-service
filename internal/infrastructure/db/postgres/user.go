package postgres

import (
	"context"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"time"
)

type User struct {
	Id           string    `db:"id"`
	FirstName    string    `db:"first_name"`
	LastName     *string   `db:"last_name"`
	Email        string    `db:"email"`
	Phone        string    `db:"phone"`
	CreatedAt    time.Time `db:"created_at"`
	HashPassword string    `db:"hash_password"`
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user User

	sql := "SELECT id, first_name, last_name, email, phone, created_at, hash_password FROM users.user WHERE email=$1"
	err := pgxscan.Get(ctx, r.connection(ctx), &user, sql, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(domain.RecordNotFoundError, "failed to find user")
		}
		return nil, err
	}

	return UserModelToDomain(user), nil
}

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var user User

	sql := "SELECT id, first_name, last_name, email, phone, created_at, hash_password FROM users WHERE phone=$1"
	err := pgxscan.Select(ctx, r.connection(ctx), &user, sql, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(domain.RecordNotFoundError, "failed to find user")
		}
		return nil, err
	}

	return UserModelToDomain(user), nil
}

func (r *Repository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	var user User

	sql := "SELECT id, first_name, last_name, email, phone, created_at, hash_password FROM users.user WHERE id=$1"
	err := pgxscan.Get(ctx, r.connection(ctx), &user, sql, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(domain.RecordNotFoundError, "failed to find user")
		}
		return nil, err
	}

	return UserModelToDomain(user), nil
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	sql := "INSERT INTO users.user(id, first_name, email, phone, created_at, hash_password) VALUES($1, $2, $3, $4, $5, $6)"
	_, err := r.connection(ctx).Exec(ctx, sql, user.Id, user.FirstName, user.Email, user.Phone, user.CreatedAt, user.HashPassword)
	if err != nil {
		return err
	}

	return nil
}
