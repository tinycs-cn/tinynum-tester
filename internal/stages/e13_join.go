package stages

import (
	"fmt"
	"time"

	"github.com/tinycs-cn/tester-utils/runner"
	"github.com/tinycs-cn/tester-utils/test_case_harness"
	"github.com/tinycs-cn/tester-utils/tester_definition"
	"github.com/tinycs-cn/tinynum-tester/internal/helpers"
)

func e13JoinTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "join-and-transform",
		Timeout:     30 * time.Second,
		TestFunc:    testE13Join,
		CompileStep: autoCompileStep("TestS13", "test_s13"),
	}
}

func testE13Join(harness *test_case_harness.TestCaseHarness) error {
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
		// concatenate
		{"concat_axis0", "[[1.0, 2.0], [3.0, 4.0], [5.0, 6.0], [7.0, 8.0]]", "concat axis=0 values"},
		{"concat_axis0_shape", "[4, 2]", "concat axis=0 shape"},
		{"concat_axis1", "[[1.0, 2.0, 5.0, 6.0], [3.0, 4.0, 7.0, 8.0]]", "concat axis=1 values"},
		{"concat_axis1_shape", "[2, 4]", "concat axis=1 shape"},
		{"concat_diff_size", "[[1.0, 2.0], [3.0, 4.0], [9.0, 10.0]]", "concat different sizes"},
		{"concat_diff_size_shape", "[3, 2]", "concat different sizes shape"},
		{"concat_1d", "[1.0, 2.0, 3.0, 4.0, 5.0]", "concat 1D"},
		// stack
		{"stack_axis0", "[[1.0, 2.0, 3.0], [4.0, 5.0, 6.0]]", "stack axis=0 values"},
		{"stack_axis0_shape", "[2, 3]", "stack axis=0 shape"},
		{"stack_axis1", "[[1.0, 4.0], [2.0, 5.0], [3.0, 6.0]]", "stack axis=1 values"},
		{"stack_axis1_shape", "[3, 2]", "stack axis=1 shape"},
		{"stack_2d_shape", "[2, 2, 2]", "stack 2D shape"},
		// pad
		{"pad_rows", "[[0.0, 0.0, 0.0], [1.0, 2.0, 3.0], [4.0, 5.0, 6.0], [0.0, 0.0, 0.0]]", "pad rows"},
		{"pad_rows_shape", "[4, 3]", "pad rows shape"},
		{"pad_cols_value", "[[-1.0, 1.0, 2.0, 3.0, -1.0, -1.0], [-1.0, 4.0, 5.0, 6.0, -1.0, -1.0]]", "pad cols with value"},
		{"pad_cols_shape", "[2, 6]", "pad cols shape"},
		{"pad_1d", "[0.0, 0.0, 1.0, 2.0, 3.0, 0.0]", "pad 1D"},
		// flip
		{"flip_axis0", "[[4.0, 5.0, 6.0], [1.0, 2.0, 3.0]]", "flip axis=0"},
		{"flip_axis1", "[[3.0, 2.0, 1.0], [6.0, 5.0, 4.0]]", "flip axis=1"},
		{"flip_1d", "[5.0, 4.0, 3.0, 2.0, 1.0]", "flip 1D"},
		// errors
		{"concat_shape_error", "ERROR", "concat shape mismatch error"},
		{"pad_dim_error", "ERROR", "pad wrong padWidth length error"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E13 tests passed!")
	return nil
}
