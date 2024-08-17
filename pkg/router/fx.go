package router

import (
	"go.uber.org/fx"
)

// Module -
var Module = fx.Options(
	fx.Provide(NewRouter),
	fx.Provide(NewHTTPServer),
)
