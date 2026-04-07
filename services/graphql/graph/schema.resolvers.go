package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ybotet/pz8-pipelineCICD-go/services/graphql/graph/generated"
	"github.com/ybotet/pz8-pipelineCICD-go/services/graphql/graph/model"
	"github.com/ybotet/pz8-pipelineCICD-go/shared/models"
)

func (r *Resolver) Tasks(ctx context.Context) ([]*models.Task, error) {
	r.Logger.Info("GraphQL Query: tasks")
	return r.TaskRepo.GetAll(ctx)
}

func (r *Resolver) Task(ctx context.Context, id string) (*models.Task, error) {
	r.Logger.Infof("GraphQL Query: task id=%s", id)
	return r.TaskRepo.GetByID(ctx, id)
}

func (r *Resolver) CreateTask(ctx context.Context, input model.CreateTaskInput) (*models.Task, error) {
	r.Logger.Info("GraphQL Mutation: createTask")
	
	now := time.Now().Format(time.RFC3339)
	task := &models.Task{
		ID:          uuid.New().String(),
		Title:       input.Title,
		Description: getString(input.Description),
		Done:        false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	err := r.TaskRepo.Create(ctx, task)
	if err != nil {
		r.Logger.WithError(err).Error("Failed to create task")
		return nil, err
	}
	return task, nil
}

func (r *Resolver) UpdateTask(ctx context.Context, id string, input model.UpdateTaskInput) (*models.Task, error) {
	r.Logger.Infof("GraphQL Mutation: updateTask id=%s", id)
	
	task, err := r.TaskRepo.GetByID(ctx, id)
	if err != nil || task == nil {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	
	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Done != nil {
		task.Done = *input.Done
	}
	task.UpdatedAt = time.Now().Format(time.RFC3339)
	
	err = r.TaskRepo.Update(ctx, task)
	if err != nil {
		r.Logger.WithError(err).Error("Failed to update task")
		return nil, err
	}
	return task, nil
}

func (r *Resolver) DeleteTask(ctx context.Context, id string) (bool, error) {
	r.Logger.Infof("GraphQL Mutation: deleteTask id=%s", id)
	
	err := r.TaskRepo.Delete(ctx, id)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (r *Resolver) Mutation() generated.MutationResolver { return r }
func (r *Resolver) Query() generated.QueryResolver { return r }

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}