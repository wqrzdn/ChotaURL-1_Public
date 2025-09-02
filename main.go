package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"shawty-master/handlers"
	"shawty-master/storages"
	"github.com/joho/godotenv"
)

func main() {
	
	
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found or could not be loaded. Relying on system environment variables.")
	}

	
	redisDSN := os.Getenv("REDIS_DSN")
	if redisDSN == "" {
		
		redisDSN = "redis://localhost:6379/0"
	}

	redisStorage, err := storages.NewRedis(redisDSN)
	if err != nil {
		log.Fatal("Failed to initialize Redis storage:", err)
	}
	log.Println("Successfully connected to Redis.")

	
	http.Handle("/shorten", handlers.EncodeHandler(redisStorage))
	http.Handle("/api/shorten", handlers.EncodeHandler(redisStorage))

	
	http.Handle("/dec/", handlers.DecodeHandler(redisStorage))
	http.Handle("/api/dec/", handlers.DecodeHandler(redisStorage))
	http.Handle("/red/", handlers.RedirectHandler(redisStorage))
	http.Handle("/api/red/", handlers.RedirectHandler(redisStorage))
	fs := http.FileServer(http.Dir("public"))
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		
		if r.URL.Path == "/" {
			fs.ServeHTTP(w, r)
			return
		}
		
		
		code := strings.TrimPrefix(r.URL.Path, "/")
		
		
		if strings.Contains(code, ".") || 
		   strings.HasPrefix(code, "api/") ||
		   strings.HasPrefix(code, "shorten") ||
		   strings.HasPrefix(code, "dec/") ||
		   strings.HasPrefix(code, "red/") {
			fs.ServeHTTP(w, r)
			return
		}
		
		
		url, err := redisStorage.Load(code)
		if err != nil {
			
			fs.ServeHTTP(w, r)
			return
		}
		
		
		http.Redirect(w, r, url, 301)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
