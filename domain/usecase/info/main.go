package usecase

import (
	"context"
	"go-rest-boilerplate/domain/entity"
	infoRepository "go-rest-boilerplate/domain/repository/info"
)

type InfoUsecase interface {
	FindAllInfo(ctx context.Context) (out entity.Infos, err error)
}

type infoUsecase struct {
	infoRepository infoRepository.InfoRepository
}

func NewInfoUsecase(
	infoRepository infoRepository.InfoRepository,
) InfoUsecase {
	return &infoUsecase{
		infoRepository: infoRepository,
	}
}
