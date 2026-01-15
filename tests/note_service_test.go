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
	id := uuid.New()

	expected := []note.Note{
		{ID: id, Title: "Note 1"},
	}

	repo := &mockRepository{
		getAllNotesFn: func() ([]note.Note, error) {
			return expected, nil
		},
	}

	service := note.NewService(repo)

	result, err := service.GetAllNotes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestService_GetNote(t *testing.T) {
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

	service := note.NewService(repo)

	note, err := service.GetNote("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if note != expected {
		t.Fatalf("expected %v, got %v", expected, note)
	}
}

func TestService_CreateNote(t *testing.T) {
	repo := &mockRepository{
		createNoteFn: func(note *note.Note) error {
			return nil
		},
	}

	service := note.NewService(repo)

	err := service.CreateNote(&note.Note{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_UpdateNote(t *testing.T) {
	repo := &mockRepository{
		updateNoteFn: func(id string, note *note.Note) error {
			return nil
		},
	}

	service := note.NewService(repo)

	err := service.UpdateNote("1", &note.Note{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_DeleteNote(t *testing.T) {
	repo := &mockRepository{
		deleteNoteFn: func(id string) error {
			return nil
		},
	}

	service := note.NewService(repo)

	err := service.DeleteNote("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_DeleteNote_Error(t *testing.T) {
	expectedErr := errors.New("delete failed")

	repo := &mockRepository{
		deleteNoteFn: func(id string) error {
			return expectedErr
		},
	}

	service := note.NewService(repo)

	err := service.DeleteNote("1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != expectedErr {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
