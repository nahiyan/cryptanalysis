package injector

import (
	services "benchmark/internal/config/services"
	services7 "benchmark/internal/database/services"
	services1 "benchmark/internal/encoder/services"
	services2 "benchmark/internal/error/services"
	services3 "benchmark/internal/filesystem/services"
	services4 "benchmark/internal/pipeline/services"
	services5 "benchmark/internal/schema/services"
	services8 "benchmark/internal/solution/services"
	services6 "benchmark/internal/solver/services"
	do "github.com/samber/do"
)

func New() *do.Injector {
	injector := do.New()
	do.Provide(injector, services.NewConfigService)
	do.Provide(injector, services1.NewEncoderService)
	do.Provide(injector, services2.NewErrorService)
	do.Provide(injector, services3.NewFilesystemService)
	do.Provide(injector, services4.NewPipelineService)
	do.Provide(injector, services5.NewSchemaService)
	do.Provide(injector, services6.NewSolverService)
	do.Provide(injector, services7.NewDatabaseService)
	do.Provide(injector, services8.NewSolutionService)
	return injector
}
