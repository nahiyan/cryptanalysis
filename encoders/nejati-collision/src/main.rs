mod differential;
mod path;
mod rules;
mod sha256;

use std::env;
use std::io::{self, BufWriter};
use std::path::PathBuf;
use std::process::ExitCode;

use path::DifferentialPath;
use rules::Rules;

struct Options {
    rounds: Option<usize>,
    xor: bool,
    diff_desc: bool,
    free_start: bool,
    random_input_difference: bool,
    differential_path: Option<PathBuf>,
    adder: String,
}

impl Default for Options {
    fn default() -> Self {
        Self {
            rounds: None,
            xor: false,
            diff_desc: false,
            free_start: false,
            random_input_difference: false,
            differential_path: None,
            adder: "espresso".into(),
        }
    }
}

fn main() -> ExitCode {
    match run() {
        Ok(()) => ExitCode::SUCCESS,
        Err(error) => {
            eprintln!("error: {error}");
            ExitCode::FAILURE
        }
    }
}

fn run() -> Result<(), String> {
    let options = parse_args()?;
    let rounds = options
        .rounds
        .ok_or_else(|| "number of rounds is required; use -r or --rounds".to_string())?;
    if !(16..=64).contains(&rounds) {
        return Err("SHA-256 collision rounds must be between 16 and 64".into());
    }
    if !options.diff_desc {
        return Err("--diff_desc is required for collision encoding".into());
    }
    if options.adder != "espresso" {
        return Err(format!(
            "adder '{}' is not supported by the collision encoder; use espresso",
            options.adder
        ));
    }
    if options.random_input_difference {
        return Err("--rand_input_diff was never implemented by the legacy encoder; provide a differential path instead".into());
    }
    let path = options
        .differential_path
        .ok_or_else(|| "--diff_const_file is required with --diff_desc".to_string())?;
    let path = DifferentialPath::read(&path, rounds)?;
    let rules = Rules::embedded()?;
    let formula = differential::encode(rounds, options.xor, options.free_start, &path, &rules);
    let stdout = io::stdout();
    formula
        .write_dimacs_with_order(BufWriter::new(stdout.lock()), rounds)
        .map_err(|error| format!("failed to write DIMACS: {error}"))
}

fn parse_args() -> Result<Options, String> {
    let mut options = Options::default();
    let mut args = env::args().skip(1);
    while let Some(argument) = args.next() {
        let value = |args: &mut std::iter::Skip<env::Args>, flag: &str| {
            args.next()
                .ok_or_else(|| format!("{flag} requires a value"))
        };
        match argument.as_str() {
            "-h" | "--help" => {
                print_help();
                std::process::exit(0);
            }
            "--xor" => options.xor = true,
            "--diff_desc" => options.diff_desc = true,
            "--free_start" => options.free_start = true,
            "--rand_input_diff" => options.random_input_difference = true,
            "-r" | "--rounds" => {
                let raw = value(&mut args, &argument)?;
                options.rounds = Some(raw.parse().map_err(|_| {
                    format!("{argument} expects a non-negative integer, got '{raw}'")
                })?);
            }
            "-A" | "--adder_type" => options.adder = value(&mut args, &argument)?,
            "-d" | "--diff_const_file" => {
                options.differential_path = Some(value(&mut args, &argument)?.into())
            }
            _ => return Err(format!("unknown argument '{argument}'; use --help")),
        }
    }
    Ok(options)
}

fn print_help() {
    println!(
        "Nejati SHA-256 collision SAT encoder\n\
         \nUSAGE:\n  nejati_collision_encoder --rounds <16..64> --diff_desc --diff_const_file <path> [OPTIONS]\n\
         \nOPTIONS:\n  -r, --rounds <N>                 Number of SHA-256 rounds\n  -d, --diff_const_file <PATH>     Differential-path constraints\n      --diff_desc                  Enable differential encoding (required)\n      --free_start                 Leave both chaining values unconstrained\n      --xor                        Emit native XOR clauses\n  -A, --adder_type <espresso>      Multi-operand adder encoding\n  -h, --help                       Print this help"
    );
}
