package sql

import (
	"context"

	"gorm.io/gorm"
)

type sqlCommandRepository[M GormModel[R], R any] struct {
	orm *gorm.DB
}

func (repository sqlCommandRepository[M, R]) Insert(
	ctx context.Context,
	entity *R,
) error {
	var start M
	model := start.FromResponse(entity).(M)
	err := repository.orm.WithContext(ctx).Create(&model).Error
	if err != nil {
		return err
	}
	*entity = *model.ToResponse()
	return nil
}

func (repository sqlCommandRepository[M, R]) InsertMany(
	ctx context.Context,
	entities []*R,
) error {
	var start M
	models := make([]M, len(entities))
	for i, entity := range entities {
		models[i] = start.FromResponse(entity).(M)
	}
	err := repository.orm.WithContext(ctx).Create(&models).Error
	if err != nil {
		return err
	}
	for i, model := range models {
		*entities[i] = *model.ToResponse()
	}
	return nil
}

func (repository sqlCommandRepository[M, R]) DeleteWithSpec(
	ctx context.Context,
	entity *R,
	specifications ...Specification,
) error {
	var start M
	model := start.FromResponse(entity).(M)
	preWarm := repository.orm.WithContext(ctx)
	for _, s := range specifications {
		preWarm = preWarm.Where(s.GetQuery(), s.GetValues()...)
	}
	return preWarm.Delete(model).Error
}

func (repository sqlCommandRepository[M, R]) DeleteByID(
	ctx context.Context,
	id uint,
) error {
	var start M
	err := repository.orm.WithContext(ctx).Delete(&start, &id).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository sqlCommandRepository[M, R]) Update(
	ctx context.Context,
	entity *R,
) error {
	var start M
	model := start.FromResponse(entity).(M)
	err := repository.orm.WithContext(ctx).Save(&model).Error
	if err != nil {
		return err
	}
	*entity = *model.ToResponse()
	return nil
}

func NewSQLCommandRepository[M GormModel[R], R any](
	orm *gorm.DB,
) ISQLCommandRepository[M, R] {
	return &sqlCommandRepository[M, R]{orm: orm}
}
