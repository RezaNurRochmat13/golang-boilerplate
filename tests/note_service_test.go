package tests

import (
	"errors"
	"golang-boilerplate-example/module/note"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

type mockRepository struct {
	getAllNotesFn func() ([]note.Note, error)
	getNoteFn     func(id string) (note.Note, error)
	createNoteFn  func(note *note.Note) error
	updateNoteFn  func(id string, note *note.Note) error
	deleteNoteFn  func(id string) error
}

func (m *mockRepository) GetAllNotes() ([]note.Note, error) {
	return m.getAllNotesFn()
}

func (m *mockRepository) GetNote(id string) (note.Note, error) {
	return m.getNoteFn(id)
}

func (m *mockRepository) CreateNote(note *note.Note) error {
	return m.createNoteFn(note)
}

func (m *mockRepository) UpdateNote(id string, note *note.Note) error {
	return m.updateNoteFn(id, note)
}

func (m *mockRepository) DeleteNote(id string) error {
	return m.deleteNoteFn(id)
}

func TestService_GetAllNotes(t *testing.T) {
	redisClient := setupRedis(t)

	id := uuid.New()

	expected := []note.Note{
		{ID: id, Title: "Note 1"},
	}

	repo := &mockRepository{
		getAllNotesFn: func() ([]note.Note, error) {
			return expected, nil
		},
	}

	service := note.NewService(repo, redisClient)

	result, err := service.GetAllNotes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestService_GetNote(t *testing.T) {
	redisClient := setupRedis(t)

	id := uuid.New()

	expected := note.Note{ID: id, Title: "Test"}

	repo := &mockRepository{
		getNoteFn: func(id string) (note.Note, error) {
			if id != "1" {
				t.Fatalf("unexpected id: %s", id)
			}
			return expected, nil
		},
	}

	service := note.NewService(repo, redisClient)

	note, err := service.GetNote("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if note != expected {
		t.Fatalf("expected %v, got %v", expected, note)
	}
}

func TestService_CreateNote_Success(t *testing.T) {
	redisClient := setupRedis(t)

	repo := &mockRepository{
		createNoteFn: func(note *note.Note) error {
			return nil
		},
	}

	service := note.NewService(repo, redisClient)

	err := service.CreateNote(&note.Note{
		Title: "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_CreateNote_InvalidPayload(t *testing.T) {
	redisClient := setupRedis(t)

	repo := &mockRepository{}

	service := note.NewService(repo, redisClient)

	err := service.CreateNote(&note.Note{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, note.ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}

func TestService_UpdateNote_Success(t *testing.T) {
	redisClient := setupRedis(t)

	repo := &mockRepository{
		updateNoteFn: func(id string, note *note.Note) error {
			return nil
		},
	}

	service := note.NewService(repo, redisClient)

	err := service.UpdateNote("1", &note.Note{
		Title: "updated",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_UpdateNote_NotFound(t *testing.T) {
	redisClient := setupRedis(t)

	repo := &mockRepository{
		updateNoteFn: func(id string, input *note.Note) error {
			return note.ErrNoteNotFound
		},
	}

	service := note.NewService(repo, redisClient)

	err := service.UpdateNote("1", &note.Note{
		Title: "updated",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, note.ErrNoteNotFound) {
		t.Fatalf("expected ErrNoteNotFound, got %v", err)
	}
}

func TestService_UpdateNote_InvalidPayload(t *testing.T) {
	redisClient := setupRedis(t)

	repo := &mockRepository{}

	service := note.NewService(repo, redisClient)

	err := service.UpdateNote("1", &note.Note{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, note.ErrInvalidPayload) {
		t.Fatalf("expected ErrInvalidPayload, got %v", err)
	}
}

func TestService_DeleteNote_Success(t *testing.T) {
	redisClient := setupRedis(t)

	repo := &mockRepository{
		deleteNoteFn: func(id string) error {
			return nil
		},
	}

	service := note.NewService(repo, redisClient)

	err := service.DeleteNote("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_DeleteNote_NotFound(t *testing.T) {
	redisClient := setupRedis(t)

	repo := &mockRepository{
		deleteNoteFn: func(id string) error {
			return note.ErrNoteNotFound
		},
	}

	service := note.NewService(repo, redisClient)

	err := service.DeleteNote("1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, note.ErrNoteNotFound) {
		t.Fatalf("expected ErrNoteNotFound, got %v", err)
	}
}

func TestService_DeleteNote_OtherError(t *testing.T) {
	repoErr := errors.New("db connection lost")
	redisClient := setupRedis(t)

	repo := &mockRepository{
		deleteNoteFn: func(id string) error {
			return repoErr
		},
	}

	service := note.NewService(repo, redisClient)

	err := service.DeleteNote("1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// harus tetap bisa di-unwarp
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}
