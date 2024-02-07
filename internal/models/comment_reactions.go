package models

type CommentReaction struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"`
	CommentID int  `json:"comment_id"`
	Status    bool `json:"status"`
}

type CommentReactionDTO struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"`
	CommentID int  `json:"comment_id"`
	Status    bool `json:"status"`
}

type CommentReactionRepo interface {
	CreateCommentsReactions(reaction *CommentReactionDTO) error
	GetReactionByUserIDAndCommentID(userID, commentID int) (*CommentReaction, error)
	GetReactionsByCommentID(commentID int) ([]*CommentReaction, error)
	DeleteCommentsReactions(commentID int) error
	GetVotesByCommentID(commentID int) ([]*CommentReactionDTO, error)
}

type CommentReactionService interface {
	CreateCommentsReactions(reaction *CommentReactionDTO) error
	GetLikesAndDislikes(commentID int) (int, int, error)
}
