package app

type App struct {
	sessionStore *SessionStore
}

func NewApp() *App {

	sessionStore := NewSessionStore()
	return &App{
		sessionStore: sessionStore,
	}
}
