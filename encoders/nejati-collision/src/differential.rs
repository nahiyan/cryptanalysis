use nejati_sat::{Formula, Var, Word};

use crate::path::DifferentialPath;
use crate::rules::{Rule, Rules};
use crate::sha256::{Block, ROUND_CONSTANTS};

const ZERO: Word = [0; 32];

pub fn encode(
    rounds: usize,
    xor: bool,
    free_start: bool,
    path: &DifferentialPath,
    rules: &Rules,
) -> Formula {
    let initial_block = !free_start;
    let f = Block::encode(rounds, initial_block, "f", 0, xor);
    let mut g = Block::encode(
        rounds,
        initial_block,
        "g",
        f.formula.variable_count() as i32,
        xor,
    );
    g.formula.name("full_collision", i32::from(!free_start));

    let mut da = vec![ZERO; rounds + 4];
    let mut de = vec![ZERO; rounds + 4];
    let mut dw = vec![ZERO; rounds];
    for index in 0..rounds + 4 {
        da[index] = new_diff(&mut g.formula, &format!("DA_{index}"));
        de[index] = new_diff(&mut g.formula, &format!("DE_{index}"));
        basic_rules(&mut g.formula, &da[index], &f.a[index], &g.a[index]);
        basic_rules(&mut g.formula, &de[index], &f.e[index], &g.e[index]);
        if index < rounds {
            dw[index] = new_diff(&mut g.formula, &format!("DW_{index}"));
            basic_rules(&mut g.formula, &dw[index], &f.w[index], &g.w[index]);
        }
    }

    let zero = g.formula.new_vars(3);
    g.formula.name("zero_g", zero[0]);
    for variable in zero.iter().take(3) {
        g.formula.add_clause(vec![-variable]);
    }

    for index in 0..rounds + 4 {
        if index >= 4 {
            let round = index - 4;
            fix_word(
                &mut g.formula,
                &path.w[round],
                &dw[round],
                &f.w[round],
                &g.w[round],
            );
        }
        fix_word(
            &mut g.formula,
            &path.a[index],
            &da[index],
            &f.a[index],
            &g.a[index],
        );
        fix_word(
            &mut g.formula,
            &path.e[index],
            &de[index],
            &f.e[index],
            &g.e[index],
        );
    }

    let mut ds0 = vec![ZERO; rounds];
    let mut ds1 = vec![ZERO; rounds];
    let mut dw_carry = vec![ZERO; rounds];
    let mut dw_carry_high = vec![ZERO; rounds];
    for round in 16..rounds {
        ds0[round] = new_diff(&mut g.formula, &format!("Ds0_{round}"));
        ds1[round] = new_diff(&mut g.formula, &format!("Ds1_{round}"));
        basic_rules(&mut g.formula, &ds0[round], &f.s0[round], &g.s0[round]);
        basic_rules(&mut g.formula, &ds1[round], &f.s1[round], &g.s1[round]);

        propagate_small_sigma(
            &mut g.formula,
            &rules.xor3,
            &dw[round - 15],
            &ds0[round],
            7,
            18,
            3,
        );
        propagate_small_sigma(
            &mut g.formula,
            &rules.xor3,
            &dw[round - 2],
            &ds1[round],
            17,
            19,
            10,
        );

        dw_carry_high[round] = new_diff(&mut g.formula, &format!("Dadd.W.r1_{round}"));
        dw_carry[round] = new_diff(&mut g.formula, &format!("Dadd.W.r0_{round}"));
        basic_rules(
            &mut g.formula,
            &dw_carry[round],
            &f.w_carry[round],
            &g.w_carry[round],
        );
        basic_rules(
            &mut g.formula,
            &dw_carry_high[round],
            &f.w_carry_high[round],
            &g.w_carry_high[round],
        );
        diff_add(
            &mut g.formula,
            &dw[round],
            &[&dw[round - 16], &ds0[round], &dw[round - 7], &ds1[round]],
            &dw_carry[round],
            Some(&dw_carry_high[round]),
        );
    }

    let mut dsigma0 = vec![ZERO; rounds];
    let mut dsigma1 = vec![ZERO; rounds];
    let mut dchoose = vec![ZERO; rounds];
    let mut dmajority = vec![ZERO; rounds];
    let mut dt = vec![ZERO; rounds];
    let mut dk = vec![ZERO; rounds];
    let mut dt_carry = vec![ZERO; rounds];
    let mut dt_carry_high = vec![ZERO; rounds];
    let mut de_carry = vec![ZERO; rounds];
    let mut da_carry = vec![ZERO; rounds];
    let mut da_carry_high = vec![ZERO; rounds];

    for round in 0..rounds {
        dsigma0[round] = new_diff(&mut g.formula, &format!("Dsigma0_{round}"));
        dsigma1[round] = new_diff(&mut g.formula, &format!("Dsigma1_{round}"));
        basic_rules(
            &mut g.formula,
            &dsigma0[round],
            &f.sigma0[round],
            &g.sigma0[round],
        );
        basic_rules(
            &mut g.formula,
            &dsigma1[round],
            &f.sigma1[round],
            &g.sigma1[round],
        );
        propagate_rotate_xor3(
            &mut g.formula,
            &rules.xor3,
            &da[round + 3],
            &dsigma0[round],
            [2, 13, 22],
        );
        propagate_rotate_xor3(
            &mut g.formula,
            &rules.xor3,
            &de[round + 3],
            &dsigma1[round],
            [6, 11, 25],
        );

        dchoose[round] = new_diff(&mut g.formula, &format!("Dif_{round}"));
        basic_rules(
            &mut g.formula,
            &dchoose[round],
            &f.choose[round],
            &g.choose[round],
        );
        propagate_boolean(
            &mut g.formula,
            &rules.choose,
            [&de[round + 3], &de[round + 2], &de[round + 1]],
            &dchoose[round],
        );

        dmajority[round] = new_diff(&mut g.formula, &format!("Dmaj_{round}"));
        basic_rules(
            &mut g.formula,
            &dmajority[round],
            &f.majority[round],
            &g.majority[round],
        );
        propagate_boolean(
            &mut g.formula,
            &rules.majority,
            [&da[round + 3], &da[round + 2], &da[round + 1]],
            &dmajority[round],
        );

        dt_carry_high[round] = new_diff(&mut g.formula, &format!("Dadd.T.r1_{round}"));
        dt_carry[round] = new_diff(&mut g.formula, &format!("Dadd.T.r0_{round}"));
        dt[round] = new_diff(&mut g.formula, &format!("DT_{round}"));
        basic_rules(
            &mut g.formula,
            &dt_carry_high[round],
            &f.t_carry_high[round],
            &g.t_carry_high[round],
        );
        basic_rules(
            &mut g.formula,
            &dt_carry[round],
            &f.t_carry[round],
            &g.t_carry[round],
        );
        basic_rules(&mut g.formula, &dt[round], &f.t[round], &g.t[round]);
        dk[round] = new_diff(&mut g.formula, &format!("DK_{round}"));
        let _constant = ROUND_CONSTANTS[round];
        for variable in dk[round] {
            g.formula.add_clause(vec![-variable]);
        }
        diff_add(
            &mut g.formula,
            &dt[round],
            &[
                &de[round],
                &dsigma1[round],
                &dchoose[round],
                &dk[round],
                &dw[round],
            ],
            &dt_carry[round],
            Some(&dt_carry_high[round]),
        );

        de_carry[round] = new_diff(&mut g.formula, &format!("Dadd.E.r0_{round}"));
        basic_rules(
            &mut g.formula,
            &de_carry[round],
            &f.e_carry[round],
            &g.e_carry[round],
        );
        diff_add(
            &mut g.formula,
            &de[round + 4],
            &[&da[round], &dt[round]],
            &de_carry[round],
            None,
        );

        da_carry_high[round] = new_diff(&mut g.formula, &format!("Dadd.A.r1_{round}"));
        da_carry[round] = new_diff(&mut g.formula, &format!("Dadd.A.r0_{round}"));
        basic_rules(
            &mut g.formula,
            &da_carry[round],
            &f.a_carry[round],
            &g.a_carry[round],
        );
        basic_rules(
            &mut g.formula,
            &da_carry_high[round],
            &f.a_carry_high[round],
            &g.a_carry_high[round],
        );
        diff_add(
            &mut g.formula,
            &da[round + 4],
            &[&dt[round], &dsigma0[round], &dmajority[round]],
            &da_carry[round],
            Some(&da_carry_high[round]),
        );
    }

    g.formula.append(&f.formula);
    g.formula
}

