package injector

import (
	services "benchmark/internal/config/services"
	services12 "benchmark/internal/cuber/services"
	services7 "benchmark/internal/database/services"
	services1 "benchmark/internal/encoder/services"
	services2 "benchmark/internal/error/services"
	services3 "benchmark/internal/filesystem/services"
	services11 "benchmark/internal/marshalling/services"
	services4 "benchmark/internal/pipeline/services"
	services10 "benchmark/internal/random/services"
	services5 "benchmark/internal/schema/services"
	services9 "benchmark/internal/slurm/services"
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
	do.Provide(injector, services9.NewSlurmService)
	do.Provide(injector, services10.NewRandomService)
	do.Provide(injector, services11.NewMarshallingService)
	do.Provide(injector, services12.NewCuberService)
	return injector
}
