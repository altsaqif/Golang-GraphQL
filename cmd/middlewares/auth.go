package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/altsaqif/graphql-go/cmd/configs"
	"github.com/altsaqif/graphql-go/cmd/helpers"
	"github.com/altsaqif/graphql-go/cmd/models"
)

// Key type for context key
type Key string

// KeyUser represents context key for user
const KeyUser Key = "user"

// Middleware function for authentication
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve token from header
		header := r.Header.Get("Authorization")
		// token, _ := r.Cookie("auth_token")
		// if err != nil {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// Allow unauthenticated users in
		if header == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Validate token (example: check if it's valid JWT)
		email, err := helpers.ParseToken(header)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		// Simulate fetching user from database based on email
		user := models.User{Email: email}
		userID, err := models.GetUserIdByEmail(configs.DB, email)
		if err != nil {
			// next.ServeHTTP(w, r)
			http.Error(w, "Email not found", http.StatusNotFound)
			return
		}

		// Create user object
		user.ID = userID

		// Add user to context
		ctx := context.WithValue(r.Context(), KeyUser, &user)

		// Call the next handler with our new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *models.User {
	raw, ok := ctx.Value(KeyUser).(*models.User)
	if !ok {
		log.Println("Error : ", ok)
	}
	fmt.Println(raw)
	return raw
}
