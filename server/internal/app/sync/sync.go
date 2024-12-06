package syncmanager

import (
	"sync"

	"google.golang.org/grpc"
)

// SyncManager управляет потоками для пользователей.
type SyncManager[T any] struct {
	mu      sync.Mutex
	streams map[string][]grpc.ServerStream // Потоки по userID
}

// NewSyncManager создает новый экземпляр SyncManager.
func NewSyncManager[T any]() *SyncManager[T] {
	return &SyncManager[T]{
		streams: make(map[string][]grpc.ServerStream),
	}
}

// AddStream добавляет поток для userID.
func (sm *SyncManager[T]) AddStream(userID string, stream grpc.ServerStream) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.streams[userID] = append(sm.streams[userID], stream)
}

// RemoveStream удаляет поток для userID.
func (sm *SyncManager[T]) RemoveStream(userID string, stream grpc.ServerStream) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	streams := sm.streams[userID]
	for i, s := range streams {
		if s == stream {
			sm.streams[userID] = append(streams[:i], streams[i+1:]...)
			break
		}
	}
	if len(sm.streams[userID]) == 0 {
		delete(sm.streams, userID)
	}
}

// Broadcast отправляет сообщение всем потокам userID.
func (sm *SyncManager[T]) Broadcast(userID string, message T) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for _, stream := range sm.streams[userID] {
		stream.SendMsg(message) // Отправляем сообщение в поток
	}
}
