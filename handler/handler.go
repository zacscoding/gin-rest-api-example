package handler

import "gin-rest-api-example/repository"

type Handler struct {
	userRepo repository.UserRepository
}


func (h *Handler) GetUser() error {
	return nil
}