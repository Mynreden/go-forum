package comment

import (
	"database/sql"
	"fmt"
	"forum/internal/domain"
)

type CommentStorage struct {
	db *sql.DB
}

func NewCommentStorage(db *sql.DB) *CommentStorage {
	return &CommentStorage{db: db}
}

func (s CommentStorage) CreateComment(comment *domain.Comment) error {
	query := `INSERT INTO comments (comment, post_id, user_id, userName, created_at) VALUES (?, ?, ?, ?, ?)`
	result, err := s.db.Exec(
		query,
		comment.Content,
		comment.PostID,
		comment.AuthorID,
		comment.AuthorName,
		comment.CreatedAt,
	)
	if err != nil {
		return err
	}

	if _, err := result.LastInsertId(); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (s CommentStorage) GetAllByPostID(postID int) ([]*domain.Comment, error) {
	query := `SELECT id, comment, post_id, user_id, userName, created_at FROM comments WHERE post_id = ?`
	rows, err := s.db.Query(query, postID)
	if err != nil {
		fmt.Println(err, "Error getting comments")
		return nil, err
	}
	defer rows.Close()

	comments := make([]*domain.Comment, 0)
	for rows.Next() {
		comment := new(domain.Comment)
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.PostID,
			&comment.AuthorID,
			&comment.AuthorName,
			&comment.CreatedAt,
		)
		if err != nil {
			fmt.Println(err, "Error getting comments")

			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {

		return nil, err
	}

	return comments, nil
}

func (s *CommentStorage) GetCommentByID(id int) (*domain.Comment, error) {
	query := `SELECT id, comment, post_id, user_id, userName, created_at FROM comments WHERE id = ?`
	row := s.db.QueryRow(query, id)

	comment := &domain.Comment{}
	err := row.Scan(&comment.ID,
		&comment.Content,
		&comment.PostID,
		&comment.AuthorID,
		&comment.AuthorName,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
