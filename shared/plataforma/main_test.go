package plataforma_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/avelinobego/gestao_pedagogica/internal/testdb"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	database, err := testdb.Start(ctx)
	if err != nil {
		cancel()
		fmt.Fprintf(os.Stderr, "start platform test database: %v\n", err)
		os.Exit(1)
	}
	if err := plataforma.Initialize(ctx, "postgres", database.DSN); err != nil {
		cancel()
		fmt.Fprintf(os.Stderr, "initialize platform test database: %v\n", err)
		terminateCtx, terminateCancel := context.WithTimeout(context.Background(), 30*time.Second)
		terminateErr := database.Terminate(terminateCtx)
		terminateCancel()
		if terminateErr != nil {
			fmt.Fprintf(os.Stderr, "terminate platform test database: %v\n", terminateErr)
		}
		os.Exit(1)
	}
	cancel()

	code := m.Run()
	if err := plataforma.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "close platform test database connection: %v\n", err)
		code = 1
	}
	terminateCtx, terminateCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer terminateCancel()
	if err := database.Terminate(terminateCtx); err != nil {
		fmt.Fprintf(os.Stderr, "terminate platform test database: %v\n", err)
		code = 1
	}
	terminateCancel()
	os.Exit(code)
}
