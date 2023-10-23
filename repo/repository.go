package repo

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"trueconf-refactor/model"
)

const store = `users.json`

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (repo Repository) SearchUsers(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(store)
	s := model.UserStore{}
	_ = json.Unmarshal(f, &s)

	render.JSON(w, r, s.List)
	return
}

func (repo Repository) CreateUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(store)
	s := model.UserStore{}
	_ = json.Unmarshal(f, &s)

	request := model.CreateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, model.ErrInvalidRequest(err))
		return
	}

	s.Increment++
	u := model.User{
		CreatedAt:   time.Now(),
		DisplayName: request.DisplayName,
		Email:       request.DisplayName,
	}

	id := strconv.Itoa(s.Increment)
	s.List[id] = u

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(store, b, fs.ModePerm)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func (repo Repository) GetUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(store)
	s := model.UserStore{}
	_ = json.Unmarshal(f, &s)

	id := chi.URLParam(r, "id")

	render.JSON(w, r, s.List[id])
}

func (repo Repository) UpdateUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(store)
	s := model.UserStore{}
	_ = json.Unmarshal(f, &s)

	request := model.UpdateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, model.ErrInvalidRequest(err))
		return
	}

	id := chi.URLParam(r, "id")

	if _, ok := s.List[id]; !ok {
		_ = render.Render(w, r, model.ErrInvalidRequest(model.UserNotFound))
		return
	}

	u := s.List[id]
	u.DisplayName = request.DisplayName
	s.List[id] = u

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(store, b, fs.ModePerm)

	render.Status(r, http.StatusNoContent)
}

func (repo Repository) DeleteUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(store)
	s := model.UserStore{}
	_ = json.Unmarshal(f, &s)

	id := chi.URLParam(r, "id")

	if _, ok := s.List[id]; !ok {
		_ = render.Render(w, r, model.ErrInvalidRequest(model.UserNotFound))
		return
	}

	delete(s.List, id)

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(store, b, fs.ModePerm)

	render.Status(r, http.StatusNoContent)
}
