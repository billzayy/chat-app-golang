package routes

import (
	"net/http"

	"github.com/billzayy/chat-golang/internal/pkg"
	"github.com/billzayy/chat-golang/internal/routes/auth"
	"github.com/billzayy/chat-golang/internal/routes/messages"
	"github.com/billzayy/chat-golang/internal/routes/user"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Group APIs
	mux.Handle("/api/auth/", http.StripPrefix("/api/auth", AuthMux()))
	mux.Handle("/api/users/", http.StripPrefix("/api/users", UserMux()))
	mux.Handle("/api/messages/", http.StripPrefix("/api/messages", MessageMux()))

	mux.HandleFunc("/ws", pkg.ServeWs)

	return mux
}

func AuthMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/signup", auth.SignUpRoute)
	mux.HandleFunc("/login", auth.LoginRoute)
	mux.HandleFunc("/logout", auth.LogOutRoute)

	return mux
}

func UserMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", user.UserRoute)

	return mux
}

func MessageMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/{receiverId}", messages.GetMessageRoute)
	mux.HandleFunc("/send/{receiverId}", messages.SendMessageRoute)

	return mux
}
