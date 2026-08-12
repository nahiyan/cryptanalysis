use std::collections::BTreeMap;
use std::io::{self, Write};

use crate::templates::adder_template;

pub type Var = i32;
pub type Word = [Var; 32];

#[derive(Clone, Debug)]
pub struct Clause {
    pub literals: Vec<Var>,
    pub xor: bool,
}

#[derive(Default)]
pub struct Formula {
    next_var: Var,
    clauses: Vec<Clause>,
    names: BTreeMap<String, Var>,
    use_xor_clauses: bool,
}

impl Formula {
    pub fn new(use_xor_clauses: bool) -> Self {
        Self {
            use_xor_clauses,
            ..Self::default()
        }
    }

    pub fn with_start(use_xor_clauses: bool, first_variable: Var) -> Self {
        Self {
            next_var: first_variable,
            use_xor_clauses,
            ..Self::default()
        }
    }

    pub fn variable_count(&self) -> usize {
        self.next_var as usize
    }

    pub fn clause_count(&self) -> usize {
        self.clauses.len()
    }

    pub fn new_vars(&mut self, count: usize) -> Vec<Var> {
        (0..count)
            .map(|_| {
                self.next_var += 1;
                self.next_var
            })
            .collect()
    }

    pub fn new_word(&mut self, name: Option<String>) -> Word {
        let vars: Word = self
            .new_vars(32)
            .try_into()
            .expect("a word always contains 32 variables");
        if let Some(name) = name {
            self.names.insert(name, vars[0]);
        }
        vars
    }

    pub fn name(&mut self, name: impl Into<String>, first_variable: Var) {
        self.names.insert(name.into(), first_variable);
    }

    pub fn append(&mut self, other: &Self) {
        self.next_var = self.next_var.max(other.next_var);
        self.clauses.extend(other.clauses.iter().cloned());
        self.names.extend(
            other
                .names
                .iter()
                .map(|(name, variable)| (name.clone(), *variable)),
        );
    }

    pub fn add_clause(&mut self, literals: Vec<Var>) {
        self.push_clause(literals, false);
    }

    fn push_clause(&mut self, literals: Vec<Var>, xor: bool) {
        assert!(
            literals.iter().all(|literal| *literal != 0),
            "a clause cannot contain variable zero"
        );
        self.clauses.push(Clause { literals, xor });
    }

    pub fn fixed_value(&mut self, word: &[Var], value: u32) {
        self.fixed_value_bits(word, value, word.len());
    }

    pub fn fixed_value_bits(&mut self, word: &[Var], value: u32, bits: usize) {
        assert!(bits <= word.len());
        for (bit, variable) in word.iter().copied().take(bits).enumerate() {
            let literal = if (value >> bit) & 1 == 1 {
                variable
            } else {
                -variable
            };
            self.add_clause(vec![literal]);
        }
    }

    pub fn rotate_left(&self, input: &Word, positions: usize) -> Word {
        std::array::from_fn(|bit| input[(bit + 32 - positions % 32) % 32])
    }

    pub fn rotate_right(&self, input: &Word, positions: usize) -> Word {
        self.rotate_left(input, 32 - positions % 32)
    }

    pub fn equivalent(&mut self, output: &Word, input: &Word) {
        for bit in 0..32 {
            self.add_clause(vec![-output[bit], input[bit]]);
            self.add_clause(vec![output[bit], -input[bit]]);
        }
    }

    pub fn choose(&mut self, output: &Word, selector: &Word, yes: &Word, no: &Word) {
        for bit in 0..32 {
            self.add_clause(vec![-output[bit], selector[bit], no[bit]]);
            self.add_clause(vec![-output[bit], -selector[bit], yes[bit]]);
            self.add_clause(vec![output[bit], selector[bit], -no[bit]]);
            self.add_clause(vec![output[bit], -selector[bit], -yes[bit]]);
        }
    }

    pub fn majority3(&mut self, output: &Word, a: &Word, b: &Word, c: &Word) {
        for bit in 0..32 {
            self.add_clause(vec![-output[bit], a[bit], b[bit]]);
            self.add_clause(vec![-output[bit], a[bit], c[bit]]);
            self.add_clause(vec![-output[bit], b[bit], c[bit]]);
            self.add_clause(vec![output[bit], -b[bit], -c[bit]]);
            self.add_clause(vec![output[bit], -a[bit], -c[bit]]);
            self.add_clause(vec![output[bit], -a[bit], -b[bit]]);
        }
    }

