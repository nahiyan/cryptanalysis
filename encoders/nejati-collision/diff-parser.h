#ifndef _DIFF_PATH_PARSER_H_
#define _DIFF_PATH_PARSER_H_

#include <vector>
#include <string>
using namespace std;

void parse_diff_path(int rounds, FILE* f, vector<string>& A, vector<string>& E, vector<string>& W);

#endif
