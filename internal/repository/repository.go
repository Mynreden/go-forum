package repository

import (
	"database/sql"
	"forum/internal/domain"
	"forum/internal/repository/category"
	"forum/internal/repository/comment"
	"forum/internal/repository/commentsReaction"
	"forum/internal/repository/post"
	"forum/internal/repository/postReaction"
	"forum/internal/repository/session"
	"forum/internal/repository/user"
	// "forum/internal/repository/postReaction"
)

type Repository struct {
	CommentRepo         domain.CommentRepo
	CommentReactionRepo domain.CommentReactionRepo
	PostRepo            domain.PostRepo
	UserRepo            domain.UserRepo
	SessionRepo         domain.SessionRepo
	CategoryRepo        domain.CategoryRepo
	PostReactionRepo    domain.PostReactionRepo
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		PostRepo:            post.NewPostStorage(db),
		UserRepo:            user.NewUserStorage(db),
		CommentRepo:         comment.NewCommentStorage(db),
		SessionRepo:         session.NewSessionStorage(db),
		CategoryRepo:        category.NewCategoryStorage(db),
		PostReactionRepo:    postReaction.NewPostReactionStorage(db),
		CommentReactionRepo: commentsReaction.NewCommentsReactionsStorage(db),
	}
}
