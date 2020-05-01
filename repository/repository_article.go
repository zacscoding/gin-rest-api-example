package repository

import (
	"gin-rest-api-example/database/models"
	"github.com/jinzhu/gorm"
)

type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository returns a new article repository
func NewArticleRepository(db *gorm.DB) *articleRepository {
	return &articleRepository{
		db: db,
	}
}

func (a *articleRepository) SaveArticle(article *models.Article) error {
	tx := a.db.Begin()

	if len(article.Tags) != 0 {
		var tags []models.Tag
		for _, tag := range article.Tags {
			var t models.Tag
			err := a.db.FirstOrCreate(&t, models.Tag{
				Name: tag.Name,
			}).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			tags = append(tags, t)
		}
		article.Tags = tags
	}

	err := tx.Save(article).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (a *articleRepository) SaveOne(data interface{}) error {
	return a.db.Save(data).Error
}
