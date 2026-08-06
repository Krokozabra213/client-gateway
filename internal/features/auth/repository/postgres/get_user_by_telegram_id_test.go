package postgresRepo

import (
	"context"
	"testing"
	"time"

	d "github.com/GargantuaLabs/client-gateway/internal/features/auth/domain"
	"github.com/GargantuaLabs/client-gateway/internal/features/shared/types"
	"github.com/GargantuaLabs/client-gateway/migrations"
	testutils "github.com/Krokozabra213/gargantua_common/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetUserByTelegramID(t *testing.T) {
	t.Run("successfully get user", func(t *testing.T) {
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
			Gender:            d.GenderMale,
			Role:              d.RoleRegular,
			IsVerified:        true,
		}

		_, err = repo.CreateUser(ctx, user)
		require.NoError(t, err)

		gotUser, err := repo.GetUserByTelegramID(ctx, tgID)
		require.NoError(t, err)
		require.Equal(t, user.TelegramID, gotUser.TelegramID)
		require.Equal(t, user.TelegramFirstName, gotUser.TelegramFirstName)
		require.Equal(t, user.Gender, gotUser.Gender)
		require.Equal(t, user.IsVerified, gotUser.IsVerified)
	})

	t.Run("fail if user doesn't exist", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		repo := &Repository{pool: pool}

		tgID, err := types.NewTelegramID(123)
		require.NoError(t, err)

		_, err = repo.GetUserByTelegramID(ctx, tgID)
		require.ErrorIs(t, ErrUserNotFound, err)
	})
}
