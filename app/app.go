package app

import (
	"github.com/go-chi/chi"
	chim "github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

type Server struct {
	Host string
}

func (s *Server) Setup() *chi.Mux {
	corsHandler := cors.AllowAll()

	r := chi.NewRouter()
	r.Use(
		corsHandler.Handler,
		chim.Logger,
		chim.Recoverer,
	)

	r.Group(func(r chi.Router) {
		r.Use(WithLogger)

		r.NotFound(s.Handler)
	})

	return r
}
