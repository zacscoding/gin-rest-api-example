package handler


type Handler struct {
	userRepo repository.UserRepository
}


func (h *Handler) GetUser() error {
	return nil
}