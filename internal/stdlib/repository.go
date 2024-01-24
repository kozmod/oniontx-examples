package stdlib

import (
	"context"
	"fmt"

	oniontx "github.com/kozmod/oniontx/stdlib"

	"github.com/kozmod/oniontx-examples/internal/utils"
)

type TextRepository struct {
	transactor    *oniontx.Transactor
	errorExpected bool
}

func NewTextRepository(transactor *oniontx.Transactor, errorExpected bool) *TextRepository {
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
		return fmt.Errorf("stdlib repository: %w", err)
	}
	return nil
}
