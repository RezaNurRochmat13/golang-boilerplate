package note

import "gorm.io/gorm"

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
	err := r.DB.Find(&notes).Error
	return notes, err
}

func (r *Repository) GetNote(id string) (Note, error) {
	var note Note
	err := r.DB.First(&note, "id = ?", id).Error
	return note, err
}

func (r *Repository) CreateNote(note *Note) error {
	return r.DB.Create(note).Error
}

func (r *Repository) UpdateNote(id string, note *Note) error {
	var existing Note

	if err := r.DB.First(&existing, "id = ?", id).Error; err != nil {
		return err
	}

	// copy fields explicitly
	existing.Title = note.Title
	existing.Content = note.Content
	existing.Text = note.Text

	return r.DB.Save(&existing).Error
}

func (r *Repository) DeleteNote(id string) error {
	result := r.DB.Delete(&Note{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
