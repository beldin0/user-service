package routing

import (
	"net/http"

	"github.com/beldin0/users/src/logging"
	"github.com/beldin0/users/src/routeutils"
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
	initialPath := r.URL.Path
	var head string
	head, r.URL.Path = routeutils.ShiftPath(r.URL.Path)
	handler, ok := h.handlerMap[head]
	if ok {
		handler.ServeHTTP(w, r)
		return
	}
	logging.NewLogger().Sugar().With("path", initialPath).Info("unable to route incoming request")
	http.Error(w, "Not Found", http.StatusNotFound)
}
