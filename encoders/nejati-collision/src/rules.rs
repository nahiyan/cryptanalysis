#[derive(Clone)]
pub struct Rule {
    pub input: String,
    pub output: String,
}

pub struct Rules {
    pub choose: Vec<Rule>,
    pub majority: Vec<Rule>,
    pub xor3: Vec<Rule>,
}

impl Rules {
    pub fn embedded() -> Result<Self, String> {
        let mut result = Self {
            choose: Vec::new(),
            majority: Vec::new(),
            xor3: Vec::new(),
        };
        for (line_number, line) in include_str!("../prop_rules.db").lines().enumerate() {
            let mut fields = line.split_whitespace();
            let kind = fields
                .next()
                .ok_or_else(|| format!("empty rule at line {}", line_number + 1))?;
            let input = fields
                .next()
                .ok_or_else(|| format!("missing rule input at line {}", line_number + 1))?;
            let output = fields
                .next()
                .ok_or_else(|| format!("missing rule output at line {}", line_number + 1))?;
            let rule = Rule {
                input: input.into(),
                output: output.into(),
            };
            match kind {
                "ch" => result.choose.push(rule),
                "maj" => result.majority.push(rule),
                "xor3" => result.xor3.push(rule),
                "add2" | "add3" | "add4" | "add5" | "add6" | "add7" => {}
                _ => {
                    return Err(format!(
                        "unknown rule kind '{kind}' at line {}",
                        line_number + 1
                    ));
                }
            }
        }
        Ok(result)
    }
}
