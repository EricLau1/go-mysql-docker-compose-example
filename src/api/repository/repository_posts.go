package repository

import "api/models"

type PostRepository interface {
	Save(models.Post) (models.Post, error)
	FindAll() ([]models.Post, error)
	FindById(uint64) (models.Post, error)
	Update(uint64, models.Post) (int64, error)
	Delete(post_id uint64, user_id uint32) (int64, error)
}
