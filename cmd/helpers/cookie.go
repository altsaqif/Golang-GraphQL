package helpers

import (
	"context"
	"log"
	"net/http"
	"time"
)

var (
	// Secure cookie hash keys
	hashKey  = []byte("long-hash-key-for-securecookie")
	blockKey = []byte("long-block-key-for-securecookie")
)

// Function to set cookie
func SetCookie(ctx context.Context, token string) {

	// Mendapatkan ResponseWriter dari konteks
	w := ctx.Value("http.ResponseWriter")
	if w == nil {
		log.Println("Failed to get ResponseWriter from context : ResponseWriter is nil")
		return
	}

	// Mengonversi ResponseWriter ke http.ResponseWriter
	writer, ok := w.(http.ResponseWriter)
	if !ok {
		log.Println("Failed to convert ResponseWriter to http.ResponseWriter")
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour), // Expires in 24 hours
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}

	http.SetCookie(writer, cookie)
}

func ClearCookie(ctx context.Context) {
	if r, ok := ctx.Value("http").(*http.Request); ok {
		w := r.Context().Value("http").(http.ResponseWriter) // Mengakses objek respons
		clearCookie := http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		}
		http.SetCookie(w, &clearCookie) // Mengatur cookie kosong untuk membersihkan token
	}
}
