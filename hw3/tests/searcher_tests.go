package tests

import (
	// "context"
	"github.com/stretchr/testify/assert"
	// "io/ioutil"
	// "net/http"
	// "strings"
	"testing"
	mock "github.com/stretchr/testify/mock"
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

func TestNewFileSearcher (t *testingT) {
	maxDepth := uint64(3)
	r := new(Searcher)
	cr := NewCrawler(r, maxDepth)
	assert.NotNil(t, cr)
	assert.Equal(t, maxDepth, cr.MaxDepth)

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


