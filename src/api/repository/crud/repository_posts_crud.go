package crud

import (
	"api/models"
	"api/utils/channels"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type repositoryPostsCRUD struct {
	db *gorm.DB
}

func NewRepositoryPostsCRUD(db *gorm.DB) *repositoryPostsCRUD {
	return &repositoryPostsCRUD{db}
}

func (r *repositoryPostsCRUD) Save(post models.Post) (models.Post, error) {
	var err error
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		err = r.db.Debug().Model(&models.Post{}).Create(&post).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return post, nil
	}
	return models.Post{}, err
}

func (r *repositoryPostsCRUD) FindAll() ([]models.Post, error) {
	var err error
	posts := []models.Post{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		err = r.db.Debug().Model(&models.Post{}).Limit(100).Find(&posts).Error
		if err != nil {
			ch <- false
			return
		}
		if len(posts) > 0 {
			for i, _ := range posts {
				err = r.db.Debug().Model(&models.User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
				if err != nil {
					ch <- false
					return
				}
			}
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return posts, nil
	}
	return nil, err
}

func (r *repositoryPostsCRUD) FindById(pid uint64) (models.Post, error) {
	var err error
	post := models.Post{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		err = r.db.Debug().Model(&models.Post{}).Where("id = ?", pid).Take(&post).Error
		if err != nil {
			ch <- false
			return
		}
		if post.ID != 0 {
			err = r.db.Debug().Model(&models.User{}).Where("id = ?", post.AuthorID).Take(&post.Author).Error
			if err != nil {
				ch <- false
				return
			}
		}
		ch <- true
	}(done)
	if channels.OK(done) {
		return post, nil
	}
	return models.Post{}, err
}

func (r *repositoryPostsCRUD) Update(pid uint64, post models.Post) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.Post{}).Where("id = ?", pid).Take(&models.Post{}).UpdateColumns(
			map[string]interface{}{
				"title":      post.Title,
				"content":    post.Content,
				"updated_at": time.Now(),
			},
		)
		ch <- true
	}(done)
	if channels.OK(done) {
		if rs.Error != nil {
			if gorm.IsRecordNotFoundError(rs.Error) {
				return 0, errors.New("Post not found")
			}
			return 0, rs.Error
		}
		return rs.RowsAffected, nil
	}
	return 0, rs.Error
}

func (r *repositoryPostsCRUD) Delete(pid uint64, uid uint32) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&models.Post{}).Delete(&models.Post{})
		ch <- true
	}(done)
	if channels.OK(done) {
		if rs.Error != nil {
			if gorm.IsRecordNotFoundError(rs.Error) {
				return 0, errors.New("Post not found")
			}
			return 0, rs.Error
		}
		return rs.RowsAffected, nil
	}
	return 0, rs.Error
}
