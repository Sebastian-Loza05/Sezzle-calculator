package httpapi

import "net/http"

func NewRouter() http.Handler {
	handler := NewHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /calculate", handler.HandleCalculate)
	mux.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", "POST")
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "use POST for calculations")
	})

	return mux
}
