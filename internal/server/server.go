package server

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"

	"byteload-agent/pkg/system"
)

// Config holds the server configuration
type Config struct {
	Port string
	Auth AuthConfig
}

type AuthConfig struct {
	Enabled  bool
	Username string
	Password string
}

// Server represents the HTTP server
type Server struct {
	config Config
}

// New creates a new server instance
func New(config Config) *Server {
	if config.Port == "" {
		config.Port = "3000"
	}
	return &Server{config: config}
}

// basicAuth middleware for HTTP basic authentication
func (s *Server) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip auth if not enabled
		if !s.config.Auth.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		// Get credentials from request
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(user), []byte(s.config.Auth.Username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(s.config.Auth.Password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// jsonResponse writes a JSON response
func jsonResponse(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	http.HandleFunc("/system", s.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		var services []string
		if servicesParam := r.URL.Query().Get("services"); servicesParam != "" {
			if err := json.Unmarshal([]byte(servicesParam), &services); err != nil {
				jsonResponse(w, nil, err)
				return
			}
		}
		data, err := system.GetSystemData(services)
		jsonResponse(w, data, err)
	}))

	http.HandleFunc("/os", s.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		data, err := system.GetOSData()
		jsonResponse(w, data, err)
	}))

	http.HandleFunc("/cpu", s.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		data, err := system.GetCPUData()
		jsonResponse(w, data, err)
	}))

	http.HandleFunc("/memory", s.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		data, err := system.GetMemoryData()
		jsonResponse(w, data, err)
	}))

	http.HandleFunc("/storage", s.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		data, err := system.GetStorageData()
		jsonResponse(w, data, err)
	}))

	http.HandleFunc("/services", s.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		var services []string
		if servicesParam := r.URL.Query().Get("services"); servicesParam != "" {
			if err := json.Unmarshal([]byte(servicesParam), &services); err != nil {
				jsonResponse(w, nil, err)
				return
			}
		}
		data, err := system.GetServicesData(services)
		jsonResponse(w, data, err)
	}))

	http.HandleFunc("/docker", s.basicAuth(func(w http.ResponseWriter, r *http.Request) {
		data, err := system.GetDockerContainersData()
		jsonResponse(w, data, err)
	}))
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.setupRoutes()
	return http.ListenAndServe(":"+s.config.Port, nil)
}
