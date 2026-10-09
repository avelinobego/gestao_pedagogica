package shared_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/avelinobego/gestao_pedagogica/shared"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("DATABASE_DRIVER", "")
	t.Setenv("DATABASE_DSN", "")
	t.Setenv("HTTP_ADDR", "")
	dir := t.TempDir()
	config := []byte("database:\n  driver: postgres\n  dsn: test-dsn\nserver:\n  address: 127.0.0.1:9090\n")
	if err := os.WriteFile(filepath.Join(dir, "gestao.yaml"), config, 0o600); err != nil {
		t.Fatal(err)
	}

	loaded, err := shared.LoadConfig(dir)
	if err != nil {
		t.Fatalf("LoadConfig() retornou erro: %v", err)
	}
	if loaded.Database.Driver != "postgres" {
		t.Errorf("driver = %q; esperado postgres", loaded.Database.Driver)
	}
	if loaded.Database.Dsn != "test-dsn" {
		t.Errorf("dsn = %q; esperado test-dsn", loaded.Database.Dsn)
	}
	if loaded.Server.Address != "127.0.0.1:9090" {
		t.Errorf("server address = %q; esperado 127.0.0.1:9090", loaded.Server.Address)
	}
}

func TestLoadConfigEnderecoPadrao(t *testing.T) {
	t.Setenv("DATABASE_DRIVER", "")
	t.Setenv("DATABASE_DSN", "")
	t.Setenv("HTTP_ADDR", "")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "gestao.yaml"), []byte("database: {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	loaded, err := shared.LoadConfig(dir)
	if err != nil {
		t.Fatalf("LoadConfig() retornou erro: %v", err)
	}
	if loaded.Server.Address != "127.0.0.1:8080" {
		t.Errorf("server address = %q; esperado 127.0.0.1:8080", loaded.Server.Address)
	}
}

func TestLoadConfigHTTPAddrSobrescreveYAML(t *testing.T) {
	t.Setenv("DATABASE_DRIVER", "")
	t.Setenv("DATABASE_DSN", "")
	t.Setenv("HTTP_ADDR", "127.0.0.1:9091")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "gestao.yaml"), []byte("server:\n  address: 127.0.0.1:9090\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	loaded, err := shared.LoadConfig(dir)
	if err != nil {
		t.Fatalf("LoadConfig() retornou erro: %v", err)
	}
	if loaded.Server.Address != "127.0.0.1:9091" {
		t.Errorf("server address = %q; esperado override 127.0.0.1:9091", loaded.Server.Address)
	}
}

func TestLoadConfigErroQuandoArquivoAusente(t *testing.T) {
	_, err := shared.LoadConfig(t.TempDir())
	if err == nil {
		t.Fatal("LoadConfig() deveria falhar quando não há arquivo de configuração")
	}
}
