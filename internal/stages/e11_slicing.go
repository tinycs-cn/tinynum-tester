package stages

import (
	"fmt"
	"time"

	"github.com/tinycs-cn/tester-utils/runner"
	"github.com/tinycs-cn/tester-utils/test_case_harness"
	"github.com/tinycs-cn/tester-utils/tester_definition"
	"github.com/tinycs-cn/tinynum-tester/internal/helpers"
)

func e11SlicingTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "slicing-and-views",
		Timeout:     30 * time.Second,
		TestFunc:    testE11Slicing,
		CompileStep: autoCompileStep("TestS11", "test_s11"),
	}
}

func testE11Slicing(harness *test_case_harness.TestCaseHarness) error {
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
		// slice
		{"slice_2d", "[[2.0, 3.0], [5.0, 6.0]]", "2D slice [0:2, 1:3]"},
		{"slice_2d_shape", "[2, 2]", "2D slice shape"},
		{"slice_step", "[0.0, 2.0, 4.0]", "1D slice step=2"},
		{"slice_step_shape", "[3]", "1D slice step shape"},
		{"slice_all", "[[4.0, 5.0, 6.0]]", "Slice.all() on axis"},
		{"slice_view_shared", "99.0", "view shares memory"},
		{"slice_3d", "[[[1.0, 2.0]], [[5.0, 6.0]]]", "3D slice"},
		{"slice_3d_shape", "[2, 1, 2]", "3D slice shape"},
		// expandDims
		{"expand_dims_0", "[1, 2, 3]", "expandDims(0) shape"},
		{"expand_dims_1", "[2, 1, 3]", "expandDims(1) shape"},
		{"expand_dims_last", "[2, 3, 1]", "expandDims(last) shape"},
		{"expand_dims_view", "99.0", "expandDims shares memory"},
		// squeeze
		{"squeeze_axis", "[[1.0], [2.0], [3.0]]", "squeeze(0) value"},
		{"squeeze_axis_shape", "[3, 1]", "squeeze(0) shape"},
		{"squeeze_all", "[1.0, 2.0, 3.0]", "squeeze() value"},
		{"squeeze_all_shape", "[3]", "squeeze() shape"},
		// errors
		{"squeeze_error", "ERROR", "squeeze non-1 axis"},
		{"slice_range_error", "ERROR", "slice wrong number of ranges"},
		{"error_step_neg", "ERROR", "slice with negative step throws"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E11 tests passed!")
	return nil
}
