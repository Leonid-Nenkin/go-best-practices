package main_tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestNewFileSearcher (t *testing T) {
	res, err := main.NewFileSearcher()
	if err != nil {
		t.Errorf("....")
	}

	expected := 300
	if res != expected {
		t.Errorf("...")
	}
}

func TestListDirectory (t *testing T) {

}

func TesTfindFiles (t *testing T) {

}


