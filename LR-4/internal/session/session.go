package session

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type Session struct {
	ClientID    uint   `json:"client_id"`
	Username    string `json:"username"`
	IsModerator bool   `json:"is_moderator"`
}

type Manager struct {
	client *redis.Client
	ctx    context.Context
}

func NewSessionManager() *Manager {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       0,
	})

	ctx := context.Background()

	// Проверяем подключение
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("⚠️ Redis connection failed: %v\n", err)
	} else {
		fmt.Printf("✅ Redis client initialized successfully\n")
	}

	return &Manager{
		client: client,
		ctx:    ctx,
	}
}

func (m *Manager) CreateSession(sessionID string, session Session, expiration time.Duration) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return m.client.Set(m.ctx, "session:"+sessionID, data, expiration).Err()
}

func (m *Manager) GetSession(sessionID string) (*Session, error) {
	data, err := m.client.Get(m.ctx, "session:"+sessionID).Result()
	if err != nil {
		return nil, err
	}

	var session Session
	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (m *Manager) DeleteSession(sessionID string) error {
	return m.client.Del(m.ctx, "session:"+sessionID).Err()
}

func (m *Manager) GetAllSessions() (map[string]Session, error) {
	keys, err := m.client.Keys(m.ctx, "session:*").Result()
	if err != nil {
		return nil, err
	}

	sessions := make(map[string]Session)
	for _, key := range keys {
		data, err := m.client.Get(m.ctx, key).Result()
		if err != nil {
			continue
		}

		var session Session
		if json.Unmarshal([]byte(data), &session) == nil {
			sessions[key] = session
		}
	}

	return sessions, nil
}
