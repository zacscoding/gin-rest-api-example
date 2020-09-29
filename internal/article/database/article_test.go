package database

import (
	accountDB "gin-rest-api-example/internal/account/database"
	accountModel "gin-rest-api-example/internal/account/model"
	"gin-rest-api-example/internal/article/model"
	"gin-rest-api-example/internal/database"
	"gin-rest-api-example/pkg/logging"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"testing"
	"time"
)

var dUser = accountModel.Account{
	Username: "user1",
	Email:    "user1@gmail.com",
	Password: "password",
}

type DBSuite struct {
	suite.Suite
	db        ArticleDB
	accountDB accountDB.AccountDB
	originDB  *gorm.DB
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DBSuite))
}

func (s *DBSuite) SetupSuite() {
	logging.SetLevel(zapcore.FatalLevel)
	s.originDB = database.NewTestDatabase(s.T(), true)
	s.db = &articleDB{db: s.originDB}
	s.accountDB = accountDB.NewAccountDB(s.originDB)
}

func (s *DBSuite) SetupTest() {
	s.NoError(database.DeleteRecordAll(s.T(), s.originDB, []string{
		"comments", "id > 0",
		"article_tags", "article_id > 0",
		"tags", "id > 0",
		"articles", "id > 0",
		"accounts", "id > 0",
	}))
	s.NoError(s.accountDB.Save(nil, &dUser))
}

func (s *DBSuite) TestSaveArticle() {
	// given
	s.NoError(s.originDB.Create(&model.Tag{Name: "tag1"}).Error)
	article := newArticle("title1", "title1", "body", dUser, []string{"tag1", "tag2"})

	// when
	now := time.Now()
	err := s.db.SaveArticle(nil, article)

	// then
	s.NoError(err)
	find, err := s.db.FindArticleBySlug(nil, article.Slug)
	s.NoError(err)
	s.NotEqual(0, find.ID)
	s.Equal(article.Slug, find.Slug)
	s.Equal(article.Title, find.Title)
	s.Equal(article.Body, find.Body)
	s.WithinDuration(now, find.CreatedAt, time.Second)
	s.WithinDuration(now, find.UpdatedAt, time.Second)
	s.Equal(int64(0), find.DeletedAtUnix)
	s.Equal(article.Author, dUser)
	s.assertArticleTag(find, []string{"tag1", "tag2"})
}

func (s *DBSuite) TestSaveArticle_WithSameSlugAfterDeleted() {
	// given
	article := newArticle("title1", "title1", "body", dUser, []string{"tag1", "tag2"})
	s.NoError(s.db.SaveArticle(nil, article))
	s.NoError(s.db.DeleteArticleBySlug(nil, dUser.ID, article.Slug))

	article2 := newArticle(article.Slug, article.Title, article.Body, dUser, []string{})

	// when
	err := s.db.SaveArticle(nil, article2)

	// then
	s.NoError(err)
}

func (s *DBSuite) TestSaveArticle_FailIfDuplicateSlug() {
	// given
	article := newArticle("title1", "title1", "body", dUser, []string{"tag1"})
	s.NoError(s.db.SaveArticle(nil, article))

	// when
	article2 := newArticle(article.Slug, article.Title, article.Body, dUser, nil)
	err := s.db.SaveArticle(nil, article2)

	// then
	s.Error(err)
	s.Equal(database.ErrKeyConflict, err)
}

func (s *DBSuite) TestFindArticleBySlug() {
	// given
	article := newArticle("title1", "title1", "body", dUser, []string{"tag1"})
	//now := time.Now()
	s.NoError(s.db.SaveArticle(nil, article))

	// when
	find, err := s.db.FindArticleBySlug(nil, article.Slug)

	// then
	s.NoError(err)
	s.assertArticle(article, find)
}

func (s *DBSuite) TestFindArticleBySlug_FailIfNotExist() {
	// when
	find, err := s.db.FindArticleBySlug(nil, "not-exist-slug")

	// then
	s.Nil(find)
	s.Error(err)
	s.Equal(database.ErrNotFound, err)
}

func (s *DBSuite) TestFindArticleBySlug_FailIfDeleted() {
	// given
	article := newArticle("title1", "title1", "body", dUser, []string{"tag1"})
	s.NoError(s.db.SaveArticle(nil, article))
	_, err := s.db.FindArticleBySlug(nil, article.Slug)
	s.NoError(err)
	s.NoError(s.db.DeleteArticleBySlug(nil, dUser.ID, article.Slug))

	// when
	find, err := s.db.FindArticleBySlug(nil, article.Slug)

	// then
	s.Nil(find)
	s.Error(err)
	s.Equal(database.ErrNotFound, err)
}

