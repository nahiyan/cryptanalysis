use nejati_sat::{Formula, Word};

const ROUND_CONSTANTS: [u32; 64] = [
    0x428a_2f98,
    0x7137_4491,
    0xb5c0_fbcf,
    0xe9b5_dba5,
    0x3956_c25b,
    0x59f1_11f1,
    0x923f_82a4,
    0xab1c_5ed5,
    0xd807_aa98,
    0x1283_5b01,
    0x2431_85be,
    0x550c_7dc3,
    0x72be_5d74,
    0x80de_b1fe,
    0x9bdc_06a7,
    0xc19b_f174,
    0xe49b_69c1,
    0xefbe_4786,
    0x0fc1_9dc6,
    0x240c_a1cc,
    0x2de9_2c6f,
    0x4a74_84aa,
    0x5cb0_a9dc,
    0x76f9_88da,
    0x983e_5152,
    0xa831_c66d,
    0xb003_27c8,
    0xbf59_7fc7,
    0xc6e0_0bf3,
    0xd5a7_9147,
    0x06ca_6351,
    0x1429_2967,
    0x27b7_0a85,
    0x2e1b_2138,
    0x4d2c_6dfc,
    0x5338_0d13,
    0x650a_7354,
    0x766a_0abb,
    0x81c2_c92e,
    0x9272_2c85,
    0xa2bf_e8a1,
    0xa81a_664b,
    0xc24b_8b70,
    0xc76c_51a3,
    0xd192_e819,
    0xd699_0624,
    0xf40e_3585,
    0x106a_a070,
    0x19a4_c116,
    0x1e37_6c08,
    0x2748_774c,
    0x34b0_bcb5,
    0x391c_0cb3,
    0x4ed8_aa4a,
    0x5b9c_ca4f,
    0x682e_6ff3,
    0x748f_82ee,
    0x78a5_636f,
    0x84c8_7814,
    0x8cc7_0208,
    0x90be_fffa,
    0xa450_6ceb,
    0xbef9_a3f7,
    0xc671_78f2,
];

pub struct Sha256Encoding {
    pub formula: Formula,
    pub message: [Word; 16],
    output: [Word; 8],
}

impl Sha256Encoding {
    pub fn build(rounds: usize, use_xor_clauses: bool) -> Self {
        assert!((16..=64).contains(&rounds));
        let mut formula = Formula::new(use_xor_clauses);

        let mut schedule = Vec::with_capacity(rounds);
        for word in 0..rounds {
            schedule.push(formula.new_word(Some(format!("w{word}"))));
        }
        let message: [Word; 16] = schedule[..16]
            .try_into()
            .expect("SHA-256 always has 16 input words");
        let chain: [Word; 8] = std::array::from_fn(|_| formula.new_word(None));
        let output: [Word; 8] =
            std::array::from_fn(|word| formula.new_word(Some(format!("hash{word}"))));

        let mut a = vec![chain[3], chain[2], chain[1], chain[0]];
        let mut e = vec![chain[7], chain[6], chain[5], chain[4]];
        for _ in 0..rounds {
            a.push(formula.new_word(None));
            e.push(formula.new_word(None));
        }

        for word in 16..rounds {
            let small_sigma0 = formula.new_word(None);
            let small_sigma1 = formula.new_word(None);

            let rotate1 = formula.rotate_right(&schedule[word - 15], 7);
            let rotate2 = formula.rotate_right(&schedule[word - 15], 18);
            formula.xor2(&small_sigma0[29..], &rotate1[29..], &rotate2[29..]);
            formula.xor3(
                &small_sigma0[..29],
                &rotate1[..29],
                &rotate2[..29],
                &schedule[word - 15][3..],
            );

            let rotate1 = formula.rotate_right(&schedule[word - 2], 17);
            let rotate2 = formula.rotate_right(&schedule[word - 2], 19);
            formula.xor2(&small_sigma1[22..], &rotate1[22..], &rotate2[22..]);
            formula.xor3(
                &small_sigma1[..22],
                &rotate1[..22],
                &rotate2[..22],
                &schedule[word - 2][10..],
            );

            let expanded = schedule[word];
            formula.add4(
                &expanded,
                &schedule[word - 16],
                &small_sigma0,
                &schedule[word - 7],
                &small_sigma1,
            );
        }

        let mut constants = Vec::with_capacity(rounds);
        for value in ROUND_CONSTANTS.iter().copied().take(rounds) {
            let word = formula.new_word(None);
            formula.fixed_value(&word, value);
            constants.push(word);
        }

        for (word, value) in chain.iter().zip([
            0x6a09_e667,
            0xbb67_ae85,
            0x3c6e_f372,
            0xa54f_f53a,
            0x510e_527f,
            0x9b05_688c,
            0x1f83_d9ab,
            0x5be0_cd19,
        ]) {
            formula.fixed_value(word, value);
        }

        for round in 0..rounds {
            let sigma0 = formula.new_word(None);
            let sigma1 = formula.new_word(None);
            big_sigma0(&mut formula, &sigma0, &a[round + 3]);
            big_sigma1(&mut formula, &sigma1, &e[round + 3]);

            let choose = formula.new_word(None);
            let majority = formula.new_word(None);
            formula.choose(&choose, &e[round + 3], &e[round + 2], &e[round + 1]);
            formula.majority3(&majority, &a[round + 3], &a[round + 2], &a[round + 1]);

            let temporary = formula.new_word(None);
            formula.add5(
                &temporary,
                &e[round],
                &sigma1,
                &choose,
                &constants[round],
                &schedule[round],
            );
            formula.add2(&e[round + 4], &a[round], &temporary);
            formula.add3(&a[round + 4], &temporary, &sigma0, &majority);
        }

        formula.add2(&output[0], &chain[0], &a[rounds + 3]);
        formula.add2(&output[1], &chain[1], &a[rounds + 2]);
        formula.add2(&output[2], &chain[2], &a[rounds + 1]);
        formula.add2(&output[3], &chain[3], &a[rounds]);
        formula.add2(&output[4], &chain[4], &e[rounds + 3]);
        formula.add2(&output[5], &chain[5], &e[rounds + 2]);
        formula.add2(&output[6], &chain[6], &e[rounds + 1]);
        formula.add2(&output[7], &chain[7], &e[rounds]);

        Self {
            formula,
            message,
            output,
        }
    }

