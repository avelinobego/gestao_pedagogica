package plataforma_test

import (
	"context"
	"testing"
	"time"

	"github.com/avelinobego/gestao_pedagogica/shared"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
	"github.com/jmoiron/sqlx"
)

func inicializarBancoDeTeste(t *testing.T) {
	t.Helper()
	config, err := shared.LoadConfig("$HOME/.config", ".")
	if err != nil {
		t.Skipf("configuração do banco de teste indisponível: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := plataforma.Initialize(ctx, config.Database.Driver, config.Database.Dsn); err != nil {
		t.Skipf("banco de teste indisponível: %v", err)
	}
	t.Cleanup(func() {
		if err := plataforma.Close(); err != nil {
			t.Errorf("falha ao fechar conexão de teste: %v", err)
		}
	})
}

func TestConexao(t *testing.T) {
	inicializarBancoDeTeste(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := plataforma.GetDB(ctx)
	if err != nil {
		t.Fatalf("GetDB() retornou erro: %v", err)
	}
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("PingContext() retornou erro: %v", err)
	}
}

func TestTransacao(t *testing.T) {
	inicializarBancoDeTeste(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	handler := func(tx *sqlx.Tx, ctx context.Context) error {
		var valor int
		return tx.GetContext(ctx, &valor, "SELECT 1")
	}
	if err := plataforma.Executar(ctx, handler); err != nil {
		t.Fatalf("Executar() retornou erro: %v", err)
	}
}
