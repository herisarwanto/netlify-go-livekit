package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		http.Error(w, "LiveKit API key/secret not set", http.StatusInternalServerError)
		return
	}

	claims := jwt.MapClaims{
		"iss": apiKey,
		"sub": req.Identity,
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
		"grants": map[string]interface{}{
			"roomJoin":    true,
			"room":        req.Room,
			"canPublish":  true,
			"canSubscribe": true,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(apiSecret))
	resp := Response{}
	if err != nil {
		resp.Success = false
		resp.Error = "Failed to generate token"
	} else {
		resp.Success = true
		resp.Token = signed
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
	fmt.Println("Listening on port", port)
	http.ListenAndServe(":"+port, nil)
} 