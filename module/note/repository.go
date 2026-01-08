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
	err := r.DB.First(&note, "id = ?", id).Error
	if err != nil {
		return err
	}

	return r.DB.Save(note).Error
}

func (r *Repository) DeleteNote(id string) error {
	return r.DB.Delete(&Note{}, "id = ?", id).Error
}
