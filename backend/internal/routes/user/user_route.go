package user

import (
	"net/http"

	"github.com/billzayy/chat-golang/internal/db/userdb"
	"github.com/billzayy/chat-golang/internal/handlers"
	"github.com/billzayy/chat-golang/internal/pkg"
	"github.com/billzayy/chat-golang/internal/pkg/middleware"
)

func UserRoute(w http.ResponseWriter, r *http.Request) {
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

	userId, err := pkg.ReadToken(cookie.Value)

	if err != nil {
		handlers.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	checkProtect, err := middleware.ProtectRoute(userId)

	if err != nil {
		panic(err)
	}

	if checkProtect {
		result, err := userdb.GetAllUsers(userId)

		if err != nil {
			switch err.Error() {
			case "user not found":
				handlers.Response(w, http.StatusNotFound, err.Error())
				return
			default:
				handlers.Response(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		handlers.Response(w, http.StatusOK, result)
		return
	} else {
		handlers.Response(w, http.StatusNotAcceptable, "Protected !!!")
		return
	}
}
