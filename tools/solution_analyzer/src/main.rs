use core::panic;
use std::fs;

use clap::{arg, Arg, Command};

enum Base {
    Binary,
    Hexadecimal,
}

// Definition of the summarization context
struct SummarizationContext<'a> {
    variable_ranges: &'a mut Vec<(u32, u32)>,
    variable_bytes: &'a mut Vec<u8>,
    base: Base,
}

impl SummarizationContext<'_> {
    // Analyze the line
    fn analyze_line(&mut self, line: &str) {
        // Ignore headers
        if line.starts_with("s SATISFIABLE") || line.starts_with("SAT") {
            return;
        }

        let pieces = line.split_whitespace();
        for piece in pieces {
            // Ignore v and 0
            if piece.eq("v") || piece.eq("0") {
                continue;
            }

            let variable: u32 = i32::unsigned_abs(piece.parse().unwrap());
            let value: u8 = if piece.starts_with('-') { 0 } else { 1 };

            for (start, end) in self.variable_ranges.iter() {
                if variable >= *start && variable <= *end {
                    self.variable_bytes.push(value)
                }
            }
        }
    }

    // Generic analysis of the entire solution, invoking the line analyzer
    fn analyze(&mut self, content: String) {
        let mut line: String = String::from("");
        for c in content.chars() {
            // End of line
            if c == '\n' {
                self.analyze_line(line.as_str());
                line = String::from("")
            } else {
                line.push(c)
            }
        }
    }

    // Dump the summarization
    fn dump(&self) {
        // Output the analysis
        let mut i = 0;
        for (start, end) in self.variable_ranges.iter() {
            println!("Range: {} - {}", start, end);
            for (j, _) in (*start..*end + 1).enumerate() {
                print!("{}", self.variable_bytes[i]);
                i += 1;

                if (j + 1) % 8 == 0 {
                    println!()
                }
            }
            println!()
        }

        // TODO: Output in hex
        // println!();
        // for i in 0..self.variable_ranges.len() + 1 {
        //     println!("Range");
        //     let value = &self.variable_bytes[(8 * i)..(8 * i + 8)];
        //     println!("{}", hex::encode(value));
        // }
    }
}

// Definition of the normalization context
struct NormalizationContext<'a> {
    variables: &'a mut Vec<i64>,
}

impl NormalizationContext<'_> {
    // Analyze the line
    fn analyze_line(&mut self, line: &str) {
        // Ignore headers
        if line.starts_with("s SATISFIABLE") || line.starts_with("SAT") || line.starts_with("c") {
            return;
        }

        let pieces = line.split_whitespace();
        for piece in pieces {
            let variable: i64 = match piece.parse() {
                Ok(x) => x,
                Err(_) => 0,
            };

            if variable == 0 {
                continue;
            }

            self.variables.push(variable)
        }
    }

    // Generic analysis of the entire solution, invoking the line analyzer
    fn analyze(&mut self, content: String) {
        let mut line: String = String::from("");
        for c in content.chars() {
            // End of line
            if c == '\n' {
                self.analyze_line(line.as_str());
                line = String::from("")
            } else {
                line.push(c)
            }
        }
    }

    fn dump(self) {
        println!("SAT");
        for variable in self.variables {
            print!("{} ", variable)
        }
        print!("0 ");
        println!()
    }
}

fn main() {
    // Process the CLI arguments
    let matches = Command::new("Solution Analyzer")
        .version("1.0")
        .author("Nahiyan Alamgir")
        .about("Analyzes SAT Solver solution to normalize the result or represent it as binary.")
        .arg(Arg::new("solution_file").required(true))
        .subcommand(
            Command::new("summarize")
                .about(
                    "Summarizes the values of the variables. It dumps the values as binary by default.",
                )
                .arg(arg!(-v --variables <RANGES> "Ranges of variables, starting from 1, with each range separated by a comma. Example: 1-64,512-1024.").required(true))
                .arg(arg!(--bin "Print the values of the variables as binary."))
                .arg(arg!(--hex "Print the values of the variables as hexadecimal."))
        )
        .subcommand(
            Command::new("normalize").about("Normalizes the solution")
        )
        .get_matches();

    // Solution
    let solution_file = matches.get_one::<String>("solution_file").unwrap();
    let content = match fs::read_to_string(solution_file) {
        Err(_) => {
            panic!("Failed to open the solutions file.");
        }
        Ok(content) => content,
    };

    // Handle the commands
    if let Some(matches_) = matches.subcommand_matches("summarize") {
        let hex = if matches_.get_flag("hex") {
            true
        } else {
            false
        };
        let binary = if matches_.get_flag("bin") {
            true
        } else if !hex {
            true
        } else {
            false
        };

        // Initialize the context
        let mut context = SummarizationContext {
            variable_bytes: &mut vec![],
            variable_ranges: &mut vec![],
            // variables: &mut vec![],
            // mode: &mut Mode::Summarize,
            base: if binary {
                Base::Binary
            } else {
                Base::Hexadecimal
            },
        };

        // Get the variabels ranges
        let variables_ = matches_.get_one::<String>("variables").unwrap();
        let variable_ranges = variables_.split(",");
        for variable_range in variable_ranges {
            let pieces = variable_range.split("-");
            if pieces.clone().count() != 2 {
                panic!("Invalid variable range provided as as argument");
            }

            let pieces_iter = &mut pieces.map(|x| x.parse().unwrap()).into_iter();
            let start: u32 = pieces_iter.next().unwrap();
            let end: u32 = pieces_iter.next().unwrap();

            context.variable_ranges.push((start, end));
        }

        // Analyze the solution
        context.analyze(content);
        context.dump();
    } else if let Some(_) = matches.subcommand_matches("normalize") {
        let mut context = NormalizationContext {
            variables: &mut vec![],
        };

        context.analyze(content);
        context.dump()
    }
}
