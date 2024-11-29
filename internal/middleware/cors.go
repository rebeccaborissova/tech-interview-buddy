package middleware

import "net/http"

func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        // Enable CORS via HTTP response header
        writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie") // Added Cookie
        writer.Header().Set("Access-Control-Allow-Credentials", "true")
        writer.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")

        // Handle preflight requests
        if request.Method == "OPTIONS" {
            writer.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(writer, request)
    })
}