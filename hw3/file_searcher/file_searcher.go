package file_searcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

type FileSearcher struct {
	logger *zap.Logger
}

func NewFileSearcher(logger *zap.Logger) *FileSearcher {
	return &FileSearcher{
		logger: logger,
	}
}

func (f *FileSearcher) listDirectory(ctx context.Context, dir string,
	depth int, maint uint32) ([]FileInfo, error) {
	if depth <= 0 {
		return nil, nil
	}

	select {
	case <-ctx.Done():
		f.logger.Info("context is done, skipping dir", zap.String("dir", dir))
		return nil, nil
	default:
		// time.Sleep(time.Second * 10)
		var result []FileInfo
		res, err := os.ReadDir(dir)
		if err != nil {
			f.logger.Error("could not read dir", zap.Error(err),
				zap.String("dir", dir))
			return nil, err
		}

		for _, entry := range res {
			path := filepath.Join(dir, entry.Name())

			if atomic.LoadUint32(&maint) == 1 {
				fmt.Println(path, depth)
			}
			if atomic.LoadUint32(&maint) == 0 {
				depth += 2
			}

			if entry.IsDir() {
				child, err := f.listDirectory(ctx, path, depth-1, maint)
				if err != nil {
					return result, err
				}
				result = append(result, child...)
			} else {
				info, err := entry.Info()
				if err != nil {
					return result, nil
				}
				result = append(result, fileInfo{info, path})
			}
		}
		return result, err
	}
}

func (f *FileSearcher) FindFiles(ctx context.Context, ext string, maxDepth int, maint uint32) (FileList, error) {
	wd, err := os.Getwd()

	if err != nil {
		f.logger.Error("Could not get work directory", zap.Error(err))
		return nil, err
	}

	files, err := f.listDirectory(ctx, wd, maxDepth, maint)

	if err != nil {
		f.logger.Error("Could not get list of files", zap.Error(err))
		if len(files) == 0 {
			return nil, err
		}
		f.logger.Warn("Could not get list of files", zap.Error(err))
	}

	fl := make(FileList, len(files))
	for _, file := range files {
		fileExt := filepath.Ext(file.Name())
		f.logger.Debug("Compare extensions", zap.String("target_ext", ext),
			zap.String("current_ext", fileExt))

		if fileExt == ext {
			f.logger.Debug("Compare extensions", zap.String("target_ext", ext),
				zap.String("current", fileExt))

			fl[file.Name()] = TargetFile{
				Name: file.Name(),
				Path: file.Path(),
			}
		}
	}
	return fl, err
}
