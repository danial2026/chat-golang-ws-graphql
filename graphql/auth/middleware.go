package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"chat-golang-react/chat/common/decoders"
)

var UserCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			tokenStr := header

			tokenParts := strings.Split(tokenStr, " ")
			if len(tokenParts) != 2 {
				next.ServeHTTP(w, r)
				return
			}

			tokenClaims, err := decoders.ParseToken(r.Context(), tokenParts[1])
			if err != nil {
				log.Println("Error during checking token:", err)
				return
			}
			tokenClaims = decoders.GenerateContextAuth(tokenClaims, tokenParts[1])

			// put it in context
			ctx := context.WithValue(r.Context(), UserCtxKey, &tokenClaims)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
