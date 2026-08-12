mod md4;
mod sha1;
mod sha256;

use std::env;
use std::io::{self, BufWriter};
use std::process::ExitCode;
use std::time::{SystemTime, UNIX_EPOCH};

use md4::Md4Encoding;
use sha1::Sha1Encoding;
use sha256::Sha256Encoding;

#[derive(Debug)]
struct Options {
    rounds: Option<usize>,
    function: String,
    analysis: String,
    adder: String,
    target: String,
    use_xor: bool,
    print_target: bool,
    fixed_bits: usize,
    dobbertin: bool,
    relaxed_bits: usize,
}

impl Default for Options {
    fn default() -> Self {
        Self {
            rounds: None,
            function: "sha1".into(),
            analysis: "preimage".into(),
            adder: "espresso".into(),
            target: "random".into(),
            use_xor: false,
            print_target: false,
            fixed_bits: 0,
            dobbertin: false,
            relaxed_bits: 32,
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
    if options.analysis != "preimage" {
        return Err("the Rust encoder currently supports preimage analysis only".into());
    }
    if options.adder != "espresso" {
        return Err(format!(
            "adder '{}' has not been migrated yet; use espresso",
            options.adder
        ));
    }
    if options.fixed_bits > 512 {
        return Err("fixed input bits must be between 0 and 512".into());
    }
    if !(1..=32).contains(&options.relaxed_bits) {
        return Err("Dobbertin relaxed bits must be between 1 and 32".into());
    }

    match options.function.as_str() {
        "md4" => run_md4(options, rounds),
        "sha1" => run_sha1(options, rounds),
        "sha256" => run_sha256(options, rounds),
        function => Err(format!("unsupported hash function '{function}'")),
    }
}

fn run_md4(options: Options, rounds: usize) -> Result<(), String> {
    if !(1..=48).contains(&rounds) {
        return Err("MD4 rounds must be between 1 and 48".into());
    }
    if options.dobbertin && rounds < 27 {
        return Err("Dobbertin constraints require at least 27 MD4 rounds".into());
    }
    let message = (options.target == "random").then(random_message);
    let target = match &message {
        Some(message) => md4::compression_state(message, rounds),
        None => parse_target::<4>(&options.target)?,
    };
    if maybe_print_target(&options, message.as_ref(), &target)? {
        return Ok(());
    }
    require_message_for_fix(&options, message.as_ref())?;

    let mut encoding = Md4Encoding::build(
        rounds,
        options.use_xor,
        options.dobbertin,
        options.relaxed_bits,
    );
    encoding.fix_output(target);
    if let Some(message) = &message {
        encoding.fix_message_prefix(message, options.fixed_bits);
    }
    write_formula(&encoding.formula)
}

fn run_sha1(options: Options, rounds: usize) -> Result<(), String> {
    reject_dobbertin(&options, "SHA-1")?;
    if !(16..=80).contains(&rounds) {
        return Err("SHA-1 rounds must be between 16 and 80".into());
    }
    let message = (options.target == "random").then(random_message);
    let target = match &message {
        Some(message) => sha1::compression(message, rounds),
        None => parse_target::<5>(&options.target)?,
    };
    if maybe_print_target(&options, message.as_ref(), &target)? {
        return Ok(());
    }
    require_message_for_fix(&options, message.as_ref())?;

    let mut encoding = Sha1Encoding::build(rounds, options.use_xor);
    encoding.fix_output(target);
    if let Some(message) = &message {
        encoding.fix_message_prefix(message, options.fixed_bits);
    }
    write_formula(&encoding.formula)
}

fn run_sha256(options: Options, rounds: usize) -> Result<(), String> {
    reject_dobbertin(&options, "SHA-256")?;
    if !(16..=64).contains(&rounds) {
        return Err("SHA-256 rounds must be between 16 and 64".into());
    }
    let message = (options.target == "random").then(random_message);
    let target = match &message {
        Some(message) => sha256::compression(message, rounds),
        None => parse_target::<8>(&options.target)?,
    };
    if maybe_print_target(&options, message.as_ref(), &target)? {
        return Ok(());
    }
    require_message_for_fix(&options, message.as_ref())?;

    let mut encoding = Sha256Encoding::build(rounds, options.use_xor);
    encoding.fix_output(target);
    if let Some(message) = &message {
        encoding.fix_message_prefix(message, options.fixed_bits);
    }
    write_formula(&encoding.formula)
}

fn reject_dobbertin(options: &Options, function: &str) -> Result<(), String> {
    if options.dobbertin || options.relaxed_bits != 32 {
        Err(format!(
            "Dobbertin options apply only to MD4, not {function}"
        ))
    } else {
        Ok(())
    }
}

fn require_message_for_fix(options: &Options, message: Option<&[u32; 16]>) -> Result<(), String> {
    if options.fixed_bits > 0 && message.is_none() {
        Err("--fix requires --target random; an explicit target has no source message".into())
    } else {
        Ok(())
    }
}

fn maybe_print_target<const N: usize>(
    options: &Options,
    message: Option<&[u32; 16]>,
    target: &[u32; N],
) -> Result<bool, String> {
    if !options.print_target {
        return Ok(false);
    }
    let message = message.ok_or_else(|| {
        "--print_target requires --target random so a message is available".to_string()
    })?;
    println!("{}", format_words(message));
    println!("{}", format_words(target));
    Ok(true)
}

fn format_words(words: &[u32]) -> String {
    words
        .iter()
        .map(|word| format!("{word:08x}"))
        .collect::<Vec<_>>()
        .join(" ")
}

fn write_formula(formula: &nejati_sat::Formula) -> Result<(), String> {
    let stdout = io::stdout();
    let mut output = BufWriter::new(stdout.lock());
    formula
        .write_dimacs(&mut output)
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
            "--xor" => options.use_xor = true,
            "--print_target" => options.print_target = true,
            "--dobbertin" => options.dobbertin = true,
            "-r" | "--rounds" => {
                options.rounds = Some(parse_number(&value(&mut args, &argument)?, &argument)?)
            }
            "-f" | "--function" => options.function = value(&mut args, &argument)?,
            "-a" | "--analysis" => options.analysis = value(&mut args, &argument)?,
            "-A" | "--adder_type" => options.adder = value(&mut args, &argument)?,
            "-t" | "--target" => options.target = value(&mut args, &argument)?,
            "-F" | "--fix" => {
                options.fixed_bits = parse_number(&value(&mut args, &argument)?, &argument)?
            }
            "-b" | "--bits" => {
                options.relaxed_bits = parse_number(&value(&mut args, &argument)?, &argument)?
            }
            _ => return Err(format!("unknown argument '{argument}'; use --help")),
        }
    }
    Ok(options)
}

