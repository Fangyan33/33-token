package store

import "time"

type AccountSubscriptionState struct {
	ID                         string
	AccountID                  string
	CurrentOrderID             string
	CurrentPlanPriceSnapshotID string
	Status                     string
	BillingPeriodStart         *time.Time
	BillingPeriodEnd           *time.Time
	ActivatedAt                *time.Time
	ExpiredAt                  *time.Time
	PausedAt                   *time.Time
	UpdatedAt                  time.Time
}

type AccountCycleSummary struct {
	ID                 string
	AccountID          string
	BillingPeriodStart time.Time
	BillingPeriodEnd   time.Time
	QuotaTotal         int64
	QuotaUsed          int64
	QuotaRemaining     int64
	InputTokensTotal   int64
	OutputTokensTotal  int64
	TotalTokensTotal   int64
	Status             string
	UpdatedAt          time.Time
}

type UsageEvent struct {
	ID                 string
	RequestID          string
	AccountID          string
	APIKeyID           string
	Protocol           string
	ModelName          string
	UpstreamProvider   string
	UpstreamModelID    string
	RequestStartedAt   time.Time
	RequestFinishedAt  time.Time
	LatencyMS          int64
	ResultStatus       string
	ErrorType          string
	InputTokens        int64
	OutputTokens       int64
	TotalTokens        int64
	BillingPeriodStart *time.Time
	BillingPeriodEnd   *time.Time
	CreatedAt          time.Time
}

type BillingEvent struct {
	ID                   string
	IdempotencyKey       string
	AccountID            string
	UsageEventID         string
	BillingPeriodStart   *time.Time
	BillingPeriodEnd     *time.Time
	SettlementType       string
	QuotaDelta           int64
	BeforeQuotaRemaining int64
	AfterQuotaRemaining  int64
	ResultStatus         string
	FailureReason        string
	CreatedAt            time.Time
	SettledAt            *time.Time
}

type SummaryRepository interface {
	GetSubscriptionState(accountID string) (AccountSubscriptionState, error)
	GetCycleSummary(accountID string) (AccountCycleSummary, error)
	AppendUsageEvent(event UsageEvent) (UsageEvent, error)
	AppendBillingEvent(event BillingEvent) (BillingEvent, error)
	UpdateCycleSummary(summary AccountCycleSummary) (AccountCycleSummary, error)
}
