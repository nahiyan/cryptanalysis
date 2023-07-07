#include "diff-parser.h"
#include "retrieve_table.h"
#include "sha256x.h"
#include "util.h"
#include <assert.h>
#include <ctime>
#include <getopt.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
using namespace std;

/* config options */
int cfg_use_xor_clauses;
Formula::MultiAdderType cfg_multi_adder_type;
int cfg_diff_desc;
int cfg_diff_impl;
int cfg_rand_inp_diff;
string cfg_diff_const_file;

void diff_xor2(Formula& f, int r, int a, int b) { }

void diff_xor3(Formula& f, int r, int a, int b, int c) { }

void diff_xor4(Formula& f, int r, int a, int b, int c, int d) { }

void collision(int rounds)
{
    SHA256 f(rounds), g(rounds);
    f.cnf.formulaName = "f";
    g.cnf.formulaName = "g";
    if (cfg_use_xor_clauses) {
        f.cnf.setUseXORClauses();
        g.cnf.setUseXORClauses();
    }
    if (cfg_multi_adder_type != Formula::MAT_NONE) {
        f.cnf.setMultiAdderType(cfg_multi_adder_type);
        g.cnf.setMultiAdderType(cfg_multi_adder_type);
    }

    f.encode();

    g.cnf.setVarID(f.cnf.getVarCnt());
    g.encode();
    g.cnf.AddFormula(f.cnf);

    if (cfg_diff_desc) {
        /* Differential Path Variables */
        int DA[70][32][4], DE[70][32][4], DW[70][32][4];
        for (int i = 0; i < rounds + 4; i++) {
            g.cnf.newVarsD2(&DA[i][0][0], 32, 4, "DA_" + to_string(i));
            g.cnf.newVarsD2(&DE[i][0][0], 32, 4, "DE_" + to_string(i));
            g.cnf.diffVar(&DA[i][0][0], f.A[i], g.A[i]);
            g.cnf.diffVar(&DE[i][0][0], f.E[i], g.E[i]);
            if (i < rounds) {
                g.cnf.newVarsD2(&DW[i][0][0], 32, 4, "DW_" + to_string(i));
                g.cnf.diffVar(&DW[i][0][0], f.w[i], g.w[i]);
            }
        }

        // Support for built-in differential characters
        vector<string> A, E, W;
        if (cfg_diff_const_file == "HARD_CODED") {
            retrieve_table(rounds, A, E, W);
        } else {
            FILE* diff_const_file = fopen(cfg_diff_const_file.c_str(), "r");
            parse_diff_path(rounds, diff_const_file, A, E, W);
            fclose(diff_const_file);
        }

        /* Fixing the differences from the initial path */
        for (int i = -4; i < rounds; i++) {
            if (i >= 0) {
                for (int j = 0; j < 32; j++)
                    if (W[i][31 - j] == '-') {
                        g.cnf.fixedValue(&DW[i][j][0], 1, 1);
                        g.cnf.fixedValue(&DW[i][j][1], 0, 1);
                        g.cnf.fixedValue(&DW[i][j][2], 0, 1);
                        g.cnf.fixedValue(&DW[i][j][3], 1, 1);
                    } else if (W[i][31 - j] == 'x') {
                        g.cnf.fixedValue(&DW[i][j][0], 0, 1);
                        g.cnf.fixedValue(&DW[i][j][1], 1, 1);
                        g.cnf.fixedValue(&DW[i][j][2], 1, 1);
                        g.cnf.fixedValue(&DW[i][j][3], 0, 1);
                    }
            }
            for (int j = 0; j < 32; j++) {
                if (A[i + 4][31 - j] == '-') {
                    g.cnf.fixedValue(&DA[i + 4][j][0], 1, 1);
                    g.cnf.fixedValue(&DA[i + 4][j][1], 0, 1);
                    g.cnf.fixedValue(&DA[i + 4][j][2], 0, 1);
                    g.cnf.fixedValue(&DA[i + 4][j][3], 1, 1);
                } else if (A[i + 4][31 - j] == 'x') {
                    g.cnf.fixedValue(&DA[i + 4][j][0], 0, 1);
                    g.cnf.fixedValue(&DA[i + 4][j][1], 1, 1);
                    g.cnf.fixedValue(&DA[i + 4][j][2], 1, 1);
                    g.cnf.fixedValue(&DA[i + 4][j][3], 0, 1);
                }

                if (E[i + 4][31 - j] == '-') {
                    g.cnf.fixedValue(&DE[i + 4][j][0], 1, 1);
                    g.cnf.fixedValue(&DE[i + 4][j][1], 0, 1);
                    g.cnf.fixedValue(&DE[i + 4][j][2], 0, 1);
                    g.cnf.fixedValue(&DE[i + 4][j][3], 1, 1);
                } else if (E[i + 4][31 - j] == 'x') {
                    g.cnf.fixedValue(&DE[i + 4][j][0], 0, 1);
                    g.cnf.fixedValue(&DE[i + 4][j][1], 1, 1);
                    g.cnf.fixedValue(&DE[i + 4][j][2], 1, 1);
                    g.cnf.fixedValue(&DE[i + 4][j][3], 0, 1);
                }
            }
        }

        /* Differential propagation over message expansion */
        int Ds0[64][32][4], Ds1[64][32][4];
        int Dwcarry[64][32][4], DwCarry[64][32][4];
        for (int i = 16; i < rounds; i++) {
            g.cnf.newVarsD2(&Ds0[i][0][0], 32, 4, "Ds0_" + to_string(i));
            g.cnf.newVarsD2(&Ds1[i][0][0], 32, 4, "Ds1_" + to_string(i));
            g.cnf.diffVar(&Ds0[i][0][0], f.s0[i], g.s0[i]);
            g.cnf.diffVar(&Ds1[i][0][0], f.s1[i], g.s1[i]);

            // s0 = (w[i-15] >>> 7) XOR (w[i-15] >>> 18) XOR (w[i-15] >> 3)
            // int r1[32], r2[32];
            // g.cnf.rotr(r1, DW[i - 15], 7);
            // g.cnf.rotr(r2, DW[i - 15], 18);
            // g.cnf.xor2(Ds0[i] + 29, r1 + 29, r2 + 29, 3);
            // g.cnf.xor3(Ds0[i], r1, r2, DW[i - 15] + 3, 29);

            // s1 = (w[i-2] >>> 17) XOR (w[i-2] >>> 19) XOR (w[i-2] >> 10)
            // g.cnf.rotr(r1, DW[i - 2], 17);
            // g.cnf.rotr(r2, DW[i - 2], 19);
            // g.cnf.xor2(Ds1[i] + 22, r1 + 22, r2 + 22, 10);
            // g.cnf.xor3(Ds1[i], r1, r2, DW[i - 2] + 10, 22);

            g.cnf.newVarsD2(&Dwcarry[i][0][0], 32, 4, "Dw_carry_" + to_string(i));
            g.cnf.newVarsD2(&DwCarry[i][0][0], 32, 4, "Dw_Carry_" + to_string(i));
            g.cnf.diffVar(&Dwcarry[i][0][0], f.wcarry[i], g.wcarry[i]);
            g.cnf.diffVar(&DwCarry[i][0][0], f.wCarry[i], g.wCarry[i]);

            // w[i] = w[i-16] + s0 + w[i-7] + s1
            // g.cnf.diff_add(DW[i], DW[i - 16], Ds0[i], Dwcarry[i], DwCarry[i],
            //                DW[i - 7], Ds1[i]);
        }

        /* Differential propagation for round function */
        int Dsigma0[64][32][4], Dsigma1[64][32][4];
        int Df1[64][32][4], Df2[64][32][4];
        int DT[70][32][4];
        int Dr0carry[64][32][4], Dr0Carry[64][32][4];
        int DK[64][32][4];
        int Dr1carry[64][32][4];
        int Dr2carry[64][32][4], Dr2Carry[64][32][4];

        for (int i = 0; i < rounds; i++) {
            // sigma0 = Sigma0(A[i+3])
            // sigma1 = Sigma1(E[i+3])
            g.cnf.newVarsD2(&Dsigma0[i][0][0], 32, 4, "Dsigma0_" + to_string(i));
            g.cnf.newVarsD2(&Dsigma1[i][0][0], 32, 4, "Dsigma1_" + to_string(i));
            g.cnf.diffVar(&Dsigma0[i][0][0], f.sigma0[i], g.sigma0[i]);
            g.cnf.diffVar(&Dsigma1[i][0][0], f.sigma1[i], g.sigma1[i]);

            // g.Sigma0(Dsigma0[i], DA[i + 3]);
            // g.Sigma1(Dsigma1[i], DE[i + 3]);

            // f1 = IF(E[i+3], E[i+2], E[i+1])
            // f2 = MAJ(A[i+3], A[i+2], A[i+1])
            g.cnf.newVarsD2(&Df1[i][0][0], 32, 4, "Df1_" + to_string(i));
            g.cnf.newVarsD2(&Df2[i][0][0], 32, 4, "Df2_" + to_string(i));
            g.cnf.diffVar(&Df1[i][0][0], f.f1[i], g.f1[i]);
            g.cnf.diffVar(&Df2[i][0][0], f.f2[i], g.f2[i]);

            // if (cfg_diff_impl) {
            //   for (int j = 0; j < 32; j++) {
            //     // MAJ: xxx -> x
            //     g.cnf.addClause(
            //         {-DA[i + 3][j], -DA[i + 2][j], -DA[i + 1][j], Df2[i][j]});

            //     // MAJ: --- -> -
            //     g.cnf.addClause(
            //         {DA[i + 3][j], DA[i + 2][j], DA[i + 1][j], -Df2[i][j]});

            //     // IF: -xx -> x
            //     g.cnf.addClause(
            //         {DE[i + 3][j], -DE[i + 2][j], -DE[i + 1][j], Df1[i][j]});

            //     // IF: --- -> -
            //     g.cnf.addClause(
            //         {DE[i + 3][j], DE[i + 2][j], DE[i + 1][j], -Df1[i][j]});
            //   }
            // }

            // T = E[i] + sigma1 + f1 + k[i] + w[i]
            g.cnf.newVarsD2(&DT[i][0][0], 32, 4, "DT_" + to_string(i));
            g.cnf.diffVar(&DT[i][0][0], f.T[i], g.T[i]);
            g.cnf.newVarsD2(&Dr0carry[i][0][0], 32, 4, "Dr0_carry_" + to_string(i));
            g.cnf.newVarsD2(&Dr0Carry[i][0][0], 32, 4, "Dr0_Carry_" + to_string(i));
            g.cnf.diffVar(&Dr0carry[i][0][0], f.r0carry[i], g.r0carry[i]);
            g.cnf.diffVar(&Dr0Carry[i][0][0], f.r0Carry[i], g.r0Carry[i]);
            g.cnf.newVarsD2(&DK[i][0][0], 32, 4, "DK_" + to_string(i));
            g.cnf.fixedValueD2(&DK[i][0][0], 0, 32, 4);

            // g.cnf.diff_add(DT[i], DE[i], Dsigma1[i], Dr0carry[i], Dr0Carry[i], Df1[i],
            //  DK[i], DW[i]);

            // E[i+4] = A[i] + T
            g.cnf.newVarsD2(&Dr1carry[i][0][0], 32, 4, "Dr1_carry_" + to_string(i));
            g.cnf.diffVar(&Dr1carry[i][0][0], f.r1carry[i], g.r1carry[i]);

            // g.cnf.diff_add(DE[i + 4], DA[i], DT[i], Dr1carry[i]);

            // A[i+4] = T + sigma0 + f2
            g.cnf.newVarsD2(&Dr2carry[i][0][0], 32, 4, "Dr2_carry_" + to_string(i));
            g.cnf.newVarsD2(&Dr2Carry[i][1][0], 31, 4, "Dr2_Carry_" + to_string(i));
            g.cnf.diffVar(&Dr2carry[i][0][0], f.r2carry[i], g.r2carry[i]);
            g.cnf.diffVar(&Dr2Carry[i][1][0], f.r2Carry[i] + 1, g.r2Carry[i] + 1, 31);

            // g.cnf.diff_add(DA[i + 4], DT[i], Dsigma0[i], Dr2carry[i], Dr2Carry[i],
            //  Df2[i]);
        }
    } else {
        /* Inputs should be different */
        if (cfg_rand_inp_diff) {
            /* Force the inputs to differ at a random bit */
            int idx = rand() % 512;
            int r = idx / 32;
            int s = idx % 32;
            g.cnf.neq(&f.w[r][s], &g.w[r][s], 1);
        } else {
            /* Force the inputs to be different */
            int tmp[16][32];
            for (int i = 0; i < 16; i++) {
                g.cnf.newVars(tmp[i], 32);
                g.cnf.xor2(tmp[i], f.w[i], g.w[i]);
            }
            vector<int> v;
            for (int i = 0; i < 16; i++)
                for (int j = 0; j < 32; j++)
                    v.push_back(tmp[i][j]);
            g.cnf.addClause(v);
        }

        /* Outputs should be the same */
        for (int i = 0; i < 8; i++)
            g.cnf.eq(f.out[i], g.out[i]);
    }

    g.cnf.dimacs();
}

