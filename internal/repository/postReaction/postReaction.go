package postReaction

import (
	"database/sql"
	"forum/internal/models"
)

type PostReactionRepo struct {
	db *sql.DB
}

func NewPostReactionStorage(db *sql.DB) *PostReactionRepo {
	return &PostReactionRepo{db}
}

func (repo *PostReactionRepo) CreatePostReaction(reaction *models.PostReactionDTO) error {

	query := "INSERT INTO postsReactions (user_id, post_id, reaction) VALUES (?, ?, ?)"

	// Выполняем запрос на вставку данных
	_, err := repo.db.Exec(
		query,
		reaction.UserID,
		reaction.PostID,
		reaction.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostReactionRepo) GetPostReactionsByPostID(postID int) ([]*models.PostReaction, error) {
	rows, err := repo.db.Query("SELECT reaction FROM postsReactions WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []*models.PostReaction

	for rows.Next() {
		// Инициализация экземпляра PostReaction перед использованием
		reaction := &models.PostReaction{}

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

func (repo *PostReactionRepo) GetPostsReactionsByUserID(userID int) ([]*models.PostReaction, error) {

	rows, err := repo.db.Query("SELECT post_id, reaction FROM postsReactions WHERE user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reactions []*models.PostReaction

	for rows.Next() {
		var reaction *models.PostReaction

		err := rows.Scan(
			&reaction.PostID,
			&reaction.Status,
		)

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

func (repo *PostReactionRepo) GetAllPostReactions() ([]*models.PostReaction, error) {

	rows, err := repo.db.Query("SELECT user_id, post_id, reaction FROM postsReactions")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []*models.PostReaction

	for rows.Next() {
		var reaction models.PostReaction
		err := rows.Scan(
			&reaction.UserID,
			&reaction.PostID,
			&reaction.Status,
		)
		if err != nil {
			return nil, err
		}

		reactions = append(reactions, &reaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}

func (repo *PostReactionRepo) GetReactionByUserIDAndPostID(userID, postID int) (*models.PostReaction, error) {
	// Запрос с использованием WHERE для фильтрации по userID и postID
	row := repo.db.QueryRow("SELECT id, reaction FROM postsReactions WHERE user_id = ? AND post_id = ?", userID, postID)

	var reaction models.PostReaction
	// Сканирование данных в структуру
	err := row.Scan(
		&reaction.ID,
		&reaction.Status,
	)
	if err != nil {
		// При отсутствии реакции возвращаем nil, ошибку можно проверить через errors.Is(err, sql.ErrNoRows)
		return nil, err
	}

	return &reaction, nil
}

func (repo *PostReactionRepo) DeletePostReactionByID(reactionID int) error {
	_, err := repo.db.Exec("DELETE FROM postsReactions WHERE id = ?", reactionID)
	return err
}
