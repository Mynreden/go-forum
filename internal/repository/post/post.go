package post

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"log"
	"time"
)

type PostStorage struct {
	db *sql.DB
}

func NewPostStorage(db *sql.DB) *PostStorage {
	return &PostStorage{db: db}
}

var ErrRecordNotFound = errors.New("post not found")

func (s *PostStorage) CreatePost(p *models.Post) (int, error) {
	query := `INSERT INTO posts (title, content, author_id, authorname, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at`

	args := []interface{}{p.Title, p.Content, p.AuthorID, p.AuthorName, p.CreatedAt, p.UpdatedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return 0, err
		// return 0, err
	}

	for _, category := range p.Categories {

		query = `INSERT INTO PostCategories (post_id, category_name) VALUES ($1, $2)`
		_, err = s.db.ExecContext(ctx, query, p.ID, category.Name)
		if err != nil {
			return 0, err
		}
	}

	return p.ID, nil
}

func (s *PostStorage) CreatePostWithImage(p *models.Post) (int, error) {
	query := `INSERT INTO posts (title, content, author_id, authorname, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id, created_at, updated_at`

	args := []interface{}{p.Title, p.Content, p.AuthorID, p.AuthorName, p.CreatedAt, p.UpdatedAt, p.ImagePath}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return 0, err
	}

	for _, category := range p.Categories {

		query = `INSERT INTO PostCategories (post_id, category_name) VALUES ($1, $2)`
		_, err = s.db.ExecContext(ctx, query, p.ID, category.Name)
		if err != nil {
			return 0, err
		}
	}

	query = `INSERT INTO images (post_id, image_path) VALUES ($1, $2)`
	_, err = s.db.ExecContext(ctx, query, p.ID, p.ImagePath)
	if err != nil {
		return 0, err
	}

	return p.ID, nil
}

func (s *PostStorage) DeletePost(post *models.Post) error {
	return nil
}

func (s *PostStorage) UpdatePost(post *models.Post) error {
	return nil
}

func (s *PostStorage) GetAllPosts(offset, limit int) ([]*models.Post, error) {
	query := `SELECT * FROM posts ORDER BY id DESC LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Println(err)

		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		post := models.Post{}

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.AuthorName,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			log.Println(err)

			return nil, err
		}
		err = s.getAllPostCategories(ctx, &post)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
			default:
				return nil, err

			}
		}
		err = s.getAllIMG(ctx, &post)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				post.ImagePath = ""

			default:
				return nil, err

			}
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostStorage) getAllIMG(ctx context.Context, post *models.Post) error {
	query := `SELECT image_path FROM images WHERE post_id = $1`
	img_row := s.db.QueryRowContext(ctx, query, post.ID)

	err := img_row.Scan(&post.ImagePath)
	if err != nil {
		return err
	}
	if err := img_row.Err(); err != nil {
		log.Println(err)
		return err
	}
	return err
}

func (s *PostStorage) getAllPostCategories(ctx context.Context, post *models.Post) error {
	query := `SELECT c.category_name FROM categories c
	JOIN PostCategories pc ON c.category_name = pc.category_name
	WHERE pc.post_id = $1`
	category_rows, err := s.db.QueryContext(ctx, query, post.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	for category_rows.Next() {
		category := models.Category{}

		err := category_rows.Scan(&category.Name)
		if err != nil {
			log.Println(err)
			return err
		}

		post.Categories = append(post.Categories, &category)
	}
	if err := category_rows.Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *PostStorage) GetPostByID(id int) (*models.Post, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT * FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	post := &models.Post{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.AuthorName,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Get post categories
	query = `SELECT c.category_name FROM categories c
	JOIN PostCategories pc ON c.category_name = pc.category_name
	WHERE pc.post_id = $1`

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		category := models.Category{}

		err := rows.Scan(&category.Name)
		if err != nil {
			return nil, err
		}

		post.Categories = append(post.Categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Get post image
	query = `SELECT image_path FROM images WHERE post_id = $1`
	row := s.db.QueryRowContext(ctx, query, id)

	err = row.Scan(&post.ImagePath)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			post.ImagePath = ""

			return post, nil
		default:
			// fmt.Println(err, "assssds")

			return nil, err
		}
	}

	return post, nil
}

func (s *PostStorage) GetLikedPosts(id int, offset int, limit int) ([]*models.Post, error) {
	query := `SELECT p.id, p.title, p.content, p.author_id, p.authorname, p.created_at, p.updated_at FROM posts p
	JOIN postsReactions a ON p.id = a.post_id
	WHERE a.user_id = $1 AND a.reaction = 1
	ORDER BY p.id DESC
	LIMIT $3 OFFSET $4`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, id, limit, offset)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		post := models.Post{}

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.AuthorName,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			fmt.Println(err)

			return nil, err
		}
		err = s.getAllPostCategories(ctx, &post)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
			default:
				return nil, err

			}
		}
		err = s.getAllIMG(ctx, &post)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				post.ImagePath = ""

			default:
				return nil, err

			}
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return posts, nil
}

func (s *PostStorage) GetPostsByCategory(category string, offset int, limit int) ([]*models.Post, error) {
	query := `SELECT p.* FROM posts p
	JOIN PostCategories pc ON p.id = pc.post_id
	WHERE pc.category_name = $1
	ORDER BY p.id DESC
	LIMIT $3 OFFSET $4`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, category, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		post := models.Post{}

		err := rows.Scan(&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.AuthorName,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		err = s.getAllPostCategories(ctx, &post)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
			default:
				return nil, err

			}
		}
		err = s.getAllIMG(ctx, &post)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				post.ImagePath = ""

			default:
				return nil, err

			}
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostStorage) GetPostsByAuthor(author int, offset int, limit int) ([]*models.Post, error) {
	query := `SELECT * FROM posts WHERE author_id = $1
	ORDER BY id DESC
	LIMIT $3 OFFSET $4`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, author, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		post := models.Post{}

		err := rows.Scan(&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.AuthorName,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		err = s.getAllPostCategories(ctx, &post)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
			default:
				return nil, err

			}
		}
		err = s.getAllIMG(ctx, &post)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				post.ImagePath = ""

			default:
				return nil, err

			}
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
