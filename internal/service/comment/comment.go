package comment

import (
	"forum/internal/domain"
	"time"
)

type CommentService struct {
	repo domain.CommentRepo
}

func NewCommentService(repo domain.CommentRepo) *CommentService {
	return &CommentService{repo}
}

func (c *CommentService) CreateComment(commentDTO *domain.CreateCommentDTO) error {
	comment := &domain.Comment{
		Content:    commentDTO.Content,
		AuthorID:   commentDTO.AuthorID,
		AuthorName: commentDTO.AuthorName,
		PostID:     commentDTO.PostID,
		CreatedAt:  time.Now(),
	}

	err := c.repo.CreateComment(comment)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommentService) GetAllByPostID(postID int) ([]*domain.Comment, error) {
	comments, err := c.repo.GetAllByPostID(postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentService) GetCommentByID(id int) (*domain.Comment, error) {
	return c.repo.GetCommentByID(id)
}

// soon===========

// func (c *CommentService) DeleteComment(id int) error {
// 	return nil
// }

// func (c *CommentService) UpdateComment(comment *domain.Comment) error {
// 	return nil
// }
