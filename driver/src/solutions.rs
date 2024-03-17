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
      path: String,
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
            // println!("{}", String::from_utf8(value.stdout).unwrap());
            let lines: Vec<String> = value
              .stdout
              .lines()
              .take(3)
              .map(|line| line.unwrap())
              .collect();
            assert!(lines.len() == 3);

            let sol = Solution {
              path: path.to_string(),
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
    let mut file = fs::File::create("solutions.csv").expect("Failed to write solutions.csv");
    writeln!(file, "names,exit_codes,process_times,time_limits").unwrap();
    for _ in 0..tasks_count {
      if let Some(Solution {
        path,
        exit_code,
        process_time,
        time_limit,
      }) = rx.recv().await
      {
        let name: &str = Path::new(&path).file_name().unwrap().to_str().unwrap();
        writeln!(
          file,
          "{},{},{:.2},{}",
          name, exit_code, process_time, time_limit
        )
        .unwrap();
      } else {
        break;
      }
    }
    file.flush().unwrap();
  });
}
