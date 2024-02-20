package memory

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
	"web-frame/session"
)

var (
	errorKeyNotFound        = errors.New("session: key 找不到")
	errorKeySessionNotFound = errors.New("session: session 找不到")
)

type Store struct {
	sessions   *cache.Cache
	expiration time.Duration
}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		sessions:   cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	sess := &Session{
		values: sync.Map{},
		id:     id,
	}
	s.sessions.Set(id, sess, s.expiration)
	return sess, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	val, ok := s.sessions.Get(id)
	if !ok {
		return errors.New("session: 该 id 对应 session 不存在")
	}
	s.sessions.Set(id, val, s.expiration)
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.sessions.Delete(id)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	sess, ok := s.sessions.Get(id)
	if !ok {
		return nil, errorKeySessionNotFound
	}
	return sess.(*Session), nil
}

type Session struct {
	values sync.Map
	id     string
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	val, ok := s.values.Load(key)
	if !ok {
		return nil, errorKeyNotFound
	}
	return val, nil
}

func (s *Session) Set(ctx context.Context, key string, val any) error {
	s.values.Store(key, val)
	return nil
}

func (s *Session) ID() string {
	return s.id
}
