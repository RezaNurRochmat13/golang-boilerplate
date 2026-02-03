package note

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrNoteNotFound = errors.New("Note not found")

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) GetAllNotes() ([]Note, error) {
	var notes []Note
	if err := r.DB.Find(&notes).Error; err != nil {
		return nil, fmt.Errorf("repo: get all notes: %w", err)
	}

	return notes, nil
}

func (r *Repository) GetNote(id string) (Note, error) {
	var note Note

	if err := r.DB.First(&note, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Note{}, ErrNoteNotFound
		}
		return Note{}, fmt.Errorf("repo: get note %s: %w", id, err)
	}

	return note, nil
}

func (r *Repository) CreateNote(note *Note) error {
	if err := r.DB.Create(note).Error; err != nil {
		return fmt.Errorf("repo: create note: %w", err)
	}

	return nil
}

func (r *Repository) UpdateNote(id string, note *Note) error {
	var existing Note

	if err := r.DB.First(&existing, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoteNotFound
		}
		return fmt.Errorf("repo: find note %s for update: %w", id, err)
	}

	// copy fields explicitly
	existing.Title = note.Title
	existing.Content = note.Content
	existing.Text = note.Text

	if err := r.DB.Save(&existing).Error; err != nil {
		return fmt.Errorf("repo: update note %s: %w", id, err)
	}

	return nil
}

func (r *Repository) DeleteNote(id string) error {
	result := r.DB.Delete(&Note{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("repo: delete note %s: %w", id, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}

	return nil
}
