#include <algorithm>
#include <cstdio>
#include <iostream>
#include <map>
#include <numeric>
#include <set>
#include <string>
#include <vector>

using namespace std;

map<char, vector<char>> symbols = { { '?', { 'u', 'n', '1', '0' } },
    { '-', { '1', '0' } },
    { 'x', { 'u', 'n' } },
    { '0', { '0' } },
    { 'u', { 'u' } },
    { 'n', { 'n' } },
    { '1', { '1' } },
    { '3', { '0', 'u' } },
    { '5', { '0', 'n' } },
    { '7', { '0', 'u', 'n' } },
    { 'A', { 'u', '1' } },
    { 'B', { '1', 'u', '0' } },
    { 'C', { 'n', '1' } },
    { 'D', { '0', 'n', '1' } },
    { 'E', { 'u', 'n', '1' } } };
map<char, set<char>> symbols_set = { { '?', { 'u', 'n', '1', '0' } },
    { '-', { '1', '0' } },
    { 'x', { 'u', 'n' } },
    { '0', { '0' } },
    { 'u', { 'u' } },
    { 'n', { 'n' } },
    { '1', { '1' } },
    { '3', { '0', 'u' } },
    { '5', { '0', 'n' } },
    { '7', { '0', 'u', 'n' } },
    { 'A', { 'u', '1' } },
    { 'B', { '1', 'u', '0' } },
    { 'C', { 'n', '1' } },
    { 'D', { '0', 'n', '1' } },
    { 'E', { 'u', 'n', '1' } } };

// Function to calculate the Cartesian product of multiple vectors of
// characters
vector<string> cartesian_product(vector<vector<char>> input)
{
    vector<string> result;
    int numVectors = input.size();
    vector<int> indices(numVectors, 0);

    while (true) {
        string currentProduct;
        for (int i = 0; i < numVectors; ++i)
            currentProduct.push_back(input[i][indices[i]]);

        result.push_back(currentProduct);

        int j = numVectors - 1;
        while (j >= 0 && indices[j] == int(input[j].size()) - 1) {
            indices[j] = 0;
            j--;
        }

        if (j < 0)
            break;

        indices[j]++;
    }

    return result;
}

vector<string> cartesian_product(vector<char> input, int repeat)
{
    vector<vector<char>> inputs;
    for (int i = 0; i < repeat; i++) {
        inputs.push_back(input);
    }
    return cartesian_product(inputs);
}

vector<int> add(vector<int> inputs)
{
    int sum = accumulate(inputs.begin(), inputs.end(), 0);
    return { sum >> 2 & 1, sum >> 1 & 1, sum & 1 };
}
vector<int> ch(vector<int> inputs)
{
    int x = inputs[0], y = inputs[1], z = inputs[2];
    return { (x & y) ^ (x & z) ^ z };
}
vector<int> maj(vector<int> inputs)
{
    int x = inputs[0], y = inputs[1], z = inputs[2];
    return { (x & y) ^ (y & z) ^ (x & z) };
}

vector<int> xor_(vector<int> inputs)
{
    int value = 0;
    for (auto& input: inputs) {
        value ^= input;
    }
    return { value };
}

string propagate(vector<int> (*func)(vector<int> inputs), string inputs, int outputs_size = 1)
{
    auto conforms_to = [](char c1, char c2) {
        auto c1_chars = symbols[c1], c2_chars = symbols[c2];
        for (auto& c : c1_chars)
            if (find(c2_chars.begin(), c2_chars.end(), c) == c2_chars.end())
                return false;
        return true;
    };

    vector<vector<char>> iterables_list;
    for (auto& input : inputs) {
        auto it = symbols.find(input);
        if (it != symbols.end())
            iterables_list.push_back(it->second);
    }

    set<char> possibilities[outputs_size];
    auto combos = cartesian_product(iterables_list);
    for (auto& combo : combos) {
        vector<int> inputs_f, inputs_g;
        for (auto& c : combo) {
            switch (c) {
            case 'u':
                inputs_f.push_back(1);
                inputs_g.push_back(0);
                break;
            case 'n':
                inputs_f.push_back(0);
                inputs_g.push_back(1);
                break;
            case '1':
                inputs_f.push_back(1);
                inputs_g.push_back(1);
                break;
            case '0':
                inputs_f.push_back(0);
                inputs_g.push_back(0);
                break;
            }
        }

        vector<int> outputs_f, outputs_g;
        outputs_f = func(inputs_f);
        outputs_g = func(inputs_g);

        vector<char> outputs;
        for (int i = 0; i < outputs_size; i++) {
            int x = outputs_f[i], x_ = outputs_g[i];
            outputs.push_back(x == 1 && x_ == 1 ? '1' : x == 1 && x_ == 0 ? 'u'
                    : x == 0 && x_ == 1                                   ? 'n'
                                                                          : '0');
        }

        for (int i = 0; i < outputs_size; i++)
            possibilities[i].insert((outputs[i]));
    }

    auto gc_from_set = [](set<char>& set) {
        // assert (set.size () > 0);
        for (auto& entry : symbols_set)
            if (set == entry.second)
                return entry.first;
        return '#';
    };

    string propagated_output = "";
    for (auto& p : possibilities)
        propagated_output += gc_from_set(p);

    return propagated_output;
}

int main()
{
    auto combos = cartesian_product({ '-', 'x', 'u', 'n', '0', '1', '?' }, 8);
    for (auto& combo : combos) {
        cout << "add ";
        cout << combo << " " << propagate(add, combo, 3) << endl;
    }
    exit(0);

    vector<char> gcs;
    for (auto it = symbols.begin(); it != symbols.end(); it++)
        gcs.push_back(it->first);

    auto combos3 = cartesian_product(gcs, 3);
    for (auto& combo : combos3) {
        cout << "ch"
             << " ";
        cout << combo << " " << propagate(ch, combo, 1) << endl;
    }
    for (auto& combo : combos3) {
        cout << "maj"
             << " ";
        cout << combo << " " << propagate(maj, combo, 1) << endl;
    }
    for (auto& combo : combos3) {
        cout << "xor3"
             << " ";
        cout << combo << " " << propagate(xor_, combo, 1) << endl;
    }

    // for (int i = 2; i <= 7; i++) {
    //     if (i == 3)
    //         continue;
    //     auto combos = cartesian_product({ '-', 'x' }, i);
    //     for (auto& combo : combos) {
    //         cout << "xor" << i
    //              << " ";
    //         cout << combo << " " << propagate(xor_, combo, 1) << endl;
    //     }
    // }
    for (int i = 2; i <= 7; i++) {
        auto combos = cartesian_product({ '-', 'x' }, i);
        for (auto& combo : combos) {
            cout << "add" << i
                 << " ";
            cout << combo << " " << propagate(add, combo, 3) << endl;
        }
    }
}