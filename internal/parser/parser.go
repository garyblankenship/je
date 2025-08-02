package parser

import (
	"errors"
	"fmt"
	"strings"
)

type OperatorType int

const (
	OpAssignString    OperatorType = iota // =
	OpAssignJSON                          // :=
	OpAssignFile                          // @
	OpAssignJSONFile                      // :@
	OpAppendArray                         // []=
	OpAppendArrayJSON                     // []:=
	OpArrayMap                            // [].key=
	OpArrayMapJSON                        // [].key:=
)

type Assignment struct {
	Path     string
	Operator OperatorType
	Value    string
}

func ParseAssignments(args []string) ([]Assignment, error) {
	var assignments []Assignment

	for _, arg := range args {
		assignment, err := parseAssignment(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid assignment %q: %w", arg, err)
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

func parseAssignment(arg string) (Assignment, error) {
	// Check for array map operators (they contain [].)
	if idx := strings.Index(arg, "[]."); idx > 0 {
		// Find the actual operator after [].
		remaining := arg[idx+3:]
		if opIdx := strings.Index(remaining, ":="); opIdx > 0 {
			return Assignment{
				Path:     arg[:idx+3+opIdx], // Include path up to operator
				Operator: OpArrayMapJSON,
				Value:    remaining[opIdx+2:],
			}, nil
		}
		if opIdx := strings.Index(remaining, "="); opIdx > 0 {
			return Assignment{
				Path:     arg[:idx+3+opIdx], // Include path up to operator
				Operator: OpArrayMap,
				Value:    remaining[opIdx+1:],
			}, nil
		}
	}

	// Check for array append operators (they contain [])
	if idx := strings.Index(arg, "[]:="); idx > 0 {
		return Assignment{
			Path:     arg[:idx+2], // Include []
			Operator: OpAppendArrayJSON,
			Value:    arg[idx+4:],
		}, nil
	}
	if idx := strings.Index(arg, "[]="); idx > 0 {
		return Assignment{
			Path:     arg[:idx+2], // Include []
			Operator: OpAppendArray,
			Value:    arg[idx+3:],
		}, nil
	}

	// Check for file operators
	if idx := strings.Index(arg, ":@"); idx > 0 {
		return Assignment{
			Path:     arg[:idx],
			Operator: OpAssignJSONFile,
			Value:    arg[idx+2:],
		}, nil
	}
	if idx := strings.Index(arg, "@"); idx > 0 && !strings.Contains(arg[:idx], "=") {
		return Assignment{
			Path:     arg[:idx],
			Operator: OpAssignFile,
			Value:    arg[idx+1:],
		}, nil
	}

	// Check for JSON operator
	if idx := strings.Index(arg, ":="); idx > 0 {
		return Assignment{
			Path:     arg[:idx],
			Operator: OpAssignJSON,
			Value:    arg[idx+2:],
		}, nil
	}

	// Check for string operator
	if idx := strings.Index(arg, "="); idx > 0 {
		return Assignment{
			Path:     arg[:idx],
			Operator: OpAssignString,
			Value:    arg[idx+1:],
		}, nil
	}

	return Assignment{}, errors.New("no valid operator found")
}

// ParsePath breaks a path into segments for navigation
func ParsePath(path string) []string {
	// Handle escaped dots
	path = strings.ReplaceAll(path, `\.`, "\x00")
	segments := strings.Split(path, ".")

	// Restore escaped dots
	for i, seg := range segments {
		segments[i] = strings.ReplaceAll(seg, "\x00", ".")
	}

	return segments
}