    pub fn xor2(&mut self, output: &[Var], a: &[Var], b: &[Var]) {
        assert_eq!(output.len(), a.len());
        assert_eq!(output.len(), b.len());
        for bit in 0..output.len() {
            if self.use_xor_clauses {
                self.push_clause(vec![-output[bit], a[bit], b[bit]], true);
            } else {
                self.add_clause(vec![-output[bit], -a[bit], -b[bit]]);
                self.add_clause(vec![output[bit], -a[bit], b[bit]]);
                self.add_clause(vec![output[bit], a[bit], -b[bit]]);
                self.add_clause(vec![-output[bit], a[bit], b[bit]]);
            }
        }
    }

    pub fn xor3(&mut self, output: &[Var], a: &[Var], b: &[Var], c: &[Var]) {
        assert_eq!(output.len(), a.len());
        assert_eq!(output.len(), b.len());
        assert_eq!(output.len(), c.len());
        for bit in 0..output.len() {
            if self.use_xor_clauses {
                self.push_clause(vec![-output[bit], a[bit], b[bit], c[bit]], true);
            } else {
                self.add_clause(vec![output[bit], -a[bit], -b[bit], -c[bit]]);
                self.add_clause(vec![-output[bit], -a[bit], -b[bit], c[bit]]);
                self.add_clause(vec![-output[bit], -a[bit], b[bit], -c[bit]]);
                self.add_clause(vec![output[bit], -a[bit], b[bit], c[bit]]);
                self.add_clause(vec![-output[bit], a[bit], -b[bit], -c[bit]]);
                self.add_clause(vec![output[bit], a[bit], -b[bit], c[bit]]);
                self.add_clause(vec![output[bit], a[bit], b[bit], -c[bit]]);
                self.add_clause(vec![-output[bit], a[bit], b[bit], c[bit]]);
            }
        }
    }

    pub fn xor4(&mut self, output: &Word, a: &Word, b: &Word, c: &Word, d: &Word) {
        for bit in 0..32 {
            if self.use_xor_clauses {
                self.push_clause(vec![-output[bit], a[bit], b[bit], c[bit], d[bit]], true);
            } else {
                self.add_clause(vec![-output[bit], -a[bit], -b[bit], -c[bit], -d[bit]]);
                self.add_clause(vec![output[bit], -a[bit], -b[bit], -c[bit], d[bit]]);
                self.add_clause(vec![output[bit], -a[bit], -b[bit], c[bit], -d[bit]]);
                self.add_clause(vec![-output[bit], -a[bit], -b[bit], c[bit], d[bit]]);
                self.add_clause(vec![output[bit], -a[bit], b[bit], -c[bit], -d[bit]]);
                self.add_clause(vec![-output[bit], -a[bit], b[bit], -c[bit], d[bit]]);
                self.add_clause(vec![-output[bit], -a[bit], b[bit], c[bit], -d[bit]]);
                self.add_clause(vec![output[bit], -a[bit], b[bit], c[bit], d[bit]]);
                self.add_clause(vec![output[bit], a[bit], -b[bit], -c[bit], -d[bit]]);
                self.add_clause(vec![-output[bit], a[bit], -b[bit], -c[bit], d[bit]]);
                self.add_clause(vec![-output[bit], a[bit], -b[bit], c[bit], -d[bit]]);
                self.add_clause(vec![output[bit], a[bit], -b[bit], c[bit], d[bit]]);
                self.add_clause(vec![-output[bit], a[bit], b[bit], -c[bit], -d[bit]]);
                self.add_clause(vec![output[bit], a[bit], b[bit], -c[bit], d[bit]]);
                self.add_clause(vec![output[bit], a[bit], b[bit], c[bit], -d[bit]]);
                self.add_clause(vec![-output[bit], a[bit], b[bit], c[bit], d[bit]]);
            }
        }
    }

    pub fn xor_many_bit(&mut self, output: Var, inputs: &[Var]) {
        assert!((2..=7).contains(&inputs.len()));
        if self.use_xor_clauses {
            let mut clause = Vec::with_capacity(inputs.len() + 1);
            clause.push(-output);
            clause.extend_from_slice(inputs);
            self.push_clause(clause, true);
            return;
        }

        for assignment in 0..(1usize << inputs.len()) {
            let parity = assignment.count_ones() & 1 == 1;
            for output_value in [false, true] {
                if output_value == parity {
                    continue;
                }
                let mut clause = Vec::with_capacity(inputs.len() + 1);
                clause.push(if output_value { -output } else { output });
                for (index, variable) in inputs.iter().copied().enumerate() {
                    clause.push(if assignment >> index & 1 == 1 {
                        -variable
                    } else {
                        variable
                    });
                }
                self.add_clause(clause);
            }
        }
    }

