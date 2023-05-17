package usecase

import (
	"context"
	"testing"
	"time"

	"go-rest-boilerplate/data"
	infoRepo "go-rest-boilerplate/domain/repository/info"
	infoRepoMocks "go-rest-boilerplate/domain/repository/info/mocks"

	"go-rest-boilerplate/domain/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestFindAllInfo(t *testing.T) {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	ctx := data.NewContext(context.Background(), sqlxDB)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		mocks   mockDeps
		mocksDB func(mock sqlmock.Sqlmock)
		args    args
		want    entity.Infos
		wantErr bool
		error   error
	}{
		{
			name: "Success",
			mocks: mockDeps{
				CreateInfoRepository: func(ctrl *gomock.Controller) infoRepo.InfoRepository {
					mock := infoRepoMocks.NewMockInfoRepository(ctrl)
					mock.EXPECT().FindAll(ctx, gomock.Any()).
						Return(&entity.Infos{
							entity.Info{
								Key:       "mock_key1",
								Value:     "mock_val1",
								CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							},
							entity.Info{
								Key:       "mock_key2",
								Value:     "mock_val2",
								CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							},
						}, nil)
					return mock

				},
			},
			mocksDB: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				sqlMock.ExpectCommit()
			},
			args: args{
				ctx: ctx,
			},
			want: entity.Infos{
				entity.Info{
					Key:       "mock_key1",
					Value:     "mock_val1",
					CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				entity.Info{
					Key:       "mock_key2",
					Value:     "mock_val2",
					CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
			error:   nil,
		},
	}

	for _, test := range tests {
		test.mocksDB(sqlMock)
		t.Run(test.name, func(t *testing.T) {
			h := initTest(t, &test.mocks)
			defer h.done()

			got, err := h.usecase.FindAllInfo(test.args.ctx)
			if (err != nil) != test.wantErr {
				t.Errorf("Detail() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if test.wantErr {
				assert.Equal(t, err, test.error)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, got, test.want)
			}
		})
	}
}
