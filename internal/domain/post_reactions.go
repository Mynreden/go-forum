package domain

type PostReaction struct {
	ID     int  `json:"id"`
	UserID int  `json:"user_id"`
	PostID int  `json:"post_id"`
	Status bool `json:"status"`
}

type PostReactionDTO struct {
	ID     int  `json:"id"`
	UserID int  `json:"user_id"`
	PostID int  `json:"post_id"`
	Status bool `json:"status"`
}

type PostReactionRepo interface {
	CreatePostReaction(*PostReactionDTO) error
	GetPostReactionsByPostID(postID int) ([]*PostReaction, error)
	GetPostsReactionsByUserID(userID int) ([]*PostReaction, error)
	GetReactionByUserIDAndPostID(userID, postID int) (*PostReaction, error)
	DeletePostReactionByID(reactionID int) error
}

type PostReactionService interface {
	CreatePostReaction(reaction *PostReactionDTO) error
	GetAllPostReactionsByPostID(posts []*Post) error
	PutReactionsToPost(post *Post) error
}
