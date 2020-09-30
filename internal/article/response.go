package article

import (
	"gin-rest-api-example/internal/article/model"
	"time"
)

type ArticleResponse struct {
	Article Article `json:"article"`
}

type ArticlesResponse struct {
	Article       []Article `json:"articles"`
	ArticlesCount int64     `json:"articlesCount"`
}

type Article struct {
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tagList"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    Author    `json:"author"`
}

type CommentResponse struct {
	Comment Comment `json:"comment"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    Author    `json:"author"`
}

type Author struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

// NewArticlesResponse converts article models and total count to ArticlesResponse
func NewArticlesResponse(articles []*model.Article, total int64) *ArticlesResponse {
	var a []Article
	for _, article := range articles {
		a = append(a, NewArticleResponse(article).Article)
	}

	return &ArticlesResponse{
		Article:       a,
		ArticlesCount: total,
	}
}

// NewArticleResponse converts article model to ArticleResponse
func NewArticleResponse(a *model.Article) *ArticleResponse {
	var tags []string
	for _, tag := range a.Tags {
		tags = append(tags, tag.Name)
	}

	return &ArticleResponse{
		Article: Article{
			Slug:      a.Slug,
			Title:     a.Title,
			Body:      a.Body,
			Tags:      tags,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
			Author: Author{
				Username: a.Author.Username,
				Bio:      a.Author.Bio,
				Image:    a.Author.Image,
			},
		},
	}
}

// NewCommentsResponse converts article comment models to CommentsResponse
func NewCommentsResponse(comments []*model.Comment) *CommentsResponse {
	var commentsRes []Comment
	for _, comment := range comments {
		commentsRes = append(commentsRes, NewCommentResponse(comment).Comment)
	}
	return &CommentsResponse{
		Comments: commentsRes,
	}
}

// NewCommentResponse converts article comment model to CommentResponse
func NewCommentResponse(comment *model.Comment) *CommentResponse {
	return &CommentResponse{
		Comment: Comment{
			ID:        comment.ID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Body:      comment.Body,
			Author: Author{
				Username: comment.Author.Username,
				Bio:      comment.Author.Bio,
				Image:    comment.Author.Image,
			},
		},
	}
}
