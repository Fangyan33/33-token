package app

import (
	"net/http"
	"time"

	billinghttp "github.com/33-token/model-api-platform/services/controlplane/internal/billing/http"
	"github.com/33-token/model-api-platform/services/controlplane/internal/httpx"
)

func NewServer() *http.Server {
	mux := http.NewServeMux()
	billingHandler := billinghttp.NewHandler(billinghttp.NewBillingService())

	mux.HandleFunc("/healthz", httpx.Healthz)
	mux.HandleFunc("/api/v1/auth/", httpx.Placeholder("auth"))
	mux.HandleFunc("/api/v1/accounts/", httpx.Placeholder("accounts"))
	mux.Handle("/api/v1/billing/", billingHandler)
	mux.HandleFunc("/api/v1/admin/", httpx.Placeholder("admin"))
	mux.HandleFunc("/api/v1/gateway/", httpx.Placeholder("gateway"))

	return &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
