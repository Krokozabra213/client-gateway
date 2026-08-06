package postgresRepo

import (
	"context"
	"testing"
	"time"

	d "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"github.com/GargantuaLabs/client-gateway/migrations"
	testutils "github.com/Krokozabra213/gargantua_common/pkg/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("successfully create user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		repo := &Repository{pool: pool}

		tgID, err := types.NewTelegramID(12345)
		require.NoError(t, err)

		user := d.User{
			TelegramID:        tgID,
			TelegramFirstName: d.Username("Alice"),
			BirthDate:         time.Now().AddDate(-30, 0, 0),
			Gender:            d.GenderFemale,
			Role:              d.RoleRegular,
			IsVerified:        true,
		}

		userID, err := repo.CreateUser(ctx, user)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, uuid.UUID(userID))
	})

	t.Run("fail on duplicate telegram_id", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		repo := &Repository{pool: pool}

		tgID, err := types.NewTelegramID(54321)
		require.NoError(t, err)

		user1 := d.User{
			TelegramID:        tgID,
			TelegramFirstName: d.Username("Bot"),
			BirthDate:         time.Now().AddDate(-25, 0, 0),
			Gender:            d.GenderMale,
			IsVerified:        false,
			Role:              d.RoleRegular,
		}
		user2 := d.User{
			TelegramID:        tgID,
			TelegramFirstName: d.Username("Charlie"),
			BirthDate:         time.Now().AddDate(-20, 0, 0),
			Gender:            d.GenderMale,
			IsVerified:        true,
			Role:              d.RoleRegular,
		}

		_, err = repo.CreateUser(ctx, user1)
		require.NoError(t, err)

		_, err = repo.CreateUser(ctx, user2)
		require.Error(t, err)
	})

	t.Run("fail on invalid gender", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		repo := &Repository{pool: pool}

		tgID, err := types.NewTelegramID(12345)
		require.NoError(t, err)

		user := d.User{
			TelegramID:        tgID,
			TelegramFirstName: d.Username("test"),
			BirthDate:         time.Now().AddDate(-30, 0, 0),
			Gender:            "unknown",
			IsVerified:        true,
			Role:              d.RoleRegular,
		}

		_, err = repo.CreateUser(ctx, user)
		require.Error(t, err)
	})

	t.Run("fail on invalid role", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		repo := &Repository{pool: pool}

		tgID, err := types.NewTelegramID(12345)
		require.NoError(t, err)

		user := d.User{
			TelegramID:        tgID,
			TelegramFirstName: d.Username("test"),
			Role:              "invalid role",
			BirthDate:         time.Now().AddDate(-30, 0, 0),
			Gender:            d.GenderMale,
			IsVerified:        true,
		}

		_, err = repo.CreateUser(ctx, user)
		require.Error(t, err)
	})
}
