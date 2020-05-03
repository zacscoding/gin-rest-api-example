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

func (a *articleRepository) SaveOne(data interface{}) error {
	return a.db.Save(data).Error
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

func (a *articleRepository) FindArticleBySlug(slug string) (*models.Article, error) {
	var m models.Article
	err := a.db.Where(&models.Article{Slug: slug}).Preload("Favorites").Preload("Tags").Preload("Author").Find(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (a *articleRepository) FindArticles(p Pageable) ([]models.Article, int, error) {
	var (
		articles []models.Article
		count    int
	)
	a.db.Model(&articles).Count(&count)
	err := a.db.Preload("Favorites").Preload("Tags").Preload("Author").Offset(p.Offset).Limit(p.Limit).Order("created_at desc").Find(&articles).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	return articles, count, nil
}

func (a *articleRepository) FindArticlesByTag(tag string, p Pageable) ([]models.Article, int, error) {
	var (
		t        models.Tag
		articles []models.Article
		count    int
	)
	err := a.db.Where(&models.Tag{Name: tag}).First(&t).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	a.db.Model(&t).Preload("Favorites").Preload("Tags").Preload("Author").Offset(p.Offset).Limit(p.Limit).Order("created_at desc").Association("Articles").Find(&articles)
	count = a.db.Model(&t).Association("Articles").Count()
	return articles, count, nil
}

func (a *articleRepository) FindArticlesByAuthor(username string, p Pageable) ([]models.Article, int, error) {
	var (
		u        models.User
		articles []models.Article
		count    int
	)
	err := a.db.Where(&models.User{Username: username}).First(&u).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	a.db.Where(&models.Article{AuthorID: u.ID}).Preload("Favorites").Preload("Tags").Preload("Author").Offset(p.Offset).Limit(p.Limit).Order("created_at desc").Find(&articles)
	a.db.Where(&models.Article{AuthorID: u.ID}).Model(&models.Article{}).Count(&count)
	return articles, count, nil
}

func (a *articleRepository) FindArticlesByFavoritedUsername(username string, p Pageable) ([]models.Article, int, error) {
	var (
		u        models.User
		articles []models.Article
		count    int
	)
	err := a.db.Where(&models.User{Username: username}).First(&u).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	a.db.Model(&u).Preload("Favorites", "user_id", u.ID).Preload("Tags").Preload("Author").Offset(p.Offset).Limit(p.Limit).Order("created_at desc").Find(&articles)
	a.db.Where(&models.ArticleFavorite{UserID: u.ID}).Model(&models.ArticleFavorite{}).Count(&count)
	return articles, count, nil
}

func (a *articleRepository) UpdateFavorite(articleFavorite *models.ArticleFavorite) error {
	return a.db.Create(articleFavorite).Error
}
