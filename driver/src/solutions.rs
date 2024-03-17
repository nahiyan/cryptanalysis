use std::{
  fs,
  io::{BufRead, Write},
  path::Path,
  process,
};
use tokio::sync::mpsc;

pub fn parse_sols(input_dir: String) {
  let parse_script = include_bytes!("../scripts/parse_log.sh");
  let mut script_file = fs::File::create("parse_log.sh").unwrap();
  script_file
    .write(parse_script)
    .expect("Failed to write script.");

  // Parse the solutions concurrently
  let rt = tokio::runtime::Runtime::new().unwrap();
  rt.block_on(async {
    struct Solution {
      name: String,
      exit_code: i32,
      process_time: f32,
      time_limit: i32,
    }

    let (tx, mut rx) = mpsc::channel::<Solution>(100);
    let mut tasks_count = 0;
    for file in fs::read_dir(input_dir).unwrap() {
      let tx_ = tx.clone();
      tokio::spawn(async move {
        let path_ = file.unwrap().path();
        let path = path_.to_str().unwrap();
        let output = process::Command::new("sh")
          .args(["parse_log.sh", path])
          .output();
        match output {
          Ok(value) => {
            let lines: Vec<String> = value
              .stdout
              .lines()
              .take(3)
              .map(|line| line.unwrap())
              .collect();
            assert!(lines.len() == 3);

            let sol = Solution {
              name: String::from(Path::new(path).file_name().unwrap().to_str().unwrap()),
              exit_code: if !lines[0].is_empty() {
                lines[0].parse::<i32>().unwrap()
              } else {
                0
              },
              process_time: if !lines[1].is_empty() {
                lines[1].parse::<f32>().unwrap()
              } else {
                0.0
              },
              time_limit: if !lines[2].is_empty() {
                lines[2].parse::<i32>().unwrap()
              } else {
                0
              },
            };

            tx_.send(sol).await.unwrap();
          }
          Err(_) => panic!("Failed to parse log."),
        }
      });
      tasks_count += 1;
    }

    // Get the solutions through the channel
    let mut sols: Vec<Solution> = Vec::new();
    for _ in 0..tasks_count {
      if let Some(sol) = rx.recv().await {
        sols.push(sol);
      } else {
        break;
      }
    }
    sols.sort_by(|a, b| b.name.cmp(&a.name));

    let mut sols_file = fs::File::create("solutions.csv").expect("Failed to write solutions.csv");
    writeln!(sols_file, "names,exit_codes,process_times,time_limits").unwrap();
    for Solution {
      name,
      exit_code,
      process_time,
      time_limit,
    } in sols
    {
      writeln!(
        sols_file,
        "{},{},{:.2},{}",
        name, exit_code, process_time, time_limit
      )
      .unwrap();
    }
    sols_file.flush().unwrap();
  });
}
