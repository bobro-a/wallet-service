package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"wallet/internal/config"
	"wallet/internal/handlers"
	"wallet/internal/repo"
	"wallet/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type App struct {
	ctx context.Context
	cfg *config.Config
	db  *sqlx.DB
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Printf("port: %s, address: %s", cfg.DB.Port, cfg.DB.Host)
	db, err := sqlx.Connect("pgx", GetDbURL(cfg.DB.Password, cfg.DB.Host, cfg.DB.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return &App{
		ctx: ctx,
		cfg: cfg,
		db:  db,
	}, nil
}

func (a *App) Start() error {
	defer func() {
		_ = a.db.Close()
	}()

	m, err := migrate.New(a.cfg.Migrate.Path, GetDbURL(a.cfg.DB.Password, a.cfg.DB.Host, a.cfg.DB.Port))
	if err != nil {
		return fmt.Errorf("failed to init migration: %w", err)
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		return fmt.Errorf("failed to run migration: %w", err)
	}
	repo := repo.NewWalletRepo(a.db, a.cfg.DB.Name)
	service := service.NewService(repo)
	h := handlers.NewHandler(service)
	r := chi.NewRouter()
	r.Get("/api/v1/wallets/{wallet_uuid}", h.GetAmount)
	r.Post("/api/v1/wallet", h.ChangeAmount)
	err = http.ListenAndServe(a.cfg.Server.Addr+":"+a.cfg.Server.Port, r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		return err
	}
	return nil
}

func GetDbURL(password, host, port string) string {
	return fmt.Sprintf("postgres://postgres:%s@%s:%s/postgres?sslmode=disable", password, host, port)
}
