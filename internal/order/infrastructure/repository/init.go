package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/moeryomenko/saga/internal/order/config"
)

// module as singleton.
var pool *pgxpool.Pool = nil

func Init(cfg *config.Config) func(context.Context) error {
	return func(ctx context.Context) error {
		dbConfig, err := pgxpool.ParseConfig(fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Name, cfg.Database.Password,
		))
		if err != nil {
			return err
		}
		dbConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
			c.ConnInfo().RegisterDataType(pgtype.DataType{
				Value: &shopspring.Numeric{},
				Name:  "decimal",
				OID:   pgtype.NumericOID,
			})
			return nil
		}
		pool, err = pgxpool.ConnectConfig(ctx, dbConfig)
		if err != nil {
			return err
		}
		return nil
	}
}

func Close(ctx context.Context) error {
	if pool != nil {
		pool.Close()
	}

	return nil
}
