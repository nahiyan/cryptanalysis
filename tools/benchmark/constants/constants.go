package constants

const (
	ArgCounterChain  = "cc"
	ArgDotMatrix     = "dm"
	ArgEspresso      = "esp"
	ArgCryptoMiniSat = "cms"
	ArgKissat        = "ks"
	ArgCadical       = "cdc"
	ArgGlucoseSyrup  = "gs"
	ArgMapleSat      = "ms"
	ArgXnfSat        = "xnf"
	CryptoMiniSat    = "cryptominisat"
	Kissat           = "kissat"
	Cadical          = "cadical"
	Glucose          = "glucose"
	MapleSat         = "maplesat"
	XnfSat           = "xnfsat"

	BenchmarkLogFileName    = "benchmark.csv"
	VerificationLogFileName = "verification.csv"
	ValidResultsLogFileName = "valid_results.csv"
	JobsDirPath             = "./jobs/"
	ResultsDirPath          = "./results/"
	EncodingsDirPath        = ResultsDirPath + "encodings/"
	LogsDirPath             = ResultsDirPath + "logs/"
	SolutionsDirPath        = ResultsDirPath + "solutions/"

	ErrOneJobScheduleFailed = "failed to schedule one of the jobs"
)

const (
	Valid        = "valid"
	Invalid      = "invalid"
	Undetermined = "undetermined"
)
