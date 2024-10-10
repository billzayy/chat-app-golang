package auth

import (
	"net/http"

	"github.com/billzayy/chat-golang/internal/handlers"
)

func LogOutRoute(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "jwt",
		Value:  "",
		MaxAge: 0,
	}

	http.SetCookie(w, &cookie)

	handlers.Response(w, http.StatusOK, "Logged out successfully")
}
