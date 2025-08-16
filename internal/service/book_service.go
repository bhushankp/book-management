package service

import (
	"context"

	"book-management/internal/models"
	pkgerr "book-management/internal/pkg/errors"
	"book-management/internal/pkg/validator"
	"book-management/internal/repository"
)

type BookService interface {
	Create(ctx context.Context, in *models.Book) (*models.Book, *pkgerr.AppError)
	Get(ctx context.Context, id uint) (*models.Book, *pkgerr.AppError)
	List(ctx context.Context, page, pageSize int) ([]models.Book, int64, *pkgerr.AppError)
	Update(ctx context.Context, id uint, in *models.Book) (*models.Book, *pkgerr.AppError)
	Delete(ctx context.Context, id uint) *pkgerr.AppError
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(r repository.BookRepository) BookService {
	return &bookService{repo: r}
}

func (s *bookService) Create(ctx context.Context, in *models.Book) (*models.Book, *pkgerr.AppError) {
	if err := validator.ValidateStruct(in); err != nil {
		return nil, pkgerr.E(pkgerr.ErrInvalidInput, err.Error(), err)
	}
	return s.repo.Create(ctx, in)
}

func (s *bookService) Get(ctx context.Context, id uint) (*models.Book, *pkgerr.AppError) {
	return s.repo.GetByID(ctx, id)
}

func (s *bookService) List(ctx context.Context, page, pageSize int) ([]models.Book, int64, *pkgerr.AppError) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, pageSize, offset)
}

func (s *bookService) Update(ctx context.Context, id uint, in *models.Book) (*models.Book, *pkgerr.AppError) {
	if err := validator.ValidateStruct(in); err != nil {
		return nil, pkgerr.E(pkgerr.ErrInvalidInput, err.Error(), err)
	}
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	b.Title, b.Author, b.ISBN = in.Title, in.Author, in.ISBN
	if err2 := s.repo.Update(ctx, b); err2 != nil {
		return nil, err2
	}
	return b, nil
}

func (s *bookService) Delete(ctx context.Context, id uint) *pkgerr.AppError {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
