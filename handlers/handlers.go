// package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"shawty-master/storages"
)

type ShortenRequest struct {
	URL    string `json:"url"`
	Custom string `json:"custom,omitempty"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url,omitempty"`
	Error    string `json:"error,omitempty"`
}

func EncodeHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		
		var req ShortenRequest
		
		if r.Header.Get("Content-Type") == "application/json" {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ShortenResponse{Error: "Invalid JSON"})
				return
			}
		} else {
			
			req.URL = r.PostFormValue("url")
			req.Custom = r.PostFormValue("custom")
		}

		if req.URL == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ShortenResponse{Error: "URL is required"})
			return
		}

		
		if req.Custom != "" {
			if len(req.Custom) < 3 || len(req.Custom) > 20 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ShortenResponse{Error: "Custom name must be 3-20 characters"})
				return
			}
			
			
			validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
			if !validName.MatchString(req.Custom) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ShortenResponse{Error: "Custom name can only contain letters, numbers, hyphens, and underscores"})
				return
			}
		}

		var code string
		var err error

		if req.Custom != "" {
			code, err = storage.SaveWithCustom(req.URL, req.Custom)
			if err != nil {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(ShortenResponse{Error: err.Error()})
				return
			}
		} else {
			code = storage.Save(req.URL)
		}

		var code string
		var err error

		if req.Custom != "" {
			code, err = storage.SaveWithCustom(req.URL, req.Custom)
			if err != nil {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(ShortenResponse{Error: err.Error()})
				return
			}
		} else {
			code = storage.Save(req.URL)
		}

		var scheme string
		if r.Header.Get("X-Forwarded-Proto") == "http" || strings.HasPrefix(r.Host, "localhost") {
			scheme = "http"
		} else {
			scheme = "https"
		}
		
		shortURL := scheme + "://" + r.Host + "/" + code
		
		json.NewEncoder(w).Encode(ShortenResponse{ShortURL: shortURL})
	}

	return http.HandlerFunc(handleFunc)
}

func DecodeHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/dec/"):]

		url, err := storage.Load(code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
			return
		}

		w.Write([]byte(url))
	}

	return http.HandlerFunc(handleFunc)
}

func RedirectHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/red/"):]

		url, err := storage.Load(code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
			return
		}

		http.Redirect(w, r, string(url), 301)
	}

	return http.HandlerFunc(handleFunc)
}
