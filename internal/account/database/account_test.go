package database

import (
	"gin-rest-api-example/internal/account/model"
	"gin-rest-api-example/internal/database"
	"gin-rest-api-example/pkg/logging"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type DBSuite struct {
	suite.Suite
	db       AccountDB
	originDB *gorm.DB
}

func (s *DBSuite) SetupSuite() {
	logging.SetLevel(zapcore.FatalLevel)
	s.originDB = database.NewTestDatabase(s.T(), true)
	s.db = &accountDB{
		db: s.originDB,
	}
}

func (s *DBSuite) SetupTest() {
	s.originDB.Where("id > 0").Delete(&model.Account{})
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DBSuite))
}

func (s *DBSuite) TestSave() {
	// given
	acc := model.Account{
		Username: "user1",
		Password: "pass1",
		Email:    "user@gmail.com",
	}

	// when
	now := time.Now()
	err := s.db.Save(nil, &acc)

	// then
	s.NoError(err)
	find, err := s.db.FindByEmail(nil, acc.Email)
	s.NoError(err)
	s.Equal(acc.ID, find.ID)
	s.Equal(acc.Username, find.Username)
	s.Equal(acc.Email, find.Email)
	s.Equal(acc.Password, find.Password)
	s.Equal(acc.Bio, find.Bio)
	s.Equal(acc.Image, find.Image)
	s.WithinDuration(now, find.CreatedAt, time.Second)
	s.WithinDuration(now, find.UpdatedAt, time.Second)
	s.False(find.Disabled)
}

func (s *DBSuite) TestSave_ErrorIfExistEmail() {
	// given
	email := "user1@email.com"
	err := s.db.Save(nil, &model.Account{
		Username: "user1",
		Email:    email,
		Password: "pass",
	})
	s.NoError(err)

	// when
	err = s.db.Save(nil, &model.Account{
		Username: "user2",
		Email:    email,
		Password: "pass2",
	})

	// then
	s.Error(err)
	s.Equal(database.ErrKeyConflict, err)
}

func (s *DBSuite) TestUpdate() {
	// given
	acc := model.Account{
		Username: "user1",
		Email:    "user@gmail.com",
		Password: "pass1",
	}
	s.NoError(s.db.Save(nil, &acc))

	updated := &model.Account{
		Username: "updated-user1",
		Email:    "updated-email@gamil.com",
		Bio:      "updated-bio",
		Image:    "updated-image",
	}

	// when
	err := s.db.Update(nil, acc.Email, updated)

	// then
	s.NoError(err)

	find, err := s.db.FindByEmail(nil, acc.Email)
	s.NoError(err)
	// unchanged fields
	s.Equal(acc.Email, find.Email)
	s.Equal(acc.Password, find.Password)

	// updated fields
	s.Equal(updated.Username, find.Username)
	s.Equal(updated.Bio, find.Bio)
	s.Equal(updated.Image, find.Image)
}

func (s *DBSuite) TestUpdate_FailIfNotExist() {
	// when
	err := s.db.Update(nil, "unknown@emai.com", &model.Account{
		Username: "updated",
	})

	// then
	s.Error(err)
	s.Equal(database.ErrNotFound, err)
}

func (s *DBSuite) TestFindByEmail() {
	// given
	now := time.Now()
	acc := model.Account{
		Username: "user1",
		Email:    "user@gmail.com",
		Password: "pass1",
	}
	s.NoError(s.db.Save(nil, &acc))

	// when
	find, err := s.db.FindByEmail(nil, acc.Email)

	// then
	s.NoError(err)
	s.Equal(acc.Username, find.Username)
	s.Equal(acc.Email, find.Email)
	s.Equal(acc.Password, find.Password)
	s.Empty(find.Bio)
	s.Empty(find.Image)
	s.WithinDuration(now, find.CreatedAt, time.Second)
	s.WithinDuration(now, find.UpdatedAt, time.Second)
	s.False(find.Disabled)
}

func (s *DBSuite) TestFindByEmail_ErrorIfNotExist() {
	// when
	find, err := s.db.FindByEmail(nil, "unknown@email.com")

	// then
	s.Nil(find)
	s.Error(err)
	s.Equal(database.ErrNotFound, err)
}
