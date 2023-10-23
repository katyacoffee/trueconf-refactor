package service

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"trueconf-refactor/model"
	"trueconf-refactor/repo"
)

type IRepository interface {
	GetAllUsers() (*model.UserStore, error)
	CreateUser(request model.CreateUserRequest) (string, error)
	GetUser(id string) (model.User, error)
	UpdateUser(request model.UpdateUserRequest, id string) error
	DeleteUser(id string) error
}

type HttpHandlers struct {
	repository IRepository
}

func NewHttpHandlers() *HttpHandlers {
	return &HttpHandlers{
		repository: repo.NewRepository(),
	}
}

func (h *HttpHandlers) SearchUsers(w http.ResponseWriter, r *http.Request) {
	s, err := h.repository.GetAllUsers()
	if err != nil {
		fmt.Printf("ERROR in SearchUsers: %s", err.Error())
	}

	render.JSON(w, r, s.List)
	return
}

func (h *HttpHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	request := model.CreateUserRequest{}
	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, model.ErrInvalidRequest(err))
		return
	}

	id, err := h.repository.CreateUser(request)
	if err != nil {
		fmt.Printf("ERROR in CreateUser: %s", err.Error())
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func (h *HttpHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.repository.GetUser(id)
	if err != nil {
		fmt.Printf("ERROR in GetUser: %s", err.Error())
	}

	render.JSON(w, r, user)
}

func (h *HttpHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	request := model.UpdateUserRequest{}
	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, model.ErrInvalidRequest(err))
		return
	}

	id := chi.URLParam(r, "id")

	err := h.repository.UpdateUser(request, id)
	if err != nil {
		_ = render.Render(w, r, model.ErrInvalidRequest(err))
		fmt.Printf("ERROR in UpdateUser: %s", err.Error())
	}

	render.Status(r, http.StatusNoContent)
}

func (h *HttpHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.repository.DeleteUser(id)
	if err != nil {
		_ = render.Render(w, r, model.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}
