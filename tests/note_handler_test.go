package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func createNote(t *testing.T, app *fiber.App, title string) string {
	t.Helper()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/notes",
		bytes.NewBuffer([]byte(`{"title":"`+title+`"}`)),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("create note failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// fetch notes to get ID
	getReq := httptest.NewRequest(http.MethodGet, "/api/notes", nil)
	getResp, _ := app.Test(getReq)

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	_ = json.NewDecoder(getResp.Body).Decode(&result)

	if len(result.Data) == 0 {
		t.Fatal("no note created")
	}

	return result.Data[len(result.Data)-1].ID
}

func TestGetAllNotes_Integration(t *testing.T) {
	app := SetupTestApp()

	createNote(t, app, "Note 1")

	req := httptest.NewRequest(http.MethodGet, "/api/notes", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var result struct {
		Data []interface{} `json:"data"`
	}

	_ = json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Data) == 0 {
		t.Fatal("expected notes, got empty")
	}
}

func TestGetNoteByID_Integration(t *testing.T) {
	app := SetupTestApp()

	id := createNote(t, app, "Detail Note")

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/notes/"+id,
		nil,
	)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetNote_InvalidID(t *testing.T) {
	app := SetupTestApp()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/notes/invalid-id",
		nil,
	)

	resp, _ := app.Test(req)

	if resp.StatusCode == fiber.StatusOK {
		t.Fatal("expected error for invalid id")
	}
}

func TestCreateNote_Integration(t *testing.T) {
	app := SetupTestApp()

	payload := map[string]string{
		"title": "Integration Note",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/notes",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestCreateNote_InvalidBody(t *testing.T) {
	app := SetupTestApp()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/notes",
		bytes.NewBuffer([]byte(`invalid-json`)),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUpdateNote_Integration(t *testing.T) {
	app := SetupTestApp()

	id := createNote(t, app, "Old Title")

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/notes/"+id,
		bytes.NewBuffer([]byte(`{"title":"Updated Title"}`)),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestUpdateNoteInvalidPayload_Integration(t *testing.T) {
	app := SetupTestApp()

	id := createNote(t, app, "Old Title")

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/notes/"+id,
		bytes.NewBuffer([]byte(`{"title":""}`)),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestDeleteNote_Integration(t *testing.T) {
	app := SetupTestApp()

	id := createNote(t, app, "Delete Me")

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/notes/"+id,
		nil,
	)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// edge case: delete again → should fail
	resp2, _ := app.Test(req)
	if resp2.StatusCode == fiber.StatusOK {
		t.Fatal("expected error when deleting non-existent note")
	}
}

func TestDeleteNoteInvalidID_Integration(t *testing.T) {
	app := SetupTestApp()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/notes/invalid-id",
		nil,
	)

	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}
