package stdlib

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	oniontx "github.com/kozmod/oniontx/stdlib"
	"github.com/stretchr/testify/assert"

	"github.com/kozmod/oniontx-examples/internal/utils"
)

const (
	textRecord = "text_A"
)

func Test_UseCase(t *testing.T) {
	var (
		db = ConnectDB(t)
	)

	t.Cleanup(func() {
		err := db.Close()
		assert.NoError(t, err)
	})

	t.Run("success_create", func(t *testing.T) {
		var (
			ctx         = context.Background()
			transactor  = oniontx.NewTransactor(db)
			repositoryA = NewTextRepository(transactor, false)
			repositoryB = NewTextRepository(transactor, false)
			useCase     = NewUseCase(repositoryA, repositoryB, transactor)
		)

		err := useCase.CreateTextRecords(ctx, textRecord)
		assert.NoError(t, err)

		{
			records, err := GetTextRecords(db)
			assert.NoError(t, err)
			assert.Len(t, records, 2)
			for _, record := range records {
				assert.Equal(t, textRecord, record)
			}
		}

		t.Cleanup(func() {
			err = ClearDB(db)
			assert.NoError(t, err)
		})
	})
	t.Run("error_and_rollback", func(t *testing.T) {
		var (
			ctx         = context.Background()
			transactor  = oniontx.NewTransactor(db)
			repositoryA = NewTextRepository(transactor, false)
			repositoryB = NewTextRepository(transactor, true)
			useCase     = NewUseCase(repositoryA, repositoryB, transactor)
		)

		err := useCase.CreateTextRecords(ctx, textRecord)
		assert.Error(t, err)
		assert.ErrorIs(t, err, utils.ErrExpected)

		{
			records, err := GetTextRecords(db)
			assert.NoError(t, err)
			assert.Len(t, records, 0)
		}

		t.Cleanup(func() {
			err = ClearDB(db)
			assert.NoError(t, err)
		})
	})
}

func Test_UseCases(t *testing.T) {
	var (
		db = ConnectDB(t)
	)

	t.Cleanup(func() {
		err := db.Close()
		assert.NoError(t, err)
	})

	t.Run("single_repository", func(t *testing.T) {
		t.Run("success_create", func(t *testing.T) {
			var (
				ctx         = context.Background()
				transactor  = oniontx.NewTransactor(db)
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
				records, err := GetTextRecords(db)
				assert.NoError(t, err)
				assert.Len(t, records, 4)
				for _, record := range records {
					assert.Equal(t, textRecord, record)
				}
			}

			t.Cleanup(func() {
				err = ClearDB(db)
				assert.NoError(t, err)
			})
		})
		t.Run("error_and_rollback", func(t *testing.T) {
			var (
				ctx         = context.Background()
				transactor  = oniontx.NewTransactor(db)
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
				records, err := GetTextRecords(db)
				assert.NoError(t, err)
				assert.Len(t, records, 0)
			}

			t.Cleanup(func() {
				err = ClearDB(db)
				assert.NoError(t, err)
			})
		})
	})
}

func ConnectDB(t *testing.T) *sql.DB {
	connConfig, err := pgx.ParseConfig(utils.ConnectionString)
	assert.NoError(t, err)

	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sql.Open("pgx", connStr)
	assert.NoError(t, err)

	err = db.Ping()
	assert.NoError(t, err)

	return db
}

func ClearDB(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE text;")
	if err != nil {
		return fmt.Errorf("clear DB: %w", err)
	}
	return nil
}

func GetTextRecords(db *sql.DB) ([]string, error) {
	row, err := db.Query("SELECT val FROM text;")
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
