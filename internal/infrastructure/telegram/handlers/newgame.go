package handlers

import "RIP-Peroni/blood_guess/internal/application/ports"

type NewGameHandler struct {
	*BaseHandler
	createGameInput ports.CreateGameInput
}
