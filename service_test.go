package server_template

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ORIGIN_VALID = "0.0.0.0:4200"
const ORIGIN_INVALID = "0.0.0.0:4201"

func TestValidOrigin(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", ORIGIN_VALID)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInValidOrigin(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", ORIGIN_INVALID)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestStatusNotFound(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/notfound", nil)
	req.Header.Set("Origin", ORIGIN_VALID)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

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
