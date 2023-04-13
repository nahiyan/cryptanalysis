package injector

import (
	services23 "benchmark/internal/combined_logs/services"
	services15 "benchmark/internal/command/services"
	services "benchmark/internal/config/services"
	services14 "benchmark/internal/cube_selector/services"
	services11 "benchmark/internal/cuber/services"
	services12 "benchmark/internal/cubeset/services"
	services1 "benchmark/internal/encoder/services"
	services13 "benchmark/internal/encoding/services"
	services2 "benchmark/internal/error/services"
	services3 "benchmark/internal/filesystem/services"
	services22 "benchmark/internal/log/services"
	services10 "benchmark/internal/marshalling/services"
	services19 "benchmark/internal/md4/services"
	services20 "benchmark/internal/md5/services"
	services4 "benchmark/internal/pipeline/services"
	services9 "benchmark/internal/random/services"
	services5 "benchmark/internal/schema/services"
	services21 "benchmark/internal/sha256/services"
	services18 "benchmark/internal/simplification/services"
	services17 "benchmark/internal/simplifier/services"
	services8 "benchmark/internal/slurm/services"
	services7 "benchmark/internal/solution/services"
	services6 "benchmark/internal/solver/services"
	services16 "benchmark/internal/summarizer/services"
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
	do.Provide(injector, services7.NewSolutionService)
	do.Provide(injector, services8.NewSlurmService)
	do.Provide(injector, services9.NewRandomService)
	do.Provide(injector, services10.NewMarshallingService)
	do.Provide(injector, services11.NewCuberService)
	do.Provide(injector, services12.NewCubesetService)
	do.Provide(injector, services13.NewEncodingService)
	do.Provide(injector, services14.NewCubeSelectorService)
	do.Provide(injector, services15.NewCommandService)
	do.Provide(injector, services16.NewSummarizerService)
	do.Provide(injector, services17.NewSimplifierService)
	do.Provide(injector, services18.NewSimplificationService)
	do.Provide(injector, services19.NewMd4Service)
	do.Provide(injector, services20.NewMd5Service)
	do.Provide(injector, services21.NewSha256Service)
	do.Provide(injector, services22.NewLogService)
	do.Provide(injector, services23.NewCombinedLogsService)
	return injector
}
