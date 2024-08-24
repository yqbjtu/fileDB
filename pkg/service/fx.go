package service

import "go.uber.org/fx"

// Module -
var Module = fx.Options(
	fx.Provide(NewCellHistoryService),
	fx.Provide(NewCellStatusService),
	fx.Provide(NewCellCompileQueueService),
	fx.Provide(NewCellCvsService),
)
