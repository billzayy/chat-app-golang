package messages

import (
	"net/http"

	"github.com/billzayy/chat-golang/internal/db/messagedb"
	"github.com/billzayy/chat-golang/internal/handlers"
	"github.com/billzayy/chat-golang/internal/pkg"
)

func GetMessageRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("jwt")

	if err != nil {
		handlers.Response(w, http.StatusNotFound, http.ErrNoCookie)
		return
	}

	if r.Method != http.MethodGet {
		handlers.Response(w, http.StatusMethodNotAllowed, "Invalid Method")
		return
	}

	cookieId, err := pkg.ReadToken(cookie.Value)

	if err != nil {
		handlers.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	pathVariableId := r.PathValue("receiverId")

	data, err := messagedb.GetMessage(cookieId, pathVariableId)

	if err != nil {
		switch err.Error() {
		case "conversation not found":
			handlers.Response(w, http.StatusNotFound, data)
			return
		default:
			handlers.Response(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	handlers.Response(w, http.StatusOK, data)
}
