package repository

import "gin-rest-api-example/database/models"

type UserRepository interface {
	// =================================
	// User CRUD
	// =================================
	// Save save a new given user
	Save(user *models.User) error

	// Update update given user's information
	Update(user *models.User) error

	// FindByEmail returns a user given email
	FindByEmail(email string) (*models.User, error)

	// FindByUsername returns a user given username
	FindByUsername(username string) (*models.User, error)

	// =================================
	// Follows
	// =================================

	// UpdateFollow given following user follow to given follower user
	UpdateFollow(following *models.User, follower *models.User) error

	// UpdateUnFollow given following user un follow to given follower user
	UpdateUnFollow(following *models.User, follower *models.User) error

	// CountFollows returns given user's counts of (following, followers)
	CountFollows(user *models.User) (int, int, error)

	// IsFollowing returns a true if follow, otherwise false
	IsFollowing(user *models.User, follower *models.User) (bool, error)

	// FindFollowers returns follower users given user
	FindFollowers(user *models.User) ([]*models.User, error)

	// FindFollowing returns following users given user
	FindFollowing(user *models.User) ([]*models.User, error)
}

type ArticleRepository interface {
	// SaveOne save given interface
	SaveOne(data interface{}) error

	// =================================
	// Articles
	// =================================

	// SaveArticle save given article with tags
	SaveArticle(article *models.Article) error

	// FindArticleBySlug returns a article with given slug or nil if empty
	FindArticleBySlug(slug string) (*models.Article, error)

	// FindArticles return articles,count given pageable
	FindArticles(p Pageable) ([]models.Article, int, error)

	// FindArticlesByTag returns articles, count given pageable and tag name
	FindArticlesByTag(tag string, p Pageable) ([]models.Article, int, error)

	// FindArticlesByAuthor returns articles, count given pageable and author name
	FindArticlesByAuthor(username string, p Pageable) ([]models.Article, int, error)

	// FindArticlesByFavoritedUsername returns articles, count given pageable and who favorited by username
	FindArticlesByFavoritedUsername(username string, p Pageable) ([]models.Article, int, error)

	// =================================
	// Article favorites
	// =================================

	// UpdateFavorite update favorite by given user and article
	UpdateFavorite(articleFavorite *models.ArticleFavorite) error

	// =================================
	// Article comments
	// =================================

}
