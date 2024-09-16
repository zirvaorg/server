package routes

import (
	"github.com/zirvaorg/ratelimit"
	"github.com/zirvaorg/ratelimit/memstore"
	"net/http"
	"server/internal/logic"
	"server/internal/msg"
)

func keyFunc(r *http.Request) string {
	return r.RemoteAddr
}

func Registrar(mux *http.ServeMux, store *memstore.MemStore) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	mux.Handle("/registrar", ratelimit.Middleware(store, handler, keyFunc))
}
