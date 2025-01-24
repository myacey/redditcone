package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/myacey/redditclone/internal/handlers"
	"github.com/myacey/redditclone/internal/service"
	"github.com/myacey/redditclone/internal/token"

	"go.uber.org/zap"
)

type Server struct {
	Logger  *zap.SugaredLogger
	Service service.ServiceInterface
	Router  *mux.Router
	Handler *handlers.Handler
}

func NewServer(logger *zap.SugaredLogger, service service.ServiceInterface, tokenMaker token.TokenMaker) *Server {
	server := Server{
		Logger:  logger,
		Service: service,
	}

	server.Handler = handlers.NewHandler(server.Service, server.Logger, tokenMaker) // Создаем Handler ПОСЛЕ инициализации Service и JWTMaker

	server.configureRouter()

	return &server
}

func (s *Server) Start() {
	s.Logger.Info("start listening :8080")
	if err := http.ListenAndServe(":8080", s.Router); err != nil {
		s.Logger.Fatal(err)
	}
}

func (s *Server) configureRouter() {
	s.Router = mux.NewRouter()

	s.Router.Use(func(h http.Handler) http.Handler { return s.Handler.LoggingMiddleware(h, s.Logger) })
	s.addStaticToRouter()

	// protected, need auth
	protected := s.Router.PathPrefix("/api").Subrouter()
	protected.Use(func(h http.Handler) http.Handler { return s.Handler.AuthMiddleware(h, s.Logger) })
	protected.HandleFunc("/posts", s.Handler.AddPost).Methods("POST")
	protected.HandleFunc("/post/{id}", s.Handler.AddComment).Methods("POST")
	protected.HandleFunc("/post/{id}", s.Handler.DeletePost).Methods("DELETE")
	protected.HandleFunc("/post/{postID}/{commentID}", s.Handler.DeleteComment).Methods("DELETE")
	protected.HandleFunc("/post/{id}/unvote", s.Handler.UnvotePost).Methods("GET")
	protected.HandleFunc("/post/{id}/upvote", s.Handler.VotePost).Methods("GET")
	protected.HandleFunc("/post/{id}/downvote", s.Handler.DownvotePost).Methods("GET")

	s.Router.HandleFunc("/api/register", s.Handler.RegisterUser).Methods("POST")
	s.Router.HandleFunc("/api/login", s.Handler.LoginUser).Methods("POST")

	s.Router.HandleFunc("/api/posts/", s.Handler.GetPosts).Methods("GET")

	s.Router.HandleFunc("/api/post/{id}", s.Handler.GetPost).Methods("GET")

	s.Router.HandleFunc("/api/posts/{category}", s.Handler.GetPostsByCategory).Methods("GET")

	s.Router.HandleFunc("/api/user/{username}", s.Handler.GetUserPosts).Methods("GET")

}

func (s *Server) addStaticToRouter() {
	staticFileDirectory := http.Dir("./static")
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))

	// Статика по префиксу /static/
	s.Router.PathPrefix("/static/").Handler(staticFileHandler)

	// Корневой маршрут для отдачи index.html
	s.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/html/index.html")
	})

	// Для всех остальных маршрутов возвращать index.html
	s.Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/html/index.html")
		s.Logger.Infow("Unknown request",
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
		)
	})
}
