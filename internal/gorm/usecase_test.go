package stdlib

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kozmod/oniontx-examples/internal/utils"
)

const (
	textRecord = "text_A"
)

func Test_UseCase_CreateTextRecords(t *testing.T) {
	var (
		db = ConnectDB(t)
	)
	t.Run("success_create", func(t *testing.T) {
		var (
			ctx         = context.Background()
			transactor  = NewGormTransactor(db)
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
				assert.Equal(t, Text{Val: textRecord}, record)
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
			transactor  = NewGormTransactor(db)
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

func Test_UseCase_CreateText(t *testing.T) {
	var (
		db = ConnectDB(t)

		text = Text{
			Val: textRecord,
		}
	)
	t.Run("success_create", func(t *testing.T) {
		var (
			ctx         = context.Background()
			transactor  = NewGormTransactor(db)
			repositoryA = NewTextRepository(transactor, false)
			repositoryB = NewTextRepository(transactor, false)
			useCase     = NewUseCase(repositoryA, repositoryB, transactor)
		)

		err := useCase.CreateText(ctx, text)
		assert.NoError(t, err)

		{
			records, err := GetTextRecords(db)
			assert.NoError(t, err)
			assert.Len(t, records, 2)
			for _, record := range records {
				assert.Equal(t, Text{Val: textRecord}, record)
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
			transactor  = NewGormTransactor(db)
			repositoryA = NewTextRepository(transactor, false)
			repositoryB = NewTextRepository(transactor, true)
			useCase     = NewUseCase(repositoryA, repositoryB, transactor)
		)

		err := useCase.CreateText(ctx, text)
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
	t.Run("single_repository", func(t *testing.T) {
		t.Run("success_create", func(t *testing.T) {
			var (
				ctx         = context.Background()
				transactor  = NewGormTransactor(db)
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
					assert.Equal(t, Text{Val: textRecord}, record)
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
				transactor  = NewGormTransactor(db)
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

func ConnectDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(postgres.Open(utils.ConnectionString), &gorm.Config{})
	assert.NoError(t, err)
	return db
}

func ClearDB(db *gorm.DB) error {
	ex := db.Exec(`TRUNCATE TABLE text;`)
	if ex.Error != nil {
		return fmt.Errorf("clear DB: %w", ex.Error)
	}
	return nil
}

func GetTextRecords(db *gorm.DB) ([]Text, error) {
	var texts []Text
	db = db.Find(&texts)
	if err := db.Error; err != nil {
		return nil, fmt.Errorf("get `text` records: %w", err)
	}
	return texts, nil
}