fn new_diff(formula: &mut Formula, name: &str) -> Word {
    formula.new_word(Some(format!("{name}_g")))
}

fn basic_rules(formula: &mut Formula, difference: &Word, first: &Word, second: &Word) {
    formula.xor2(difference, first, second);
}

fn fix_word(
    formula: &mut Formula,
    description: &str,
    difference: &Word,
    first: &Word,
    second: &Word,
) {
    for (bit, character) in description.bytes().rev().enumerate() {
        fix_bit(formula, character, difference[bit], first[bit], second[bit]);
    }
}

fn fix_bit(formula: &mut Formula, character: u8, difference: Var, first: Var, second: Var) {
    match character {
        b'?' => {}
        b'-' => formula.add_clause(vec![-difference]),
        b'x' => formula.add_clause(vec![difference]),
        b'u' => {
            formula.add_clause(vec![first]);
            formula.add_clause(vec![-second]);
            formula.add_clause(vec![difference]);
        }
        b'n' => {
            formula.add_clause(vec![-first]);
            formula.add_clause(vec![second]);
            formula.add_clause(vec![difference]);
        }
        b'1' => {
            formula.add_clause(vec![first]);
            formula.add_clause(vec![second]);
            formula.add_clause(vec![-difference]);
        }
        b'0' => {
            formula.add_clause(vec![-first]);
            formula.add_clause(vec![-second]);
            formula.add_clause(vec![-difference]);
        }
        _ => unreachable!("the path parser validates differential characters"),
    }
}

