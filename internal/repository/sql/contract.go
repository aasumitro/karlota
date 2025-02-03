package sql

import "context"

type (
	GormModel[R any] interface {
		ToResponse() *R
		FromResponse(response *R) interface{}
	}

	ISQLCommandRepository[M, R any] interface {
		Insert(ctx context.Context, entity *R) error
		InsertMany(ctx context.Context, entities []*R) error
		DeleteWithSpec(ctx context.Context, entity *R, specifications ...Specification) error
		DeleteByID(ctx context.Context, id uint) error
		Update(ctx context.Context, entity *R) error
	}

	ISQLQueryRepository[M, R any] interface {
		Count(
			ctx context.Context,
			specifications ...Specification,
		) (i int64, err error)
		FindAll(
			ctx context.Context,
		) ([]*R, error)
		FindWithLimit(
			ctx context.Context,
			limit int, offset int,
			specifications ...Specification,
		) ([]*R, error)
		FindByID(
			ctx context.Context,
			id uint,
		) (*R, error)
	}
)
