package postReaction

import (
	"database/sql"
	"errors"
	"forum/internal/models"
)

type PostReactionService struct {
	repo models.PostReactionRepo
}

func NewPostReactionService(repo models.PostReactionRepo) *PostReactionService {
	return &PostReactionService{repo}
}

func (s *PostReactionService) CreatePostReaction(reaction *models.PostReactionDTO) error {
	r, err := s.repo.GetReactionByUserIDAndPostID(reaction.UserID, reaction.PostID)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	} else {
		s.repo.DeletePostReactionByID(r.ID)
		if r.Status == reaction.Status {
			return nil
		}
	}

	return s.repo.CreatePostReaction(reaction)
}

func (s *PostReactionService) GetAllPostReactionsByPostID(posts []*models.Post) error {
	for i, r := range posts {
		reactions, err := s.repo.GetPostReactionsByPostID(r.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		var like, dislike int
		for _, r := range reactions {
			// switch r.Status {
			// case 1:
			// 	like++
			// case 0:
			// 	dislike++

			// }
			if !r.Status {
				dislike++
			} else {
				like++
			}
		}
		posts[i].Likes = like
		posts[i].Dislikes = dislike

	}

	return nil
}

func (s *PostReactionService) PutReactionsToPost(post *models.Post) error {
	reactions, err := s.repo.GetPostReactionsByPostID(post.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	var like, dislike int
	for _, r := range reactions {
		if !r.Status {
			dislike++
		} else {
			like++
		}
	}
	post.Likes = like
	post.Dislikes = dislike
	return nil
}
