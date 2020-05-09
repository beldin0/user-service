package routing

import (
	"net/http"
	"path"
	"strings"

	"github.com/beldin0/users/src/userhandler"
	"github.com/jmoiron/sqlx"
)

type routing struct {
	handlerMap map[string]http.Handler
}

// NewRouting returns a routing instance
func NewRouting(db *sqlx.DB) http.Handler {
	handlerMap := make(map[string]http.Handler)
	handlerMap["users"] = userhandler.New(db)
	return &routing{
		handlerMap: handlerMap,
	}
}

func (h *routing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	handler, ok := h.handlerMap[head]
	if ok {
		handler.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
