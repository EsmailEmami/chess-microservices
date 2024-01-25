package service

import (
	"context"

	"github.com/esmailemami/chess/shared/database/psql"
	"github.com/esmailemami/chess/shared/errs"
	"github.com/google/uuid"
)

type BaseService[T any] struct {
}

func (b *BaseService[T]) Get(ctx context.Context, id uuid.UUID) (*T, error) {
	db := psql.DBContext(ctx)

	var model T

	if err := db.Find(&model, "id=?", id).Error; err != nil {
		return nil, errs.NotFoundErr()
	}

	return &model, nil
}

func (b *BaseService[T]) Create(ctx context.Context, model *T) error {
	db := psql.DBContext(ctx)

	if err := db.Create(model).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (b *BaseService[T]) Update(ctx context.Context, model *T) error {
	db := psql.DBContext(ctx)

	if err := db.Save(model).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}

func (b *BaseService[T]) Delete(ctx context.Context, model *T) error {
	db := psql.DBContext(ctx)

	if err := db.Delete(model).Error; err != nil {
		return errs.InternalServerErr().WithError(err)
	}

	return nil
}
