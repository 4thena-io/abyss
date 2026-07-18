package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/config"
)

type SetupHandler struct {
	configPath string
	configured bool
	defaults   *config.Config
	restartCh  chan<- struct{}
}

func NewSetupHandler(configPath string, configured bool, defaults *config.Config, restartCh chan<- struct{}) *SetupHandler {
	return &SetupHandler{configPath: configPath, configured: configured, defaults: defaults, restartCh: restartCh}
}

// Status reports whether the instance is configured and, while it isn't,
// echoes back the non-secret values already set via env vars/config file so
// the setup wizard can pre-fill them instead of showing blank fields.
func (h *SetupHandler) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := map[string]any{"configured": h.configured}
	if !h.configured && h.defaults != nil {
		res["defaults"] = map[string]string{
			"db_type":     h.defaults.Database.Type,
			"db_path":     h.defaults.Database.Path,
			"db_host":     h.defaults.Database.Host,
			"db_port":     h.defaults.Database.Port,
			"db_user":     h.defaults.Database.User,
			"db_name":     h.defaults.Database.Name,
			"forge_type":  h.defaults.Forge.Type,
			"forge_host":  h.defaults.Forge.Host,
			"forge_owner": h.defaults.Forge.Owner,
			"ci_type":     h.defaults.CI.Type,
			"ci_host":     h.defaults.CI.Host,
			"client_id":   h.defaults.Auth.ClientID,
		}
	}
	_ = json.NewEncoder(w).Encode(res)
}

type setupRequest struct {
	// Database
	DBType     string `json:"db_type"`
	DBPath     string `json:"db_path"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`

	// Forge
	ForgeType   string `json:"forge_type"`
	ForgeHost   string `json:"forge_host"`
	ForgeToken  string `json:"forge_token"`
	ForgeOwner  string `json:"forge_owner"`
	ForgeBranch string `json:"forge_branch"`

	// OAuth
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackURL  string `json:"callback_url"`

	// CI
	CIType  string `json:"ci_type"`
	CIHost  string `json:"ci_host"`
	CIToken string `json:"ci_token"`
}

func (h *SetupHandler) Configure(w http.ResponseWriter, r *http.Request) {
	if h.configured {
		response.Conflict(w, "instance is already configured")
		return
	}

	var req setupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.DBType == "" {
		req.DBType = "sqlite"
	}
	if req.DBType == "sqlite" && req.DBPath == "" {
		req.DBPath = "./abyss.db"
	}

	if req.ForgeType == "" || req.ForgeHost == "" || req.ClientID == "" || req.ClientSecret == "" || req.CIType == "" {
		response.BadRequest(w, "forge_type, forge_host, client_id, client_secret, and ci_type are required")
		return
	}
	if (req.DBType == "postgres" || req.DBType == "mysql") && (req.DBHost == "" || req.DBUser == "" || req.DBName == "") {
		response.BadRequest(w, "db_host, db_user, and db_name are required for postgres and mysql")
		return
	}
	if (req.CIType == "woodpecker" || req.CIType == "drone") && (req.CIHost == "" || req.CIToken == "") {
		response.BadRequest(w, "ci_host and ci_token are required for woodpecker and drone")
		return
	}

	path := h.configPath
	if path == "" {
		path = config.DefaultPath()
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		response.InternalError(w)
		return
	}
	if err := os.WriteFile(path, []byte(buildConfigYAML(req)), 0600); err != nil {
		response.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "restarting"})

	// Trigger restart after the response has been flushed.
	go func() {
		time.Sleep(200 * time.Millisecond)
		h.restartCh <- struct{}{}
	}()
}

func buildConfigYAML(r setupRequest) string {
	return fmt.Sprintf(`server:
  host: "0.0.0.0"
  port: "8000"

%s

forge:
  type: %s
  host: %s
  token: %s
  owner: %s
  branch: %s

ci:
  type: %s
  host: %s
  token: %s

auth:
  client_id: %s
  client_secret: %s
  callback_url: %s

logging:
  level: "info"
  format: "pretty"
  file: ""
`,
		buildDatabaseYAML(r),
		ys(r.ForgeType), ys(r.ForgeHost), ys(r.ForgeToken), ys(r.ForgeOwner), ys(forgeBranch(r)),
		ys(r.CIType), ys(r.CIHost), ys(r.CIToken),
		ys(r.ClientID), ys(r.ClientSecret), ys(r.CallbackURL),
	)
}

func buildDatabaseYAML(r setupRequest) string {
	switch r.DBType {
	case "postgres", "mysql":
		port := r.DBPort
		if port == "" {
			if r.DBType == "postgres" {
				port = "5432"
			} else {
				port = "3306"
			}
		}
		return fmt.Sprintf(`database:
  type: %s
  host: %s
  port: %s
  user: %s
  password: %s
  name: %s`,
			ys(r.DBType), ys(r.DBHost), ys(port), ys(r.DBUser), ys(r.DBPassword), ys(r.DBName))
	default:
		return fmt.Sprintf(`database:
  type: "sqlite"
  path: %s`, ys(r.DBPath))
	}
}

// ys wraps a value in a double-quoted YAML string, escaping internal quotes.
func ys(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

func forgeBranch(r setupRequest) string {
	if r.ForgeBranch == "" {
		return "main"
	}
	return r.ForgeBranch
}
