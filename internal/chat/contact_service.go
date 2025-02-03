package chat

import (
	"context"

	"karlota.aasumitro.id/internal/model/entity"
	sqlRepo "karlota.aasumitro.id/internal/repository/sql"
)

type contactService struct {
	sqlQueryRepository sqlRepo.ISQLQueryRepository[*entity.User, entity.User]
}

func (service contactService) List(
	ctx context.Context,
	requestBy string,
) ([]*entity.User, error) {
	limit, offset := -1, -1 // we will load all for now
	users, err := service.sqlQueryRepository.FindWithLimit(
		ctx, limit, offset, sqlRepo.NotEqual("email", requestBy))
	if err != nil {
		return nil, err
	}
	return users, nil
}

func NewContactService(
	sqlQueryRepository sqlRepo.ISQLQueryRepository[*entity.User, entity.User],
) IContactService {
	return &contactService{sqlQueryRepository: sqlQueryRepository}
}
