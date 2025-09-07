package auth

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sudo-init-do/okies-backend/pkg/util"
)

type Handler struct {
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := util.DecodeJSON(r, &req); err != nil {
		util.JSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	hashed, _ := util.HashPassword(req.Password)

	var id int64
	err := h.DB.QueryRow(
		context.Background(),
		`INSERT INTO users (email, password, role, created_at, updated_at)
         VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`,
		req.Email, hashed, req.Role).Scan(&id)

	if err != nil {
		util.JSONError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	access, refresh, _ := util.GenerateTokens(id, req.Role)
	util.JSON(w, http.StatusCreated, map[string]any{
		"user_id": id,
		"access":  access,
		"refresh": refresh,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := util.DecodeJSON(r, &req); err != nil {
		util.JSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	var id int64
	var hashed, role string
	err := h.DB.QueryRow(
		context.Background(),
		`SELECT id, password, role FROM users WHERE email=$1`,
		req.Email).Scan(&id, &hashed, &role)

	if err == sql.ErrNoRows {
		util.JSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	} else if err != nil {
		util.JSONError(w, http.StatusInternalServerError, "query error")
		return
	}

	if !util.CheckPasswordHash(req.Password, hashed) {
		util.JSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	access, refresh, _ := util.GenerateTokens(id, role)
	util.JSON(w, http.StatusOK, map[string]any{
		"user_id": id,
		"access":  access,
		"refresh": refresh,
	})
}
