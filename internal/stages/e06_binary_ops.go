package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
	"github.com/bootcraft-cn/tinynum-tester/internal/helpers"
)

func e06BinaryOpsTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "binary-ops",
		Timeout:     30 * time.Second,
		TestFunc:    testE06BinaryOps,
		CompileStep: autoCompileStep("TestE06", "test_e06"),
	}
}

func testE06BinaryOps(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir
	lang := harness.DetectedLang

	r := runner.Run(workDir, lang.RunCmd, lang.RunArgs...).
		WithTimeout(10 * time.Second).
		WithLogger(logger).
		Execute().
		Exit(0)

	if err := r.Error(); err != nil {
		return fmt.Errorf("test driver failed: %v", err)
	}

	results := helpers.ParseStructuredOutput(string(r.Result().Stdout))

	tests := []struct {
		name     string
		expected string
		label    string
	}{
		// Arithmetic (NDArray)
		{"add_toString", "[[3.0, 6.0, 9.0], [12.0, 15.0, 18.0]]", "add() element-wise"},
		{"sub_toString", "[[1.0, 2.0, 3.0], [4.0, 5.0, 6.0]]", "sub() element-wise"},
		{"mul_toString", "[[2.0, 8.0, 18.0], [32.0, 50.0, 72.0]]", "mul() element-wise"},
		{"div_toString", "[[2.0, 2.0, 2.0], [2.0, 2.0, 2.0]]", "div() element-wise"},
		{"pow_toString", "[8.0, 9.0, 1.0]", "pow() element-wise"},
		{"maximum_toString", "[4.0, 5.0, 3.0]", "maximum() element-wise"},
		// Arithmetic (scalar)
		{"add_scalar", "[15.0, 25.0, 35.0]", "add(scalar)"},
		{"sub_scalar", "[5.0, 15.0, 25.0]", "sub(scalar)"},
		{"mul_scalar", "[20.0, 40.0, 60.0]", "mul(scalar)"},
		{"div_scalar", "[1.0, 2.0, 3.0]", "div(scalar)"},
		// Comparisons (NDArray)
		{"eq_toString", "[1.0, 0.0, 0.0, 1.0]", "eq() element-wise"},
		{"neq_toString", "[0.0, 1.0, 1.0, 0.0]", "neq() element-wise"},
		{"gt_toString", "[0.0, 0.0, 1.0, 0.0]", "gt() element-wise"},
		{"gte_toString", "[1.0, 0.0, 1.0, 1.0]", "gte() element-wise"},
		{"lt_toString", "[0.0, 1.0, 0.0, 0.0]", "lt() element-wise"},
		{"lte_toString", "[1.0, 1.0, 0.0, 1.0]", "lte() element-wise"},
		// Comparisons (scalar)
		{"eq_scalar", "[0.0, 0.0, 1.0, 0.0, 0.0]", "eq(scalar)"},
		{"neq_scalar", "[1.0, 1.0, 0.0, 1.0, 1.0]", "neq(scalar)"},
		{"gt_scalar", "[0.0, 0.0, 0.0, 1.0, 1.0]", "gt(scalar)"},
		{"gte_scalar", "[0.0, 0.0, 1.0, 1.0, 1.0]", "gte(scalar)"},
		{"lt_scalar", "[1.0, 1.0, 0.0, 0.0, 0.0]", "lt(scalar)"},
		{"lte_scalar", "[1.0, 1.0, 1.0, 0.0, 0.0]", "lte(scalar)"},
		// Invariants
		{"binary_independent", "1.0", "binary returns new array (original unchanged)"},
		{"binary_transposed", "[[11.0, 24.0], [32.0, 45.0], [53.0, 66.0]]", "binary works on transposed views"},
		{"shape_mismatch", "ERROR", "shape mismatch throws error"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E06 tests passed!")
	return nil
}
