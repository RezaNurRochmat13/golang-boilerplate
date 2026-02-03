package note

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidPayload = errors.New("invalid note payload")
)

type Service struct {
	repo RepositoryContract
}

func NewService(repo RepositoryContract) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAllNotes() ([]Note, error) {
	notes, err := s.repo.GetAllNotes()

	if err != nil {
		return nil, fmt.Errorf("service: get all notes: %w", err)
	}

	return notes, nil
}

func (s *Service) GetNote(id string) (Note, error) {
	note, err := s.repo.GetNote(id)

	if err != nil {
		if errors.Is(err, ErrNoteNotFound) {
			return Note{}, ErrNoteNotFound
		}
		return Note{}, fmt.Errorf("service: get note %s: %w", id, err)
	}

	return note, nil
}

func (s *Service) CreateNote(note *Note) error {
	if note.Title == "" {
		return ErrInvalidPayload
	}

	if err := s.repo.CreateNote(note); err != nil {
		return fmt.Errorf("service: create note: %w", err)
	}

	return nil
}

func (s *Service) UpdateNote(id string, note *Note) error {
	if note.Title == "" {
		return ErrInvalidPayload
	}

	if err := s.repo.UpdateNote(id, note); err != nil {
		if errors.Is(err, ErrNoteNotFound) {
			return ErrNoteNotFound
		}
		return fmt.Errorf("service: update note %s: %w", id, err)
	}

	return nil
}

func (s *Service) DeleteNote(id string) error {
	if err := s.repo.DeleteNote(id); err != nil {
		if errors.Is(err, ErrNoteNotFound) {
			return ErrNoteNotFound
		}
		return fmt.Errorf("service: delete note %s: %w", id, err)
	}

	return nil
}
