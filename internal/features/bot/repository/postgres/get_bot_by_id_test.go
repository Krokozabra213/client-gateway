package postgresRepo

import (
	"context"
	"testing"
	"time"

	authDomain "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	authR "github.com/GargantuaLabs/client-gateway/internal/features/auth/repository/postgres"
	"github.com/GargantuaLabs/client-gateway/internal/features/bot/domain"

	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"github.com/GargantuaLabs/client-gateway/migrations"
	testutils "github.com/Krokozabra213/gargantua_common/pkg/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetBotByID(t *testing.T) {
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

		id := domain.FromUuid(uuid.New())
		botTgID, _ := types.NewTelegramID(123)

		err = repo.CreateBot(ctx, &domain.Bot{
			ID:         id,
			TelegramID: botTgID,
			Token:      func() types.NonEmptyString { v, _ := types.NewNonEmptyString("token"); return v }(),
			Name:       func() types.NonEmptyString { v, _ := types.NewNonEmptyString("name"); return v }(),
			UserID:     userID,
		})
		require.NoError(t, err)

		bot, err := repo.GetBotByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, bot.ID)
	})

	t.Run("fail if bot doesn't exist", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		repo := &Repository{pool: pool}

		_, err := repo.GetBotByID(ctx, domain.FromUuid(uuid.New()))
		require.ErrorIs(t, ErrBotNotFound, err)
	})
}
