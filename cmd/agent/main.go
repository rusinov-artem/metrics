package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/rusinov-artem/metrics/agent"
	"github.com/rusinov-artem/metrics/agent/client"
	"github.com/rusinov-artem/metrics/cmd/agent/config"

	"net/http"
	_ "net/http/pprof"
)

var runAgent = func(cfg *config.Config) {
	ctx := context.Background()
	client := client.New(fmt.Sprintf("http://%s", cfg.Address), fetchPublicKey(cfg.CryptoKey))
	client.Key = cfg.Key

	logger, _ := zap.NewDevelopment()
	logger.Info("config", zap.Any("config", cfg))

	agent.New(
		client,
		time.Second*time.Duration(cfg.PollInterval),
		time.Second*time.Duration(cfg.ReportInterval),
		cfg.RateLimit,
	).Run(ctx)
}

func NewAgent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Agent of best metrics project ever",
		Short: "Run agent to send metrics",
		Long:  "Run agent to send metrics",
	}

	cfg := config.New(cmd)

	cmd.Run = func(*cobra.Command, []string) {
		runAgent(cfg)
	}

	return cmd
}

func main() {
	go http.ListenAndServe(":9999", nil)
	err := NewAgent().Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func fetchPublicKey(publicKeyFile string) *rsa.PublicKey {
	if publicKeyFile == "" {
		return nil
	}

	fmt.Printf("%s will be used to encrypt data\n", publicKeyFile)

	f, err := os.Open(publicKeyFile)
	if err != nil {
		fmt.Printf("unable to open file %q: %v\n", publicKeyFile, err)
		return nil
	}

	pemData, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("unable to read file %q: %v\n", publicKeyFile, err)
		return nil
	}

	block, _ := pem.Decode(pemData)
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		fmt.Printf("unable to parse publicKey from %q: %v\n", publicKeyFile, err)
		return nil
	}

	return publicKey
}
