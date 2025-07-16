package auth

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"os"
)

func InitGoogleProvider() error {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := os.Getenv("CALLBACK_URL")
	sessionSecret := os.Getenv("SESSION_SECRET")

	if clientID == "" || clientSecret == "" || callbackURL == "" || sessionSecret == "" {
		return errors.New("missing one or more required environment variables: GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, CALLBACK_URL, SESSION_SECRET")
	}

	goth.UseProviders(
		google.New(
			clientID,
			clientSecret,
			callbackURL,
			"email", "profile",
		),
	)

	store := sessions.NewCookieStore([]byte(sessionSecret))
	store.MaxAge(86400 * 30)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false

	gothic.Store = store

	return nil
}
