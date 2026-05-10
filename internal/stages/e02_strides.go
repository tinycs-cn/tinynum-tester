package stages

import (
	"fmt"
	"time"

	"github.com/tinycs-cn/tester-utils/runner"
	"github.com/tinycs-cn/tester-utils/test_case_harness"
	"github.com/tinycs-cn/tester-utils/tester_definition"
	"github.com/tinycs-cn/tinynum-tester/internal/helpers"
)

func e02StridesTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "strides-and-indexing",
		Timeout:     30 * time.Second,
		TestFunc:    testS02Strides,
		CompileStep: autoCompileStep("TestS02", "test_s02"),
	}
}

func testS02Strides(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir
	lang := harness.DetectedLang

	// Run test driver
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
		// computeStrides
		{"strides_2d", "3,1", "computeStrides([2,3]) == [3,1]"},
		{"strides_3d", "20,5,1", "computeStrides([3,4,5]) == [20,5,1]"},
		{"strides_1d", "1", "computeStrides([5]) == [1]"},

		// get 2D
		{"get_2d_00", "1.0", "get(0,0) == 1.0"},
		{"get_2d_02", "3.0", "get(0,2) == 3.0"},
		{"get_2d_10", "4.0", "get(1,0) == 4.0"},
		{"get_2d_12", "6.0", "get(1,2) == 6.0"},

		// get 3D
		{"get_3d_000", "1.0", "3D get(0,0,0) == 1.0"},
		{"get_3d_123", "24.0", "3D get(1,2,3) == 24.0"},
		{"get_3d_012", "7.0", "3D get(0,1,2) == 7.0"},

		// set
		{"set_get", "99.0", "set(99,1,1) then get(1,1) == 99.0"},
		{"set_toString", "[[1.0, 2.0, 3.0], [4.0, 99.0, 6.0]]", "toString reflects set"},

		// isContiguous
		{"isContiguous_fresh", "true", "freshly created array isContiguous"},

		// error
		{"error_wrong_indices", "ERROR", "get with wrong index count throws"},
		{"error_out_of_bounds", "ERROR", "get with out-of-bounds index throws"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All S02 tests passed!")
	return nil
}
