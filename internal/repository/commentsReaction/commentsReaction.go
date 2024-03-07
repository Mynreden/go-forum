package commentsReaction

import (
	"database/sql"
	"forum/internal/domain"
)

type CommentsReactionsStorage struct {
	db *sql.DB
}

func NewCommentsReactionsStorage(db *sql.DB) *CommentsReactionsStorage {
	return &CommentsReactionsStorage{db}
}

func (repo *CommentsReactionsStorage) CreateCommentsReactions(reaction *domain.CommentReactionDTO) error {
	query := "INSERT INTO commentsReactions (user_id, comment_id, reaction) VALUES (?, ?, ?)"

	_, err := repo.db.Exec(
		query,
		reaction.UserID,
		reaction.CommentID,
		reaction.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *CommentsReactionsStorage) DeleteCommentsReactions(commentID int) error {
	_, err := repo.db.Exec("DELETE FROM commentsReactions WHERE id = ?", commentID)
	return err
}

func (repo *CommentsReactionsStorage) GetReactionByUserIDAndCommentID(userID, commentID int) (*domain.CommentReaction, error) {

	var reaction domain.CommentReaction
	row := repo.db.QueryRow("SELECT id, reaction FROM commentsReactions WHERE user_id = ? AND comment_id = ?", userID, commentID)
	// Сканирование данных в структуру
	err := row.Scan(
		&reaction.ID,
		&reaction.Status,
	)
	if err != nil {
		// fmt.Println("================================================================")
		// При отсутствии реакции возвращаем nil, ошибку можно проверить через errors.Is(err, sql.ErrNoRows)
		return nil, err
	}

	return &reaction, nil
}

func (repo *CommentsReactionsStorage) GetReactionsByCommentID(commentID int) ([]*domain.CommentReaction, error) {
	rows, err := repo.db.Query("SELECT reaction FROM commentsReactions WHERE comment_id = ?", commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []*domain.CommentReaction

	for rows.Next() {
		// Инициализация экземпляра PostReaction перед использованием
		reaction := &domain.CommentReaction{}

		err := rows.Scan(&reaction.Status)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, reaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}

func (s *CommentsReactionsStorage) GetVotesByCommentID(commentID int) ([]*domain.CommentReactionDTO, error) {
	var votes []*domain.CommentReactionDTO
	rows, err := s.db.Query("SELECT * FROM commentsReactions WHERE comment_id = $1", commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v domain.CommentReactionDTO
		err := rows.Scan(
			&v.ID,
			&v.UserID,
			&v.CommentID,
			&v.Status,
		)
		if err != nil {
			return nil, err
		}
		votes = append(votes, &v)
	}

	return votes, nil
}