    pub fn add2(&mut self, output: &Word, a: &Word, b: &Word) {
        self.add_operands(output, &[a, b]);
    }

    pub fn add3(&mut self, output: &Word, a: &Word, b: &Word, c: &Word) {
        self.add_operands(output, &[a, b, c]);
    }

    pub fn add4(&mut self, output: &Word, a: &Word, b: &Word, c: &Word, d: &Word) {
        self.add_operands(output, &[a, b, c, d]);
    }

    pub fn add5(&mut self, output: &Word, a: &Word, b: &Word, c: &Word, d: &Word, e: &Word) {
        self.add_operands(output, &[a, b, c, d, e]);
    }

    fn add_operands(&mut self, output: &Word, operands: &[&Word]) {
        let mut columns = vec![Vec::<Var>::new(); 37];
        for bit in 0..32 {
            for operand in operands {
                columns[bit].push(operand[bit]);
            }

            let carry_bits = floor_log2(columns[bit].len());
            let mut sum = Vec::with_capacity(carry_bits + 1);
            sum.push(output[bit]);
            sum.extend(self.new_vars(carry_bits));

            for carry in 1..sum.len() {
                columns[bit + carry].push(sum[carry]);
            }

            self.emit_adder_template(&columns[bit], &sum);
        }
    }

    pub fn add_with_carries(
        &mut self,
        output: &Word,
        operands: &[&Word],
        carry_low: &Word,
        carry_high: Option<&Word>,
    ) {
        assert!((2..=5).contains(&operands.len()));
        assert_eq!(carry_high.is_some(), operands.len() > 2);
        let mut columns = vec![Vec::<Var>::new(); 37];
        for bit in 0..32 {
            for operand in operands {
                columns[bit].push(operand[bit]);
            }

            let mut sum = vec![output[bit], carry_low[bit]];
            columns[bit + 1].push(carry_low[bit]);
            if columns[bit].len() > 3 {
                let carry_high = carry_high.expect("wide additions require a high carry");
                sum.push(carry_high[bit]);
                columns[bit + 2].push(carry_high[bit]);
            }
            self.emit_adder_template(&columns[bit], &sum);
        }
    }

    fn emit_adder_template(&mut self, lhs: &[Var], rhs_lsb: &[Var]) {
        let template = adder_template(lhs.len(), rhs_lsb.len()).unwrap_or_else(|| {
            panic!(
                "unsupported embedded Espresso adder shape ({},{})",
                lhs.len(),
                rhs_lsb.len()
            )
        });

        for relative_clause in template {
            let mut clause = Vec::with_capacity(relative_clause.len());
            for relative_literal in *relative_clause {
                let position = relative_literal.unsigned_abs() as usize - 1;
                let variable = if position < lhs.len() {
                    lhs[position]
                } else {
                    let rhs_position = rhs_lsb.len() - 1 - (position - lhs.len());
                    rhs_lsb[rhs_position]
                };
                clause.push(if *relative_literal < 0 {
                    -variable
                } else {
                    variable
                });
            }
            self.add_clause(clause);
        }
    }

    pub fn write_dimacs(&self, mut output: impl Write) -> io::Result<()> {
        writeln!(
            output,
            "p cnf {} {}",
            self.variable_count(),
            self.clause_count()
        )?;
        for clause in &self.clauses {
            if clause.xor {
                write!(output, "x ")?;
            }
            for literal in &clause.literals {
                write!(output, "{literal} ")?;
            }
            writeln!(output, "0")?;
        }
        for (name, first_variable) in &self.names {
            writeln!(output, "c {name} {first_variable}")?;
        }
        Ok(())
    }

    pub fn write_dimacs_with_order(&self, mut output: impl Write, rounds: usize) -> io::Result<()> {
        self.write_dimacs(&mut output)?;
        writeln!(output, "c order {rounds}")
    }
}

fn floor_log2(value: usize) -> usize {
    assert!(value > 0);
    usize::BITS as usize - 1 - value.leading_zeros() as usize
}
