package routes

import (
	"net/http"
	"server/internal/logic"
	"server/internal/msg"
)

func Registrar(mux *http.ServeMux) {
	mux.HandleFunc("GET /registrar", func(w http.ResponseWriter, r *http.Request) {
		t := r.URL.Query().Get("t")
		z := r.URL.Query().Get("z")

		if t != logic.TempRegistrarToken || logic.CheckAuthFile() {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if z == "" {
			w.WriteHeader(http.StatusBadRequest)
			logic.Output("error", msg.RegistrarEnterPortal)
			return
		}

		ok, err := logic.Registrar(t)
		if err != nil || !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		logic.Output("ok", msg.RegistrarOk)
	})
}
