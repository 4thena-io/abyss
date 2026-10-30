package service

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type ServiceHealth struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type HealthStatus struct {
	Database *ServiceHealth `json:"database,omitempty"`
	Forge    ServiceHealth  `json:"forge"`
	CI       ServiceHealth  `json:"ci"`
}

// hostedCI types have no self-hosted endpoint to ping.
var hostedCI = map[string]bool{
	"github-actions": true,
	"gitlab-ci":      true,
}

type HealthService struct {
	forgeHost    string
	ciType       string
	ciHost       string
	client       *http.Client
	pingDatabase func(context.Context) error
	checkForge   func(context.Context) error
}

func NewHealthService(forgeHost, ciType, ciHost string) *HealthService {
	return &HealthService{
		forgeHost: forgeHost,
		ciType:    ciType,
		ciHost:    ciHost,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

// WithDependencies installs real checks before the service is published.
func (s *HealthService) WithDependencies(database, forge func(context.Context) error) *HealthService {
	s.pingDatabase, s.checkForge = database, forge
	return s
}

func (s *HealthService) Check(ctx context.Context) HealthStatus {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var status HealthStatus
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if s.checkForge != nil {
			status.Forge = dependencyHealth(ctx, s.checkForge)
		} else {
			status.Forge = s.ping(ctx, s.forgeHost)
		}
	}()
	go func() { defer wg.Done(); status.CI = s.checkCI(ctx) }()
	if s.pingDatabase != nil {
		wg.Add(1)
		go func() { defer wg.Done(); health := dependencyHealth(ctx, s.pingDatabase); status.Database = &health }()
	}
	wg.Wait()
	return status
}

func dependencyHealth(ctx context.Context, check func(context.Context) error) ServiceHealth {
	if check(ctx) != nil {
		return ServiceHealth{Error: "temporarily unavailable"}
	}
	return ServiceHealth{OK: true}
}

func (s *HealthService) Ready(ctx context.Context) (HealthStatus, bool) {
	status := s.Check(ctx)
	ready := status.Database != nil && status.Database.OK
	if s.forgeHost != "" {
		ready = ready && status.Forge.OK
	}
	if s.ciHost != "" || hostedCI[s.ciType] {
		ready = ready && status.CI.OK
	}
	return status, ready
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
		return ServiceHealth{OK: false, Error: "invalid endpoint"}
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return ServiceHealth{OK: false, Error: "unreachable"}
	}
	_ = resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ServiceHealth{Error: "unexpected HTTP status"}
	}
	return ServiceHealth{OK: true}
}
