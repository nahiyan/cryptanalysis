use std::fs;
use std::path::Path;

pub struct DifferentialPath {
    pub a: Vec<String>,
    pub e: Vec<String>,
    pub w: Vec<String>,
}

impl DifferentialPath {
    pub fn read(path: &Path, rounds: usize) -> Result<Self, String> {
        let contents = fs::read_to_string(path)
            .map_err(|error| format!("failed to read {}: {error}", path.display()))?;
        let lines: Vec<_> = contents
            .lines()
            .filter(|line| !line.trim().is_empty())
            .collect();
        if lines.len() < rounds + 4 {
            return Err(format!(
                "{} contains {} path rows; {} are required",
                path.display(),
                lines.len(),
                rounds + 4
            ));
        }

        let mut result = Self {
            a: Vec::with_capacity(rounds + 4),
            e: Vec::with_capacity(rounds + 4),
            w: Vec::with_capacity(rounds),
        };
        for (index, line) in lines.iter().take(rounds + 4).enumerate() {
            let fields: Vec<_> = line.split_whitespace().collect();
            let expected_round = index as isize - 4;
            let round = fields
                .first()
                .ok_or_else(|| format!("path row {} is empty", index + 1))?
                .parse::<isize>()
                .map_err(|_| format!("path row {} has an invalid round", index + 1))?;
            if round != expected_round {
                return Err(format!(
                    "path row {} describes round {round}, expected {expected_round}",
                    index + 1
                ));
            }
            let required = if expected_round < 0 { 5 } else { 7 };
            if fields.len() < required
                || fields.get(1) != Some(&"A:")
                || fields.get(3) != Some(&"E:")
            {
                return Err(format!("path row {} has an invalid format", index + 1));
            }
            validate_word(fields[2], index + 1, "A")?;
            validate_word(fields[4], index + 1, "E")?;
            result.a.push(fields[2].to_string());
            result.e.push(fields[4].to_string());
            if expected_round >= 0 {
                if fields.get(5) != Some(&"W:") {
                    return Err(format!("path row {} is missing W:", index + 1));
                }
                validate_word(fields[6], index + 1, "W")?;
                result.w.push(fields[6].to_string());
            }
        }
        Ok(result)
    }
}

fn validate_word(word: &str, row: usize, field: &str) -> Result<(), String> {
    if word.len() != 32 {
        return Err(format!("path row {row} {field} must contain 32 characters"));
    }
    if let Some(character) = word
        .chars()
        .find(|character| !matches!(character, '-' | 'x' | 'u' | 'n' | '0' | '1' | '?'))
    {
        return Err(format!(
            "path row {row} {field} contains unsupported character '{character}'"
        ));
    }
    Ok(())
}