fn propagate_small_sigma(
    formula: &mut Formula,
    rules: &[Rule],
    input: &Word,
    output: &Word,
    rotate_a: usize,
    rotate_b: usize,
    shift: usize,
) {
    for bit in 0..32 {
        let inputs = [
            input[(bit + rotate_a) % 32],
            input[(bit + rotate_b) % 32],
            if bit + shift < 32 {
                input[bit + shift]
            } else {
                0
            },
        ];
        for rule in rules {
            if !rule.input.bytes().all(|c| matches!(c, b'-' | b'x' | b'0'))
                || !rule.output.bytes().all(|c| matches!(c, b'-' | b'x'))
                || matches!(rule.input.as_bytes()[0], b'0')
                || matches!(rule.input.as_bytes()[1], b'0')
                || (inputs[2] == 0) != matches!(rule.input.as_bytes()[2], b'0')
            {
                continue;
            }
            impose_rule(formula, &inputs, &[output[bit]], rule);
        }
    }
}

fn propagate_rotate_xor3(
    formula: &mut Formula,
    rules: &[Rule],
    input: &Word,
    output: &Word,
    rotations: [usize; 3],
) {
    for bit in 0..32 {
        let inputs = rotations.map(|rotation| input[(bit + rotation) % 32]);
        for rule in rules {
            if rule.input.bytes().all(|c| matches!(c, b'-' | b'x')) {
                impose_rule(formula, &inputs, &[output[bit]], rule);
            }
        }
    }
}

fn propagate_boolean(formula: &mut Formula, rules: &[Rule], inputs: [&Word; 3], output: &Word) {
    for bit in 0..32 {
        let variables = [inputs[0][bit], inputs[1][bit], inputs[2][bit]];
        for rule in rules {
            if rule.input.bytes().all(|c| matches!(c, b'-' | b'x')) {
                impose_rule(formula, &variables, &[output[bit]], rule);
            }
        }
    }
}

fn impose_rule(formula: &mut Formula, inputs: &[Var], outputs: &[Var], rule: &Rule) {
    let mut antecedent = Vec::new();
    for (variable, character) in inputs.iter().copied().zip(rule.input.bytes()) {
        if variable == 0 || character == b'?' {
            continue;
        }
        antecedent.push(if character == b'-' {
            variable
        } else {
            -variable
        });
    }
    for (variable, character) in outputs.iter().copied().zip(rule.output.bytes()) {
        if character == b'?' {
            continue;
        }
        let mut clause = antecedent.clone();
        clause.push(if character == b'-' {
            -variable
        } else {
            variable
        });
        formula.add_clause(clause);
    }
}

fn diff_add(
    formula: &mut Formula,
    output: &Word,
    operands: &[&Word],
    carry_low: &Word,
    carry_high: Option<&Word>,
) {
    let operand_count = operands.len();
    for bit in 0..32 {
        let mut inputs: Vec<_> = operands.iter().map(|operand| operand[bit]).collect();
        if bit > 0 {
            inputs.push(carry_low[bit - 1]);
        }
        if bit > 1 && ((operand_count == 3 && bit >= 3) || operand_count > 3) {
            inputs.push(
                carry_high.expect("wide differential addition requires a high carry")[bit - 2],
            );
        }
        diff_compressor(
            formula,
            output[bit],
            &inputs,
            carry_low[bit],
            carry_high.map(|carry| carry[bit]),
        );
    }
}

fn diff_compressor(
    formula: &mut Formula,
    output: Var,
    inputs: &[Var],
    carry_low: Var,
    carry_high: Option<Var>,
) {
    formula.xor_many_bit(output, inputs);
    let positives = || inputs.to_vec();
    let negatives = || inputs.iter().map(|variable| -*variable).collect::<Vec<_>>();
    match inputs.len() {
        2 => {
            let mut clause = positives();
            clause.push(-carry_low);
            formula.add_clause(clause);
        }
        3 => {
            let mut clause = positives();
            clause.push(-carry_low);
            formula.add_clause(clause);
            let mut clause = negatives();
            clause.push(carry_low);
            formula.add_clause(clause);
        }
        4 => {
            for carry in [
                carry_low,
                carry_high.expect("four-input compressor requires high carry"),
            ] {
                let mut clause = positives();
                clause.push(-carry);
                formula.add_clause(clause);
            }
        }
        5 => {
            for carry in [
                carry_low,
                carry_high.expect("five-input compressor requires high carry"),
            ] {
                let mut clause = positives();
                clause.push(-carry);
                formula.add_clause(clause);
            }
            let mut clause = negatives();
            clause.push(-carry_low);
            formula.add_clause(clause);
        }
        6 => {
            for carry in [
                carry_low,
                carry_high.expect("six-input compressor requires high carry"),
            ] {
                let mut clause = positives();
                clause.push(-carry);
                formula.add_clause(clause);
            }
        }
        7 => {
            let high = carry_high.expect("seven-input compressor requires high carry");
            for carry in [carry_low, high] {
                let mut clause = positives();
                clause.push(-carry);
                formula.add_clause(clause);
                let mut clause = negatives();
                clause.push(carry);
                formula.add_clause(clause);
            }
        }
        count => panic!("unsupported differential compressor width {count}"),
    }
}
