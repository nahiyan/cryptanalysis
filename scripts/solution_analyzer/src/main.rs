use std::env;
use std::fs;
use std::str::FromStr;
use std::vec;

// Definition of the context
struct Context<'a> {
    variable_ranges: &'a mut Vec<(u32, u32)>,
    variable_bytes: &'a mut Vec<u8>,
}

impl Context<'_> {
    // Analyze the line
    fn analyze_line(&mut self, line: &str) {
        // Ensure that the line starts with 'v'
        if !line.starts_with('v') {
            return;
        }

        let pieces = line.split_whitespace();
        for piece in pieces {
            if !piece.eq("v") && !piece.eq("0") {
                let variable: u32 = i32::unsigned_abs(FromStr::from_str(piece).unwrap());
                let value: u8 = if piece.starts_with('-') { 0 } else { 1 };

                for (start, end) in self.variable_ranges.iter() {
                    if variable >= *start && variable <= *end {
                        self.variable_bytes.push(value)
                    }
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
}

fn main() {
    // Initialize the context
    let mut context = Context {
        variable_bytes: &mut vec![],
        variable_ranges: &mut vec![],
    };

    // Process the CLI arguments
    let args: Vec<String> = env::args().collect();
    if args.len() == 1 {
        panic!("Missing argument: The solutions file, which holds the output of the SAT solver.");
    }

    if args.len() > 2 {
        for arg in &args[2..] {
            let pieces = arg.split('-');
            if pieces.clone().count() != 2 {
                panic!("Invalid variable range provided as as argument");
            }

            let pieces_int: Vec<u32> = pieces.map(|x| FromStr::from_str(x).unwrap()).collect();

            let start = pieces_int[0];
            let end = pieces_int[1];

            // if (end - start + 1) % 8 != 0 {
            //     panic!("One of the ranges has a width not divisible by 8!")
            // }

            context.variable_ranges.push((start, end));
        }
    }

    // Analyze the solution
    let encodings_file = &args[1];
    match fs::read_to_string(encodings_file) {
        Err(_) => {
            println!("Failed to open the solutions file.");
        }
        Ok(content) => {
            context.analyze(content);
        }
    }

    // Output the analysis
    let mut i = 0;
    for (start, end) in context.variable_ranges.iter() {
        println!("Range: {} - {}", start, end);
        for (j, _) in (*start..*end + 1).enumerate() {
            print!("{}", context.variable_bytes[i]);
            i += 1;

            if (j + 1) % 8 == 0 {
                println!()
            }
        }
        println!()
    }

    // TODO: Output in hex
    // println!();
    // for i in 0..context.variable_ranges.len() + 1 {
    //     println!("Range");
    //     let value = &context.variable_bytes[(8 * i)..(8 * i + 8)];
    //     println!("{}", hex::encode(value));
    // }
}
