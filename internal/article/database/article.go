package database

import (
	"context"
	"gin-rest-api-example/internal/article/model"
	"gin-rest-api-example/internal/database"
	"gin-rest-api-example/pkg/logging"
	"gorm.io/gorm"
	"time"
)

type IterateArticleCriteria struct {
	Tags   []string
	Author string
	Offset uint
	Limit  uint
}

//go:generate mockery --name ArticleDB --filename article_mock.go
type ArticleDB interface {
	// SaveArticle saves a given article with tags.
	// if not exist tags, then save a new tag
	SaveArticle(ctx context.Context, article *model.Article) error

	// FindArticleBySlug returns a article with given slug
	// database.ErrNotFound error is returned if not exist
	FindArticleBySlug(ctx context.Context, slug string) (*model.Article, error)

	// FindArticles returns article list with given criteria and total count
	FindArticles(ctx context.Context, criteria IterateArticleCriteria) ([]*model.Article, int64, error)

	// DeleteArticleBySlug deletes a article with given slug
	// and returns nil if success to delete, otherwise returns an error
	DeleteArticleBySlug(ctx context.Context, authorId uint, slug string) error

	// SaveComment saves a comment with given article slug and comment
	SaveComment(ctx context.Context, slug string, comment *model.Comment) error

	// FindComments returns all comments with given article slug
	FindComments(ctx context.Context, slug string) ([]*model.Comment, error)

	// DeleteCommentById deletes a comment with given article slug and comment id
	// database.ErrNotFound error is returned if not exist
	DeleteCommentById(ctx context.Context, authorId uint, slug string, id uint) error

	// DeleteComments deletes all comment with given author id and slug
	// and returns deleted records count
	DeleteComments(ctx context.Context, authorId uint, slug string) (int64, error)
}

type articleDB struct {
	db *gorm.DB
}

func (a *articleDB) SaveArticle(ctx context.Context, article *model.Article) error {
	logger := logging.FromContext(ctx)
	db := database.FromContext(ctx, a.db)
	logger.Debugw("article.db.SaveArticle", "article", article)

	for _, tag := range article.Tags {
		if err := db.WithContext(ctx).FirstOrCreate(&tag, "name = ?", tag.Name).Error; err != nil {
			logger.Errorw("article.db.SaveArticle failed to first or save tag", "err", err)
			return err
		}
	}

	if err := db.WithContext(ctx).Create(article).Error; err != nil {
		logger.Errorw("article.db.SaveArticle failed to save article", "err", err)
		if database.IsKeyConflictErr(err) {
			return database.ErrKeyConflict
		}
		return err
	}
	return nil
}

func (a *articleDB) FindArticleBySlug(ctx context.Context, slug string) (*model.Article, error) {
	logger := logging.FromContext(ctx)
	db := database.FromContext(ctx, a.db)
	logger.Debugw("article.db.FindArticleBySlug", "slug", slug)

	var ret model.Article
	// 1) load article with author
	// SELECT articles.*, accounts.*
	// FROM `articles` LEFT JOIN `accounts` `Author` ON `articles`.`author_id` = `Author`.`id`
	// WHERE slug = "title1" AND deleted_at_unix = 0 ORDER BY `articles`.`id` LIMIT 1
	err := db.Joins("Author").
		First(&ret, "slug = ? AND deleted_at_unix = 0", slug).Error
	// 2) load tags
	if err == nil {
		// SELECT * from tags JOIN article_tags ON article_tags.tag_id = tags.id AND article_tags.article_id = ?
		err = db.Model(&ret).Association("Tags").Find(&ret.Tags)
	}

	if err != nil {
		logger.Errorw("failed to find article", "err", err)
		if database.IsRecordNotFoundErr(err) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}
	return &ret, nil
}

func (a *articleDB) FindArticles(ctx context.Context, criteria IterateArticleCriteria) ([]*model.Article, int64, error) {
	logger := logging.FromContext(ctx)
	db := database.FromContext(ctx, a.db)
	logger.Debugw("article.db.FindArticles", "criteria", criteria)

	chain := db.Table("articles a").Where("deleted_at_unix = 0")
	if len(criteria.Tags) != 0 {
		chain = chain.Where("t.name IN ?", criteria.Tags)
	}
	if criteria.Author != "" {
		chain = chain.Where("au.username = ?", criteria.Author)
	}
	if len(criteria.Tags) != 0 {
		chain = chain.Joins("LEFT JOIN article_tags ats on ats.article_id = a.id").
			Joins("LEFT JOIN tags t on t.id = ats.tag_id")
	}
	if criteria.Author != "" {
		chain = chain.Joins("LEFT JOIN accounts au on au.id = a.author_id")
	}

	// get total count
	var totalCount int64
	err := chain.Distinct("a.id").Count(&totalCount).Error
	if err != nil {
		logger.Error("failed to get total count", "err", err)
	}

	// get article ids
	rows, err := chain.Select("DISTINCT(a.id) id").
		Offset(int(criteria.Offset)).
		Limit(int(criteria.Limit)).
		Order("a.id DESC").
		Rows()
	if err != nil {
		logger.Error("failed to read article ids", "err", err)
		return nil, 0, err
	}
	var ids []uint
	for rows.Next() {
		var id uint
		err := rows.Scan(&id)
		if err != nil {
			logger.Error("failed to scan id from id rows", "err", err)
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	// get articles with author by ids
	var ret []*model.Article
	if len(ids) == 0 {
		return []*model.Article{}, totalCount, nil
	}
	err = db.Joins("Author").
		Where("articles.id IN (?)", ids).
		Order("articles.id DESC").
		Find(&ret).Error
	if err != nil {
		logger.Error("failed to find article by ids", "err", err)
		return nil, 0, err
	}

	// get tags by article ids
	ma := make(map[uint]*model.Article)
	for _, r := range ret {
		ma[r.ID] = r
	}
	type ArticleTag struct {
		model.Tag
		ArticleId uint
	}
	batchSize := 100 // TODO : config
	for i := 0; i < len(ret); i += batchSize {
		var at []*ArticleTag
		last := i + batchSize
		if last > len(ret) {
			last = len(ret)
		}

		err = db.Table("tags").
			Where("article_tags.article_id IN (?)", ids[i:last]).
			Joins("LEFT JOIN article_tags ON article_tags.tag_id = tags.id").
			Select("tags.*, article_tags.article_id article_id").
			Find(&at).Error

		if err != nil {
			logger.Error("failed to load tags by article ids", "articleIds", ids[i:last], "err", err)
			return nil, 0, err
		}
		for _, tag := range at {
			a := ma[tag.ArticleId]
			a.Tags = append(a.Tags, &tag.Tag)
		}
	}
	return ret, totalCount, nil
}

func (a *articleDB) DeleteArticleBySlug(ctx context.Context, authorId uint, slug string) error {
	logger := logging.FromContext(ctx)
	db := database.FromContext(ctx, a.db)
	logger.Debugw("article.db.DeleteArticleBySlug", "slug", slug)

	chain := db.Model(&model.Article{}).
		Where("slug = ?", slug).
		Where("author_id = ?", authorId).
		Update("deleted_at_unix", time.Now().Unix())
	if chain.Error != nil {
		logger.Errorw("failed to delete an article", "err", chain.Error)
		return chain.Error
	}
	if chain.RowsAffected == 0 {
		logger.Error("failed to delete an article because not found")
		return database.ErrNotFound
	}
	return nil
}
