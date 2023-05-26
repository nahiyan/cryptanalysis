package injector

import (
	services23 "cryptanalysis/internal/combined_logs/services"
	services15 "cryptanalysis/internal/command/services"
	services "cryptanalysis/internal/config/services"
	services14 "cryptanalysis/internal/cube_selector/services"
	services11 "cryptanalysis/internal/cuber/services"
	services12 "cryptanalysis/internal/cubeset/services"
	services1 "cryptanalysis/internal/encoder/services"
	services13 "cryptanalysis/internal/encoding/services"
	services2 "cryptanalysis/internal/error/services"
	services3 "cryptanalysis/internal/filesystem/services"
	services22 "cryptanalysis/internal/log/services"
	services10 "cryptanalysis/internal/marshalling/services"
	services19 "cryptanalysis/internal/md4/services"
	services20 "cryptanalysis/internal/md5/services"
	services4 "cryptanalysis/internal/pipeline/services"
	services9 "cryptanalysis/internal/random/services"
	services5 "cryptanalysis/internal/schema/services"
	services21 "cryptanalysis/internal/sha256/services"
	services18 "cryptanalysis/internal/simplification/services"
	services17 "cryptanalysis/internal/simplifier/services"
	services8 "cryptanalysis/internal/slurm/services"
	services7 "cryptanalysis/internal/solution/services"
	services6 "cryptanalysis/internal/solver/services"
	services16 "cryptanalysis/internal/summarizer/services"
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
