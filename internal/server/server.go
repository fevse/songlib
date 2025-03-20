package server

import (
	"context"
	"net"
	"net/http"

	_ "github.com/fevse/songlib/docs"
	"github.com/fevse/songlib/internal/app"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	server *http.Server
	app    *app.SongLibApp
}

func NewServer(app *app.SongLibApp, host, port string) *Server {
	return &Server{
		server: &http.Server{
			Addr: net.JoinHostPort(host, port),
		},
		app: app}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.Handle("POST /songs", s.CreateSong())
	mux.Handle("GET /songs", s.GetSongs())
	mux.Handle("PUT /songs/{id}", s.UpdateSong())
	mux.Handle("DELETE /songs/{id}", s.DeleteSong())
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	s.server.Handler = mux
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil

}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
