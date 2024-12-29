package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/cmd/server/config"
	"github.com/rusinov-artem/metrics/server"
	"github.com/rusinov-artem/metrics/server/handler"
	"github.com/rusinov-artem/metrics/server/middleware"
	"github.com/rusinov-artem/metrics/server/migration"
	"github.com/rusinov-artem/metrics/server/router"
	"github.com/rusinov-artem/metrics/server/storage"

	"github.com/spf13/cobra"

	"net/http"
	_ "net/http/pprof"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

var runServer = func(cfg *config.Config) {
	var err error
	var metricsStorage handler.MetricsStorage
	var dbpool *pgxpool.Pool

	logger, _ := zap.NewDevelopment()
	metricsStorage, destructor := storage.NewBufferedFileStorage(
		logger,
		cfg.FileStoragePath,
		cfg.Restore,
		cfg.StoreInterval,
	)
	defer destructor()
	metricsStorageFactory := func() handler.MetricsStorage {
		return metricsStorage
	}

	if cfg.DatabaseDsn != "" {
		dbpool, err = pgxpool.New(context.Background(), cfg.DatabaseDsn)
		if err != nil {
			logger.Error("unable to connect to database", zap.Error(err))
		} else {
			migration.Migrate(logger, dbpool)
		}
		defer func() {
			dbpool.Reset()
			dbpool.Close()
		}()

		metricsStorageFactory = func() handler.MetricsStorage {
			return storage.NewPgxStorage(dbpool)
		}
	}

	logger = logger.With(zap.Any("config", cfg))

	privateKey := fetchPrivateKey(cfg.CryptoKey)

	handler := handler.New(logger, metricsStorageFactory, dbpool)
	router := router.New(chi.NewRouter()).SetHandler(handler)
	// router.AddMiddleware(middleware.Sign(logger, cfg.Key))
	if privateKey != nil {
		router.AddMiddleware(middleware.Decrypt(privateKey, logger))
	}
	router.AddMiddleware(middleware.Logger(logger))
	router.AddMiddleware(middleware.GzipEncoder())
	server.New(logger, router.Mux(), cfg.Address).Run()

}

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Server of best metrics project ever",
		Short: "Run server on port 8080",
		Long:  "Run metrics collector server on port 8080",
	}

	cfg := config.New(cmd)

	cmd.Run = func(*cobra.Command, []string) {
		runServer(cfg)
	}

	return cmd
}

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	go http.ListenAndServe(":9898", nil)

	err := NewServerCmd().Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func fetchPrivateKey(privateKeyFile string) *rsa.PrivateKey {
	if privateKeyFile == "" {
		return nil
	}

	fmt.Printf("%s will be used to decrypt data\n", privateKeyFile)

	f, err := os.Open(privateKeyFile)
	if err != nil {
		fmt.Printf("unable to open file %q: %v\n", privateKeyFile, err)
		return nil
	}

	pemData, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("unable to read file %q: %v\n", privateKeyFile, err)
		return nil
	}

	block, _ := pem.Decode(pemData)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Printf("unable to parse privateKey from %q: %v\n", privateKeyFile, err)
		return nil
	}

	return privateKey
}
