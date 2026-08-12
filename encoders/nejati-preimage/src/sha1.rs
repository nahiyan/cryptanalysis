use nejati_sat::{Formula, Word};

const ROUND_CONSTANTS: [u32; 4] = [0x5a82_7999, 0x6ed9_eba1, 0x8f1b_bcdc, 0xca62_c1d6];

pub struct Sha1Encoding {
    pub formula: Formula,
    pub message: [Word; 16],
    output: [Word; 5],
}

impl Sha1Encoding {
    pub fn build(rounds: usize, use_xor_clauses: bool) -> Self {
        assert!((16..=80).contains(&rounds));
        let mut formula = Formula::new(use_xor_clauses);

        let initial_message: [Word; 16] =
            std::array::from_fn(|word| formula.new_word(Some(format!("w{word}"))));
        let mut schedule = initial_message.to_vec();
        let mut expansion_temporary = Vec::with_capacity(rounds.saturating_sub(16));
        for word in 16..rounds {
            expansion_temporary.push(formula.new_word(Some(format!("w{word}"))));
            schedule.push([0; 32]);
        }

        let chain: [Word; 5] = std::array::from_fn(|_| formula.new_word(None));
        let output: [Word; 5] =
            std::array::from_fn(|word| formula.new_word(Some(format!("hash{word}"))));

        let mut state = vec![
            formula.rotate_left(&chain[4], 2),
            formula.rotate_left(&chain[3], 2),
            formula.rotate_left(&chain[2], 2),
            chain[1],
            chain[0],
        ];
        for _ in 0..rounds {
            state.push(formula.new_word(None));
        }

        for word in 16..rounds {
            let temporary = expansion_temporary[word - 16];
            formula.xor4(
                &temporary,
                &schedule[word - 3],
                &schedule[word - 8],
                &schedule[word - 14],
                &schedule[word - 16],
            );
            schedule[word] = formula.rotate_left(&temporary, 1);
        }

        let constants: [Word; 4] = std::array::from_fn(|index| {
            let word = formula.new_word(None);
            formula.fixed_value(&word, ROUND_CONSTANTS[index]);
            word
        });

        formula.fixed_value(&chain[0], 0x6745_2301);
        formula.fixed_value(&chain[1], 0xefcd_ab89);
        formula.fixed_value(&chain[2], 0x98ba_dcfe);
        formula.fixed_value(&chain[3], 0x1032_5476);
        formula.fixed_value(&chain[4], 0xc3d2_e1f0);

        for round in 0..rounds {
            let previous_a = formula.rotate_left(&state[round + 4], 5);
            let b = state[round + 3];
            let c = formula.rotate_left(&state[round + 2], 30);
            let d = formula.rotate_left(&state[round + 1], 30);
            let e = formula.rotate_left(&state[round], 30);

            let function = formula.new_word(None);
            if round < 20 {
                formula.choose(&function, &b, &c, &d);
            } else if round < 40 {
                formula.xor3(&function, &b, &c, &d);
            } else if round < 60 {
                formula.majority3(&function, &b, &c, &d);
            } else {
                formula.xor3(&function, &b, &c, &d);
            }

            formula.add5(
                &state[round + 5],
                &previous_a,
                &function,
                &e,
                &constants[round / 20],
                &schedule[round],
            );
        }

        let c = formula.rotate_left(&state[rounds + 2], 30);
        let d = formula.rotate_left(&state[rounds + 1], 30);
        let e = formula.rotate_left(&state[rounds], 30);
        formula.add2(&output[0], &chain[0], &state[rounds + 4]);
        formula.add2(&output[1], &chain[1], &state[rounds + 3]);
        formula.add2(&output[2], &chain[2], &c);
        formula.add2(&output[3], &chain[3], &d);
        formula.add2(&output[4], &chain[4], &e);

        Self {
            formula,
            message: initial_message,
            output,
        }
    }

    pub fn fix_output(&mut self, target: [u32; 5]) {
        for (word, value) in self.output.iter().zip(target) {
            self.formula.fixed_value(word, value);
        }
    }

    pub fn fix_message_prefix(&mut self, message: &[u32; 16], bits: usize) {
        fix_message_prefix(&mut self.formula, &self.message, message, bits);
    }
}

pub fn compression(message: &[u32; 16], rounds: usize) -> [u32; 5] {
    let mut schedule = vec![0u32; rounds.max(16)];
    schedule[..16].copy_from_slice(message);
    for word in 16..rounds {
        schedule[word] =
            (schedule[word - 3] ^ schedule[word - 8] ^ schedule[word - 14] ^ schedule[word - 16])
                .rotate_left(1);
    }

    let initial = [
        0x6745_2301u32,
        0xefcd_ab89,
        0x98ba_dcfe,
        0x1032_5476,
        0xc3d2_e1f0,
    ];
    let [mut a, mut b, mut c, mut d, mut e] = initial;
    for (round, schedule_word) in schedule.iter().copied().enumerate().take(rounds) {
        let (function, constant) = if round < 20 {
            ((b & c) | ((!b) & d), ROUND_CONSTANTS[0])
        } else if round < 40 {
            (b ^ c ^ d, ROUND_CONSTANTS[1])
        } else if round < 60 {
            ((b & c) | (b & d) | (c & d), ROUND_CONSTANTS[2])
        } else {
            (b ^ c ^ d, ROUND_CONSTANTS[3])
        };
        let temporary = a
            .rotate_left(5)
            .wrapping_add(function)
            .wrapping_add(e)
            .wrapping_add(constant)
            .wrapping_add(schedule_word);
        e = d;
        d = c;
        c = b.rotate_left(30);
        b = a;
        a = temporary;
    }
    [
        initial[0].wrapping_add(a),
        initial[1].wrapping_add(b),
        initial[2].wrapping_add(c),
        initial[3].wrapping_add(d),
        initial[4].wrapping_add(e),
    ]
}

fn fix_message_prefix(
    formula: &mut Formula,
    variables: &[Word; 16],
    message: &[u32; 16],
    bits: usize,
) {
    for bit in 0..bits {
        let word = bit / 32;
        let offset = bit % 32;
        let variable = variables[word][offset];
        formula.add_clause(vec![if (message[word] >> offset) & 1 == 1 {
            variable
        } else {
            -variable
        }]);
    }
}
