package store

import (
	"encoding/json"
	"time"
)

type UpstreamCredentialRef struct {
	ID            string
	Provider      string
	CredentialKey string
	DisplayName   string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ModelRoute struct {
	ID                      string
	Protocol                string
	PublicModelName         string
	UpstreamProvider        string
	UpstreamModelID         string
	UpstreamCredentialRefID string
	Status                  string
	RequestAdapterType      string
	ResponseAdapterType     string
	Priority                int
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type AdminAuditLog struct {
	ID             string
	OperatorID     string
	OperatorType   string
	ActionType     string
	TargetType     string
	TargetID       string
	BeforeSnapshot json.RawMessage
	AfterSnapshot  json.RawMessage
	Reason         string
	CreatedAt      time.Time
}

type AdminRepository interface {
	ListModelRoutes() ([]ModelRoute, error)
	ListCredentialRefs() ([]UpstreamCredentialRef, error)
	AppendAuditLog(log AdminAuditLog) (AdminAuditLog, error)
}
