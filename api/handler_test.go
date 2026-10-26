package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeDB struct{ err error }

func (f fakeDB) Ping(context.Context) error { return f.err }

func get(t *testing.T, db pinger, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	newMux(db).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestHealth(t *testing.T) {
	// la liveness reste OK même si la base est indisponible
	rec := get(t, fakeDB{err: errors.New("down")}, "/api/health")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"ok"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestDBHealthUp(t *testing.T) {
	rec := get(t, fakeDB{}, "/api/health/db")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"database":"ok"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestDBHealthDown(t *testing.T) {
	rec := get(t, fakeDB{err: errors.New("connection refused")}, "/api/health/db")
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"database":"down"`) {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
