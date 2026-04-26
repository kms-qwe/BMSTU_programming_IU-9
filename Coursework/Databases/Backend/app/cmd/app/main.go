package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"coffee-shop-backend/app/internal/api"
	employeeapi "coffee-shop-backend/app/internal/api/employee"
	ingredientapi "coffee-shop-backend/app/internal/api/ingredient"
	inventoryapi "coffee-shop-backend/app/internal/api/inventory"
	menuitemapi "coffee-shop-backend/app/internal/api/menuitem"
	orderapi "coffee-shop-backend/app/internal/api/order"
	reportapi "coffee-shop-backend/app/internal/api/report"
	shiftapi "coffee-shop-backend/app/internal/api/shift"
	"coffee-shop-backend/app/internal/config"
	employeerepo "coffee-shop-backend/app/internal/infra/repos/employee"
	ingredientrepo "coffee-shop-backend/app/internal/infra/repos/ingredient"
	inventoryrepo "coffee-shop-backend/app/internal/infra/repos/inventory"
	menuitemrepo "coffee-shop-backend/app/internal/infra/repos/menuitem"
	orderrepo "coffee-shop-backend/app/internal/infra/repos/order"
	reportrepo "coffee-shop-backend/app/internal/infra/repos/report"
	shiftrepo "coffee-shop-backend/app/internal/infra/repos/shift"
	employeeservice "coffee-shop-backend/app/internal/service/employee"
	ingredientservice "coffee-shop-backend/app/internal/service/ingredient"
	inventoryservice "coffee-shop-backend/app/internal/service/inventory"
	menuitemservice "coffee-shop-backend/app/internal/service/menuitem"
	orderservice "coffee-shop-backend/app/internal/service/order"
	reportservice "coffee-shop-backend/app/internal/service/report"
	shiftservice "coffee-shop-backend/app/internal/service/shift"
	"coffee-shop-backend/app/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()
	logger := utils.NewLogger(cfg.AppEnv)
	utils.SetLogger(logger)

	db, err := utils.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer closeDB(logger, db)

	txManager := utils.NewTxManager(db)

	employeeRepo := employeerepo.New(txManager)
	shiftRepo := shiftrepo.New(txManager)
	ingredientRepo := ingredientrepo.New(txManager)
	menuRepo := menuitemrepo.New(txManager)
	orderRepo := orderrepo.New(txManager)
	invRepo := inventoryrepo.New(txManager)
	repRepo := reportrepo.New(txManager)

	shiftSvc := shiftservice.New(shiftRepo, employeeRepo, txManager)
	employeeSvc := employeeservice.New(employeeRepo)
	ingredientSvc := ingredientservice.New(ingredientRepo, shiftRepo, invRepo, txManager)
	menuItemSvc := menuitemservice.New(menuRepo, ingredientRepo)
	orderSvc := orderservice.New(shiftRepo, menuRepo, ingredientRepo, orderRepo, invRepo, txManager)
	inventorySvc := inventoryservice.New(invRepo, shiftRepo, ingredientRepo, txManager)
	reportSvc := reportservice.New(repRepo)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	router := api.NewRouter(
		engine,
		shiftapi.NewAPI(shiftSvc),
		employeeapi.NewAPI(employeeSvc),
		ingredientapi.NewAPI(ingredientSvc),
		menuitemapi.NewAPI(menuItemSvc),
		orderapi.NewAPI(orderSvc),
		inventoryapi.NewAPI(inventorySvc),
		reportapi.NewAPI(reportSvc),
	)
	router.Register()

	server := &http.Server{
		Addr:              cfg.HTTP.Address(),
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("starting http server", "address", cfg.HTTP.Address(), "env", cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down http server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shutdown http server gracefully", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}

func closeDB(logger *slog.Logger, db *pgxpool.Pool) {
	db.Close()
	logger.Info("database connection closed")
}
