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

func TestRepository_GetBotCountByUserID(t *testing.T) {
	t.Parallel()

	t.Run("successfully get count", func(t *testing.T) {
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

		token, _ := types.NewNonEmptyString("token")
		name, _ := types.NewNonEmptyString("name")
		botTgID, _ := types.NewTelegramID(123)

		err = repo.CreateBot(ctx, &domain.Bot{
			ID:         domain.FromUuid(uuid.New()),
			TelegramID: botTgID,
			Token:      token,
			Name:       name,
			UserID:     userID,
		})
		require.NoError(t, err)

		count, err := repo.GetBotCountByUserID(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("returns 0 when user does not exist", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		pool := testutils.SetupTestDB(t, migrations.Files)
		repo := &Repository{pool}

		nonExistentUserID := authDomain.UserID(uuid.New())

		count, err := repo.GetBotCountByUserID(ctx, nonExistentUserID)

		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

}
