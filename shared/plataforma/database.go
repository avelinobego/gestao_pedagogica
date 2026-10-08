package plataforma

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/avelinobego/gestao_pedagogica/shared"
	"time"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var conn *sqlx.Conn

type Handler func(*sqlx.Tx, context.Context) error

func init() {
	conf, err := shared.LoadConfig("$HOME/.config",
		".")
	if err != nil {
		panic(err)
	}
	db = sqlx.MustConnect(conf.Database.Driver, conf.Database.Dsn)
}

func CincoSegudos() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(time.Second*5))
}

func TrintaSegudos() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(time.Second*30))
}

func UmMinuto() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(time.Minute*1))
}

func Executar(ctx context.Context, handlers ...Handler) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("falha ao iniciar transação: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			if errrb := tx.Rollback(); errrb != nil && !errors.Is(errrb, sql.ErrTxDone) {
				err = fmt.Errorf("falha ao executar Rollback: %w: %w", errrb, err)
			}
		}
	}()

	for _, fn := range handlers {
		if err = ctx.Err(); err != nil {
			return err
		}
		if err = fn(tx, ctx); err != nil {
			return fmt.Errorf("falha ao executar função: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("falha ao executar Commit: %w", err)
	}
	return nil
}

func GetConn(ctx context.Context) (result *sqlx.Conn, err error) {
	if conn == nil {
		conn, err = db.Connx(ctx)
	}
	result = conn
	return
}

func CloseConn() {
	if conn != nil {
		_ = conn.Close()
		conn = nil
	}
}
