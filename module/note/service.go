package note

type Service struct {
	repo RepositoryContract
}

func NewService(repo RepositoryContract) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAllNotes() ([]Note, error) {
	return s.repo.GetAllNotes()
}

func (s *Service) GetNote(id string) (Note, error) {
	return s.repo.GetNote(id)
}

func (s *Service) CreateNote(note *Note) error {
	return s.repo.CreateNote(note)
}

func (s *Service) UpdateNote(id string, note *Note) error {
	return s.repo.UpdateNote(id, note)
}

func (s *Service) DeleteNote(id string) error {
	return s.repo.DeleteNote(id)
}
