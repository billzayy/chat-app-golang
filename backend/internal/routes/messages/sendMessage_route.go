package messages

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/billzayy/chat-golang/internal/db/messagedb"
	"github.com/billzayy/chat-golang/internal/handlers"
	"github.com/billzayy/chat-golang/internal/pkg"
	"github.com/billzayy/chat-golang/internal/types"
)

func SendMessageRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("jwt")

	if err != nil {
		handlers.Response(w, http.StatusNotFound, http.ErrNoCookie)
		return
	}

	if r.Method != http.MethodPost {
		handlers.Response(w, http.StatusMethodNotAllowed, "Invalid Method")
		return
	}

	cookieId, err := pkg.ReadToken(cookie.Value)

	if err != nil {
		handlers.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	pathVariableId := r.PathValue("receiverId")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		handlers.Response(w, http.StatusBadRequest, "Error Request Body")
		return
	}

	var message types.InputMessage

	err = json.Unmarshal(body, &message)

	if err != nil {
		handlers.Response(w, http.StatusBadRequest, "Error on Unmarshal Data")
		return
	}

	data, err := messagedb.SendMessage(cookieId, pathVariableId, message)

	if err != nil {
		switch err.Error() {
		case "conversation not found":
			handlers.Response(w, http.StatusNotFound, err.Error())
			return
		default:
			handlers.Response(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	handlers.Response(w, http.StatusOK, data)
}
