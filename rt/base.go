package rt

import (
	"net/http"
)

var PingRouter = Router{
	Name:   "ping",
	Path:   "ping",
	Method: METHOD.GET,
	Handler: func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == "GET" || r.Method == "HEAD" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong"))
			return nil
		}
		return nil
	},
}
