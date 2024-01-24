package gorm

import (
	"context"
	"database/sql"

	"github.com/kozmod/oniontx"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func (t *DB) Commit(_ context.Context) error {
	tx := t.DB
	tx.Commit()
	return nil
}

func (t *DB) Rollback(_ context.Context) error {
	tx := t.DB
	tx.Rollback()
	return nil
}

func (t *DB) BeginTx(_ context.Context, opts ...oniontx.Option[*sql.TxOptions]) (*DB, error) {
	var txOptions sql.TxOptions
	for _, opt := range opts {
		opt.Apply(&txOptions)
	}
	b := t.DB
	tx := b.Begin(&txOptions)
	return &DB{DB: tx}, nil
}

//goland:noinspection GoNameStartsWithPackageName
type GormTransactor struct {
	*oniontx.Transactor[*DB, *DB, *sql.TxOptions]
}

func NewGormTransactor(db *gorm.DB) *GormTransactor {
	d := DB{DB: db}
	co := oniontx.NewContextOperator[*DB, *DB](&d)
	tr := oniontx.NewTransactor[*DB, *DB, *sql.TxOptions](&d, co)
	return &GormTransactor{
		Transactor: tr,
	}
}

func (t *GormTransactor) GetExecutor(ctx context.Context) *gorm.DB {
	tx, ok := t.TryGetTx(ctx)
	if !ok {
		tx = t.TxBeginner()
	}
	return tx.DB
}
