package graph

import (
	"github.com/sirupsen/logrus"
	"github.com/ybotet/pz8-pipelineCICD-go/services/graphql/internal/repository"
)

type Resolver struct {
	TaskRepo *repository.PostgresTaskRepository
	Logger   *logrus.Logger
}

func NewResolver(taskRepo *repository.PostgresTaskRepository, log *logrus.Logger) *Resolver {
	return &Resolver{
		TaskRepo: taskRepo,
		Logger:   log,
	}
}