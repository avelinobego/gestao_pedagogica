package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avelinobego/gestao_pedagogica/shared"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
	"github.com/avelinobego/gestao_pedagogica/web"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf, err := shared.LoadConfig("$HOME/.config", ".")
	if err != nil {
		return err
	}
	if err := plataforma.Initialize(context.Background(), conf.Database.Driver, conf.Database.Dsn); err != nil {
		return err
	}
	defer func() {
		if err := plataforma.Close(); err != nil {
			log.Printf("falha ao fechar banco de dados: %v", err)
		}
	}()

	handler, err := web.NewHandler(web.MatriculaService{}, web.AcademicoService{})
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:              conf.Server.Address,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("falha ao iniciar servidor HTTP: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("falha ao encerrar servidor HTTP: %w", err)
		}
		if err := <-serverErr; !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("falha ao encerrar servidor HTTP: %w", err)
		}
	}
	return nil
}
