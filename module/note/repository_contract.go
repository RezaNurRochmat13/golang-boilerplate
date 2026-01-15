package note

type RepositoryContract interface {
	GetAllNotes() ([]Note, error)
	GetNote(id string) (Note, error)
	CreateNote(note *Note) error
	UpdateNote(id string, note *Note) error
	DeleteNote(id string) error
}