void display_usage()
{
    printf("USAGE: ./main {number_of_rounds}\n"
           "  --help or -h                             Prints this message\n"
           "  --xor                                    Use XOR clauses (default: "
           "off)\n"
           "  --adder_type or -A {two_operand | counter_chain | espresso | "
           "dot_matrix}\n"
           "                                           Specifies the type of "
           "multi operand addition encoding (default: espresso)\n"
           "  --rounds or -r {int(16..80)}             Number of rounds in your "
           "function\n"
           "  --diff_desc                              Adds differential "
           "description\n"
           "  --diff_impl                              Adds differential "
           "implication for MAJ, IF, XOR2, and XOR3\n"
           "  --diff_const_file or -d {file_path}      Path to the differential "
           "constraints file\n"
           "  --rand_input_diff                        Randomly pick a bit in "
           "input to be different (for collision)\n");
}

int main(int argc, char** argv)
{
    unsigned long seed = time(NULL);

    /* Arguments default values */
    cfg_use_xor_clauses = 0;
    cfg_multi_adder_type = Formula::MAT_NONE;
    cfg_diff_desc = 0;
    cfg_diff_impl = 0;
    cfg_rand_inp_diff = 0;
    cfg_diff_const_file = "";
    int rounds = -1;

    struct option long_options[] = {
        /* flag options */
        { "xor", no_argument, &cfg_use_xor_clauses, 1 },
        { "diff_desc", no_argument, &cfg_diff_desc, 1 },
        { "diff_impl", no_argument, &cfg_diff_impl, 1 },
        { "rand_input_diff", no_argument, &cfg_rand_inp_diff, 1 },
        /* valued options */
        { "rounds", required_argument, 0, 'r' },
        { "adder_type", required_argument, 0, 'A' },
        { "diff_const_file", required_argument, 0, 'd' },
        { "help", no_argument, 0, 'h' },
        { 0, 0, 0, 0 }
    };

    /* Process command line */
    int c, option_index;
    while ((c = getopt_long(argc, argv, "r:A:d:h", long_options,
                &option_index))
        != -1) {
        switch (c) {
        case 0:
            /* If this option set a flag, do nothing else now. */
            if (long_options[option_index].flag != 0)
                break;
            printf("option %s", long_options[option_index].name);
            if (optarg)
                printf(" with arg %s", optarg);
            printf("\n");
            break;

        case 'r':
            rounds = atoi(optarg);
            break;

        case 'A':
            cfg_multi_adder_type = strcmp(optarg, "two_operand") == 0 ? Formula::TWO_OPERAND
                : strcmp(optarg, "counter_chain") == 0                ? Formula::COUNTER_CHAIN
                : strcmp(optarg, "espresso") == 0                     ? Formula::ESPRESSO
                : strcmp(optarg, "dot_matrix") == 0                   ? Formula::DOT_MATRIX
                                                                      : Formula::MAT_NONE;
            if (cfg_multi_adder_type == Formula::MAT_NONE) {
                fprintf(stderr, "Invalid or missing multi-adder type!\nUse -h to see "
                                "the optionsi\n");
                return 1;
            }
            break;

        case 'd':
            cfg_diff_const_file = string(optarg);
            break;

        case 'h':
            display_usage();
            return 1;

        case '?':
            return 1;

        default:
            abort();
        }
    }

    /* Check for argument consistency */
    if (cfg_diff_desc && cfg_diff_const_file == "") {
        fprintf(stderr, "Differential description flag is set, but no path to the "
                        "differential constraints file is provided!\n");
        return 1;
    }

    if (rounds == -1) {
        fprintf(stderr, "Number of rounds is required! Use -r or --rounds\n");
        return 1;
    }

    srand(seed);
    srand48(rand());

    collision(rounds);

    return 0;
}
