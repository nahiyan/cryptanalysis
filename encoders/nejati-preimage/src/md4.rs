use nejati_sat::{Formula, Word};

const ROTATIONS: [usize; 48] = [
    3, 7, 11, 19, 3, 7, 11, 19, 3, 7, 11, 19, 3, 7, 11, 19, 3, 5, 9, 13, 3, 5, 9, 13, 3, 5, 9, 13,
    3, 5, 9, 13, 3, 9, 11, 15, 3, 9, 11, 15, 3, 9, 11, 15, 3, 9, 11, 15,
];

const MESSAGE_INDEX: [usize; 48] = [
    0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0, 4, 8, 12, 1, 5, 9, 13, 2, 6, 10, 14,
    3, 7, 11, 15, 0, 8, 4, 12, 2, 10, 6, 14, 1, 9, 5, 13, 3, 11, 7, 15,
];

pub struct Md4Encoding {
    pub formula: Formula,
    pub message: [Word; 16],
    output: [Word; 4],
}

impl Md4Encoding {
    pub fn build(
        rounds: usize,
        use_xor_clauses: bool,
        dobbertin: bool,
        relaxed_bits: usize,
    ) -> Self {
        assert!((1..=48).contains(&rounds));
        assert!((1..=32).contains(&relaxed_bits));
        assert!(!dobbertin || rounds >= 27);

        let mut formula = Formula::new(use_xor_clauses);
        let message: [Word; 16] =
            std::array::from_fn(|word| formula.new_word(Some(format!("w{word}"))));
        let chain: [Word; 4] = std::array::from_fn(|_| formula.new_word(None));
        let output: [Word; 4] =
            std::array::from_fn(|word| formula.new_word(Some(format!("hash{word}"))));

        formula.fixed_value(&chain[0], 0x6745_2301);
        formula.fixed_value(&chain[1], 0xefcd_ab89);
        formula.fixed_value(&chain[2], 0x98ba_dcfe);
        formula.fixed_value(&chain[3], 0x1032_5476);

        let round_constant_2 = formula.new_word(None);
        let round_constant_3 = formula.new_word(None);
        formula.fixed_value(&round_constant_2, 0x5a82_7999);
        formula.fixed_value(&round_constant_3, 0x6ed9_eba1);

        let mut q = vec![chain[0], chain[3], chain[2], chain[1]];
        for round in 0..rounds {
            let temporary = formula.new_word(None);
            let function = formula.new_word(None);
            if round < 16 {
                formula.choose(&function, &q[round + 3], &q[round + 2], &q[round + 1]);
                formula.add3(
                    &temporary,
                    &q[round],
                    &function,
                    &message[MESSAGE_INDEX[round]],
                );
            } else if round < 32 {
                formula.majority3(&function, &q[round + 3], &q[round + 2], &q[round + 1]);
                formula.add4(
                    &temporary,
                    &q[round],
                    &function,
                    &message[MESSAGE_INDEX[round]],
                    &round_constant_2,
                );
            } else {
                formula.xor3(&function, &q[round + 3], &q[round + 2], &q[round + 1]);
                formula.add4(
                    &temporary,
                    &q[round],
                    &function,
                    &message[MESSAGE_INDEX[round]],
                    &round_constant_3,
                );
            }
            q.push(formula.rotate_left(&temporary, ROTATIONS[round]));
        }

        if dobbertin {
            for index in [16usize, 20, 24, 28, 17, 21, 25, 29, 18, 22, 26, 30] {
                if index == 16 && relaxed_bits != 32 {
                    formula.fixed_value_bits(
                        &q[index][32 - relaxed_bits..],
                        u32::MAX,
                        relaxed_bits,
                    );
                } else {
                    formula.fixed_value(&q[index], u32::MAX);
                }
            }
        }

        let remainder = rounds % 4;
        formula.equivalent(&output[0], &q[rounds + (4 - remainder) % 4]);
        formula.equivalent(&output[1], &q[rounds + (3 - remainder) % 4]);
        formula.equivalent(&output[2], &q[rounds + (6 - remainder) % 4]);
        formula.equivalent(&output[3], &q[rounds + (5 - remainder) % 4]);

        Self {
            formula,
            message,
            output,
        }
    }

    pub fn fix_output(&mut self, target: [u32; 4]) {
        for (word, value) in self.output.iter().zip(target) {
            self.formula.fixed_value(word, value);
        }
    }

    pub fn fix_message_prefix(&mut self, message: &[u32; 16], bits: usize) {
        for bit in 0..bits {
            let word = bit / 32;
            let offset = bit % 32;
            let value = (message[word] >> offset) & 1;
            let variable = self.message[word][offset];
            self.formula
                .add_clause(vec![if value == 1 { variable } else { -variable }]);
        }
    }
}

pub fn compression_state(message: &[u32; 16], rounds: usize) -> [u32; 4] {
    let mut state = [0x6745_2301u32, 0xefcd_ab89, 0x98ba_dcfe, 0x1032_5476];
    for round in 0..rounds {
        let slot = (48 - round) % 4;
        let b = state[(slot + 1) % 4];
        let c = state[(slot + 2) % 4];
        let d = state[(slot + 3) % 4];
        let (function, constant) = if round < 16 {
            ((b & c) | ((!b) & d), 0)
        } else if round < 32 {
            ((b & c) | (b & d) | (c & d), 0x5a82_7999)
        } else {
            (b ^ c ^ d, 0x6ed9_eba1)
        };
        state[slot] = state[slot]
            .wrapping_add(function)
            .wrapping_add(message[MESSAGE_INDEX[round]])
            .wrapping_add(constant)
            .rotate_left(ROTATIONS[round] as u32);
    }
    state
}
