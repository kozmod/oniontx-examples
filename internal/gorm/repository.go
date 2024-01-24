package gorm

import (
	"context"
	"fmt"

	"github.com/kozmod/oniontx-examples/internal/utils"
)

type Text struct {
	Val string `gorm:"column:val"`
}

func (t *Text) TableName() string {
	return "text"
}

type TextRepository struct {
	transactor    transactor
	errorExpected bool
}

func NewTextRepository(transactor transactor, errorExpected bool) *TextRepository {
	return &TextRepository{
		transactor:    transactor,
		errorExpected: errorExpected,
	}
}

func (r *TextRepository) RawInsert(ctx context.Context, val string) error {
	if r.errorExpected {
		return utils.ErrExpected
	}
	ex := r.transactor.GetExecutor(ctx)
	ex = ex.Exec(`INSERT INTO text (val) VALUES ($1)`, val)
	if ex.Error != nil {
		return fmt.Errorf("gorm repository - raw insert: %w", ex.Error)
	}
	return nil
}

func (r *TextRepository) Insert(ctx context.Context, text Text) error {
	if r.errorExpected {
		return utils.ErrExpected
	}
	ex := r.transactor.GetExecutor(ctx)
	ex = ex.Create(text)
	if ex.Error != nil {
		return fmt.Errorf("gorm repository - raw insert: %w", ex.Error)
	}
	return nil
}
