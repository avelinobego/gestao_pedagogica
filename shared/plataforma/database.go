package plataforma

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var databaseMu sync.Mutex
var db *sqlx.DB
var conn *sqlx.Conn

type Handler func(*sqlx.Tx, context.Context) error

func Initialize(ctx context.Context, driver, dsn string) error {
	if driver == "" {
		return errors.New("driver do banco de dados não foi configurado")
	}
	if dsn == "" {
		return errors.New("DSN do banco de dados não foi configurado")
	}
	newDB, err := sqlx.ConnectContext(ctx, driver, dsn)
	if err != nil {
		return fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}
	databaseMu.Lock()
	if db != nil {
		databaseMu.Unlock()
		if closeErr := newDB.Close(); closeErr != nil {
			return fmt.Errorf("banco de dados já inicializado; falha ao fechar conexão adicional: %w", closeErr)
		}
		return errors.New("banco de dados já inicializado")
	}
	db = newDB
	databaseMu.Unlock()
	return nil
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

func Executar(ctx context.Context, handlers ...Handler) (err error) {
	database, err := obterDB()
	if err != nil {
		return err
	}
	tx, err := database.BeginTxx(ctx, nil)
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
		if fn == nil {
			return errors.New("handler transacional não pode ser nil")
		}
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

func GetDB(ctx context.Context) (*sqlx.DB, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return obterDB()
}

func obterDB() (*sqlx.DB, error) {
	databaseMu.Lock()
	defer databaseMu.Unlock()
	if db == nil {
		return nil, errors.New("banco de dados não inicializado")
	}
	return db, nil
}

func GetConn(ctx context.Context) (result *sqlx.Conn, err error) {
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	databaseMu.Lock()
	defer databaseMu.Unlock()
	if db == nil {
		return nil, errors.New("banco de dados não inicializado")
	}
	if conn == nil {
		conn, err = db.Connx(ctx)
	}
	result = conn
	return
}

func CloseConn() {
	databaseMu.Lock()
	defer databaseMu.Unlock()
	if conn != nil {
		_ = conn.Close()
		conn = nil
	}
}

func Close() error {
	databaseMu.Lock()
	defer databaseMu.Unlock()
	if conn != nil {
		_ = conn.Close()
		conn = nil
	}
	if db == nil {
		return nil
	}
	err := db.Close()
	db = nil
	if err != nil {
		return fmt.Errorf("falha ao fechar conexão com o banco: %w", err)
	}
	return nil
}
