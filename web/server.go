package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Flack74/pom/config"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
}

type TimerSession struct {
	ID          string `json:"id"`
	WorkTime    int    `json:"work_time"`
	BreakTime   int    `json:"break_time"`
	Sessions    int    `json:"sessions"`
	CurrentSess int    `json:"current_session"`
	IsRunning   bool   `json:"is_running"`
	IsPaused    bool   `json:"is_paused"`
	IsBreak     bool   `json:"is_break"`
	TimeLeft    int    `json:"time_left"`
	Profile     string `json:"profile"`
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		clients:  make(map[*websocket.Conn]bool),
	}
}

func (s *Server) Start(port int) error {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/profiles", s.handleProfiles).Methods("GET")
	api.HandleFunc("/session/start", s.handleStartSession).Methods("POST")
	api.HandleFunc("/insights/suggestions", s.handleSuggestions).Methods("GET")
	api.HandleFunc("/insights/today", s.handleTodayStats).Methods("GET")
	api.HandleFunc("/plugins", s.handlePlugins).Methods("GET")
	api.HandleFunc("/privacy/status", s.handlePrivacyStatus).Methods("GET")

	// Serve embedded HTML/JS web UI
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(getWebUI()))
	})

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("🌐 Web UI: http://localhost%s\n", addr)
	return http.ListenAndServe(addr, r)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func (s *Server) handleProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, _ := config.LoadProfiles()
	cfg, _ := config.LoadConfig()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"profiles": profiles.Profiles,
		"current":  cfg.CurrentProfile,
	})
}

func (s *Server) handleStartSession(w http.ResponseWriter, r *http.Request) {
	var req TimerSession
	json.NewDecoder(r.Body).Decode(&req)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

func (s *Server) handleSuggestions(w http.ResponseWriter, r *http.Request) {
	suggestions, _ := config.GenerateSuggestions()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

func (s *Server) handleTodayStats(w http.ResponseWriter, r *http.Request) {
	sessions, minutes, _ := config.GetTodayStats()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sessions": sessions,
		"minutes":  minutes,
		"hours":    float64(minutes) / 60,
	})
}

func (s *Server) handlePlugins(w http.ResponseWriter, r *http.Request) {
	plugins, _ := config.LoadPlugins()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plugins.Plugins)
}

func (s *Server) handlePrivacyStatus(w http.ResponseWriter, r *http.Request) {
	cfg, _ := config.LoadConfig()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"privacy_mode": cfg.PrivacyMode,
		"cloud_sync":   cfg.CloudSync,
	})
}