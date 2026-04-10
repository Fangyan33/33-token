package access

import "errors"

var ErrQuotaBlocked = errors.New("quota blocked")

// Summary 表示统一准入校验所需的最小账户摘要。
type Summary struct {
	QuotaRemaining int64
	Status         string
}

// Check 根据最小摘要判断是否允许准入。
func Check(summary Summary) error {
	if summary.Status != "active" || summary.QuotaRemaining <= 0 {
		return ErrQuotaBlocked
	}
	return nil
}
