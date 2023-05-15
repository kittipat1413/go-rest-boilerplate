package repository

import (
	"context"
	"time"

	"go-rest-boilerplate/data"
	"go-rest-boilerplate/domain/entity"

	"github.com/jinzhu/copier"
)

const tableName = "info"

type InfoRepository interface {
	FindAll(ctx context.Context, scope data.Scope) (out *entity.Infos, err error)
	FindByKey(ctx context.Context, scope data.Scope, key string) (out *entity.Info, err error)
}

type infoRepository struct{}

func NewInfoRepository() InfoRepository {
	return &infoRepository{}
}

type InfosModel []InfoModel
type InfoModel struct {
	Key       string    `db:"key"`
	Value     string    `db:"value"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *InfoModel) ToEntity() (*entity.Info, error) {
	var info entity.Info
	err := copier.Copy(&info, m)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (o *InfosModel) ToEntities() (*entity.Infos, error) {
	var infos entity.Infos
	err := copier.Copy(&infos, o)
	if err != nil {
		return nil, err
	}
	return &infos, nil
}
