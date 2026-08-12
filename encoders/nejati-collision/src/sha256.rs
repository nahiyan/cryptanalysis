use nejati_sat::{Formula, Word};

pub const ROUND_CONSTANTS: [u32; 64] = [
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
    0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
    0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
    0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
    0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
    0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
    0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
    0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
];

const ZERO: Word = [0; 32];

pub struct Block {
    pub formula: Formula,
    pub w: Vec<Word>,
    pub a: Vec<Word>,
    pub e: Vec<Word>,
    pub t: Vec<Word>,
    pub s0: Vec<Word>,
    pub s1: Vec<Word>,
    pub sigma0: Vec<Word>,
    pub sigma1: Vec<Word>,
    pub choose: Vec<Word>,
    pub majority: Vec<Word>,
    pub w_carry: Vec<Word>,
    pub w_carry_high: Vec<Word>,
    pub t_carry: Vec<Word>,
    pub t_carry_high: Vec<Word>,
    pub e_carry: Vec<Word>,
    pub a_carry: Vec<Word>,
    pub a_carry_high: Vec<Word>,
}

impl Block {
    pub fn encode(rounds: usize, initial_block: bool, suffix: &str, start: i32, xor: bool) -> Self {
        let mut formula = Formula::with_start(xor, start);
        let named_word =
            |formula: &mut Formula, name: &str| formula.new_word(Some(format!("{name}_{suffix}")));

        let w: Vec<_> = (0..rounds)
            .map(|round| named_word(&mut formula, &format!("W_{round}")))
            .collect();
        let input: [Word; 8] = std::array::from_fn(|_| formula.new_word(None));
        formula.name(format!("cv_{suffix}"), input[0][0]);
        let output: [Word; 8] =
            std::array::from_fn(|word| named_word(&mut formula, &format!("hash{word}")));

        let mut a = vec![ZERO; rounds + 4];
        let mut e = vec![ZERO; rounds + 4];
        for round in 0..rounds {
            a[round + 4] = named_word(&mut formula, &format!("A_{}", round + 4));
            e[round + 4] = named_word(&mut formula, &format!("E_{}", round + 4));
        }

        let mut s0 = vec![ZERO; rounds];
        let mut s1 = vec![ZERO; rounds];
        let mut w_carry = vec![ZERO; rounds];
        let mut w_carry_high = vec![ZERO; rounds];
        for round in 16..rounds {
            s0[round] = formula.new_word(None);
            s1[round] = formula.new_word(None);
            let rotate7 = formula.rotate_right(&w[round - 15], 7);
            let rotate18 = formula.rotate_right(&w[round - 15], 18);
            formula.name(format!("s0_{round}_{suffix}"), s0[round][0]);
            formula.xor3(
                &s0[round][..29],
                &rotate7[..29],
                &rotate18[..29],
                &w[round - 15][3..],
            );
            formula.xor2(&s0[round][29..], &rotate7[29..], &rotate18[29..]);

            let rotate17 = formula.rotate_right(&w[round - 2], 17);
            let rotate19 = formula.rotate_right(&w[round - 2], 19);
            formula.name(format!("s1_{round}_{suffix}"), s1[round][0]);
            formula.xor3(
                &s1[round][..22],
                &rotate17[..22],
                &rotate19[..22],
                &w[round - 2][10..],
            );
            formula.xor2(&s1[round][22..], &rotate17[22..], &rotate19[22..]);

            w_carry_high[round] = named_word(&mut formula, &format!("add.W.r1_{round}"));
            w_carry[round] = named_word(&mut formula, &format!("add.W.r0_{round}"));
            formula.add_with_carries(
                &w[round],
                &[&w[round - 16], &s0[round], &w[round - 7], &s1[round]],
                &w_carry[round],
                Some(&w_carry_high[round]),
            );
        }

        let constants: Vec<_> = (0..rounds)
            .map(|round| {
                let word = named_word(&mut formula, &format!("K_{round}"));
                formula.fixed_value(&word, ROUND_CONSTANTS[round]);
                word
            })
            .collect();

        if initial_block {
            for (word, value) in input.iter().zip([
                0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f, 0x9b05688c, 0x1f83d9ab,
                0x5be0cd19,
            ]) {
                formula.fixed_value(word, value);
            }
        }
        a[3] = input[0];
        a[2] = input[1];
        a[1] = input[2];
        a[0] = input[3];
        e[3] = input[4];
        e[2] = input[5];
        e[1] = input[6];
        e[0] = input[7];
        for index in 0..4 {
            formula.name(format!("A_{index}_{suffix}"), a[index][0]);
            formula.name(format!("E_{index}_{suffix}"), e[index][0]);
        }

        let mut sigma0 = vec![ZERO; rounds];
        let mut sigma1 = vec![ZERO; rounds];
        let mut choose = vec![ZERO; rounds];
        let mut majority = vec![ZERO; rounds];
        let mut t = vec![ZERO; rounds];
        let mut t_carry = vec![ZERO; rounds];
        let mut t_carry_high = vec![ZERO; rounds];
        let mut e_carry = vec![ZERO; rounds];
        let mut a_carry = vec![ZERO; rounds];
        let mut a_carry_high = vec![ZERO; rounds];

        for round in 0..rounds {
            sigma0[round] = named_word(&mut formula, &format!("sigma0_{round}"));
            sigma1[round] = named_word(&mut formula, &format!("sigma1_{round}"));
            let r2 = formula.rotate_right(&a[round + 3], 2);
            let r13 = formula.rotate_right(&a[round + 3], 13);
            let r22 = formula.rotate_right(&a[round + 3], 22);
            formula.xor3(&sigma0[round], &r2, &r13, &r22);
            let r6 = formula.rotate_right(&e[round + 3], 6);
            let r11 = formula.rotate_right(&e[round + 3], 11);
            let r25 = formula.rotate_right(&e[round + 3], 25);
            formula.xor3(&sigma1[round], &r6, &r11, &r25);

            choose[round] = named_word(&mut formula, &format!("if_{round}"));
            formula.choose(&choose[round], &e[round + 3], &e[round + 2], &e[round + 1]);
            majority[round] = named_word(&mut formula, &format!("maj_{round}"));
            formula.majority3(
                &majority[round],
                &a[round + 3],
                &a[round + 2],
                &a[round + 1],
            );

            t_carry_high[round] = named_word(&mut formula, &format!("add.T.r1_{round}"));
            t_carry[round] = named_word(&mut formula, &format!("add.T.r0_{round}"));
            t[round] = named_word(&mut formula, &format!("T_{round}"));
            formula.add_with_carries(
                &t[round],
                &[
                    &e[round],
                    &sigma1[round],
                    &choose[round],
                    &constants[round],
                    &w[round],
                ],
                &t_carry[round],
                Some(&t_carry_high[round]),
            );

            e_carry[round] = named_word(&mut formula, &format!("add.E.r0_{round}"));
            formula.add_with_carries(
                &e[round + 4],
                &[&a[round], &t[round]],
                &e_carry[round],
                None,
            );

            a_carry_high[round] = named_word(&mut formula, &format!("add.A.r1_{round}"));
            a_carry[round] = named_word(&mut formula, &format!("add.A.r0_{round}"));
            formula.add_with_carries(
                &a[round + 4],
                &[&t[round], &sigma0[round], &majority[round]],
                &a_carry[round],
                Some(&a_carry_high[round]),
            );
            formula.add_clause(vec![-a_carry_high[round][0]]);
        }

        let output_carry: [Word; 8] = std::array::from_fn(|_| formula.new_word(None));
        let final_state = [
            a[rounds + 3],
            a[rounds + 2],
            a[rounds + 1],
            a[rounds],
            e[rounds + 3],
            e[rounds + 2],
            e[rounds + 1],
            e[rounds],
        ];
        for index in 0..8 {
            formula.add_with_carries(
                &output[index],
                &[&input[index], &final_state[index]],
                &output_carry[index],
                None,
            );
        }

        Self {
            formula,
            w,
            a,
            e,
            t,
            s0,
            s1,
            sigma0,
            sigma1,
            choose,
            majority,
            w_carry,
            w_carry_high,
            t_carry,
            t_carry_high,
            e_carry,
            a_carry,
            a_carry_high,
        }
    }
}
