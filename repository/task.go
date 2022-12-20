package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	tasks := []entity.Task{}
	err := r.db.Model(&entity.Task{}).Where("user_id = ?", id).Order("id").Scan(&tasks).Error
	return tasks, err // TODO: replace this
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	err = r.db.Create(&task).Error
	return task.ID, err // TODO: replace this
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	task := entity.Task{}
	err := r.db.Model(&entity.Task{}).Where("id = ?", id).Scan(&task).Error
	return task, err // TODO: replace this
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	tasks := []entity.Task{}
	err := r.db.Model(&entity.Task{}).Where("category_id = ?", catId).Scan(&tasks).Error
	return tasks, err // TODO: replace this
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	return r.db.Model(&task).Updates(task).Error
	// return r.db.Save(&task).Error // TODO: replace this
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	return r.db.Delete(&entity.Task{}, id).Error // TODO: replace this
}
