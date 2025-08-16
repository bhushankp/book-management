package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"book-management/internal/api/v1/book/handler"
	"book-management/internal/models"
	pkgerr "book-management/internal/pkg/errors"
	"book-management/internal/pkg/validator"
	"book-management/internal/service"

	"github.com/gorilla/mux"
)

type fakeRepo struct {
	items map[uint]*models.Book
	next  uint
}

func (f *fakeRepo) Create(_ context.Context, b *models.Book) (*models.Book, *pkgerr.AppError) {
	if f.items == nil {
		f.items = map[uint]*models.Book{}
	}
	f.next++
	b.ID = f.next
	f.items[b.ID] = b
	return b, nil
}
func (f *fakeRepo) GetByID(_ context.Context, id uint) (*models.Book, *pkgerr.AppError) {
	b := f.items[id]
	if b == nil {
		return nil, pkgerr.E(pkgerr.ErrNotFound, "nf", nil)
	}
	return b, nil
}
func (f *fakeRepo) List(_ context.Context, limit, offset int) ([]models.Book, int64, *pkgerr.AppError) {
	out := make([]models.Book, 0, len(f.items))
	for _, v := range f.items {
		out = append(out, *v)
	}
	return out, int64(len(out)), nil
}
func (f *fakeRepo) Update(_ context.Context, b *models.Book) *pkgerr.AppError {
	if f.items[b.ID] == nil {
		return pkgerr.E(pkgerr.ErrNotFound, "nf", nil)
	}
	f.items[b.ID] = b
	return nil
}
func (f *fakeRepo) Delete(_ context.Context, id uint) *pkgerr.AppError {
	delete(f.items, id)
	return nil
}

func TestCreateBook(t *testing.T) {
	validator.Init()
	repo := &fakeRepo{}
	svc := service.NewBookService(repo)
	h := handler.NewHandler(svc)

	r := mux.NewRouter()
	r.HandleFunc("/v1/books", h.Create).Methods(http.MethodPost)

	body, _ := json.Marshal(map[string]string{"title": "Go in Action", "author": "John", "isbn": "1234567890123"})
	req := httptest.NewRequest(http.MethodPost, "/v1/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d, body=%s", rec.Code, rec.Body.String())
	}
}
