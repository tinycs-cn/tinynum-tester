package stages

import (
	"fmt"
	"time"

	"github.com/tensorhero-cn/tester-utils/runner"
	"github.com/tensorhero-cn/tester-utils/test_case_harness"
	"github.com/tensorhero-cn/tester-utils/tester_definition"
	"github.com/tensorhero-cn/tinynum-tester/internal/helpers"
)

func e08SumMeanTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "sum-and-mean",
		Timeout:     30 * time.Second,
		TestFunc:    testE08SumMean,
		CompileStep: autoCompileStep("TestE08", "test_e08"),
	}
}

func testE08SumMean(harness *test_case_harness.TestCaseHarness) error {
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
		// Global reduction
		{"sum_global", "21.0", "global sum"},
		{"mean_global", "3.5", "global mean"},
		// Sum axis (no keepDims)
		{"sum_axis0_nokeep", "[5.0, 7.0, 9.0]", "sum axis=0 no keepDims"},
		{"sum_axis0_shape_nokeep", "[3]", "sum axis=0 shape no keepDims"},
		{"sum_axis1_nokeep", "[6.0, 15.0]", "sum axis=1 no keepDims"},
		{"sum_axis1_shape_nokeep", "[2]", "sum axis=1 shape no keepDims"},
		// Sum axis (keepDims=true)
		{"sum_axis0_keep", "[[5.0, 7.0, 9.0]]", "sum axis=0 keepDims"},
		{"sum_axis0_shape_keep", "[1, 3]", "sum axis=0 shape keepDims"},
		{"sum_axis1_keep", "[[6.0], [15.0]]", "sum axis=1 keepDims"},
		{"sum_axis1_shape_keep", "[2, 1]", "sum axis=1 shape keepDims"},
		// Mean axis
		{"mean_axis1_nokeep", "[2.0, 5.0]", "mean axis=1 no keepDims"},
		{"mean_axis0_keep", "[[2.5, 3.5, 4.5]]", "mean axis=0 keepDims"},
		// Negative axis
		{"sum_neg_axis", "[6.0, 15.0]", "sum negative axis"},
		// 3D
		{"sum_3d_axis1", "[[4.0, 6.0], [12.0, 14.0]]", "sum 3D axis=1"},
		{"sum_3d_axis1_shape", "[2, 2]", "sum 3D axis=1 shape"},
		// Multi-axis
		{"sum_multi_axes", "[14.0, 22.0]", "sum multi-axes [0,2]"},
		{"sum_multi_axes_shape", "[2]", "sum multi-axes shape"},
		// 1D
		{"sum_1d", "60.0", "sum 1D"},
		// Error
		{"sum_axis_error", "ERROR", "sum invalid axis throws error"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E08 tests passed!")
	return nil
}
