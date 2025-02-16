package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRoot)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "text/html"
	if contentType := rr.Header().Get("Content-Type"); !strings.Contains(contentType, expected) {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expected)
	}
}

func TestHandlePrinciple(t *testing.T) {
	req, err := http.NewRequest("GET", "/principle?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlePrinciple)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "text/html"
	if contentType := rr.Header().Get("Content-Type"); !strings.Contains(contentType, expected) {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expected)
	}
}

func TestHandleCategory(t *testing.T) {
	req, err := http.NewRequest("GET", "/category?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleCategory)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "text/html"
	if contentType := rr.Header().Get("Content-Type"); !strings.Contains(contentType, expected) {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expected)
	}
}

func TestHandleCss(t *testing.T) {
	req, err := http.NewRequest("GET", "/style.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleCss)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "text/css"
	if contentType := rr.Header().Get("Content-Type"); !strings.Contains(contentType, expected) {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expected)
	}
}
