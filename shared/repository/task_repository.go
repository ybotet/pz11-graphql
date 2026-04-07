package repository

import (
	"context"

	"github.com/ybotet/pz8-pipelineCICD-go/shared/models"
)

type TaskRepository interface {
	GetAll(ctx context.Context) ([]*models.Task, error)
	GetByID(ctx context.Context, id string) (*models.Task, error)
	GetByStatus(ctx context.Context, done bool) ([]*models.Task, error)
	Create(ctx context.Context, task *models.Task) error
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
}