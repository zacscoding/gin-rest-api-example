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

	// Following given following user follow to given follower user
	Following(following *models.User, follower *models.User) error

	// UnFollowing given following user un follow to given follower user
	UnFollowing(following *models.User, follower *models.User) error

	// CountFollows returns given user's counts of (following, followers)
	CountFollows(user *models.User) (int, int, error)

	// IsFollowing returns a true if follow, otherwise false
	IsFollowing(user *models.User, follower *models.User) (bool, error)

	// FindFollowers returns follower users given user
	FindFollowers(user *models.User) ([]*models.User, error)

	// FindFollowing returns following users given user
	FindFollowing(user *models.User) ([]*models.User, error)
}
