package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	newServer().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("unexpected content type %q", got)
	}
}

func TestBookingPage(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	newServer().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if got := response.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Fatalf("unexpected content type %q", got)
	}
	if !strings.Contains(response.Body.String(), "Book a PawIt VetCare visit") {
		t.Fatal("expected booking page headline")
	}
	if !strings.Contains(response.Body.String(), "/api/v1/public/slots") {
		t.Fatal("expected booking page to load public slots")
	}
}

func TestSlots(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/public/slots", nil)

	newServer().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload struct {
		Items []map[string]string `json:"items"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode slots response: %v", err)
	}
	if len(payload.Items) != 1 {
		t.Fatalf("expected 1 slot, got %d", len(payload.Items))
	}
	if payload.Items[0]["id"] != "slot_001" {
		t.Fatalf("unexpected slot id %q", payload.Items[0]["id"])
	}
}
