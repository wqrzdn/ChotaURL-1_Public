package api

import (
	"log"
	"net/http"
	"os"
	"shawty-master/handlers"
	"shawty-master/storages"
	"github.com/joho/godotenv"
	"strings"
)

var redisStorage *storages.Redis

func init() {
	_ = godotenv.Load()
	var err error
	redisStorage, err = storages.NewRedis(os.Getenv("REDIS_DSN"))
	if err != nil {
		
		log.Printf("Failed to connect to Redis: %v", err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
    
    if redisStorage == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"error":"Database connection failed"}`))
        return
    }

    path := r.URL.Path
    originalPath := path
    
   
    if strings.HasPrefix(path, "/api") {
        path = strings.TrimPrefix(path, "/api")
    }
    
    
    log.Printf("Original path: %s, Processed path: %s", originalPath, path)
    
    switch {
    case path == "/shorten" && r.Method == http.MethodPost:
        handlers.EncodeHandler(redisStorage).ServeHTTP(w, r)
    case strings.HasPrefix(path, "/dec/") && r.Method == http.MethodGet:
       
        newReq := r.Clone(r.Context())
        newReq.URL.Path = path
        handlers.DecodeHandler(redisStorage).ServeHTTP(w, newReq)
    case strings.HasPrefix(path, "/red/") && r.Method == http.MethodGet:
        
        newReq := r.Clone(r.Context())
        newReq.URL.Path = path
        handlers.RedirectHandler(redisStorage).ServeHTTP(w, newReq)
    case path != "" && path != "/" && r.Method == http.MethodGet:
        
        code := strings.TrimPrefix(path, "/")
        
        log.Printf("Trying to load short code: %s", code)
        
       
        if strings.Contains(code, ".") || 
           strings.HasPrefix(code, "api/") ||
           strings.HasPrefix(code, "shorten") ||
           strings.HasPrefix(code, "dec/") ||
           strings.HasPrefix(code, "red/") {
            log.Printf("Skipping code %s - looks like static file or API path", code)
            http.NotFound(w, r)
            return
        }
        
        
        url, err := redisStorage.Load(code)
        if err != nil {
            log.Printf("Failed to load URL for code %s: %v", code, err)
            http.NotFound(w, r)
            return
        }
        
        log.Printf("Redirecting %s to %s", code, url)
        
        http.Redirect(w, r, url, 301)
    default:
        http.NotFound(w, r)
    }
}
