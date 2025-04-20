package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "default port",
			config: Config{
				Port: "",
			},
			expected: "9001",
		},
		{
			name: "custom port",
			config: Config{
				Port: "8080",
			},
			expected: "8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := New(tt.config)
			if server.config.Port != tt.expected {
				t.Errorf("expected port %s, got %s", tt.expected, server.config.Port)
			}
		})
	}
}

func TestBasicAuth(t *testing.T) {
	tests := []struct {
		name           string
		authEnabled    bool
		username       string
		password       string
		requestUser    string
		requestPass    string
		expectedStatus int
	}{
		{
			name:           "auth disabled",
			authEnabled:    false,
			username:       "test",
			password:       "test",
			requestUser:    "",
			requestPass:    "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "auth enabled, valid credentials",
			authEnabled:    true,
			username:       "test",
			password:       "test",
			requestUser:    "test",
			requestPass:    "test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "auth enabled, invalid credentials",
			authEnabled:    true,
			username:       "test",
			password:       "test",
			requestUser:    "wrong",
			requestPass:    "wrong",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Auth: AuthConfig{
					Enabled:  tt.authEnabled,
					Username: tt.username,
					Password: tt.password,
				},
			}
			server := New(config)

			handler := server.basicAuth(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/", nil)
			if tt.requestUser != "" && tt.requestPass != "" {
				req.SetBasicAuth(tt.requestUser, tt.requestPass)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetRequestedServices(t *testing.T) {
	tests := []struct {
		name            string
		queryParam      string
		defaultServices []string
		expected        []string
	}{
		{
			name:            "no query param",
			queryParam:      "",
			defaultServices: []string{"nginx", "mysql"},
			expected:        []string{"nginx", "mysql"},
		},
		{
			name:            "valid query param",
			queryParam:      `["redis","postgres"]`,
			defaultServices: []string{"nginx", "mysql"},
			expected:        []string{"redis", "postgres"},
		},
		{
			name:            "invalid query param",
			queryParam:      "invalid json",
			defaultServices: []string{"nginx", "mysql"},
			expected:        []string{"nginx", "mysql"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Services: ServicesConfig{
					Default: tt.defaultServices,
				},
			}
			server := New(config)

			req := httptest.NewRequest("GET", "/", nil)
			if tt.queryParam != "" {
				q := req.URL.Query()
				q.Add("services", tt.queryParam)
				req.URL.RawQuery = q.Encode()
			}

			services := server.getRequestedServices(req)

			if len(services) != len(tt.expected) {
				t.Errorf("expected %d services, got %d", len(tt.expected), len(services))
			}

			for i, service := range services {
				if service != tt.expected[i] {
					t.Errorf("expected service %s, got %s", tt.expected[i], service)
				}
			}
		})
	}
}

func TestJsonResponse(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		err            error
		expectedStatus int
	}{
		{
			name:           "success response",
			data:           map[string]string{"test": "data"},
			err:            nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error response",
			data:           nil,
			err:            json.Unmarshal([]byte("invalid"), &struct{}{}),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			jsonResponse(w, tt.data, tt.err)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response Response
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if response.Meta.Version != Version {
					t.Errorf("expected version %s, got %s", Version, response.Meta.Version)
				}
			}
		})
	}
}
