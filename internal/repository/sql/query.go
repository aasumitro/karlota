package sql

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type sqlQueryRepository[M GormModel[R], R any] struct {
	orm *gorm.DB
}

func (repository sqlQueryRepository[M, R]) Count(
	ctx context.Context,
	specifications ...Specification,
) (i int64, err error) {
	model := new(M)
	err = repository.getPreWarmDbForSelect(ctx, specifications...).
		Model(model).Count(&i).Error
	return
}

func (repository sqlQueryRepository[M, R]) FindAll(
	ctx context.Context,
) ([]*R, error) {
	return repository.FindWithLimit(ctx, -1, -1)
}

func (repository sqlQueryRepository[M, R]) FindWithLimit(
	ctx context.Context,
	limit int, offset int,
	specifications ...Specification,
) ([]*R, error) {
	var models []M
	preWarm := repository.getPreWarmDbForSelect(ctx, specifications...)
	err := preWarm.Limit(limit).Offset(offset).Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*R, 0, len(models))
	for _, row := range models {
		result = append(result, row.ToResponse())
	}
	return result, nil
}

func (repository sqlQueryRepository[M, R]) FindByID(
	ctx context.Context,
	id uint,
) (*R, error) {
	var model M
	err := repository.orm.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return new(R), err
	}
	return model.ToResponse(), nil
}

func (repository sqlQueryRepository[M, R]) getPreWarmDbForSelect(
	ctx context.Context,
	specification ...Specification,
) *gorm.DB {
	preWarm := repository.orm.WithContext(ctx)
	for _, s := range specification {
		if joinType, joinCondition := s.GetJoin(); joinCondition != "" {
			preWarm = preWarm.Joins(fmt.Sprintf("%s %s", joinType, joinCondition))
		}
		preWarm = preWarm.Where(s.GetQuery(), s.GetValues()...)
	}
	return preWarm
}

func NewSQLQueryRepository[M GormModel[R], R any](
	orm *gorm.DB,
) ISQLQueryRepository[M, R] {
	return &sqlQueryRepository[M, R]{orm: orm}
}
