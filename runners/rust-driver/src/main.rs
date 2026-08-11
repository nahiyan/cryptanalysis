use clap::{Parser, Subcommand};
use solutions::parse_sols;
use std::{fs, io::Write, process::Command};
mod solutions;

#[derive(Subcommand)]
enum Commands {
  Cube {
    #[clap(short, long, default_value = "encodings")]
    input_dir: String,
    #[clap(short, long, default_value = "cubes")]
    output_dir: String,
    #[clap(short = 'n', long, default_value = "-1")]
    cutoff_vars: i32,
    #[clap(short, long, default_value = "5")]
    decrement: i32,
    #[clap(short, long, default_value = "1")]
    min_cubes: i32,
  },
  Solutions {
    #[clap(short, long, default_value = "logs")]
    input_dir: String,
  },
}

#[derive(Parser)]
#[command(version, about, long_about = None)]
struct Cli {
  #[command(subcommand)]
  command: Commands,
}

fn main() {
  let args = Cli::parse();

  match args.command {
    Commands::Cube {
      input_dir,
      output_dir,
      cutoff_vars,
      min_cubes,
      decrement,
    } => {
      // TODO: Implement the basic features
      let cubes_script = include_bytes!("../scripts/gen_cubes.sh");
      let mut script_file = fs::File::create("cubes.sh").unwrap();
      script_file
        .write(cubes_script)
        .expect("Failed to write script.");
      assert!(cutoff_vars >= 1);

      fs::create_dir(&output_dir).unwrap();

      for file in fs::read_dir(input_dir).unwrap() {
        Command::new("sh")
          .args([
            "gen_cubes.sh",
            &file.unwrap().path().to_str().unwrap(),
            &output_dir,
            &min_cubes.to_string(),
            &decrement.to_string(),
          ])
          .spawn()
          .expect("Failed to generate cubes");
      }
    }
    Commands::Solutions { input_dir } => parse_sols(input_dir),
  }
}
