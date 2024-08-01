package pkg

import "net/http"

func ContextMiddleware(ctx RequestContext, handler func(ctx RequestContext, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(ctx, w, r)
	}
}
