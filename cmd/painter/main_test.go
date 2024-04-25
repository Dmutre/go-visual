package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dmutre/go-visual/painter"
	"github.com/Dmutre/go-visual/painter/lang"
)

func TestHttpHandler(t *testing.T) {
	var (
		opLoop painter.Loop 
		parser lang.Parser  
	)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := lang.HttpHandler(&opLoop, &parser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}
