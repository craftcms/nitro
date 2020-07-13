package middleware

import (
	"encoding/json"
	"github.com/craftcms/nitro/internal/api"
	"net/http"
)

// RequirePOST is used to abstract requiring post requests for
// handle requests
func RequirePOST(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// requires a POST request
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)

			resp, err := json.Marshal(api.Response{Success: false, Errors: []string{"method not allowed"}})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, _ = w.Write(resp)

			return
		}

		next.ServeHTTP(w, r)
	})
}
