package post

import (
	"forum/internal/domain"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type PostService struct {
	repo domain.PostRepo
}

func NewPostService(repo domain.PostRepo) *PostService {
	return &PostService{repo}
}

func (s *PostService) DeletePost(id int) error {
	return nil
}

func (p *PostService) CreatePost(postDTO *domain.CreatePostDTO) (int, error) {
	post := &domain.Post{
		Title:      postDTO.Title,
		Content:    postDTO.Content,
		AuthorID:   postDTO.Author,
		AuthorName: postDTO.AuthorName,
		Categories: postDTO.Categories,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return p.repo.CreatePost(post)
}

func (p *PostService) CreatePostWithImage(postDTO *domain.CreatePostDTO) (int, error) {
	if postDTO.ImageFile == nil {
		return p.CreatePost(postDTO)
	}

	data, err := ioutil.ReadAll(postDTO.ImageFile)
	if err != nil {
		return 0, err
	}

	fileName, err := uuid.NewV4()
	if err != nil {
		return 0, err
	}
	filePath := "ui/static/img/" + fileName.String()
	err = ioutil.WriteFile(filePath, data, 0o666)

	if err != nil {
		return 0, err
	}
	filePath = filePath[2:]
	post := &domain.Post{
		Title:      postDTO.Title,
		Content:    postDTO.Content,
		AuthorID:   postDTO.Author,
		AuthorName: postDTO.AuthorName,
		Categories: postDTO.Categories,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ImagePath:  filePath,
	}

	id, err := p.repo.CreatePostWithImage(post)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *PostService) UpdatePost(post *domain.Post) error {
	return nil
}

func (s *PostService) GetPostsByAuthorID(author int, offset int, limit int) ([]*domain.Post, error) {
	return s.repo.GetPostsByAuthor(author, offset, limit)
}

func (s *PostService) GetAllPosts(offset, limit int) ([]*domain.Post, error) {
	return s.repo.GetAllPosts(offset, limit)
}

func (p *PostService) GetPostByID(id int) (*domain.Post, error) {
	post, err := p.repo.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	if post.ImagePath == "" {

		return post, nil
	}

	post.ImagePath = ".." + strings.TrimPrefix(post.ImagePath, "ui")

	return post, nil
}

func (p *PostService) GetLikedPosts(id int, offset int, limit int) ([]*domain.Post, error) {
	return p.repo.GetLikedPosts(id, offset, limit)
}

func (p *PostService) GetPostsByCategory(category string, offset int, limit int) ([]*domain.Post, error) {
	return p.repo.GetPostsByCategory(category, offset, limit)
}