fn parse_number(value: &str, flag: &str) -> Result<usize, String> {
    value
        .parse()
        .map_err(|_| format!("{flag} expects a non-negative integer, got '{value}'"))
}

fn parse_target<const N: usize>(value: &str) -> Result<[u32; N], String> {
    if value.len() != N * 8 || !value.bytes().all(|byte| byte.is_ascii_hexdigit()) {
        return Err(format!(
            "the target must contain exactly {} hexadecimal characters",
            N * 8
        ));
    }
    let mut target = [0u32; N];
    for (index, word) in target.iter_mut().enumerate() {
        *word = u32::from_str_radix(&value[index * 8..index * 8 + 8], 16)
            .map_err(|error| format!("invalid target: {error}"))?;
    }
    Ok(target)
}

fn random_message() -> [u32; 16] {
    let mut state = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap_or_default()
        .as_nanos() as u64
        ^ u64::from(std::process::id());
    std::array::from_fn(|_| {
        state ^= state << 13;
        state ^= state >> 7;
        state ^= state << 17;
        state as u32
    })
}

fn print_help() {
    println!(
        "nejati_preimage_encoder_rs\n\
         \n\
         Usage: nejati_preimage_encoder_rs -r <rounds> [options]\n\
         \n\
         Options:\n\
           -r, --rounds <count>       Reduced-round count\n\
           -f, --function <name>      md4, sha1, or sha256\n\
           -a, --analysis <preimage>  Analysis type\n\
           -A, --adder_type <type>    Adder type (espresso currently migrated)\n\
           -t, --target <value>       random or the complete hexadecimal target\n\
           -F, --fix <0..512>         Fix a random message prefix\n\
               --xor                  Emit XOR clauses where supported\n\
               --dobbertin            Add MD4 Dobbertin constraints\n\
           -b, --bits <1..32>         Bits retained in relaxed Dobbertin q[16]\n\
               --print_target         Print random message and target, then exit\n\
           -h, --help                 Show this help"
    );
}
