package repo

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"time"

	"trueconf-refactor/model"
)

const store = `users.json`

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (repo Repository) GetAllUsers() (*model.UserStore, error) {
	f, _ := os.ReadFile(store)
	s := model.UserStore{}
	err := json.Unmarshal(f, &s)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &s, nil
}

func (repo Repository) CreateUser(request model.CreateUserRequest) (string, error) {
	allUsers, err := repo.GetAllUsers()
	if err != nil {
		return "", fmt.Errorf("GetAllUsers: %w", err)
	}

	allUsers.Increment++
	u := model.User{
		CreatedAt:   time.Now(),
		DisplayName: request.DisplayName,
		Email:       request.DisplayName,
	}

	id := strconv.Itoa(allUsers.Increment)
	allUsers.List[id] = u

	b, err := json.Marshal(allUsers)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	err = os.WriteFile(store, b, fs.ModePerm)
	if err != nil {
		return "", fmt.Errorf("WriteFile: %w", err)
	}

	return id, nil
}

func (repo Repository) GetUser(id string) (model.User, error) {
	allUsers, err := repo.GetAllUsers()
	if err != nil {
		return model.User{}, fmt.Errorf("GetAllUsers: %w", err)
	}

	return allUsers.List[id], nil
}

func (repo Repository) UpdateUser(request model.UpdateUserRequest, id string) error {
	allUsers, err := repo.GetAllUsers()
	if err != nil {
		return fmt.Errorf("GetAllUsers: %w", err)
	}

	if _, ok := allUsers.List[id]; !ok {
		return model.UserNotFound
	}

	u := allUsers.List[id]
	u.DisplayName = request.DisplayName
	allUsers.List[id] = u

	b, err := json.Marshal(allUsers)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	err = os.WriteFile(store, b, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("WriteFile: %w", err)
	}

	return nil
}

func (repo Repository) DeleteUser(id string) error {
	allUsers, err := repo.GetAllUsers()
	if err != nil {
		return fmt.Errorf("GetAllUsers: %w", err)
	}

	if _, ok := allUsers.List[id]; !ok {
		return model.UserNotFound
	}

	delete(allUsers.List, id)

	b, _ := json.Marshal(allUsers)
	_ = os.WriteFile(store, b, fs.ModePerm)

	return nil
}
