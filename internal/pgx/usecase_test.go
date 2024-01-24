package pgx

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/kozmod/oniontx-examples/internal/utils"
	"github.com/stretchr/testify/assert"
)

const (
	textRecord = "text_A"
)

func Test_UseCase_CreateTextRecords(t *testing.T) {
	var (
		globalCtx = context.Background()
		db        = ConnectDB(globalCtx, t)
	)

	t.Cleanup(func() {
		err := db.Close(globalCtx)
		assert.NoError(t, err)
	})

	t.Run("success_create", func(t *testing.T) {
		var (
			ctx         = context.Background()
			transactor  = NewPgxTransactor(db)
			repositoryA = NewTextRepository(transactor, false)
			repositoryB = NewTextRepository(transactor, false)
			useCase     = NewUseCase(repositoryA, repositoryB, transactor)
		)

		err := useCase.CreateTextRecords(ctx, textRecord)
		assert.NoError(t, err)

		{
			records, err := GetTextRecords(globalCtx, db)
			assert.NoError(t, err)
			assert.Len(t, records, 2)
			for _, record := range records {
				assert.Equal(t, textRecord, record)
			}
		}

		t.Cleanup(func() {
			err = ClearDB(globalCtx, db)
			assert.NoError(t, err)
		})
	})
	t.Run("error_and_rollback", func(t *testing.T) {
		var (
			ctx         = context.Background()
			transactor  = NewPgxTransactor(db)
			repositoryA = NewTextRepository(transactor, false)
			repositoryB = NewTextRepository(transactor, true)
			useCase     = NewUseCase(repositoryA, repositoryB, transactor)
		)

		err := useCase.CreateTextRecords(ctx, textRecord)
		assert.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrExpected)

		{
			records, err := GetTextRecords(globalCtx, db)
			assert.NoError(t, err)
			assert.Len(t, records, 0)

		}

		t.Cleanup(func() {
			err = ClearDB(globalCtx, db)
			assert.NoError(t, err)
		})
	})
}

func Test_UseCases(t *testing.T) {
	var (
		globalCtx = context.Background()
		db        = ConnectDB(globalCtx, t)
	)

	t.Cleanup(func() {
		err := db.Close(globalCtx)
		assert.NoError(t, err)
	})

	t.Run("single_repository", func(t *testing.T) {
		t.Run("success_create", func(t *testing.T) {
			var (
				ctx         = context.Background()
				transactor  = NewPgxTransactor(db)
				repositoryA = NewTextRepository(transactor, false)
				repositoryB = NewTextRepository(transactor, false)
				useCases    = NewUseCases(
					NewUseCase(repositoryA, repositoryB, transactor),
					NewUseCase(repositoryA, repositoryB, transactor),
					transactor,
				)
			)

			err := useCases.CreateTextRecords(ctx, textRecord)
			assert.NoError(t, err)

			{
				records, err := GetTextRecords(globalCtx, db)
				assert.NoError(t, err)
				assert.Len(t, records, 4)
				for _, record := range records {
					assert.Equal(t, textRecord, record)
				}
			}

			t.Cleanup(func() {
				err = ClearDB(globalCtx, db)
				assert.NoError(t, err)
			})
		})
		t.Run("error_and_rollback", func(t *testing.T) {
			var (
				ctx         = context.Background()
				transactor  = NewPgxTransactor(db)
				repositoryA = NewTextRepository(transactor, false)
				repositoryB = NewTextRepository(transactor, true)
				useCases    = NewUseCases(
					NewUseCase(repositoryA, repositoryB, transactor),
					NewUseCase(repositoryA, repositoryB, transactor),
					transactor,
				)
			)

			err := useCases.CreateTextRecords(ctx, textRecord)
			assert.Error(t, err)
			assert.ErrorIs(t, err, utils.ErrExpected)

			{
				records, err := GetTextRecords(globalCtx, db)
				assert.NoError(t, err)
				assert.Len(t, records, 0)
			}

			t.Cleanup(func() {
				err = ClearDB(globalCtx, db)
				assert.NoError(t, err)
			})
		})
	})
}

func ConnectDB(ctx context.Context, t *testing.T) *pgx.Conn {
	conn, err := pgx.Connect(ctx, utils.ConnectionString)
	assert.NoError(t, err)

	err = conn.Ping(ctx)
	assert.NoError(t, err)
	return conn
}

func ClearDB(ctx context.Context, db *pgx.Conn) error {
	_, err := db.Exec(ctx, `TRUNCATE TABLE text;`)
	if err != nil {
		return fmt.Errorf("clear DB: %w", err)
	}
	return nil
}

func GetTextRecords(ctx context.Context, db *pgx.Conn) ([]string, error) {
	row, err := db.Query(ctx, "SELECT val FROM text;")
	if err != nil {
		return nil, fmt.Errorf("get `text` records: %w", err)
	}

	var texts []string
	for row.Next() {
		var text string
		err = row.Scan(&text)
		if err != nil {
			return nil, fmt.Errorf("scan `text` records: %w", err)
		}
		texts = append(texts, text)
	}
	return texts, nil
}
