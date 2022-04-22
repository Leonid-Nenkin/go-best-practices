package file_searcher_test

import (
	// "context"
	"github.com/stretchr/testify/assert"
	// "io/ioutil"
	// "net/http"
	// "strings"
	"hw3/config"
	crw "hw3/file_searcher"
	"log"
	"os"
	"testing"

	mock "github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// создание нового crawler'a
// указатель не nil
// максимальная глубина равна заданной

// количество найденных файлов
// увеличение глубины поиска
// увеличение/уменьшение глубины поиска по сигналу

type Searcher struct {
	mock.Mock
}

func TestNewFileSearcherSmoke(t *testing.T) {
	const (
		wantExt     = ".go"
		development = "DEVELOPMENT"
		production  = "PRODUCTION"
		env         = "ENV"
	)

	maxDepth := int(3)
	cfg := config.NewConfig(maxDepth)
	// r := new(Searcher)
	curEnv := os.Getenv(env)
	var err error

	var logger *zap.Logger
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

	cr := crw.NewFileSearcher(logger)
	assert.NotNil(t, cr)
	assert.Equal(t, maxDepth, cfg.MaxDepth)

	// res, err := main.NewFileSearcher()
	// if err != nil {
	// 	t.Errorf("....")
	// }

	// expected := 300
	// if res != expected {
	// 	t.Errorf("...")
	// }
}

// func TestListDirectory (t *testing T) {

// }

// func TesTfindFiles (t *testing T) {

// }
