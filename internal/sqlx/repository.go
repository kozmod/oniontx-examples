package sqlx

import (
	"context"
	"fmt"

	"github.com/kozmod/oniontx-examples/internal/utils"
)

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

func (r *TextRepository) Insert(ctx context.Context, val string) error {
	if r.errorExpected {
		return utils.ErrExpected
	}
	ex := r.transactor.GetExecutor(ctx)
	_, err := ex.ExecContext(ctx, `INSERT INTO text (val) VALUES ($1)`, val)
	if err != nil {
		return fmt.Errorf("sqlx repository - raw insert: %w", err)
	}
	return nil
}
