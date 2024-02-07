package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/service/category"
	"forum/internal/service/comment"
	"forum/internal/service/commentReaction"
	"forum/internal/service/post"
	"forum/internal/service/postReaction"
	"forum/internal/service/session"
	"forum/internal/service/user"
	"log"
)

type Service struct {
	UserService            models.UserService
	PostService            models.PostService
	CommentService         models.CommentService
	CommentReactionService models.CommentReactionService
	SessionService         models.SessionServise
	CategoryService        models.CategoryService
	PostReactionService    models.PostReactionService
	Log                    *log.Logger
}

func NewService(repo *repository.Repository, log *log.Logger) *Service {
	return &Service{
		UserService:            user.NewUserService(repo.UserRepo),
		PostService:            post.NewPostService(repo.PostRepo),
		CommentService:         comment.NewCommentService(repo.CommentRepo),
		CommentReactionService: commentReaction.NewCommentReactionService(repo.CommentReactionRepo),
		SessionService:         session.NewSessionService(repo.SessionRepo),
		CategoryService:        category.NewCategoryService(repo.CategoryRepo),
		PostReactionService:    postReaction.NewPostReactionService(repo.PostReactionRepo),
		Log:                    log,
	}
}
