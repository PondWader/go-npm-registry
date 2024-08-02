package pkg

import (
	"net/http"
	"strings"

	"github.com/PondWader/go-npm-registry/pkg/response"
)

func ContextMiddleware(ctx RequestContext, handler RequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(ctx, w, r)
	}
}

func AuthMiddleware(ctx RequestContext, handler RequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			response.Error(w, http.StatusUnauthorized, "Invalid authorization.")
			return
		}
		key := auth[7:]

		keyFound := false
		for _, userKey := range ctx.Config.UserKeys {
			if key == userKey {
				keyFound = true
				break
			}
		}
		if !keyFound {
			response.Error(w, http.StatusUnauthorized, "Invalid authorization.")
			return
		}

		ctx.UserKey = key

		handler(ctx, w, r)
	}
}
