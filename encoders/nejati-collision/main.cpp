#include "diff-parser.h"
#include "retrieve_table.h"
#include "sha256x.h"
#include "xformula.h"
#include <assert.h>
#include <ctime>
#include <fstream>
#include <getopt.h>
#include <sstream>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <unordered_map>

using namespace std;

/* config options */
int cfg_use_xor_clauses;
Formula::MultiAdderType cfg_multi_adder_type;
int cfg_diff_desc;
int cfg_free_start;
int cfg_rand_inp_diff;
string cfg_diff_const_file;
Rules prop_rules;

unsigned rnd_const[] = {
    0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1,
    0x923f82a4, 0xab1c5ed5, 0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3,
    0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174, 0xe49b69c1, 0xefbe4786,
    0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
    0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147,
    0x06ca6351, 0x14292967, 0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13,
    0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85, 0xa2bfe8a1, 0xa81a664b,
    0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
    0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a,
    0x5b9cca4f, 0x682e6ff3, 0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208,
    0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2
};

void load_rules(string filePath)
{
    // Open the file
    ifstream inputFile(filePath);

    // Check if the file is opened successfully
    if (!inputFile.is_open()) {
        printf("Error opening the rules file\n");
        exit(1);
    }

    // Read the file line by line
    string line;
    while (getline(inputFile, line)) {
        istringstream iss(line);
        string id, input, output;
        iss >> id >> input >> output;
        unordered_map<string, string>* rules;
        if (id == "ch")
            rules = &prop_rules.ch;
        else if (id == "maj")
            rules = &prop_rules.maj;
        else if (id == "xor3")
            rules = &prop_rules.xor3;
        else if (id == "add2")
            rules = &prop_rules.add2;
        else if (id == "add3")
            rules = &prop_rules.add3;
        else if (id == "add4")
            rules = &prop_rules.add4;
        else if (id == "add5")
            rules = &prop_rules.add5;
        else if (id == "add6")
            rules = &prop_rules.add6;
        else if (id == "add7")
            rules = &prop_rules.add7;
        rules->insert({ input, output });
    }

    // Close the file
    inputFile.close();
}

void fix_4bit_starting_point(SHA256& block, char& diff, int* dx, int* x, int* x_)
{
    auto& formula = block.cnf;
    switch (diff) {
    case '-':
        formula.fixedValue(dx + 1, 0, 1);
        formula.fixedValue(dx + 2, 0, 1);
        break;
    case 'x':
        formula.fixedValue(dx + 0, 0, 1);
        formula.fixedValue(dx + 3, 0, 1);
        break;
    case '0':
        formula.fixedValue(x, 0, 1);
        formula.fixedValue(x_, 0, 1);
        formula.fixedValue(dx + 0, 1, 1);
        formula.fixedValue(dx + 1, 0, 1);
        formula.fixedValue(dx + 2, 0, 1);
        formula.fixedValue(dx + 3, 0, 1);
        break;
    case 'u':
        formula.fixedValue(x, 1, 1);
        formula.fixedValue(x_, 0, 1);
        formula.fixedValue(dx + 0, 0, 1);
        formula.fixedValue(dx + 1, 1, 1);
        formula.fixedValue(dx + 2, 0, 1);
        formula.fixedValue(dx + 3, 0, 1);
        break;
    case 'n':
        formula.fixedValue(x, 0, 1);
        formula.fixedValue(x_, 1, 1);
        formula.fixedValue(dx + 0, 0, 1);
        formula.fixedValue(dx + 1, 0, 1);
        formula.fixedValue(dx + 2, 1, 1);
        formula.fixedValue(dx + 3, 0, 1);
        break;
    case '1':
        formula.fixedValue(x, 1, 1);
        formula.fixedValue(x_, 1, 1);
        formula.fixedValue(dx + 0, 0, 1);
        formula.fixedValue(dx + 1, 0, 1);
        formula.fixedValue(dx + 2, 0, 1);
        formula.fixedValue(dx + 3, 1, 1);
        break;
    }
}

void fix_1bit_starting_point(SHA256& block, char& diff, int* dx, int* x, int* x_)
{
    auto& formula = block.cnf;
    switch (diff) {
    case '-':
        formula.fixedValue(dx, 0, 1);
        break;
    case 'x':
        formula.fixedValue(dx, 1, 1);
        break;
    case 'u':
        formula.fixedValue(x, 1, 1);
        formula.fixedValue(x_, 0, 1);
        formula.fixedValue(dx, 1, 1);
        break;
    case 'n':
        formula.fixedValue(x, 0, 1);
        formula.fixedValue(x_, 1, 1);
        formula.fixedValue(dx, 1, 1);
        break;
    case '1':
        formula.fixedValue(x, 1, 1);
        formula.fixedValue(x_, 1, 1);
        formula.fixedValue(dx, 0, 1);
        break;
    case '0':
        formula.fixedValue(x, 0, 1);
        formula.fixedValue(x_, 0, 1);
        formula.fixedValue(dx, 0, 1);
        break;
    }
}

