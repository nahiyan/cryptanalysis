package constants

const (
	ArgCounterChain            = "cc"
	ArgDotMatrix               = "dm"
	ArgCryptoMiniSat           = "cmc"
	ArgKissat                  = "ks"
	ArgCadical                 = "cdc"
	ArgGlucoseSyrup            = "gc"
	ArgMapleSat                = "ms"
	ArgXnfSat                  = "xnf"
	CRYPTOMINISAT              = "cryptominisat"
	KISSAT                     = "kissat"
	CADICAL                    = "cadical"
	GLUCOSE                    = "glucose"
	MAPLESAT                   = "maplesat"
	CRYPTOMINISAT_BIN_PATH     = "../../../sat-solvers/cryptominisat"
	KISSAT_BIN_PATH            = "../../../sat-solvers/kissat"
	CADICAL_BIN_PATH           = "../../../sat-solvers/cadical"
	GLUCOSE_BIN_PATH           = "../../../sat-solvers/glucose"
	MAPLESAT_BIN_PATH          = "../../../sat-solvers/maplesat"
	VERIFIER_BIN_PATH          = "../../encoders/saeed/crypto/verify-md4"
	SOLUTION_ANALYZER_BIN_PATH = "../solution_analyzer/target/release/solution_analyzer"

	BENCHMARK_LOG_FILE_NAME    = "benchmark.log"
	VERIFICATION_LOG_FILE_NAME = "verification.log"
	BASE_PATH                  = "../../"
	SOLUTIONS_DIR_PATH         = BASE_PATH + "solutions/saeed/"
	ENCODINGS_DIR_PATH         = BASE_PATH + "encodings/saeed/"
)
