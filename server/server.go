package server

import (
	"net/http"

	"go-rest-boilerplate/config"
	"go-rest-boilerplate/data"
	"go-rest-boilerplate/domain/domainerror"
	inforepo "go-rest-boilerplate/domain/repository/info"
	infousecase "go-rest-boilerplate/domain/usecase/info"
	"go-rest-boilerplate/presentation/controllers"
	infocontrollers "go-rest-boilerplate/presentation/controllers/info/v1"
	"go-rest-boilerplate/presentation/render"

	"github.com/felixge/httpsnoop"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type Server struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{cfg}
}

func (s *Server) Start() error {
	router := httprouter.New()

	infoRepository := inforepo.NewInfoRepository()
	infoUsecase := infousecase.NewInfoUsecase(infoRepository)

	allControllers := []controllers.ControllerInterface{
		infocontrollers.NewInfoController(infoUsecase),
	}

	if err := controllers.MountAll(router, allControllers); err != nil {
		return err
	}

	router.NotFound = http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		render.Error(resp, r, new(domainerror.NotFoundError))
	})

	db, err := data.Connect(s.cfg)
	if err != nil {
		return err
	}

	stack := s.insertCfg(
		s.logRequests(
			s.insertDB(db, router),
		),
	)

	handler := cors.AllowAll().Handler(stack)
	return http.ListenAndServe(s.cfg.ListenAddr(), handler)
}

func (s *Server) logRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		var metrics *httpsnoop.Metrics
		s.cfg.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.RequestURI)
		defer func() {
			if metrics == nil {
				return
			}

			s.cfg.Printf("%s %s %s - HTTP %d %s\n",
				r.RemoteAddr, r.Method, r.RequestURI,
				metrics.Code, metrics.Duration)
		}()

		m := httpsnoop.CaptureMetrics(handler, resp, r)
		metrics = &m
	})
}

func (s *Server) insertCfg(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(resp, config.NewRequest(r, s.cfg))
	})
}

func (s *Server) insertDB(db *sqlx.DB, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(resp, r.WithContext(data.NewContext(r.Context(), db)))
	})
}
