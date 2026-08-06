package postgresRepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	d "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

const timeout = time.Second * 4

type UserModel struct {
	ID                uuid.UUID `db:"id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	TelegramID        int64     `db:"telegram_id"`
	TelegramFirstname string    `db:"telegram_first_name"`
	Username          *string   `db:"username"`
	BirthDate         time.Time `db:"birth_date"`
	Gender            string    `db:"gender"`
	Role              string    `db:"role"`
	IsBanned          bool      `db:"is_banned"`
	IsVerified        bool      `db:"is_verified"`
}

func (u *UserModel) ToDomain() (*d.User, error) {
	telegramId, err := types.NewTelegramID(u.TelegramID)
	if err != nil {
		return nil, err
	}

	userGender, err := d.NewUserGender(u.Gender)
	if err != nil {
		return nil, err
	}

	telegramFirstName, err := d.NewUsername(u.TelegramFirstname)
	if err != nil {
		return nil, err
	}

	username, err := shared.ToNullable(u.Username, d.NewUsername)
	if err != nil {
		return nil, err
	}

	userRole, err := d.NewUserRole(u.Role)
	if err != nil {
		return nil, err
	}

	return &d.User{
		ID:                d.FromUuid(u.ID),
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
		TelegramID:        telegramId,
		TelegramFirstName: telegramFirstName,
		Username:          username,
		BirthDate:         u.BirthDate,
		Gender:            userGender,
		Role:              userRole,
		IsBanned:          u.IsBanned,
		IsVerified:        u.IsVerified,
	}, nil
}

type Repository struct {
	pool *pgxpool.Pool
}

func (r *Repository) CreateUser(ctx context.Context, user d.User) (d.UserID, error) {
	context, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	stmt := `INSERT INTO users (telegram_id, telegram_first_name, birth_date, gender, is_verified, role) VALUES($1, $2, $3, $4, $5, $6) RETURNING id;`

	var id uuid.UUID

	err := r.pool.QueryRow(context, stmt, user.TelegramID.Value(), user.TelegramFirstName, user.BirthDate, user.Gender, user.IsVerified, user.Role).Scan(&id)
	if err != nil {
		return d.UserID(uuid.Nil), err
	}

	return d.FromUuid(id), nil
}

func (r *Repository) GetUserByID(ctx context.Context, id d.UserID) (*d.User, error) {
	context, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
			SELECT
				id,
				created_at,
				updated_at,
				telegram_id,
				telegram_first_name,
				username,
				birth_date,
				gender,
				role,
				is_banned,
				is_verified
			FROM users
			WHERE id = $1;
			`

	row, err := r.pool.Query(context, query, id)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[UserModel])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user.ToDomain()
}

func (r *Repository) GetUserByTelegramID(ctx context.Context, telegramID types.TelegramID) (*d.User, error) {
	context, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
			SELECT
				id,
				created_at,
				updated_at,
				telegram_id,
				telegram_first_name,
				username,
				birth_date,
				gender,
				role,
				is_banned,
				is_verified
			FROM users
			WHERE telegram_id = $1;
			`

	row, err := r.pool.Query(context, query, telegramID.Value())
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[UserModel])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user.ToDomain()
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
