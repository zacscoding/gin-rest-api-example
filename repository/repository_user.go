package repository

import (
	"gin-rest-api-example/database/models"
	"github.com/jinzhu/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a new user repository
func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Save(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *userRepository) Update(user *models.User) error {
	return u.db.Model(user).Update(user).Error
}

func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	var m models.User
	if err := u.db.Where(&models.User{Email: email}).First(&m).Error; err != nil {
		// RecordNotFound
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (u *userRepository) FindByUsername(username string) (*models.User, error) {
	var m models.User
	if err := u.db.Where(&models.User{Username: username}).First(&m).Error; err != nil {
		// RecordNotFound
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (u *userRepository) UpdateFollow(following *models.User, follower *models.User) error {
	return u.db.Create(&models.Follow{
		FollowerID:  follower.ID,
		FollowingID: following.ID,
	}).Error
}

func (u *userRepository) UpdateUnFollow(following *models.User, follower *models.User) error {
	return u.db.Where(models.Follow{
		FollowerID:  follower.ID,
		FollowingID: following.ID,
	}).Delete(models.Follow{}).Error
}

func (u *userRepository) CountFollows(user *models.User) (int, int, error) {
	rawQuery := `
SELECT
    SUM(CASE WHEN follower_id = ? THEN 1 ELSE 0 END) as following_cnt,
    SUM(CASE WHEN following_id = ? THEN 1 ELSE 0 END) as follower_cnt
FROM
    follows
WHERE
    (follower_id = ? OR following_id = ?);
`
	type FollowCount struct {
		FollowingCnt int
		FollowerCnt  int
	}

	var result FollowCount
	err := u.db.Raw(rawQuery, user.ID, user.ID, user.ID, user.ID).Scan(&result).Error
	if err != nil {
		return 0, 0, err
	}
	return result.FollowingCnt, result.FollowerCnt, nil
}

func (u *userRepository) IsFollowing(user *models.User, follower *models.User) (bool, error) {
	var f models.Follow
	err := u.db.Where("following_id = ? AND follower_id = ?", user.ID, follower.ID).Find(&f).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// TODO : refactor to use models
func (u *userRepository) FindFollowers(user *models.User) ([]*models.User, error) {
	var followers []*models.User
	err := u.db.Table("users").
		Select("users.*").
		Joins("left join follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", user.ID).
		Order("follows.created_at DESC").
		Scan(&followers).Error
	if err != nil {
		return nil, err
	}
	return followers, nil
}

// TODO : refactor to use models
func (u *userRepository) FindFollowing(user *models.User) ([]*models.User, error) {
	var following []*models.User
	err := u.db.Table("users").
		Select("users.*").
		Joins("left join follows ON users.id = follows.following_id").
		Where("follows.follower_id = ?", user.ID).
		Order("follows.created_at DESC").
		Scan(&following).Error
	if err != nil {
		return nil, err
	}
	return following, nil
}