    pub fn fix_output(&mut self, target: [u32; 8]) {
        for (word, value) in self.output.iter().zip(target) {
            self.formula.fixed_value(word, value);
        }
    }

    pub fn fix_message_prefix(&mut self, message: &[u32; 16], bits: usize) {
        for bit in 0..bits {
            let word = bit / 32;
            let offset = bit % 32;
            let variable = self.message[word][offset];
            self.formula
                .add_clause(vec![if (message[word] >> offset) & 1 == 1 {
                    variable
                } else {
                    -variable
                }]);
        }
    }
}

fn big_sigma0(formula: &mut Formula, output: &Word, input: &Word) {
    let rotate1 = formula.rotate_right(input, 2);
    let rotate2 = formula.rotate_right(input, 13);
    let rotate3 = formula.rotate_right(input, 22);
    formula.xor3(output, &rotate1, &rotate2, &rotate3);
}

fn big_sigma1(formula: &mut Formula, output: &Word, input: &Word) {
    let rotate1 = formula.rotate_right(input, 6);
    let rotate2 = formula.rotate_right(input, 11);
    let rotate3 = formula.rotate_right(input, 25);
    formula.xor3(output, &rotate1, &rotate2, &rotate3);
}

pub fn compression(message: &[u32; 16], rounds: usize) -> [u32; 8] {
    let mut schedule = vec![0u32; rounds.max(16)];
    schedule[..16].copy_from_slice(message);
    for word in 16..rounds {
        let sigma0 = schedule[word - 15].rotate_right(7)
            ^ schedule[word - 15].rotate_right(18)
            ^ (schedule[word - 15] >> 3);
        let sigma1 = schedule[word - 2].rotate_right(17)
            ^ schedule[word - 2].rotate_right(19)
            ^ (schedule[word - 2] >> 10);
        schedule[word] = schedule[word - 16]
            .wrapping_add(sigma0)
            .wrapping_add(schedule[word - 7])
            .wrapping_add(sigma1);
    }

    let initial = [
        0x6a09_e667u32,
        0xbb67_ae85,
        0x3c6e_f372,
        0xa54f_f53a,
        0x510e_527f,
        0x9b05_688c,
        0x1f83_d9ab,
        0x5be0_cd19,
    ];
    let [mut a, mut b, mut c, mut d, mut e, mut f, mut g, mut h] = initial;
    for round in 0..rounds {
        let sigma1 = e.rotate_right(6) ^ e.rotate_right(11) ^ e.rotate_right(25);
        let choose = (e & f) ^ ((!e) & g);
        let temporary1 = h
            .wrapping_add(sigma1)
            .wrapping_add(choose)
            .wrapping_add(ROUND_CONSTANTS[round])
            .wrapping_add(schedule[round]);
        let sigma0 = a.rotate_right(2) ^ a.rotate_right(13) ^ a.rotate_right(22);
        let majority = (a & b) ^ (a & c) ^ (b & c);
        let temporary2 = sigma0.wrapping_add(majority);

        h = g;
        g = f;
        f = e;
        e = d.wrapping_add(temporary1);
        d = c;
        c = b;
        b = a;
        a = temporary1.wrapping_add(temporary2);
    }

    [
        initial[0].wrapping_add(a),
        initial[1].wrapping_add(b),
        initial[2].wrapping_add(c),
        initial[3].wrapping_add(d),
        initial[4].wrapping_add(e),
        initial[5].wrapping_add(f),
        initial[6].wrapping_add(g),
        initial[7].wrapping_add(h),
    ]
}
