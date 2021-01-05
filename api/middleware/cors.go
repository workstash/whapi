package middleware

import (
	"github.com/workstash/whapi/config"

	"net/http"
)

//Cors add headers to suport the browser CORS
func Cors(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", config.Main.API.CorsOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		return
	}
	next(w, r)
}
