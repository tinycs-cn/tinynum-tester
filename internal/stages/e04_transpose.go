package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
	"github.com/bootcraft-cn/tinynum-tester/internal/helpers"
)

func e04TransposeTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "transpose",
		Timeout:     30 * time.Second,
		TestFunc:    testE04Transpose,
		CompileStep: autoCompileStep("TestE04", "test_e04"),
	}
}

func testE04Transpose(harness *test_case_harness.TestCaseHarness) error {
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
		// 2D transpose
		{"transpose_2d_shape", "3,2", "transpose() shape == [3,2]"},
		{"transpose_2d_toString", "[[1.0, 4.0], [2.0, 5.0], [3.0, 6.0]]", "transpose() toString"},
		{"transpose_get_equiv", "2.0", "a.get(0,1) == transpose().get(1,0)"},
		{"transpose_zerocopy", "99.0", "transpose view shares data (zero-copy)"},
		{"transpose_not_contiguous", "false", "transposed array is not contiguous"},

		// N-D transpose
		{"transpose_nd_shape", "4,2,3", "transpose(2,0,1) shape == [4,2,3]"},
		{"transpose_identity_toString", "[[1.0, 2.0, 3.0], [4.0, 5.0, 6.0]]", "transpose(0,1) is identity"},

		// swapAxes
		{"swapAxes_shape", "4,3,2", "swapAxes(0,2) shape == [4,3,2]"},
		{"swapAxes_same_toString", "[[1.0, 2.0, 3.0], [4.0, 5.0, 6.0]]", "swapAxes(0,0) is identity"},

		// errors
		{"error_transpose_non2d", "EXCEPTION", "transpose() on non-2D throws"},
		{"error_invalid_axes", "EXCEPTION", "transpose with duplicate axes throws"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E04 tests passed!")
	return nil
}
