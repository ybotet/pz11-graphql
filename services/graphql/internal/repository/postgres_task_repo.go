package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ybotet/pz8-pipelineCICD-go/shared/models"
)

type PostgresTaskRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresTaskRepository(db *sql.DB, log *logrus.Logger) *PostgresTaskRepository {
	return &PostgresTaskRepository{
		db:     db,
		logger: log,
	}
}

func (r *PostgresTaskRepository) GetAll(ctx context.Context) ([]*models.Task, error) {
	query := `SELECT id, title, description, done, created_at, updated_at FROM tasks`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	query := `SELECT id, title, description, done, created_at, updated_at FROM tasks WHERE id = $1`
	task := &models.Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *PostgresTaskRepository) GetByStatus(ctx context.Context, done bool) ([]*models.Task, error) {
	query := `SELECT id, title, description, done, created_at, updated_at FROM tasks WHERE done = $1`
	rows, err := r.db.QueryContext(ctx, query, done)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks (id, title, description, done, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, task.ID, task.Title, task.Description, task.Done, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *PostgresTaskRepository) Update(ctx context.Context, task *models.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, done = $3, updated_at = $4 WHERE id = $5`
	result, err := r.db.ExecContext(ctx, query, task.Title, task.Description, task.Done, task.UpdatedAt, task.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task not found: %s", task.ID)
	}
	return nil
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task not found: %s", id)
	}
	return nil
}