package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"sync/atomic"
	"syscall"
	"time"

	"go.uber.org/zap"

	crw "hw3/file_searcher"
)

type Config struct {
	MaxDepth int
}

var (
	GitHash  = ""
	BuidTime = ""
	Version  = ""
)

func main() {
	const (
		wantExt     = ".go"
		development = "DEVELOPMENT"
		production  = "PRODUCTION"
		env         = "ENV"
	)

	var logger *zap.Logger
	curEnv := os.Getenv(env)

	var err error
	if curEnv == production {
		logCfg := zap.NewProductionConfig()
		logCfg.OutputPaths = []string{"stderr"}
		logger, err = logCfg.Build()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal("Failed to initialize logger", err)
	}

	logger.Info("starting", zap.Int("pid", os.Getpid()),
		zap.String("commit_hash", GitHash), zap.String("Buid_time", BuidTime),
		zap.String("version", Version))

	defer logger.Sync()

	cfg := Config{2}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	var maint uint32
	go func() {
		for sig := range sigCh {
			switch sig {
			case syscall.SIGUSR1:
				atomic.StoreUint32(&maint, 1)
			case syscall.SIGUSR2:
				atomic.StoreUint32(&maint, 0)
			}
		}
	}()
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2)

	waitCh := make(chan struct{})
	fileSearcher := crw.NewFileSearcher(logger)

	go func() {
		res, err := fileSearcher.FindFiles(ctx, wantExt, cfg.MaxDepth, maint)

		if err != nil {
			logger.Error("Error on search: ", zap.Error(err))
			os.Exit(2)
		}

		for _, f := range res {
			fmt.Printf("Name: %s\t\t Path: %s\n", f.Name, f.Path)
		}
		waitCh <- struct{}{}
	}()
	go func() {
		<-sigCh
		logger.Info("Signal received, terminate ...")
		cancel()
	}()

	<-waitCh
	logger.Info("Done")
}
