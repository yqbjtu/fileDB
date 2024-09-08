package schedule

import "go.uber.org/fx"

// Module -
var Module = fx.Options(
	fx.Provide(NewCellCompileScheduleService),
)
