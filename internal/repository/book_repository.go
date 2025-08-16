package repository

import (
	"context"

	"book-management/internal/models"
	pkgerr "book-management/internal/pkg/errors"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(ctx context.Context, b *models.Book) (*models.Book, *pkgerr.AppError)
	GetByID(ctx context.Context, id uint) (*models.Book, *pkgerr.AppError)
	List(ctx context.Context, limit, offset int) ([]models.Book, int64, *pkgerr.AppError)
	Update(ctx context.Context, b *models.Book) *pkgerr.AppError
	Delete(ctx context.Context, id uint) *pkgerr.AppError
}

type bookRepo struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepo{db: db}
}

func (r *bookRepo) Create(ctx context.Context, b *models.Book) (*models.Book, *pkgerr.AppError) {
	if err := r.db.WithContext(ctx).Create(b).Error; err != nil {
		return nil, pkgerr.E(pkgerr.ErrInternal, "failed to create book", err)
	}
	return b, nil
}

func (r *bookRepo) GetByID(ctx context.Context, id uint) (*models.Book, *pkgerr.AppError) {
	var out models.Book
	if err := r.db.WithContext(ctx).First(&out, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkgerr.E(pkgerr.ErrNotFound, "book not found", err)
		}
		return nil, pkgerr.E(pkgerr.ErrInternal, "failed to get book", err)
	}
	return &out, nil
}

func (r *bookRepo) List(ctx context.Context, limit, offset int) ([]models.Book, int64, *pkgerr.AppError) {
	var items []models.Book
	var count int64
	q := r.db.WithContext(ctx).Model(&models.Book{})
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, pkgerr.E(pkgerr.ErrInternal, "failed to count books", err)
	}
	if err := q.Limit(limit).Offset(offset).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, pkgerr.E(pkgerr.ErrInternal, "failed to list books", err)
	}
	return items, count, nil
}

func (r *bookRepo) Update(ctx context.Context, b *models.Book) *pkgerr.AppError {
	if err := r.db.WithContext(ctx).Save(b).Error; err != nil {
		return pkgerr.E(pkgerr.ErrInternal, "failed to update book", err)
	}
	return nil
}

func (r *bookRepo) Delete(ctx context.Context, id uint) *pkgerr.AppError {
	if err := r.db.WithContext(ctx).Delete(&models.Book{}, id).Error; err != nil {
		return pkgerr.E(pkgerr.ErrInternal, "failed to delete book", err)
	}
	return nil
}
