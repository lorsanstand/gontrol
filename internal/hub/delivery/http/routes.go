package http

import (
	"net/http"

	"github.com/lorsanstand/gontrol/internal/hub/delivery/http/handler"
	"github.com/lorsanstand/gontrol/internal/hub/store"
)

func RegisterRoutes(mux *http.ServeMux, store *store.AgentStore) {
	heartbeatHandler := handler.NewHeartbeatHandler(store)
	mux.Handle("POST /api/v1/heartbeat", heartbeatHandler)
}
