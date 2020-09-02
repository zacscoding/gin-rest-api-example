package database

import (
	"fmt"
	"gin-rest-api-example/internal/account/model"
	"gin-rest-api-example/internal/database"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type DBSuite struct {
	suite.Suite
	db       AccountDB
	originDB *gorm.DB
}

func (s *DBSuite) TestSave() {
	// TODO : temporary
	// given
	acc := model.Account{
		Username: "user1",
		Password: "pass1",
		Email:    "user@gmail.com",
	}

	// when
	err := s.db.Save(nil, &acc)

	// then
	s.NoError(err)
	var find model.Account
	s.originDB.First(&find, "id = ?", acc.ID)
	fmt.Println(">> Find account:", find.String())
}

func (s *DBSuite) SetupSuite() {
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
