package runtime

import (
	"users_api/src/infrastructure/datasource"
	"users_api/src/infrastructure/repository"
	"users_api/src/interface/web"
	"users_api/src/interface/web/controller"
	"users_api/src/usecase"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		datasource.LoadDBConfig,
		datasource.NewPgPool,
		repository.NewUserPg,
		usecase.NewUserSearchUsecase,
		controller.NewUserController,
		web.NewRouter,
	),
	fx.Invoke(Start),
)