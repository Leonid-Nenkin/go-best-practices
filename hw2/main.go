package main

import (
	"fmt"
	"os"
	"path/filepath"
	"log"
	"context"
	"time"
	"syscall"
	"os/signal"
	"sync/atomic"
	"go.uber.org/zap"
)

type TargetFile struct {
	Path string
	Name string
}

type FileList map[string]TargetFile

type FileInfo interface {
	os.FileInfo
	Path() string
}

type fileInfo struct {
	os.FileInfo
	path string
}

func (fi fileInfo) Path() string {
	return fi.path
}

func ListDirectory(ctx context.Context, dir string, depth int, maint uint32) ([]FileInfo, error) {
	if depth <= 0 {
		return nil, nil
	}

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// time.Sleep(time.Second * 10)
		var result []FileInfo
		res, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, entry := range res {
			path := filepath.Join(dir, entry.Name())

			if atomic.LoadUint32(&maint) == 1 {
				fmt.Println(path, depth)
			}
			if atomic.LoadUint32(&maint) == 0 {
				depth+=2
			}

			if entry.IsDir() {
				child, err := ListDirectory(ctx, path, depth-1, maint)
				if err != nil {
					return nil, err
				}
				result = append(result, child...)
			} else {
				info, err := entry.Info()
				if err == nil {
					result = append(result, fileInfo{info,path})
				}
			}
		}
	return result, err
	}
}

func FindFiles(ctx context.Context, ext string, maxDepth int, maint uint32) (FileList, error) {	
	wd, err := os.Getwd()
	
	if err != nil {
		return nil, err
	}

	files, err := ListDirectory(ctx, wd, maxDepth, maint)

	if err != nil {
		return nil, err
	}

	fl := make(FileList, len(files))
	for _, file := range files {
		if filepath.Ext(file.Name()) == ext {
			fl[file.Name()] = TargetFile{
				Name: file.Name(),
				Path: file.Path(),
			}
		}
	}
	return fl, nil
}

type Config struct {
	MaxDepth   int
}

var (
	GitHash = ""
	BuidTime= ""
	Version = ""
)

func main() {
	const (
		wantExt = ".go"
		development = "DEVELOPMENT"
		production = "PRODUCTION"
		env = "ENV"
	)

	var logger *zap.Logger
	curEnv  := os.Getenv(env)

	var err error
	if curEnv == production {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal("Failed to initialize logger", err)
	}
	
	logger.Info("starting", zap.Int("pid", os.Getpid()),
		zap.String("commit hash", GitHash), zap.String("Buid time", BuidTime),
		zap.String("version", Version))

	defer logger.Sync()

	cfg := Config{2}
	
	ctx := context.Background()
	ctx, cancel := context.WithTimeout (ctx, 30 * time.Second)
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
	go func() {
		res, err := FindFiles(ctx, wantExt, cfg.MaxDepth, maint)
	
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
		<- sigCh
		log.Println("Signal received, terminate")	
		cancel()
	}()
	
	<-waitCh
	log.Println("Done")
}
