package constants

const (
	ArgCounterChain         = "cc"
	ArgDotMatrix            = "dm"
	ArgCryptoMiniSat        = "cms"
	ArgKissat               = "ks"
	ArgCadical              = "cdc"
	ArgGlucoseSyrup         = "gs"
	ArgMapleSat             = "ms"
	ArgXnfSat               = "xnf"
	CryptoMiniSat           = "cryptominisat"
	Kissat                  = "kissat"
	Cadical                 = "cadical"
	Glucose                 = "glucose"
	MapleSat                = "maplesat"
	XnfSat                  = "xnfsat"
	CryptoMiniSatBinPath    = "../../../sat-solvers/cryptominisat"
	KissatBinPath           = "../../../sat-solvers/kissat"
	CadicalBinPath          = "../../../sat-solvers/cadical"
	GlucoseBinPath          = "../../../sat-solvers/glucose"
	MapleSatBinPath         = "../../../sat-solvers/maplesat"
	XnfSatBinPath           = "../../../sat-solvers/xnfsat"
	VerifierBinPath         = "../../encoders/saeed/crypto/verify-md4"
	SolutionAnalyzerBinPath = "../solution_analyzer/target/release/solution_analyzer"

	BenchmarkLogFileName    = "benchmark.log"
	VerificationLogFileName = "verification.log"
	ValidResultsLogFileName = "valid_results.log"
	ResultsDirPath          = "./results/"
	EncodingsDirPath        = ResultsDirPath + "encodings/"
	LogsDirPath             = ResultsDirPath + "logs/"
	SolutionsDirPath        = ResultsDirPath + "solutions/"
	EncoderPath             = "../../encoders/saeed/crypto/main"
)
