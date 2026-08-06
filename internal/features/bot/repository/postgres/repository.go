package postgresRepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	d "github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrBotNotFound = errors.New("bot not found")

const timeout = time.Second * 4

type BotModel struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	TelegramID int64     `db:"telegram_id"`
	TokenHash  string    `db:"token_hash"`
	Token      string    `db:"token"`
	Name       string    `db:"name"`
	UserID     uuid.UUID `db:"user_id"`
	IsActive   bool      `json:"is_active"`
}

func (b *BotModel) ToDomain() (*d.Bot, error) {
	telegramId, err := types.NewTelegramID(b.TelegramID)
	if err != nil {
		return nil, err
	}

	token, err := types.NewNonEmptyString(b.Token)
	if err != nil {
		return nil, err
	}

	name, err := types.NewNonEmptyString(b.Name)
	if err != nil {
		return nil, err
	}

	return &d.Bot{
		ID:         d.FromUuid(b.ID),
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
		TelegramID: telegramId,
		Token:      token,
		Name:       name,
		UserID:     authDomain.FromUuid(b.UserID),
		IsActive:   b.IsActive,
	}, nil
}

func botModelsToDomain(models []BotModel) ([]*d.Bot, error) {
	bots := make([]*d.Bot, 0, len(models))

	for _, model := range models {
		bot, err := model.ToDomain()
		if err != nil {
			return nil, err
		}

		bots = append(bots, bot)
	}

	return bots, nil
}

type Repository struct {
	pool *pgxpool.Pool
}

func (r *Repository) CreateBot(ctx context.Context, bot *d.Bot) error {
	context, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	stmt := `
		INSERT INTO bot (id, telegram_id, token_hash, token, name, user_id)
		VALUES (@id, @telegram_id, @token_hash, @token, @name, @user_id)
		RETURNING id;
		`

	args := pgx.NamedArgs{
		"id":          bot.ID,
		"telegram_id": bot.TelegramID.Value(),
		"token_hash":  bot.TokenHash.Value(),
		"token":       bot.Token.Value(),
		"name":        bot.Name.Value(),
		"user_id":     bot.UserID,
	}

	_, err := r.pool.Exec(context, stmt, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBotsByUserID(ctx context.Context, userID authDomain.UserID, pagination *types.Pagination) ([]*d.Bot, error) {
	context, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
			SELECT
				id,
				created_at,
				updated_at,
				telegram_id,
				token_hash,
				token,
				name,
				user_id,
				is_active
			FROM bot
			WHERE user_id = $1
			ORDER BY created_at ASC
			OFFSET $2
			LIMIT $3
			`

	rows, err := r.pool.Query(
		context,
		query,
		userID,
		pagination.Offset(),
		pagination.Limit.Value(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bots, err := pgx.CollectRows(rows, pgx.RowToStructByName[BotModel])
	if err != nil {
		return nil, err
	}

	return botModelsToDomain(bots)
}

func (r *Repository) GetBotByTokenHash(ctx context.Context, tokenHash types.NonEmptyString) (*d.Bot, error) {
	context, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
			SELECT
				id,
				created_at,
				updated_at,
				telegram_id,
				token_hash,
				token,
				name,
				user_id,
				is_active
			FROM bot
			WHERE token_hash = $1;
			`

	row, err := r.pool.Query(context, query, tokenHash.Value())
	if err != nil {
		return nil, err
	}

	bot, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[BotModel])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBotNotFound
		}

		return nil, err
	}

	return bot.ToDomain()
}

func (r *Repository) GetBotByID(ctx context.Context, botID d.BotID) (*d.Bot, error) {
	context, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
			SELECT
				id,
				created_at,
				updated_at,
				telegram_id,
				token_hash,
				token,
				name,
				user_id,
				is_active
			FROM bot
			WHERE id = $1
			`

	row, err := r.pool.Query(context, query, botID)
	if err != nil {
		return nil, err
	}

	bot, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[BotModel])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBotNotFound
		}

		return nil, err
	}

	return bot.ToDomain()
}

func (r *Repository) GetBotByUserID(ctx context.Context, userID authDomain.UserID, botID d.BotID) (*d.Bot, error) {
	context, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
			SELECT
				id,
				created_at,
				updated_at,
				telegram_id,
				token_hash,
				token,
				name,
				user_id,
				is_active
			FROM bot
			WHERE user_id = $1 AND id = $2;
			`

	row, err := r.pool.Query(context, query, userID, botID)
	if err != nil {
		return nil, err
	}

	bot, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[BotModel])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBotNotFound
		}

		return nil, err
	}

	return bot.ToDomain()
}

func (r *Repository) GetBotCountByUserID(ctx context.Context, userID authDomain.UserID) (int, error) {
	context, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query := `
			SELECT COUNT(*)
			FROM bot
			WHERE user_id = $1;
			`

	var count int
	err := r.pool.QueryRow(context, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 0, nil
	}

	return count, nil
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
