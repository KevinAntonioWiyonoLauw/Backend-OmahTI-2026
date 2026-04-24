package repository

import (
	"context"
	"database/sql"
	"fmt"

	"inventory-management/internal/models"
	"inventory-management/pkg/utils"

	"github.com/jmoiron/sqlx"
)

// Layer repository berubah ketika detail query SQL berubah,
// bukan saat kontrak HTTP endpoint berubah.
type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) (*models.Item, error)
	GetByID(ctx context.Context, id int64) (*models.Item, error)
	List(ctx context.Context) ([]models.Item, error)
	Update(ctx context.Context, item *models.Item) (*models.Item, error)
	Delete(ctx context.Context, id int64) error
}

type itemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(ctx context.Context, item *models.Item) (*models.Item, error) {
	const query = `
		INSERT INTO items (name, quantity, price)
		VALUES ($1, $2, $3)
		RETURNING id, name, quantity, price, created_at
	`

	created := models.Item{}
	if err := r.db.GetContext(ctx, &created, query, item.Name, item.Quantity, item.Price); err != nil {
		return nil, utils.NewInternalError("gagal menyimpan item", fmt.Errorf("insert item: %w", err))
	}

	return &created, nil
}

func (r *itemRepository) GetByID(ctx context.Context, id int64) (*models.Item, error) {
	const query = `
		SELECT id, name, quantity, price, created_at
		FROM items
		WHERE id = $1
	`

	item := models.Item{}
	if err := r.db.GetContext(ctx, &item, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewNotFoundError("item tidak ditemukan", err)
		}
		return nil, utils.NewInternalError("gagal mengambil item", fmt.Errorf("select item by id: %w", err))
	}

	return &item, nil
}

func (r *itemRepository) List(ctx context.Context) ([]models.Item, error) {
	const query = `
		SELECT id, name, quantity, price, created_at
		FROM items
		ORDER BY created_at DESC
	`

	items := []models.Item{}
	if err := r.db.SelectContext(ctx, &items, query); err != nil {
		return nil, utils.NewInternalError("gagal mengambil daftar item", fmt.Errorf("select items: %w", err))
	}

	return items, nil
}

func (r *itemRepository) Update(ctx context.Context, item *models.Item) (*models.Item, error) {
	const query = `
		UPDATE items
		SET name = $2, quantity = $3, price = $4
		WHERE id = $1
		RETURNING id, name, quantity, price, created_at
	`

	updated := models.Item{}
	if err := r.db.GetContext(ctx, &updated, query, item.ID, item.Name, item.Quantity, item.Price); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewNotFoundError("item tidak ditemukan", err)
		}
		return nil, utils.NewInternalError("gagal memperbarui item", fmt.Errorf("update item: %w", err))
	}

	return &updated, nil
}

func (r *itemRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM items
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return utils.NewInternalError("gagal menghapus item", fmt.Errorf("delete item: %w", err))
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return utils.NewInternalError("gagal membaca hasil penghapusan item", fmt.Errorf("rows affected: %w", err))
	}

	if rows == 0 {
		return utils.NewNotFoundError("item tidak ditemukan", nil)
	}

	return nil
}
