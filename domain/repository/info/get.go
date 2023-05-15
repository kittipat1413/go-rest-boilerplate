package repository

import (
	"context"

	"database/sql"

	"go-rest-boilerplate/data"
	"go-rest-boilerplate/domain/domainerror"
	"go-rest-boilerplate/domain/entity"
	"go-rest-boilerplate/internal"

	"github.com/Masterminds/squirrel"
)

func (r *infoRepository) FindAll(ctx context.Context, scope data.Scope) (out *entity.Infos, err error) {
	errLocation := "[domain repository/info GetAll] "
	defer internal.WrapErr(errLocation, &err)

	if scope == nil {
		scope, err = data.NewScope(ctx, nil)
		if err != nil {
			err = new(domainerror.DatabaseError).Wrap(err)
			return
		} else {
			defer scope.End(&err)
		}
	}

	sqlBuilder := data.CreateSqlBuilder().
		Select("*").
		From(tableName)

	selectSql, args, err := sqlBuilder.ToSql()
	if err != nil {
		err = new(domainerror.InternalError).Wrap(err)
		return
	}

	model := InfosModel{}
	err = scope.Select(&model, selectSql, args...)
	if err != nil && err != sql.ErrNoRows {
		err = new(domainerror.DatabaseError).Wrap(err)
		return
	}

	out, err = model.ToEntities()
	if err != nil {
		err = new(domainerror.InternalError).Wrap(err)
		return
	}
	return
}

func (b *infoRepository) FindByKey(ctx context.Context, scope data.Scope, key string) (out *entity.Info, err error) {
	errLocation := "[repository repository/info GetByKey] "
	defer internal.WrapErr(errLocation, &err)

	if scope == nil {
		scope, err = data.NewScope(ctx, nil)
		if err != nil {
			err = new(domainerror.DatabaseError).Wrap(err)
			return
		} else {
			defer scope.End(&err)
		}
	}

	selectSql, args, err := data.CreateSqlBuilder().
		Select("*").
		From(tableName).
		Where(squirrel.Eq{"key": key}).
		ToSql()
	if err != nil {
		err = new(domainerror.InternalError).Wrap(err)
		return
	}

	model := InfoModel{}
	err = scope.Get(&model, selectSql, args...)
	if err != nil && err != sql.ErrNoRows {
		err = new(domainerror.DatabaseError).Wrap(err)
		return
	} else if err == sql.ErrNoRows {
		err = new(domainerror.NotFoundError).Wrap(err)
		return
	}

	out, err = model.ToEntity()
	if err != nil {
		err = new(domainerror.InternalError).Wrap(err)
		return
	}
	return
}
