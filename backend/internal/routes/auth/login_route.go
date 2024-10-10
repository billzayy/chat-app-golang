package auth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/billzayy/chat-golang/internal/db/authdb"
	"github.com/billzayy/chat-golang/internal/handlers"
	"github.com/billzayy/chat-golang/internal/types"
)

func LoginRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		handlers.Response(w, http.StatusMethodNotAllowed, "Invalid Method")
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		handlers.Response(w, http.StatusBadRequest, "Error Request Body")
		return
	}

	var loginData types.RequestUser

	err = json.Unmarshal(body, &loginData)

	if err != nil {
		handlers.Response(w, http.StatusBadRequest, "Error on Unmarshal Data")
		return
	}

	result, err := authdb.Login(loginData, w)

	if err != nil {
		switch err.Error() {
		case "user not found":
			handlers.Response(w, http.StatusNotFound, err)
			return
		default:
			handlers.Response(w, http.StatusInternalServerError, err)
			return
		}
	}

	handlers.Response(w, http.StatusOK, result)
}
