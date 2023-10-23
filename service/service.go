package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"trueconf-refactor/repo"
)

type IRepository interface {
	SearchUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type Service struct {
	repository IRepository
}

func NewService() *Service {
	return &Service{
		repository: repo.NewRepository(),
	}
}

func (s Service) Init() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(time.Now().String()))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.repository.SearchUsers)
				r.Post("/", s.repository.CreateUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", s.repository.GetUser)
					r.Patch("/", s.repository.UpdateUser)
					r.Delete("/", s.repository.DeleteUser)
				})
			})
		})
	})

	err := http.ListenAndServe(":3333", r)
	if err != nil {
		return fmt.Errorf("ListenAndServe: %w", err)
	}
	return nil
}
