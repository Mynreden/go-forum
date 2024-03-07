package repository

import (
	"database/sql"
	"errors"
	"forum/internal/domain"
	"forum/internal/repository/category"
	categoryService "forum/internal/service/category"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCategoryStorage_GetAllCategories(t *testing.T) {
	testCases := []struct {
		name            string
		expectedRows    *sqlmock.Rows
		expectedError   error
		expectedRowsNum int
	}{
		{
			name: "Successful retrieval",
			expectedRows: sqlmock.NewRows([]string{"category_name"}).
				AddRow("Category1").
				AddRow("Category2"),
			expectedError:   nil,
			expectedRowsNum: 2,
		},
		{
			name:            "No categories found",
			expectedRows:    sqlmock.NewRows([]string{"category_name"}),
			expectedError:   nil,
			expectedRowsNum: 0,
		},
		{
			name:            "Database error",
			expectedRows:    nil,
			expectedError:   errors.New("database error"),
			expectedRowsNum: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)
			if tc.expectedError == nil {
				mock.ExpectQuery("SELECT category_name FROM categories").WillReturnRows(tc.expectedRows)
			} else {
				mock.ExpectQuery("SELECT category_name FROM categories").WillReturnError(tc.expectedError)
			}

			storage := category.NewCategoryStorage(db)
			service := categoryService.NewCategoryService(storage)
			categories, err := service.GetAllCategories()

			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.Len(t, categories, tc.expectedRowsNum)
			}
		})
	}
}

func TestCategoryStorage_CreateCategory(t *testing.T) {
	testCases := []struct {
		name            string
		categoryName    string
		expectedError   error
		expectedCreated string
	}{
		{
			name:            "Successful creation",
			categoryName:    "Category1",
			expectedError:   nil,
			expectedCreated: "Category1",
		},
		{
			name:            "Empty category name",
			categoryName:    "",
			expectedError:   errors.New("categories name is empty"),
			expectedCreated: "",
		},
		{
			name:            "Database error",
			categoryName:    "Category2",
			expectedError:   errors.New("database error"),
			expectedCreated: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			if tc.expectedError == nil {
				mock.ExpectExec("INSERT INTO categories").
					WithArgs(tc.categoryName).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("INSERT INTO categories").
					WithArgs(tc.categoryName).
					WillReturnError(tc.expectedError)
			}

			storage := category.NewCategoryStorage(db)
			service := categoryService.NewCategoryService(storage)
			cat := &domain.Category{Name: tc.categoryName}
			createdName, err := service.CreateCategory(cat)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedCreated, createdName)
		})
	}
}

func TestCategoryStorage_GetCategoryByName(t *testing.T) {
	testCases := []struct {
		name             string
		categoryName     string
		expectedError    error
		expectedCategory *domain.Category
	}{
		{
			name:          "Successful retrieval",
			categoryName:  "Category1",
			expectedError: nil,
			expectedCategory: &domain.Category{
				Name: "Category1",
			},
		},
		{
			name:             "Category not found",
			categoryName:     "NonexistentCategory",
			expectedError:    sql.ErrNoRows,
			expectedCategory: nil,
		},
		{
			name:             "Database error",
			categoryName:     "Category2",
			expectedError:    errors.New("database error"),
			expectedCategory: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			if tc.expectedError == nil {
				rows := sqlmock.NewRows([]string{"category_name"}).AddRow(tc.categoryName)
				mock.ExpectQuery("SELECT category_name FROM categories WHERE category_name = ?").
					WithArgs(tc.categoryName).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery("SELECT category_name FROM categories WHERE category_name = ?").
					WithArgs(tc.categoryName).
					WillReturnError(tc.expectedError)
			}

			storage := category.NewCategoryStorage(db)
			service := categoryService.NewCategoryService(storage)
			cat, err := service.GetCategoryByName(tc.categoryName)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedCategory, cat)
		})
	}
}
