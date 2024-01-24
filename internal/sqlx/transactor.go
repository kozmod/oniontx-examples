package sqlx

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/kozmod/oniontx"
)

type DB struct {
	*sqlx.DB
}

func (w *DB) BeginTx(ctx context.Context, opts ...oniontx.Option[*sql.TxOptions]) (*Tx, error) {
	var txOptions sql.TxOptions
	for _, opt := range opts {
		opt.Apply(&txOptions)
	}
	tx, err := w.DB.BeginTxx(ctx, &txOptions)
	return &Tx{Tx: tx}, err
}

type Tx struct {
	*sqlx.Tx
}

func (t *Tx) Commit(_ context.Context) error {
	return t.Tx.Commit()
}

func (t *Tx) Rollback(_ context.Context) error {
	return t.Tx.Rollback()
}

//goland:noinspection GoNameStartsWithPackageName
type SqlxTransactor struct {
	*oniontx.Transactor[*DB, *Tx, *sql.TxOptions]
}

func NewSqlxTransactor(db *sqlx.DB) *SqlxTransactor {
	d := DB{DB: db}
	co := oniontx.NewContextOperator[*DB, *Tx](&d)
	tr := oniontx.NewTransactor[*DB, *Tx, *sql.TxOptions](&d, co)
	return &SqlxTransactor{
		Transactor: tr,
	}
}

func (t *SqlxTransactor) GetExecutor(ctx context.Context) executor {
	tx, ok := t.TryGetTx(ctx)
	if ok {
		return tx
	}
	return t.TxBeginner()
}
