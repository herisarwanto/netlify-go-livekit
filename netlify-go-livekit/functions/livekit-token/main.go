package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/livekit/protocol/auth"
)

type Request struct {
	Identity string `json:"identity"`
	Room     string `json:"room"`
}

type Response struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	at := auth.NewAccessToken(os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET"))
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     req.Room,
	}
	at.AddGrant(grant).
		SetIdentity(req.Identity).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	resp := Response{}
	if err != nil {
		resp.Success = false
		resp.Error = "Failed to generate token"
	} else {
		resp.Success = true
		resp.Token = token
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
} 