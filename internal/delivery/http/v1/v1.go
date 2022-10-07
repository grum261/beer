package v1

import (
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func RegisterRoutes(mux *runtime.ServeMux, u *UserHandler) error {
	return mux.HandlePath("POST", "/api/v1/users", u.HandleCreateUser)
}

func write(w http.ResponseWriter, statusCode int, payload interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"errMsg":  err.Error(),
			"payload": payload,
		})
		w.WriteHeader(statusCode)
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
	w.WriteHeader(statusCode)
}

func writeInternal(w http.ResponseWriter, err error) {
	write(w, http.StatusInternalServerError, nil, err)
}

func writeUnprocessable(w http.ResponseWriter, err error) {
	write(w, http.StatusUnprocessableEntity, nil, err)
}

func writeOK(w http.ResponseWriter, payload interface{}) {
	write(w, http.StatusOK, payload, nil)
}

func writeUnauthorized(w http.ResponseWriter, err error) {
	write(w, http.StatusUnauthorized, nil, err)
}
