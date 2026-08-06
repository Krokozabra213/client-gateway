package postgresRepo

import (
	"context"
	"testing"
	"time"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"

	authR "github.com/GargantuaLabs/client-gateway/internal/features/auth/repository/postgres"

	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"github.com/GargantuaLabs/client-gateway/migrations"
	testutils "github.com/Krokozabra213/gargantua_common/pkg/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetBotByToken(t *testing.T) {
	t.Parallel()

	t.Run("successfully get bot", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		authRepo := authR.New(pool)
		repo := &Repository{pool}

		tgID, _ := types.NewTelegramID(12345)
		username, _ := authDomain.NewUsername("Alice")

		user := authDomain.User{
			TelegramID:        tgID,
			TelegramFirstName: username,
			BirthDate:         time.Now().AddDate(-30, 0, 0),
			Gender:            authDomain.GenderMale,
			Role:              authDomain.RoleRegular,
			IsVerified:        true,
		}

		userID, err := authRepo.CreateUser(ctx, user)
		require.NoError(t, err)

		token := func() types.NonEmptyString { v, _ := types.NewNonEmptyString("token"); return v }()
		botTgID, _ := types.NewTelegramID(123)

		err = repo.CreateBot(ctx, &domain.Bot{
			ID:         domain.FromUuid(uuid.New()),
			TelegramID: botTgID,
			TokenHash:  token,
			Token:      token,
			Name:       func() types.NonEmptyString { v, _ := types.NewNonEmptyString("name"); return v }(),
			UserID:     userID,
		})
		require.NoError(t, err)

		bot, err := repo.GetBotByTokenHash(ctx, token)
		require.NoError(t, err)
		require.Equal(t, bot.Token, token)
	})

	t.Run("fail if bot doesn't exist", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		repo := &Repository{pool: pool}

		token, _ := types.NewNonEmptyString("token")
		_, err := repo.GetBotByTokenHash(ctx, token)
		require.ErrorIs(t, ErrBotNotFound, err)
	})
}
