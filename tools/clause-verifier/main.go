package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const maxLogLineSize = 64 * 1024 * 1024

type clause struct {
	literals []int
	kind     string
}

func scanFile(path string, visit func(line string, lineNumber int) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), maxLogLineSize)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		if err := visit(scanner.Text(), lineNumber); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func parseSolution(path string) (map[int]bool, error) {
	assignment := make(map[int]bool)
	err := scanFile(path, func(line string, lineNumber int) error {
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] != "v" {
			return nil
		}

		for _, field := range fields[1:] {
			if field == "v" {
				continue
			}

			literal, err := strconv.Atoi(field)
			if err != nil {
				return fmt.Errorf("%s:%d: invalid solution literal %q", path, lineNumber, field)
			}
			if literal == 0 {
				continue
			}

			variable := literal
			if variable < 0 {
				variable = -variable
			}
			assignment[variable] = literal > 0
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(assignment) == 0 {
		return nil, errors.New("solution log contains no assignment lines")
	}

	return assignment, nil
}

func parseClauses(path string) ([]clause, error) {
	var clauses []clause
	err := scanFile(path, func(line string, lineNumber int) error {
		kind := ""
		body := ""
		switch {
		case strings.HasPrefix(line, "Reason clause:"):
			kind = "reason"
			body = strings.TrimSpace(strings.TrimPrefix(line, "Reason clause:"))
		case strings.HasPrefix(line, "Blocking clause:"):
			kind = "blocking"
			body = strings.TrimSpace(strings.TrimPrefix(line, "Blocking clause:"))
		default:
			return nil
		}

		parsed := clause{kind: kind}
		for _, field := range strings.Fields(body) {
			literal, err := strconv.Atoi(field)
			if err != nil {
				return fmt.Errorf("%s:%d: invalid clause literal %q", path, lineNumber, field)
			}
			if literal != 0 {
				parsed.literals = append(parsed.literals, literal)
			}
		}
		clauses = append(clauses, parsed)
		return nil
	})

	return clauses, err
}

func literalValue(literal int, assignment map[int]bool) (bool, bool, error) {
	variable := literal
	if variable < 0 {
		variable = -variable
	}
	assigned, ok := assignment[variable]
	if !ok {
		return false, false, fmt.Errorf("variable %d is missing from the solution", variable)
	}

	matches := assigned == (literal > 0)
	return matches, assigned, nil
}

func verifyClause(candidate clause, assignment map[int]bool) (bool, error) {
	for _, literal := range candidate.literals {
		matches, _, err := literalValue(literal, assignment)
		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}

	return false, nil
}

func printClause(candidate clause, assignment map[int]bool) error {
	label := "Blocking"
	if candidate.kind == "reason" {
		label = "Reason"
	}
	fmt.Printf("%s: ", label)

	for _, literal := range candidate.literals {
		matches, _, err := literalValue(literal, assignment)
		if err != nil {
			return err
		}
		mark := "𐄂"
		if matches {
			mark = "✓"
		}
		fmt.Printf("%d(%s) ", literal, mark)
	}
	fmt.Println()
	return nil
}

func solverFinished(path string) (bool, error) {
	finished := false
	err := scanFile(path, func(line string, _ int) error {
		if strings.HasPrefix(line, "c exit ") {
			finished = true
		}
		return nil
	})
	return finished, err
}

func run(solutionLogPath, solverLogPath string) error {
	assignment, err := parseSolution(solutionLogPath)
	if err != nil {
		return fmt.Errorf("read solution: %w", err)
	}

	clauses, err := parseClauses(solverLogPath)
	if err != nil {
		return fmt.Errorf("read clauses: %w", err)
	}

	counts := map[string]int{"reason": 0, "blocking": 0}
	for _, candidate := range clauses {
		counts[candidate.kind]++
		satisfied, err := verifyClause(candidate, assignment)
		if err != nil {
			return fmt.Errorf("verify %s clause: %w", candidate.kind, err)
		}
		if !satisfied {
			if err := printClause(candidate, assignment); err != nil {
				return err
			}
		}
	}

	fmt.Printf("{'reason': %d, 'blocking': %d}\n", counts["reason"], counts["blocking"])
	finished, err := solverFinished(solverLogPath)
	if err != nil {
		return fmt.Errorf("check solver status: %w", err)
	}
	if finished {
		fmt.Println("Solved")
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <solution-log> <solver-log>\n", os.Args[0])
		os.Exit(2)
	}

	if err := run(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, "clause-verifier:", err)
		os.Exit(1)
	}
}