func (s *DBSuite) TestFindArticles() {
	// given
	// User1
	// article1 - tag1, tag2 <- second itr [0]
	// article2 - tag1		 <- first itr [1]
	// article3 - tag4
	// article4 - tag3
	// article5 - tag1		 <- first itr [0]
	// article6 - tag1 (deleted)
	// User2
	// article7 - tag1
	user1 := accountModel.Account{Username: "test-user1", Email: "test-user1@gmail.com", Password: "password"}
	s.NoError(s.accountDB.Save(nil, &user1))
	article1 := newArticle("article1", "article1", "body1", user1, []string{"tag1", "tag2"})
	s.NoError(s.db.SaveArticle(nil, article1))
	article2 := newArticle("article2", "article2", "body2", user1, []string{"tag1"})
	s.NoError(s.db.SaveArticle(nil, article2))
	article3 := newArticle("article3", "article3", "body3", user1, []string{"tag4"})
	s.NoError(s.db.SaveArticle(nil, article3))
	article4 := newArticle("article4", "article4", "body4", user1, []string{"tag3"})
	s.NoError(s.db.SaveArticle(nil, article4))
	article5 := newArticle("article5", "article5", "body5", user1, []string{"tag1"})
	s.NoError(s.db.SaveArticle(nil, article5))
	article6 := newArticle("article6", "article6", "body6", user1, []string{"tag1"})
	s.NoError(s.db.SaveArticle(nil, article6))
	s.NoError(s.db.DeleteArticleBySlug(nil, user1.ID, article6.Slug))

	user2 := accountModel.Account{Username: "test-user2", Email: "test-user2@gmail.com", Password: "password"}
	s.NoError(s.accountDB.Save(nil, &user2))
	article7 := newArticle("article7", "article7", "body7", user2, []string{"tag1"})
	s.NoError(s.db.SaveArticle(nil, article7))

	criteria := IterateArticleCriteria{
		Tags:   []string{"tag1", "tag2"},
		Author: user1.Username,
		Offset: 0,
		Limit:  2,
	}

	// when : first iteration
	results, total, err := s.db.FindArticles(nil, criteria)

	// then
	s.NoError(err)
	s.Equal(int64(3), total)
	s.Equal(2, len(results))
	s.assertArticle(article5, results[0])
	s.assertArticle(article2, results[1])

	// second iteration
	criteria.Offset = criteria.Offset + uint(len(results))
	results, total, err = s.db.FindArticles(nil, criteria)

	// then
	s.NoError(err)
	s.Equal(int64(3), total)
	s.Equal(1, len(results))
	s.assertArticle(article1, results[0])
}

func (s *DBSuite) TestDeleteArticleBySlug() {
	// given
	article := newArticle("title1", "title1", "body", dUser, []string{"tag1"})
	s.NoError(s.db.SaveArticle(nil, article))

	// when
	err := s.db.DeleteArticleBySlug(nil, dUser.ID, article.Slug)

	// then
	s.NoError(err)
	find, err := s.db.FindArticleBySlug(nil, article.Slug)
	s.Nil(find)
	s.Equal(database.ErrNotFound, err)
}

func (s *DBSuite) TestDeleteArticleBySlug_FailIfNotExist() {
	// given
	article := newArticle("title1", "title1", "body", dUser, []string{"tag1"})
	s.NoError(s.db.SaveArticle(nil, article))

	cases := []struct {
		AuthorID uint
		Slug     string
	}{
		{
			AuthorID: dUser.ID,
			Slug:     "not-exist-slug",
		}, {
			AuthorID: dUser.ID + 1,
			Slug:     article.Slug,
		},
	}

	for _, tc := range cases {
		// when
		err := s.db.DeleteArticleBySlug(nil, tc.AuthorID, tc.Slug)

		// then
		s.Error(err)
		s.Equal(database.ErrNotFound, err)
	}
}

func (s *DBSuite) assertArticle(expected, actual *model.Article) {
	s.Equal(expected.Slug, actual.Slug)
	s.Equal(expected.Title, actual.Title)
	s.Equal(expected.Body, actual.Body)
	s.WithinDuration(expected.CreatedAt, actual.CreatedAt, time.Second)
	s.WithinDuration(expected.UpdatedAt, actual.UpdatedAt, time.Second)
	s.Equal(expected.DeletedAtUnix, actual.DeletedAtUnix)
	s.Equal(expected.Author.ID, actual.Author.ID)
	s.Equal(expected.Author.Email, actual.Author.Email)
	s.Equal(expected.Author.Username, actual.Author.Username)
	var tags []string
	for _, tag := range expected.Tags {
		tags = append(tags, tag.Name)
	}
	s.assertArticleTag(actual, tags)
}

func (s *DBSuite) assertArticleTag(article *model.Article, tags []string) {
	s.Equal(len(article.Tags), len(tags))
	if len(article.Tags) == 0 {
		return
	}
	m := make(map[string]struct{})
	for _, tag := range article.Tags {
		m[tag.Name] = struct{}{}
	}
	for _, tag := range tags {
		_, ok := m[tag]
		s.True(ok)
	}
}

func newArticle(slug, title, body string, author accountModel.Account, tags []string) *model.Article {
	var tagArr []*model.Tag
	for _, tag := range tags {
		tagArr = append(tagArr, &model.Tag{Name: tag})
	}
	return &model.Article{
		Slug:   slug,
		Title:  title,
		Body:   body,
		Author: author,
		Tags:   tagArr,
	}
}
