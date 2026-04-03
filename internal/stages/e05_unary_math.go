package stages

import (
	"fmt"
	"time"

	"github.com/tensorhero-cn/tester-utils/runner"
	"github.com/tensorhero-cn/tester-utils/test_case_harness"
	"github.com/tensorhero-cn/tester-utils/tester_definition"
	"github.com/tensorhero-cn/tinynum-tester/internal/helpers"
)

func e05UnaryMathTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "unary-math",
		Timeout:     30 * time.Second,
		TestFunc:    testE05UnaryMath,
		CompileStep: autoCompileStep("TestE05", "test_e05"),
	}
}

func testE05UnaryMath(harness *test_case_harness.TestCaseHarness) error {
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
		{"neg_toString", "[-1.0, 2.0, -3.0]", "neg() element-wise"},
		{"abs_toString", "[1.0, 2.0, 3.0]", "abs() element-wise"},
		{"exp_zeros", "[[1.0, 1.0, 1.0], [1.0, 1.0, 1.0]]", "exp(zeros) == ones"},
		{"log_ones", "[[0.0, 0.0], [0.0, 0.0]]", "log(ones) == zeros"},
		{"sqrt_toString", "[0.0, 1.0, 2.0, 3.0]", "sqrt([0,1,4,9])"},
		{"square_toString", "[4.0, 0.0, 9.0]", "square([-2,0,3])"},
		{"tanh_zero", "0.0", "tanh(0) == 0"},
		{"sin_zero", "0.0", "sin(0) == 0"},
		{"cos_zero", "1.0", "cos(0) == 1"},
		{"sign_toString", "[-1.0, 0.0, 1.0]", "sign([-5,0,7])"},
		{"round_toString", "[1.0, 2.0, 0.0, 2.0]", "round([1.4,1.6,-0.5,2.3])"},
		{"clip_toString", "[-2.0, -1.0, 0.0, 1.0, 2.0]", "clip(-2,2)"},
		{"pow_half", "[1.0, 2.0, 3.0]", "pow(0.5) == sqrt"},
		{"unary_independent", "1.0", "unary returns new array (original unchanged)"},
		{"unary_transposed", "[[1.0, 16.0], [4.0, 25.0], [9.0, 36.0]]", "unary works on transposed view"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E05 tests passed!")
	return nil
}
