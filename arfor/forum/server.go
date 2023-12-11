package foruum

import (
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
	}
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	infoLog.Print(" http://localhost:", *&port)

	return s.httpServer.ListenAndServe()
}

// func (s *Server) Shutdown(ctx context.Context) error {
// 	return s.httpServer.Shutdown(ctx)

// }
