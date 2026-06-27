package http

import stdhttp "net/http"

func NewRouter(handler *Handler) stdhttp.Handler {
	mux := stdhttp.NewServeMux()
	mux.HandleFunc("/weather", handler.Weather)
	return mux
}
