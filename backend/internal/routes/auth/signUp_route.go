package auth

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/billzayy/chat-golang/internal/db/authdb"
	"github.com/billzayy/chat-golang/internal/handlers"
	"github.com/billzayy/chat-golang/internal/types"
)

func SignUpRoute(w http.ResponseWriter, r *http.Request) {
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

	var user types.User

	err = json.Unmarshal(body, &user)

	if err != nil {
		handlers.Response(w, http.StatusBadRequest, "Error on Unmarshal Data")
		return
	}

	data, err := authdb.SignUp(user, w)

	if err != nil {
		handlers.Response(w, http.StatusInternalServerError, err.Error())
		return
	}

	if data == "User Existed" {
		handlers.Response(w, http.StatusFound, data)
		return
	}

	handlers.Response(w, http.StatusCreated, data)
}
