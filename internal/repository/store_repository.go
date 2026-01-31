package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/abu-umair/be-lms-go/internal/entity"
	"github.com/abu-umair/be-lms-go/pkg/database"
)

type IStoreRepository interface {
	WithTransaction(tx *sql.Tx) IStoreRepository
	CreateNewStore(ctx context.Context, store *entity.Store) error
	GetStoreById(ctx context.Context, storeId string) (*entity.Store, error)
	UpdateStore(ctx context.Context, store *entity.Store) error
	DeleteStore(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error
}

type storeRepository struct {
	db database.DatabaseQuery
}

func (ss *storeRepository) WithTransaction(tx *sql.Tx) IStoreRepository {
	return &storeRepository{
		db: tx,
	}
}

func (sr *storeRepository) CreateNewStore(ctx context.Context, store *entity.Store) error {
	_, err := sr.db.ExecContext(
		ctx,
		`INSERT INTO stores (id, name, image_file_name, address, created_at, created_by, updated_at, updated_by, deleted_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,

		store.Id,
		store.Name,
		store.ImageFileName,
		store.Address,
		store.CreatedAt,
		store.CreatedBy,
		store.UpdatedAt,
		store.UpdatedBy,
		store.DeletedBy,
	)

	if err != nil {
		return err
	}

	return nil
}

func (sr *storeRepository) GetStoreById(ctx context.Context, storeId string) (*entity.Store, error) {
	var storeEntity entity.Store

	row := sr.db.QueryRowContext(
		ctx,
		"SELECT id, name, address, image_file_name FROM stores WHERE id = $1 AND deleted_at IS NULL",
		storeId,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&storeEntity.Id,
		&storeEntity.Name,
		&storeEntity.Address,
		&storeEntity.ImageFileName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &storeEntity, nil
}

func (sr *storeRepository) UpdateStore(ctx context.Context, store *entity.Store) error {
	_, err := sr.db.ExecContext(
		ctx, "UPDATE stores SET name = $1, address = $2, image_file_name = $3, updated_at = $4, updated_by = $5 WHERE id = $6",
		store.Name,
		store.Address,
		store.ImageFileName,
		store.UpdatedAt,
		store.UpdatedBy,
		store.Id,
	)

	if err != nil {
		return nil
	}

	return nil
}

func (sr *storeRepository) DeleteStore(ctx context.Context, id string, deletedAt time.Time, deletedBy string) error {
	_, err := sr.db.ExecContext(
		ctx, "UPDATE stores SET deleted_at = $1, deleted_by = $2 WHERE id = $3",
		deletedAt,
		deletedBy,
		id,
	)

	if err != nil {
		return nil
	}

	return nil
}

func NewStoreRepository(db database.DatabaseQuery) IStoreRepository {
	return &storeRepository{db: db}
}
