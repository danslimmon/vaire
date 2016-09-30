package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func assertAuthFailed(t *testing.T, w *httptest.ResponseRecorder) bool {
	var err error

	if w.Code != 403 {
		t.Logf("Expected status code 403; got %d instead", w.Code)
		return false
	}
	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Logf("Failed to read response body; error: %s", err)
		return false
	}
	var body struct{ Error string }
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		t.Logf("Failed to parse response body; error: %s\n\nBODY:\n\n%s", err, string(bodyBytes))
		return false
	}
	if body.Error == "" {
		t.Logf("Response missing error message")
		return false
	}
	return true
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	Config.Token = "abc123"
	router := gin.New()
	router.Use(middlewareCheckAuth)
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "{\"blah\":\"foo\"}")
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, r)
	if !assertAuthFailed(t, w) {
		t.Log("Auth middleware failed to rebuff request with no Authorization header")
		t.Fail()
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/test", nil)
	r.Header["Authorization"] = []string{"Token "}
	router.ServeHTTP(w, r)
	if !assertAuthFailed(t, w) {
		t.Log("Auth middleware failed to rebuff request with malformed Authorization header")
		t.Fail()
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/test", nil)
	r.Header["Authorization"] = []string{"Token def456"}
	router.ServeHTTP(w, r)
	if !assertAuthFailed(t, w) {
		t.Log("Auth middleware failed to rebuff request with incorrect token")
		t.Fail()
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/test", nil)
	r.Header["Authorization"] = []string{"Token abc123"}
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Log("Auth middleware failed to allow request with correct token")
		t.Fail()
	}
}
