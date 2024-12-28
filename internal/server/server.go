package server

import (
	"encoding/json"
	"net/http"

	"byteload-agent/pkg/system"
)

// Server represents the HTTP server
type Server struct {
	port string
}

// New creates a new server instance
func New(port string) *Server {
	if port == "" {
		port = "3000"
	}
	return &Server{port: port}
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
	http.HandleFunc("/system", func(w http.ResponseWriter, r *http.Request) {
		var services []string
		if servicesParam := r.URL.Query().Get("services"); servicesParam != "" {
			if err := json.Unmarshal([]byte(servicesParam), &services); err != nil {
				jsonResponse(w, nil, err)
				return
			}
		}
		data, err := system.GetSystemData(services)
		jsonResponse(w, data, err)
	})

	http.HandleFunc("/os", func(w http.ResponseWriter, r *http.Request) {
		data, err := system.GetOSData()
		jsonResponse(w, data, err)
	})

	http.HandleFunc("/cpu", func(w http.ResponseWriter, r *http.Request) {
		data, err := system.GetCPUData()
		jsonResponse(w, data, err)
	})

	http.HandleFunc("/memory", func(w http.ResponseWriter, r *http.Request) {
		data, err := system.GetMemoryData()
		jsonResponse(w, data, err)
	})

	http.HandleFunc("/storage", func(w http.ResponseWriter, r *http.Request) {
		data, err := system.GetStorageData()
		jsonResponse(w, data, err)
	})

	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		var services []string
		if servicesParam := r.URL.Query().Get("services"); servicesParam != "" {
			if err := json.Unmarshal([]byte(servicesParam), &services); err != nil {
				jsonResponse(w, nil, err)
				return
			}
		}
		data, err := system.GetServicesData(services)
		jsonResponse(w, data, err)
	})

	http.HandleFunc("/docker", func(w http.ResponseWriter, r *http.Request) {
		data, err := system.GetDockerContainersData()
		jsonResponse(w, data, err)
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.setupRoutes()
	return http.ListenAndServe(":"+s.port, nil)
}
