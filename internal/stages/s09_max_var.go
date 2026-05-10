package stages

import (
	"fmt"
	"time"

	"github.com/tinycs-cn/tester-utils/runner"
	"github.com/tinycs-cn/tester-utils/test_case_harness"
	"github.com/tinycs-cn/tester-utils/tester_definition"
	"github.com/tinycs-cn/tinynum-tester/internal/helpers"
)

func e09MaxVarTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "max-var-and-friends",
		Timeout:     30 * time.Second,
		TestFunc:    testS09MaxVar,
		CompileStep: autoCompileStep("TestS09", "test_s09"),
	}
}

func testS09MaxVar(harness *test_case_harness.TestCaseHarness) error {
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
		// max / min
		{"max_axis1", "[4.0, 9.0]", "max axis=1"},
		{"max_axis0_keep", "[[3.0, 5.0, 9.0]]", "max axis=0 keepDims"},
		{"max_axis0_keep_shape", "[1, 3]", "max axis=0 keepDims shape"},
		{"min_axis0", "[1.0, 1.0, 4.0]", "min axis=0"},
		{"min_axis1_keep", "[[1.0], [1.0]]", "min axis=1 keepDims"},
		// argmax / argmin
		{"argmax_axis1", "[2.0, 2.0]", "argmax axis=1"},
		{"argmax_axis0", "[0.0, 1.0, 1.0]", "argmax axis=0"},
		{"argmin_axis1", "[1.0, 0.0]", "argmin axis=1"},
		{"argmin_axis0", "[1.0, 0.0, 0.0]", "argmin axis=0"},
		// prod
		{"prod_axis1", "[6.0, 120.0]", "prod axis=1"},
		{"prod_axis0", "[4.0, 10.0, 18.0]", "prod axis=0"},
		// var / std
		{"var_axis1", "[1.0, 1.0]", "var axis=1"},
		{"var_axis0_keep", "[[0.25, 0.25]]", "var axis=0 keepDims"},
		{"std_axis1", "[1.0, 1.0]", "std axis=1"},
		{"std_axis0", "[0.5, 0.5]", "std axis=0"},
		// countNonZero
		{"count_nonzero", "3", "countNonZero 1D"},
		{"count_nonzero_2d", "4", "countNonZero 2D"},
		// Negative axis
		{"max_neg_axis", "[4.0, 9.0]", "max negative axis"},
		// 3D
		{"max_3d", "[[4.0, 1.0], [5.0, 9.0]]", "max 3D axis=1"},
		// Error
		{"max_axis_error", "ERROR", "max invalid axis throws error"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All S09 tests passed!")
	return nil
}
