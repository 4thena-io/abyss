package service

import (
	"context"
	"net/http"
	"time"
)

type ServiceHealth struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type HealthStatus struct {
	Forge ServiceHealth `json:"forge"`
	CI    ServiceHealth `json:"ci"`
}

// hostedCI types have no self-hosted endpoint to ping.
var hostedCI = map[string]bool{
	"github-actions": true,
	"gitlab-ci":      true,
}

type HealthService struct {
	forgeHost string
	ciType    string
	ciHost    string
	client    *http.Client
}

func NewHealthService(forgeHost, ciType, ciHost string) *HealthService {
	return &HealthService{
		forgeHost: forgeHost,
		ciType:    ciType,
		ciHost:    ciHost,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *HealthService) Check(ctx context.Context) HealthStatus {
	return HealthStatus{
		Forge: s.ping(ctx, s.forgeHost),
		CI:    s.checkCI(ctx),
	}
}

func (s *HealthService) checkCI(ctx context.Context) ServiceHealth {
	if hostedCI[s.ciType] {
		return ServiceHealth{OK: true}
	}
	return s.ping(ctx, s.ciHost)
}

func (s *HealthService) ping(ctx context.Context, host string) ServiceHealth {
	if host == "" {
		return ServiceHealth{OK: false, Error: "not configured"}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host, nil)
	if err != nil {
		return ServiceHealth{OK: false, Error: err.Error()}
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return ServiceHealth{OK: false, Error: err.Error()}
	}
	resp.Body.Close()
	return ServiceHealth{OK: true}
}
