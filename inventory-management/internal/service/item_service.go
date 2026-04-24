package service

import (
	"context"
	"strings"

	"inventory-management/internal/models"
	"inventory-management/internal/repository"
	"inventory-management/pkg/utils"
)

// Layer service berubah ketika aturan bisnis berubah,
// bukan saat framework HTTP atau SQL query detail berubah.
type ItemService interface {
	Create(ctx context.Context, input models.CreateItemInput) (*models.Item, error)
	GetByID(ctx context.Context, id int64) (*models.Item, error)
	List(ctx context.Context) ([]models.Item, error)
	Update(ctx context.Context, id int64, input models.UpdateItemInput) (*models.Item, error)
	Delete(ctx context.Context, id int64) error
}

type itemService struct {
	repo repository.ItemRepository
}

func NewItemService(repo repository.ItemRepository) ItemService {
	return &itemService{repo: repo}
}

func (s *itemService) Create(ctx context.Context, input models.CreateItemInput) (*models.Item, error) {
	if err := validateCreateInput(input); err != nil {
		return nil, err
	}

	item := &models.Item{
		Name:     strings.TrimSpace(input.Name),
		Quantity: input.Quantity,
		Price:    input.Price,
	}

	return s.repo.Create(ctx, item)
}

func (s *itemService) GetByID(ctx context.Context, id int64) (*models.Item, error) {
	if id <= 0 {
		return nil, utils.NewInvalidInputError("id harus lebih besar dari 0", nil)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *itemService) List(ctx context.Context) ([]models.Item, error) {
	return s.repo.List(ctx)
}

func (s *itemService) Update(ctx context.Context, id int64, input models.UpdateItemInput) (*models.Item, error) {
	if id <= 0 {
		return nil, utils.NewInvalidInputError("id harus lebih besar dari 0", nil)
	}

	if err := validateUpdateInput(input); err != nil {
		return nil, err
	}

	item := &models.Item{
		ID:       id,
		Name:     strings.TrimSpace(input.Name),
		Quantity: input.Quantity,
		Price:    input.Price,
	}

	return s.repo.Update(ctx, item)
}

func (s *itemService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return utils.NewInvalidInputError("id harus lebih besar dari 0", nil)
	}

	return s.repo.Delete(ctx, id)
}

func validateCreateInput(input models.CreateItemInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return utils.NewInvalidInputError("nama item wajib diisi", nil)
	}
	if input.Quantity < 0 {
		return utils.NewInvalidInputError("quantity tidak boleh negatif", nil)
	}
	if input.Price < 0 {
		return utils.NewInvalidInputError("price tidak boleh negatif", nil)
	}
	return nil
}

func validateUpdateInput(input models.UpdateItemInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return utils.NewInvalidInputError("nama item wajib diisi", nil)
	}
	if input.Quantity < 0 {
		return utils.NewInvalidInputError("quantity tidak boleh negatif", nil)
	}
	if input.Price < 0 {
		return utils.NewInvalidInputError("price tidak boleh negatif", nil)
	}
	return nil
}
