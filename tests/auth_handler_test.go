package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterUser_Integration(t *testing.T) {

	app := SetupTestApp()

	payload := map[string]string{
		"email":    "test@mail.com",
		"name":     "Reja",
		"password": "password123",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status 201 got %d", resp.StatusCode)
	}
}

func TestRegisterUser_InvalidBody(t *testing.T) {

	app := SetupTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewBuffer([]byte(`invalid-json`)),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}

func TestLoginUser_Integration(t *testing.T) {

	app := SetupTestApp()

	// register user dulu
	registerPayload := map[string]string{
		"email":    "login@mail.com",
		"name":     "Login User",
		"password": "password123",
	}

	registerBody, _ := json.Marshal(registerPayload)

	registerReq := httptest.NewRequest(
		http.MethodPost,
		"/register",
		bytes.NewBuffer(registerBody),
	)

	registerReq.Header.Set("Content-Type", "application/json")

	_, _ = app.Test(registerReq)

	// login
	loginPayload := map[string]string{
		"email":    "login@mail.com",
		"password": "password123",
	}

	loginBody, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewBuffer(loginBody),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200 got %d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}

	_ = json.NewDecoder(resp.Body).Decode(&result)

	if result.Token == "" {
		t.Fatal("expected jwt token")
	}
}

func TestLoginUser_InvalidCredentials(t *testing.T) {

	app := SetupTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewBuffer([]byte(`{
			"email":"wrong@mail.com",
			"password":"wrong"
		}`)),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", resp.StatusCode)
	}
}

func TestLoginUser_InvalidBody(t *testing.T) {

	app := SetupTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/login",
		bytes.NewBuffer([]byte(`invalid-json`)),
	)

	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400 got %d", resp.StatusCode)
	}
}
