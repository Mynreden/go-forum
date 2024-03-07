package category

import (
	"database/sql"
	"errors"
	"forum/internal/domain"
)

type CategoryStorage struct {
	db *sql.DB
}

func NewCategoryStorage(db *sql.DB) *CategoryStorage {
	return &CategoryStorage{db}
}

func (s *CategoryStorage) GetAllCategories() ([]*domain.Category, error) {
	rows, err := s.db.Query("SELECT category_name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*domain.Category, 0)
	for rows.Next() {
		category := new(domain.Category)
		err := rows.Scan(&category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CategoryStorage) CreateCategory(category *domain.Category) (string, error) {
	if category.Name == "" {
		return "", errors.New("categories name is empty")
	}

	_, err := s.db.Exec("INSERT INTO categories (category_name) VALUES (?)", category.Name)
	if err != nil {
		return "", err
	}

	return category.Name, nil
}

// TODO: implement
func (s *CategoryStorage) GetCategoryByID(id int) (*domain.Category, error) {
	return nil, nil
}

func (s *CategoryStorage) GetCategoryByName(name string) (*domain.Category, error) {
	row := s.db.QueryRow("SELECT category_name FROM categories WHERE category_name = ?", name)

	category := new(domain.Category)
	err := row.Scan(&category.Name)
	if err != nil {

		return nil, err
	}

	return category, nil
}
