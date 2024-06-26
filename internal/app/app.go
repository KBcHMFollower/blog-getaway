package app

import (
	"context"
	"github.com/unrolled/render"
	"google.golang.org/grpc"
	"log/slog"
	"net/http"
	"test-plate/config"
	postdepend "test-plate/internal/dependencies/postservice"
	dependencies "test-plate/internal/dependencies/userservice"
	"test-plate/internal/handlers"
	"test-plate/internal/services"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	router     *chi.Mux
	log        *slog.Logger
	cfg        *config.Config
	srv        *http.Server
	grpcClient *grpc.ClientConn
}

func New(log *slog.Logger, cfg *config.Config) *App {
	router := chi.NewRouter()

	return &App{
		router: router,
		log:    log,
		cfg:    cfg,
	}
}

func (a *App) AddMiddlewares() *App {
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
	a.router.Use(middleware.URLFormat)

	a.log.Info("middlewares is added")

	return a
}

func (a *App) AddHandlers() *App {

	rnd := render.New()

	a.log.Info("1")

	userGrpcClient, err := dependencies.NewUsersClient(a.cfg.GrpcServices.UserService.Addr)
	postGrpcClient, err := postdepend.NewPostsClient(a.cfg.GrpcServices.PostService.Addr)

	if err != nil {
		panic(err)
	}

	a.log.Info("2")

	authService := services.NewAuthService(userGrpcClient, a.log)
	userService := services.NewUserService(userGrpcClient, a.log)
	postService := services.NewPostService(postGrpcClient, a.log)
	commService := services.NewCommentService(postGrpcClient, a.log)

	authHandler := handlers.NewAuthHandler(authService, rnd)
	userHandler := handlers.NewUsersHandler(userService, rnd, a.cfg.HttpServer.PublicUrl)
	postHandler := handlers.NewPostHandler(postService, rnd)
	commHandler := handlers.NewCommentHandler(commService, rnd)

	a.router.Route("/api", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		r.Get("/auth", authHandler.CheckAuth)
		r.Get("/images/{imgName}", userHandler.GetAvatar)
		r.Route("/users", func(r chi.Router) {

			r.Delete("/{userId}", userHandler.DeleteUser)
			r.Patch("/{userId}", userHandler.UpdateUser)
			r.Get("/{userId}", userHandler.GetUser)

			r.Get("/{userId}/subscribers", userHandler.GetSubscribers)
			r.Post("/{userId}/subscribers", userHandler.Subscribe)
			r.Delete("/{userId}/subscribers", userHandler.Unsubscribe)
			r.Get("/{userId}/subscriptions", userHandler.GetSubscriptions)

			r.Patch("/{userId}/posts/{postId}", postHandler.UpdatePost)
			r.Route("/{userId}/posts", func(r chi.Router) {
				r.Get("/", postHandler.GetUserPosts)
				r.Post("/", postHandler.CreatePost)
				r.Route("/{postId}", func(r chi.Router) {
					r.Get("/", postHandler.GetPost)
					r.Patch("/", postHandler.UpdatePost)
					r.Delete("/", postHandler.DeletePost)

					r.Post("/comments", commHandler.CreateComment)
					r.Get("/comments", commHandler.GetPostComments)
					r.Get("/comments/{commentId}", commHandler.GetComment)
					r.Delete("/comments/{commentId}", commHandler.DeleteComment)
					r.Patch("/comments/{commentId}", commHandler.UpdateComment)
				})
			})
		})
	})

	a.log.Info("handlers is added")
	return a
}

func (a *App) Run() {
	a.log.Info("server is trying to get up")

	srv := &http.Server{
		Addr:         a.cfg.HttpServer.Address,
		Handler:      a.router,
		ReadTimeout:  a.cfg.HttpServer.Timeout,
		WriteTimeout: a.cfg.HttpServer.Timeout,
		IdleTimeout:  a.cfg.HttpServer.IddleTimeout,
	}

	a.srv = srv

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			a.log.Error("failed to start server")
		}
	}()

	a.log.Info("server started with", slog.String("addr", srv.Addr))
}

func (a *App) Stop() {
	a.log.Info("server is trying to stop")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		a.log.Error("failed to stop server", err)

		return
	}
}
