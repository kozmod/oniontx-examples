package pgx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kozmod/oniontx"
)

type DB struct {
	*pgx.Conn
}

func (w *DB) BeginTx(ctx context.Context, opts ...oniontx.Option[*pgx.TxOptions]) (*Tx, error) {
	var txOptions pgx.TxOptions
	for _, opt := range opts {
		opt.Apply(&txOptions)
	}
	tx, err := w.Conn.BeginTx(ctx, txOptions)
	return &Tx{Tx: tx}, err
}

type Tx struct {
	pgx.Tx
}

func (t *Tx) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

func (t *Tx) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}

//goland:noinspection GoNameStartsWithPackageName
type PgxTransactor struct {
	*oniontx.Transactor[*DB, *Tx, *pgx.TxOptions]
}

func NewPgxTransactor(conn *pgx.Conn) *PgxTransactor {
	d := DB{Conn: conn}
	co := oniontx.NewContextOperator[*DB, *Tx](&d)
	tr := oniontx.NewTransactor[*DB, *Tx, *pgx.TxOptions](&d, co)
	return &PgxTransactor{
		Transactor: tr,
	}
}

//goland:noinspection GoExportedFuncWithUnexportedType
func (t *PgxTransactor) GetExecutor(ctx context.Context) executor {
	tx, ok := t.TryGetTx(ctx)
	if ok {
		return tx
	}
	return t.TxBeginner()
}
