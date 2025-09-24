package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"community-support-service/internal/database"
	"community-support-service/internal/models"
	"community-support-service/internal/repositories"
)

type categoryRepository struct {
	db *database.DB
}

func NewCategoryRepository(db *database.DB) repositories.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO categories (
			id, name, description, type, isActive, metadata,
			createdAt, updatedAt
		) VALUES (
			:id, :name, :description, :type, :isActive, :metadata,
			:createdAt, :updatedAt
		)`
	
	_, err := r.db.NamedExecContext(ctx, query, category)
	return err
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	query := `
		SELECT 
			id, name, description, type, isActive, metadata,
			createdAt, updatedAt
		FROM categories 
		WHERE id = $1`
	
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetAll(ctx context.Context, activeOnly bool) ([]*models.Category, error) {
	query := `
		SELECT 
			id, name, description, type, isActive, metadata,
			createdAt, updatedAt
		FROM categories`
	
	var args []interface{}
	if activeOnly {
		query += " WHERE isActive = true"
	}
	
	query += " ORDER BY name ASC"
	
	var categories []*models.Category
	err := r.db.SelectContext(ctx, &categories, query, args...)
	return categories, err
}

func (r *categoryRepository) GetByType(ctx context.Context, categoryType string, activeOnly bool) ([]*models.Category, error) {
	query := `
		SELECT 
			id, name, description, type, isActive, metadata,
			createdAt, updatedAt
		FROM categories 
		WHERE type = $1`
	
	args := []interface{}{categoryType}
	
	if activeOnly {
		query += " AND isActive = true"
	}
	
	query += " ORDER BY name ASC"
	
	var categories []*models.Category
	err := r.db.SelectContext(ctx, &categories, query, args...)
	return categories, err
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories SET 
			name = :name,
			description = :description,
			type = :type,
			isActive = :isActive,
			metadata = :metadata,
			updatedAt = :updatedAt
		WHERE id = :id`
	
	_, err := r.db.NamedExecContext(ctx, query, category)
	return err
}

func (r *categoryRepository) UpdateStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	query := `
		UPDATE categories SET 
			isActive = $1,
			updatedAt = NOW()
		WHERE id = $2`
	
	_, err := r.db.ExecContext(ctx, query, isActive, id)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}