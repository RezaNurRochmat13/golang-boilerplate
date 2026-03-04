package note

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidPayload = errors.New("invalid note payload")
)

const (
	cacheTTL         = 5 * time.Minute
	cacheKeyAllNotes = "notes:all"
	cacheKeyNoteByID = "notes:%s"
)

type Service struct {
	repo  RepositoryContract
	redis *redis.Client
}

func NewService(repo RepositoryContract, redisClient *redis.Client) *Service {
	return &Service{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *Service) GetAllNotes() ([]Note, error) {
	ctx := context.Background()

	// Fetch from redis
	cached, err := s.redis.Get(ctx, cacheKeyAllNotes).Result()
	if err == nil {
		var notes []Note
		if err := json.Unmarshal([]byte(cached), &notes); err == nil {
			return notes, nil
		}
	}

	// Fetch in DB when redis cache miss
	notes, err := s.repo.GetAllNotes()
	if err != nil {
		return nil, fmt.Errorf("service: get all notes: %w", err)
	}

	// Cache to redis
	bytes, _ := json.Marshal(notes)
	_ = s.redis.Set(ctx, cacheKeyAllNotes, bytes, cacheTTL).Err()

	return notes, nil
}

func (s *Service) GetNote(id string) (Note, error) {
	ctx := context.Background()
	key := fmt.Sprintf(cacheKeyNoteByID, id)

	// Fetch from cache
	cached, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var note Note
		if err := json.Unmarshal([]byte(cached), &note); err == nil {
			return note, nil
		}
	}

	// Fetch from cache when cache miss
	note, err := s.repo.GetNote(id)
	if err != nil {
		if errors.Is(err, ErrNoteNotFound) {
			return Note{}, ErrNoteNotFound
		}
		return Note{}, fmt.Errorf("service: get note %s: %w", id, err)
	}

	// Cache to redis
	bytes, _ := json.Marshal(note)
	_ = s.redis.Set(ctx, key, bytes, cacheTTL).Err()

	return note, nil
}

func (s *Service) CreateNote(note *Note) error {
	if note.Title == "" {
		return ErrInvalidPayload
	}

	if err := s.repo.CreateNote(note); err != nil {
		return fmt.Errorf("service: create note: %w", err)
	}

	ctx := context.Background()

	// Invalidate cache
	_ = s.redis.Del(ctx, cacheKeyAllNotes).Err()

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

	ctx := context.Background()

	// Invalidate specific note cache
	key := fmt.Sprintf(cacheKeyNoteByID, id)
	_ = s.redis.Del(ctx, key).Err()

	// Invalidate list cache
	_ = s.redis.Del(ctx, cacheKeyAllNotes).Err()

	return nil
}

func (s *Service) DeleteNote(id string) error {
	if err := s.repo.DeleteNote(id); err != nil {
		if errors.Is(err, ErrNoteNotFound) {
			return ErrNoteNotFound
		}
		return fmt.Errorf("service: delete note %s: %w", id, err)
	}

	ctx := context.Background()

	key := fmt.Sprintf(cacheKeyNoteByID, id)

	// Invalidate cache
	_ = s.redis.Del(ctx, key).Err()
	_ = s.redis.Del(ctx, cacheKeyAllNotes).Err()

	return nil
}
