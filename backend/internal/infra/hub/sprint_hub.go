package hub

import (
	"encoding/json"
	domain "main/internal/domain/models"
	"sync"

	"github.com/gorilla/websocket"
)

type WSEventType string

const (
	EventBindingUpdated  WSEventType = "binding_updated"
	EventBindingCreated  WSEventType = "binding_created"
	EventBindingDeleted  WSEventType = "binding_deleted"
	EventBindingOrdering WSEventType = "binding_ordering"
)

type BindingUpdatedEvent struct {
	ID domain.IssueBindingId `json:"id"`
}

type BindingCreatedEvent struct {
	ID domain.IssueBindingId `json:"id"`
}

type BindingDeletedEvent struct {
	ID       uint             `json:"id"`
	SprintID domain.SprintId  `json:"sprint_id"`
}

type WSEvent struct {
	Type WSEventType `json:"type"`
	Data any         `json:"data,omitempty"`
}

type SprintHub struct {
	mu    sync.RWMutex
	rooms map[domain.SprintId]map[*websocket.Conn]struct{}
}

func NewSprintHub() *SprintHub {
	return &SprintHub{
		rooms: make(map[domain.SprintId]map[*websocket.Conn]struct{}),
	}
}

func (h *SprintHub) Register(sprintId domain.SprintId, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[sprintId] == nil {
		h.rooms[sprintId] = make(map[*websocket.Conn]struct{})
	}
	h.rooms[sprintId][conn] = struct{}{}
}

func (h *SprintHub) Unregister(sprintId domain.SprintId, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, ok := h.rooms[sprintId]; ok {
		delete(room, conn)
		if len(room) == 0 {
			delete(h.rooms, sprintId)
		}
	}
}

func (h *SprintHub) Broadcast(sprintId domain.SprintId, event WSEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	conns := h.getRoomConns(sprintId)

	var deads []*websocket.Conn
	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			deads = append(deads, conn)
		}
	}

	h.deleteDeadConns(sprintId, deads)
}

func (h *SprintHub) getRoomConns(sprintId domain.SprintId) []*websocket.Conn {
	h.mu.RLock()
	room := h.rooms[sprintId]
	conns := make([]*websocket.Conn, 0, len(room))
	for conn := range room {
		conns = append(conns, conn)
	}
	h.mu.RUnlock()
	return conns
}

func (h *SprintHub) deleteDeadConns(sprintId domain.SprintId, deads []*websocket.Conn) {
	if len(deads) > 0 {
		h.mu.Lock()
		if room, ok := h.rooms[sprintId]; ok {
			for _, conn := range deads {
				delete(room, conn)
			}
		}
		h.mu.Unlock()
	}
}
