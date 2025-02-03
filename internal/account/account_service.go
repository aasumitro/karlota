package account

import (
	"context"

	"karlota.aasumitro.id/internal/common"
	"karlota.aasumitro.id/internal/model/entity"
	sqlRepo "karlota.aasumitro.id/internal/repository/sql"
)

type accountService struct {
	sqlQueryRepository   sqlRepo.ISQLQueryRepository[*entity.User, entity.User]
	sqlCommandRepository sqlRepo.ISQLCommandRepository[*entity.User, entity.User]
}

func (service *accountService) Profile(
	ctx context.Context,
	id uint,
) (interface{}, error) {
	return service.sqlQueryRepository.FindByID(ctx, id)
}

func (service *accountService) UpdateDisplayName(
	ctx context.Context,
	id uint, displayName string,
) (interface{}, error) {
	user, err := service.sqlQueryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.DisplayName = displayName
	if err := service.sqlCommandRepository.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (service *accountService) UpdatePassword(
	ctx context.Context,
	id uint, oldPwd, newPwd string,
) (interface{}, error) {
	user, err := service.sqlQueryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// validate old password
	if err := user.ValidatePassword(oldPwd); err != nil {
		return nil, common.ErrOldPasswordNotValid
	}
	// generate new password
	if err := user.MakeNewPassword(newPwd); err != nil {
		return nil, err
	}
	if err := service.sqlCommandRepository.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func NewAccountService(
	sqlQueryRepository sqlRepo.ISQLQueryRepository[*entity.User, entity.User],
	sqlCommandRepository sqlRepo.ISQLCommandRepository[*entity.User, entity.User],
) IAccountService {
	return &accountService{
		sqlCommandRepository: sqlCommandRepository,
		sqlQueryRepository:   sqlQueryRepository,
	}
}
