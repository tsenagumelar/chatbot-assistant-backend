package services

import (
	"police-assistant-backend/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SessionStore menyimpan chat history per session
type SessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// Session menyimpan history dan metadata per session
type Session struct {
	ID        string
	History   []models.OpenAIMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	sessionStore *SessionStore
	once         sync.Once
)

// GetSessionStore returns singleton instance of SessionStore
func GetSessionStore() *SessionStore {
	once.Do(func() {
		sessionStore = &SessionStore{
			sessions: make(map[string]*Session),
		}
		// Start cleanup goroutine
		go sessionStore.cleanupExpiredSessions()
	})
	return sessionStore
}

// CreateSession membuat session baru
func (s *SessionStore) CreateSession() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID := uuid.New().String()
	s.sessions[sessionID] = &Session{
		ID:        sessionID,
		History:   []models.OpenAIMessage{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return sessionID
}

// GetSession mengambil session berdasarkan ID
func (s *SessionStore) GetSession(sessionID string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	return session, exists
}

// AddMessage menambahkan pesan ke history session
func (s *SessionStore) AddMessage(sessionID string, role string, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		// Jika session tidak ada, buat baru
		session = &Session{
			ID:        sessionID,
			History:   []models.OpenAIMessage{},
			CreatedAt: time.Now(),
		}
		s.sessions[sessionID] = session
	}

	session.History = append(session.History, models.OpenAIMessage{
		Role:    role,
		Content: content,
	})
	session.UpdatedAt = time.Now()

	// Batasi history max 30 pesan (15 exchange)
	if len(session.History) > 30 {
		session.History = session.History[len(session.History)-30:]
	}

	return nil
}

// GetHistory mengambil history dari session
func (s *SessionStore) GetHistory(sessionID string) []models.OpenAIMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return []models.OpenAIMessage{}
	}

	return session.History
}

// ClearSession menghapus history session
func (s *SessionStore) ClearSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if session, exists := s.sessions[sessionID]; exists {
		session.History = []models.OpenAIMessage{}
		session.UpdatedAt = time.Now()
	}
}

// DeleteSession menghapus session
func (s *SessionStore) DeleteSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sessionID)
}

// GetSessionCount mengembalikan jumlah session aktif
func (s *SessionStore) GetSessionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.sessions)
}

// cleanupExpiredSessions membersihkan session yang sudah tidak aktif > 24 jam
func (s *SessionStore) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, session := range s.sessions {
			// Hapus session yang tidak aktif > 24 jam
			if now.Sub(session.UpdatedAt) > 24*time.Hour {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}
