package handlers

import "dblocker_logs_server/internal/repository"

type DBlockerHandler struct {
	Repo *repository.DBlockerRepository
}
