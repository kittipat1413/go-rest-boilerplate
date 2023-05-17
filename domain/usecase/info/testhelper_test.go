package usecase

import (
	"testing"

	infoRepo "go-rest-boilerplate/domain/repository/info"
	infoRepoMocks "go-rest-boilerplate/domain/repository/info/mocks"

	"github.com/golang/mock/gomock"
)

type mockDeps struct {
	CreateInfoRepository func(ctrl *gomock.Controller) infoRepo.InfoRepository
}

var defaultMockRepos = mockDeps{
	CreateInfoRepository: func(ctrl *gomock.Controller) infoRepo.InfoRepository {
		return infoRepoMocks.NewMockInfoRepository(ctrl)
	},
}

type helper struct {
	usecase InfoUsecase
	ctrl    *gomock.Controller
	done    func()
}

func initTest(t *testing.T, m *mockDeps) *helper {
	ctrl := gomock.NewController(t)
	if m.CreateInfoRepository == nil {
		m.CreateInfoRepository = defaultMockRepos.CreateInfoRepository
	}

	return &helper{
		usecase: NewInfoUsecase(
			m.CreateInfoRepository(ctrl),
		),
		ctrl: ctrl,
		done: func() {
			ctrl.Finish()
		},
	}
}
