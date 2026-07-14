package repository

import (
	"context"

	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
	"gorm.io/gorm"
)

type IDealershipRepository interface {
	List(ctx context.Context) ([]model.Dealership, error)
	Exists(ctx context.Context, dealershipID string) (bool, error)
}

type dealershipRepository struct{ db *gorm.DB }

func NewDealershipRepository(db *gorm.DB) IDealershipRepository {
	return &dealershipRepository{db: db}
}

func (r *dealershipRepository) List(ctx context.Context) ([]model.Dealership, error) {
	var rows []model.Dealership
	err := r.db.WithContext(ctx).Order("name ASC").Find(&rows).Error
	return rows, err
}

func (r *dealershipRepository) Exists(ctx context.Context, dealershipID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Dealership{}).Where("id = ?", dealershipID).Count(&count).Error
	return count > 0, err
}
