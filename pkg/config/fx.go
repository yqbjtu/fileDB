package config

import "go.uber.org/fx"

// Module -
var Module = fx.Options(
	fx.Provide(GetConfig),
)
