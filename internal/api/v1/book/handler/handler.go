package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"book-management/internal/models"
	pkgerr "book-management/internal/pkg/errors"
	"book-management/internal/pkg/logger"
	"book-management/internal/pkg/response"
	"book-management/internal/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct{ svc service.BookService }

func NewHandler(s service.BookService) *Handler { return &Handler{svc: s} }

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	size, _ := strconv.Atoi(q.Get("page_size"))
	items, total, err := h.svc.List(r.Context(), page, size)
	if err != nil {
		writeErr(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"data": items, "total": total})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseID(mux.Vars(r)["id"])
	if appErr != nil {
		writeErr(w, appErr)
		return
	}
	item, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in struct{ Title, Author, ISBN string }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, pkgerr.E(pkgerr.ErrInvalidInput, "invalid body", err))
		return
	}
	item, err := h.svc.Create(r.Context(), &models.Book{Title: in.Title, Author: in.Author, ISBN: in.ISBN})
	if err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseID(mux.Vars(r)["id"])
	if appErr != nil {
		writeErr(w, appErr)
		return
	}
	var in struct{ Title, Author, ISBN string }
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeErr(w, pkgerr.E(pkgerr.ErrInvalidInput, "invalid body", err))
		return
	}
	item, err2 := h.svc.Update(r.Context(), id, &models.Book{Title: in.Title, Author: in.Author, ISBN: in.ISBN})
	if err2 != nil {
		writeErr(w, err2)
		return
	}
	response.JSON(w, http.StatusOK, item)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, appErr := parseID(mux.Vars(r)["id"])
	if appErr != nil {
		writeErr(w, appErr)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}

func parseID(s string) (uint, *pkgerr.AppError) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil || n == 0 {
		return 0, pkgerr.E(pkgerr.ErrInvalidInput, "invalid id", err)
	}
	return uint(n), nil
}

func writeErr(w http.ResponseWriter, e *pkgerr.AppError) {
	logger.Log.Warn("http error", zap.String("code", string(e.Code)), zap.String("msg", e.Message), zap.Error(e.Err))
	response.JSON(w, pkgerr.HTTPStatus(e.Code), map[string]string{"error": e.Message, "code": string(e.Code)})
}
