package controller

import "go.uber.org/fx"

// Module -
var Module = fx.Options(
	fx.Provide(NewCvsController),
	fx.Provide(NewMiscController),
	fx.Provide(NewQueryController),
	fx.Provide(NewAdminController),
	fx.Provide(NewCompileQueueController),
)
