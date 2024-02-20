package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
	"web-frame/session"
)

var (
	errorKeyNotFound        = errors.New("session: key 找不到")
	errorKeySessionNotFound = errors.New("session: session 找不到")
)

type Store struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewStore(client redis.Cmdable) *Store {
	return &Store{
		expiration: time.Minute * 15,
		client:     client,
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	_, err := s.client.HSet(ctx, id, id, id).Result()
	if err != nil {
		return nil, err
	}
	_, err = s.client.Expire(ctx, id, s.expiration).Result()
	if err != nil {
		return nil, err
	}
	return &Session{
		id:     id,
		client: s.client,
	}, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	ok, err := s.client.Expire(ctx, id, s.expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("session: 该 id 对应 session 不存在")
	}
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	_, err := s.client.Del(ctx, id).Result()
	return err
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	cnt, err := s.client.Exists(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, errorKeySessionNotFound
	}
	return &Session{
		id:     id,
		client: s.client,
	}, nil
}

type Session struct {
	client redis.Cmdable
	id     string
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	val, err := s.client.HGet(ctx, s.id, key).Result()
	return val, err
}

func (s *Session) Set(ctx context.Context, key string, val any) error {
	const lua = `
if redis.call("exists", KEYS[1])
then
	return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
else
	return -1
end
`
	res, err := s.client.Eval(ctx, lua, []string{s.id}, key, val).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return errorKeyNotFound
	}
	return nil
}

func (s *Session) ID() string {
	return s.id
}
