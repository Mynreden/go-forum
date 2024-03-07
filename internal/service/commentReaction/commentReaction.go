package commentReaction

import (
	"database/sql"
	"errors"
	"forum/internal/domain"
)

type CommentReactionService struct {
	repo domain.CommentReactionRepo
}

func NewCommentReactionService(repo domain.CommentReactionRepo) *CommentReactionService {
	return &CommentReactionService{repo}
}

func (s *CommentReactionService) CreateCommentsReactions(reaction *domain.CommentReactionDTO) error {
	r, err := s.repo.GetReactionByUserIDAndCommentID(reaction.UserID, reaction.CommentID)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	} else {
		s.repo.DeleteCommentsReactions(r.ID)
		if r.Status == reaction.Status {
			return nil
		}
	}

	return s.repo.CreateCommentsReactions(reaction)
}

func (s *CommentReactionService) GetLikesAndDislikes(commentID int) (int, int, error) {
	votes, err := s.repo.GetVotesByCommentID(commentID)
	if err != nil {
		return 0, 0, err
	}

	var likes, dislikes int
	for _, v := range votes {
		if v.Status {
			likes++
		} else {
			dislikes++
		}
	}

	return likes, dislikes, nil
}
