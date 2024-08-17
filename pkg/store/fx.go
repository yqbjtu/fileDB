package store

import (
	"go.uber.org/fx"
)

// Module -
var Module = fx.Options(
	fx.Provide(NewPgDB),
	fx.Provide(NewCellHistoryStore),
	fx.Provide(NewCellStatusStore),
	fx.Provide(NewCellGisMetaStore),
	fx.Provide(NewCellCompileQueueStore),
)
