package usecase

import (
	"context"
	"go-rest-boilerplate/domain/domainerror"
	"go-rest-boilerplate/domain/entity"
	"go-rest-boilerplate/internal"
)

func (u infoUsecase) FindAllInfo(ctx context.Context) (out entity.Infos, err error) {
	errLocation := "[usecase info/find_all_info FindAllInfo] "
	defer internal.WrapErr(errLocation, &err)

	infos, err := u.infoRepository.FindAll(ctx, nil)
	if err != nil {
		return entity.Infos{}, new(domainerror.InternalError).Wrap(err)
	}
	return *infos, nil
}
