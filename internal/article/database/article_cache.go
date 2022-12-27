package database

import (
	"context"
	"fmt"
	"gin-rest-api-example/internal/article/model"
	"gin-rest-api-example/internal/cache"
	"gin-rest-api-example/internal/metric"
)

var _ ArticleDB = (*articleCacheDB)(nil)

const (
	cacheKeyArticleBySlug = "article-by-slug"
)

func newarticleCacheDB(cacher cache.Cacher, mp *metric.MetricsProvider, delegate ArticleDB) ArticleDB {
	return &articleCacheDB{
		cacher:   cacher,
		mp:       mp,
		delegate: delegate,
	}
}

type articleCacheDB struct {
	cacher   cache.Cacher
	mp       *metric.MetricsProvider
	delegate ArticleDB
}

func (ac *articleCacheDB) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	return ac.delegate.RunInTx(ctx, f)
}

func (ac *articleCacheDB) SaveArticle(ctx context.Context, article *model.Article) error {
	if err := ac.delegate.SaveArticle(ctx, article); err != nil {
		return err
	}
	if cache.IsCacheSkip(ctx) {
		return nil
	}
	key := ac.articleBySlugCacheKey(article.Slug)
	ac.cacher.Set(ctx, key, article)
	return nil
}

func (ac *articleCacheDB) FindArticleBySlug(ctx context.Context, slug string) (*model.Article, error) {
	if cache.IsCacheSkip(ctx) {
		return ac.delegate.FindArticleBySlug(ctx, slug)
	}

	var (
		item     model.Article
		key      = ac.articleBySlugCacheKey(slug)
		cacheHit = true
	)
	err := ac.cacher.Fetch(ctx, key, &item, func() (interface{}, error) {
		cacheHit = false
		return ac.delegate.FindArticleBySlug(ctx, slug)
	})
	if err != nil {
		return nil, err
	}
	ac.mp.RecordCache(cacheKeyArticleBySlug, cacheHit)
	return &item, nil
}

func (ac *articleCacheDB) FindArticles(ctx context.Context, criteria IterateArticleCriteria) ([]*model.Article, int64, error) {
	return ac.delegate.FindArticles(ctx, criteria)
}

func (ac *articleCacheDB) DeleteArticleBySlug(ctx context.Context, authorId uint, slug string) error {
	if err := ac.delegate.DeleteArticleBySlug(ctx, authorId, slug); err != nil {
		return err
	}
	// TODO: require tx?
	key := ac.articleBySlugCacheKey(slug)
	if exists, _ := ac.cacher.Exists(ctx, key); exists {
		ac.cacher.Delete(ctx, slug)
	}
	return nil
}

func (ac *articleCacheDB) SaveComment(ctx context.Context, slug string, comment *model.Comment) error {
	return ac.delegate.SaveComment(ctx, slug, comment)
}

func (ac *articleCacheDB) FindComments(ctx context.Context, slug string) ([]*model.Comment, error) {
	return ac.delegate.FindComments(ctx, slug)
}

func (ac *articleCacheDB) DeleteCommentById(ctx context.Context, authorId uint, slug string, id uint) error {
	return ac.delegate.DeleteCommentById(ctx, authorId, slug, id)
}

func (ac *articleCacheDB) DeleteComments(ctx context.Context, authorId uint, slug string) (int64, error) {
	return ac.delegate.DeleteComments(ctx, authorId, slug)
}

func (ac *articleCacheDB) articleBySlugCacheKey(slug string) string {
	return fmt.Sprintf("%s.%s", cacheKeyArticleBySlug, slug)
}
