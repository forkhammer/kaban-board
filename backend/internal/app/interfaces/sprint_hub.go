package interfaces

import (
	domain "main/internal/domain/models"

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
	ID        domain.IssueBindingId `json:"id"`
	AccountID domain.AccountId      `json:"account_id"`
}

type BindingCreatedEvent struct {
	ID        domain.IssueBindingId `json:"id"`
	AccountID domain.AccountId      `json:"account_id"`
}

type BindingDeletedEvent struct {
	ID       uint            `json:"id"`
	SprintID domain.SprintId `json:"sprint_id"`
}

type WSEvent struct {
	Type WSEventType `json:"type"`
	Data any         `json:"data,omitempty"`
}

type SprintHub interface {
	Register(sprintId domain.SprintId, conn *websocket.Conn)
	Unregister(sprintId domain.SprintId, conn *websocket.Conn)
	Broadcast(sprintId domain.SprintId, event WSEvent)
}
