package store

import "time"

type Account struct {
	ID                  string
	Status              string
	DisplayName         string
	DefaultContactEmail string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type UserIdentity struct {
	ID           string
	AccountID    string
	LoginEmail   string
	AuthProvider string
	Status       string
	LastLoginAt  *time.Time
	CreatedAt    time.Time
}

type APIKey struct {
	ID         string
	AccountID  string
	KeyPrefix  string
	KeyHash    string
	Status     string
	CreatedAt  time.Time
	DisabledAt *time.Time
	LastUsedAt *time.Time
}

type AccountsRepository interface {
	GetAccount(accountID string) (Account, error)
	GetUserIdentity(accountID string) (UserIdentity, error)
	ListAPIKeys(accountID string) ([]APIKey, error)
}
