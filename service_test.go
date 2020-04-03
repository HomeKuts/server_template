package server_template

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"os"
)

const ORIGIN_VALID = "0.0.0.0:4200"
const ORIGIN_INVALID = "0.0.0.0:4201"

func TestMain(m *testing.M) {
	Config()
	os.Exit(m.Run())
}

// Перевірка запиту з корректним Origin 
func TestValidOrigin(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", ORIGIN_VALID)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Перевірка запиту з некорректним Origin 
func TestInValidOrigin(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", ORIGIN_INVALID)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// Перевірка запиту з неснуючим маршрутом 
func TestStatusNotFound(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notfound", nil)
	req.Header.Set("Origin", ORIGIN_VALID)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Перевірка запиту з маршрутом: /info 
func TestInfo(t *testing.T) {
        router := setupRouter()

        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/info", nil)
		req.Header.Set("Origin", ORIGIN_VALID)
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
        
        b := w.Body.String()        
        assert.Contains(t, b, "ver")
}
