package tests

import (
	"bytes"
	"golang-boilerplate-example/internal/auth"
	"net/http"
	"net/http/httptest"
)

func authRequest(method, url string, body []byte) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")

	// generate token
	token, _ := auth.GenerateToken(1)

	req.Header.Set("Authorization", "Bearer "+token)

	return req
}
