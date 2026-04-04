package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
	"github.com/bootcraft-cn/tinynum-tester/internal/helpers"
)

func e12CreationTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "creation-and-random",
		Timeout:     30 * time.Second,
		TestFunc:    testE12Creation,
		CompileStep: autoCompileStep("TestE12", "test_e12"),
	}
}

func testE12Creation(harness *test_case_harness.TestCaseHarness) error {
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
		// arange
		{"arange_basic", "[0.0, 1.0, 2.0, 3.0, 4.0]", "arange(0,5,1) values"},
		{"arange_basic_shape", "[5]", "arange(0,5,1) shape"},
		{"arange_float_len", "4", "arange(1,2,0.3) length"},
		// linspace
		{"linspace_basic", "[0.0, 0.25, 0.5, 0.75, 1.0]", "linspace(0,1,5) values"},
		{"linspace_basic_shape", "[5]", "linspace(0,1,5) shape"},
		{"linspace_three", "[0.0, 5.0, 10.0]", "linspace(0,10,3) values"},
		{"linspace_single", "[5.0]", "linspace(5,5,1) single"},
		// eye
		{"eye_3", "[[1.0, 0.0, 0.0], [0.0, 1.0, 0.0], [0.0, 0.0, 1.0]]", "eye(3) values"},
		{"eye_3_shape", "[3, 3]", "eye(3) shape"},
		{"eye_1", "[[1.0]]", "eye(1) values"},
		// diag
		{"diag_basic", "[[3.0, 0.0, 0.0], [0.0, 5.0, 0.0], [0.0, 0.0, 7.0]]", "diag([3,5,7]) values"},
		{"diag_basic_shape", "[3, 3]", "diag([3,5,7]) shape"},
		// randn
		{"randn_shape", "[2, 3]", "randn(2,3) shape"},
		{"randn_mean_near_zero", "true", "randn(10000) mean ≈ 0"},
		// rand
		{"rand_shape", "[3, 4]", "rand(3,4) shape"},
		{"rand_in_range", "true", "rand values in [0,1)"},
		// uniform
		{"uniform_shape", "[2, 5]", "uniform(-2,3,2,5) shape"},
		{"uniform_in_range", "true", "uniform values in [-2,3)"},
		// shuffle
		{"shuffle_length", "5", "shuffle preserves length"},
		{"shuffle_elements_preserved", "true", "shuffle preserves elements"},
		// fill
		{"fill_zeros", "[[0.0, 0.0, 0.0], [0.0, 0.0, 0.0]]", "fill(0) zeros"},
		{"fill_value", "[[7.0, 7.0], [7.0, 7.0]]", "fill(7) values"},
		// eye identity property
		{"eye_identity_property", "[[1.0, 2.0, 3.0], [4.0, 5.0, 6.0], [7.0, 8.0, 9.0]]", "eye(3) × x = x"},
		// error
		{"diag_error_non1d", "ERROR", "diag with non-1D throws"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E12 tests passed!")
	return nil
}
