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

func TestRepository_CreateBot(t *testing.T) {

	t.Parallel()

	t.Run("successfully create bot", func(t *testing.T) {
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

		botTgID, _ := types.NewTelegramID(123)

		token, err := types.NewNonEmptyString("token")
		require.NoError(t, err)
		name, err := types.NewNonEmptyString("name")
		require.NoError(t, err)

		err = repo.CreateBot(ctx, &domain.Bot{
			ID:         domain.FromUuid(uuid.New()),
			TelegramID: botTgID,
			Token:      token,
			Name:       name,
			UserID:     userID,
		})
		require.NoError(t, err)
	})

	t.Run("fail if user doesn't exists", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		repo := &Repository{pool}

		tgID, _ := types.NewTelegramID(123)

		token, _ := types.NewNonEmptyString("token")
		name, _ := types.NewNonEmptyString("name")
		err := repo.CreateBot(ctx, &domain.Bot{
			ID:         domain.FromUuid(uuid.New()),
			TelegramID: tgID,
			Token:      token,
			Name:       name,
			UserID:     authDomain.FromUuid(uuid.New()),
		})
		require.Error(t, err)
	})

	t.Run("failed on duplicate telegram_id", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		authRepo := authR.New(pool)
		repo := &Repository{pool}

		tgID, _ := types.NewTelegramID(12345)

		user := authDomain.User{
			TelegramID:        tgID,
			TelegramFirstName: authDomain.Username("test"),
			BirthDate:         time.Now().AddDate(-30, 0, 0),
			Gender:            authDomain.GenderMale,
			Role:              authDomain.RoleRegular,
			IsVerified:        true,
		}

		userID, err := authRepo.CreateUser(ctx, user)
		require.NoError(t, err)

		botTgID, _ := types.NewTelegramID(123)
		bot := &domain.Bot{
			ID:         domain.FromUuid(uuid.New()),
			TelegramID: botTgID,
			Token:      func() types.NonEmptyString { v, _ := types.NewNonEmptyString("token"); return v }(),
			Name:       func() types.NonEmptyString { v, _ := types.NewNonEmptyString("name"); return v }(),
			UserID:     userID,
		}

		err = repo.CreateBot(ctx, bot)
		require.NoError(t, err)

		bot.ID = domain.FromUuid(uuid.New())
		bot.Token, _ = types.NewNonEmptyString("another token")
		err = repo.CreateBot(ctx, bot)
		require.Error(t, err)
	})

	t.Run("failed on duplicate token", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		authRepo := authR.New(pool)
		repo := &Repository{pool}

		tgID, _ := types.NewTelegramID(12345)

		user := authDomain.User{
			TelegramID:        tgID,
			TelegramFirstName: authDomain.Username("test"),
			BirthDate:         time.Now().AddDate(-30, 0, 0),
			Gender:            authDomain.GenderMale,
			Role:              authDomain.RoleRegular,
			IsVerified:        true,
		}

		userID, err := authRepo.CreateUser(ctx, user)
		require.NoError(t, err)

		botTgID, _ := types.NewTelegramID(123)

		token, _ := types.NewNonEmptyString("token")
		name, _ := types.NewNonEmptyString("name")
		bot := &domain.Bot{
			ID:         domain.FromUuid(uuid.New()),
			TelegramID: botTgID,
			Token:      token,
			Name:       name,
			UserID:     userID,
		}

		err = repo.CreateBot(ctx, bot)
		require.NoError(t, err)

		botTgID2, _ := types.NewTelegramID(123)

		bot.ID = domain.FromUuid(uuid.New())
		bot.TelegramID = botTgID2
		
		err = repo.CreateBot(ctx, bot)
		require.Error(t, err)
	})

}