void fix_starting_point(SHA256& block, char& diff, int* dx, int* x, int* x_)
{
#if IS_4bit
    fix_4bit_starting_point(block, diff, dx, x, x_);
#else
    fix_1bit_starting_point(block, diff, dx, x, x_);
#endif
}

void collision(int rounds)
{
    load_rules("prop_rules.db");

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
    if (cfg_free_start) {
        f.initialBlock = false;
        g.initialBlock = false;
    }

    f.encode();

    g.cnf.setVarID(f.cnf.getVarCnt());
    g.encode();
    g.cnf.AddFormula(f.cnf);

    assert(cfg_diff_desc);
    /* Differential Path Variables */
    int DA[70][32][4], DE[70][32][4], DW[70][32][4];
    int Ds0[64][32][4], Ds1[64][32][4];
    int Dwcarry[64][32][4], DwCarry[64][32][4];
    int Dsigma0[64][32][4], Dsigma1[64][32][4];
    int Df1[64][32][4], Df2[64][32][4];
    int DT[70][32][4];
    int Dr0carry[64][32][4], Dr0Carry[64][32][4];
    int DK[64][32][4];
    int Dr1carry[64][32][4];
    int Dr2carry[64][32][4], Dr2Carry[64][32][4];
    for (int i = 0; i < rounds + 4; i++) {
        g.cnf.newDiff(DA[i], "DA_" + to_string(i));
        g.cnf.newDiff(DE[i], "DE_" + to_string(i));
        g.cnf.basic_rules(DA[i], f.A[i], g.A[i]);
        g.cnf.basic_rules(DE[i], f.E[i], g.E[i]);
        if (i < rounds) {
            g.cnf.newDiff(DW[i], "DW_" + to_string(i));
            g.cnf.basic_rules(DW[i], f.w[i], g.w[i]);
        }
    }

#if IS_4bit
    int zero[6]; // GC '0'
    g.cnf.newVars(zero, 6, "zero");
    g.cnf.fixedValue(&zero[0], 0, 1);
    g.cnf.fixedValue(&zero[1], 0, 1);
    g.cnf.fixedValue(&zero[2], 1, 1);
    g.cnf.fixedValue(&zero[3], 0, 1);
    g.cnf.fixedValue(&zero[4], 0, 1);
    g.cnf.fixedValue(&zero[5], 0, 1);
#else
    int zero[3]; // GC '0'
    g.cnf.newVars(zero, 3, "zero");
    g.cnf.fixedValue(&zero[0], 0, 1);
    g.cnf.fixedValue(&zero[1], 0, 1);
    g.cnf.fixedValue(&zero[2], 0, 1);
#endif

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
        if (i >= 0)
            for (int j = 0; j < 32; j++)
                fix_starting_point(g, W[i][31 - j], DW[i][j], &f.w[i][j], &g.w[i][j]);

        for (int j = 0; j < 32; j++) {
            fix_starting_point(g, A[i + 4][31 - j], DA[i + 4][j], &f.A[i + 4][j], &g.A[i + 4][j]);
            fix_starting_point(g, E[i + 4][31 - j], DE[i + 4][j], &f.E[i + 4][j], &g.E[i + 4][j]);
        }
    }

    /* Differential propagation over message expansion */
    for (int i = 16; i < rounds; i++) {
        g.cnf.newDiff(Ds0[i], "Ds0_" + to_string(i));
        g.cnf.newDiff(Ds1[i], "Ds1_" + to_string(i));
        g.cnf.basic_rules(Ds0[i], f.s0[i], g.s0[i]);
        g.cnf.basic_rules(Ds1[i], f.s1[i], g.s1[i]);

        // s0 = (w[i-15] >>> 7) XOR (w[i-15] >>> 18) XOR (w[i-15] >> 3)
        {
            int inputs[3][32][4];
            for (int j = 0; j < 32; j++) {
                // Perform rotations and shifts
                for (int k = 0; k < 4; k++) {
                    inputs[0][j][k] = DW[i - 15][(j + 7) % 32][k];
                    inputs[1][j][k] = DW[i - 15][(j + 18) % 32][k];
                    if (j < 29)
                        inputs[2][j][k] = DW[i - 15][(j + 3) % 32][k];
                    else
                        inputs[2][j][k] = zero[2 + k];
                }
            }
            // Add XOR3 difference rules
            for (auto& entry : prop_rules.xor3) {
                bool skip = false;
                for (auto& c : entry.first)
                    if (c != '-' && c != 'x')
                        skip = true;
                if (skip)
                    continue;
                g.cnf.impose_rule({ &inputs[0], &inputs[1], &inputs[2] }, { &Ds0[i] }, entry);
            }
        }

        // s1 = (w[i-2] >>> 17) XOR (w[i-2] >>> 19) XOR (w[i-2] >> 10)
        {
            int inputs[3][32][4];
            for (int j = 0; j < 32; j++) {
                // Perform rotations and shifts
                for (int k = 0; k < 4; k++) {
                    inputs[0][j][k] = DW[i - 2][(j + 17) % 32][k];
                    inputs[1][j][k] = DW[i - 2][(j + 19) % 32][k];
                    if (j < 22)
                        inputs[2][j][k] = DW[i - 2][(j + 10) % 32][k];
                    else
                        inputs[2][j][k] = zero[2 + k];
                }
            }
            // Add XOR3 difference rules
            for (auto& entry : prop_rules.xor3) {
                bool skip = false;
                for (auto& c : entry.first)
                    if (c != '-' && c != 'x')
                        skip = true;
                if (skip)
                    continue;
                g.cnf.impose_rule({ &inputs[0], &inputs[1], &inputs[2] }, { &Ds1[i] }, entry);
            }
        }

        // Addition: w[i] = w[i-16] + s0 + w[i-7] + s1
        g.cnf.newDiff(DwCarry[i], "Dadd.W.r1_" + to_string(i));
        g.cnf.newDiff(Dwcarry[i], "Dadd.W.r0_" + to_string(i));
        g.cnf.basic_rules(Dwcarry[i], f.wcarry[i], g.wcarry[i]);
        g.cnf.basic_rules(DwCarry[i], f.wCarry[i], g.wCarry[i]);
        g.cnf.diff_add(prop_rules, DW[i], DW[i - 16], Ds0[i], Dwcarry[i], DwCarry[i], DW[i - 7], Ds1[i]);
    }

    /* Differential propagation for round function */
    for (int i = 0; i < rounds; i++) {
        // sigma0 = Sigma0(A[i+3])
        // sigma1 = Sigma1(E[i+3])
        g.cnf.newDiff(Dsigma0[i], "Dsigma0_" + to_string(i));
        g.cnf.newDiff(Dsigma1[i], "Dsigma1_" + to_string(i));
        g.cnf.basic_rules(Dsigma0[i], f.sigma0[i], g.sigma0[i]);
        g.cnf.basic_rules(Dsigma1[i], f.sigma1[i], g.sigma1[i]);

        // g.Sigma0(Dsigma0[i], DA[i + 3]);
        {
            int inputs[3][32][4];
            for (int j = 0; j < 32; j++) {
                // Perform rotations
                for (int k = 0; k < 4; k++) {
                    inputs[0][j][k] = DA[i + 3][(j + 2) % 32][k];
                    inputs[1][j][k] = DA[i + 3][(j + 13) % 32][k];
                    inputs[2][j][k] = DA[i + 3][(j + 22) % 32][k];
                }
            }
            // Add XOR3 difference rules
            for (auto& entry : prop_rules.xor3) {
                bool skip = false;
                for (auto& c : entry.first)
                    if (c != '-' && c != 'x')
                        skip = true;
                if (skip)
                    continue;
                g.cnf.impose_rule({ &inputs[0], &inputs[1], &inputs[2] }, { &Dsigma0[i] }, entry);
            }
        }

        // g.Sigma1(Dsigma1[i], DE[i + 3]);
        {
            int inputs[3][32][4];
            for (int j = 0; j < 32; j++) {
                // Perform rotations
                for (int k = 0; k < 4; k++) {
                    inputs[0][j][k] = DE[i + 3][(j + 6) % 32][k];
                    inputs[1][j][k] = DE[i + 3][(j + 11) % 32][k];
                    inputs[2][j][k] = DE[i + 3][(j + 25) % 32][k];
                }
            }
            // Add XOR3 difference rules
            for (auto& entry : prop_rules.xor3) {
                bool skip = false;
                for (auto& c : entry.first)
                    if (c != '-' && c != 'x')
                        skip = true;
                if (skip)
                    continue;
                g.cnf.impose_rule({ &inputs[0], &inputs[1], &inputs[2] }, { &Dsigma1[i] }, entry);
            }
        }

#if !IS_4bit
        {
            int output[32], input[32];
            for (int x = 0; x < 32; x++) {
                output[x] = Dsigma0[i][x][0];
                input[x] = DA[i + 3][x][0];
            }

            g.Sigma0(output, input);
        }
        {
            int output[32], input[32];
            for (int x = 0; x < 32; x++) {
                output[x] = Dsigma1[i][x][0];
                input[x] = DE[i + 3][x][0];
            }

            g.Sigma1(output, input);
        }
#endif

        // f1 = IF(E[i+3], E[i+2], E[i+1])
        g.cnf.newDiff(Df1[i], "Dif_" + to_string(i));
        g.cnf.basic_rules(Df1[i], f.f1[i], g.f1[i]);
        // Add IF difference rules
        for (auto& entry : prop_rules.ch) {
            bool skip = false;
            for (auto& c : entry.first)
                if (c != '-' && c != 'x')
                    skip = true;
            if (skip)
                continue;
            g.cnf.impose_rule({ &DE[i + 3], &DE[i + 2], &DE[i + 1] }, { &Df1[i] }, entry);
        }

        // f2 = MAJ(A[i+3], A[i+2], A[i+1])
        g.cnf.newDiff(Df2[i], "Dmaj_" + to_string(i));
        g.cnf.basic_rules(Df2[i], f.f2[i], g.f2[i]);
        // Add MAJ difference rules
        for (auto& entry : prop_rules.maj) {
            bool skip = false;
            for (auto& c : entry.first)
                if (c != '-' && c != 'x')
                    skip = true;
            if (skip)
                continue;
            g.cnf.impose_rule({ &DA[i + 3], &DA[i + 2], &DA[i + 1] }, { &Df2[i] }, entry);
        }

        // Addition: T = E[i] + sigma1 + f1 + k[i] + w[i]
        g.cnf.newDiff(Dr0Carry[i], "Dadd.T.r1_" + to_string(i));
        g.cnf.newDiff(Dr0carry[i], "Dadd.T.r0_" + to_string(i));
        g.cnf.newDiff(DT[i], "DT_" + to_string(i));
        g.cnf.basic_rules(Dr0Carry[i], f.r0Carry[i], g.r0Carry[i]);
        g.cnf.basic_rules(Dr0carry[i], f.r0carry[i], g.r0carry[i]);
        g.cnf.basic_rules(DT[i], f.T[i], g.T[i]);
        g.cnf.newDiff(DK[i], "DK_" + to_string(i));

        // Fix the differences
        for (int j = 0; j < 32; j++) {
#if IS_4bit
            bool is_true = rnd_const[i] >> j & 1;
            if (is_true) {
                g.cnf.addClause({ -DK[i][j][0] });
                g.cnf.addClause({ -DK[i][j][1] });
                g.cnf.addClause({ -DK[i][j][2] });
                g.cnf.addClause({ DK[i][j][3] });
            } else {
                g.cnf.addClause({ DK[i][j][0] });
                g.cnf.addClause({ -DK[i][j][1] });
                g.cnf.addClause({ -DK[i][j][2] });
                g.cnf.addClause({ -DK[i][j][3] });
            }
#else
            g.cnf.addClause({ -DK[i][j][0] });
#endif
        }
        g.cnf.diff_add(prop_rules, DT[i], DE[i], Dsigma1[i], Dr0carry[i], Dr0Carry[i], Df1[i], DK[i], DW[i]);

        // Addition: E[i+4] = A[i] + T
        g.cnf.newDiff(Dr1carry[i], "Dadd.E.r0_" + to_string(i));
        g.cnf.basic_rules(Dr1carry[i], f.r1carry[i], g.r1carry[i]);
        g.cnf.diff_add(prop_rules, DE[i + 4], DA[i], DT[i], Dr1carry[i]);

        // Addition: A[i+4] = T + sigma0 + f2
        g.cnf.newDiff(Dr2Carry[i], "Dadd.A.r1_" + to_string(i));
        g.cnf.newDiff(Dr2carry[i], "Dadd.A.r0_" + to_string(i));
        g.cnf.basic_rules(Dr2carry[i], f.r2carry[i], g.r2carry[i]);
        g.cnf.basic_rules(Dr2Carry[i], f.r2Carry[i], g.r2Carry[i]);
        g.cnf.diff_add(prop_rules, DA[i + 4], DT[i], Dsigma0[i], Dr2carry[i], Dr2Carry[i], Df2[i]);
    }

    g.cnf.dimacs(rounds);
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
           "  --free_start                             Free up the chaining value"
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
    cfg_free_start = 0;
    cfg_rand_inp_diff = 0;
    cfg_diff_const_file = "";
    int rounds = -1;

    struct option long_options[] = {
        /* flag options */
        { "xor", no_argument, &cfg_use_xor_clauses, 1 },
        { "diff_desc", no_argument, &cfg_diff_desc, 1 },
        { "free_start", no_argument, &cfg_free_start, 1 },
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
