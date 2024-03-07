package category

import (
	"forum/internal/domain"
)

type CategoryService struct {
	repo domain.CategoryRepo
}

func NewCategoryService(repo domain.CategoryRepo) *CategoryService {
	return &CategoryService{repo}
}

func (c *CategoryService) CreateCategory(category *domain.Category) (string, error) {
	name, err := c.repo.CreateCategory(category)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (c *CategoryService) GetAllCategories() ([]*domain.Category, error) {
	categories, err := c.repo.GetAllCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *CategoryService) GetCategoryByName(name string) (*domain.Category, error) {
	category, err := c.repo.GetCategoryByName(name)
	if err != nil {
		return nil, err
	}

	return category, nil
}
