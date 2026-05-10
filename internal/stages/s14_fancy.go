package stages

import (
	"fmt"
	"time"

	"github.com/tinycs-cn/tester-utils/runner"
	"github.com/tinycs-cn/tester-utils/test_case_harness"
	"github.com/tinycs-cn/tester-utils/tester_definition"
	"github.com/tinycs-cn/tinynum-tester/internal/helpers"
)

func e14FancyTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "fancy-indexing",
		Timeout:     30 * time.Second,
		TestFunc:    testS14Fancy,
		CompileStep: autoCompileStep("TestS14", "test_s14"),
	}
}

func testS14Fancy(harness *test_case_harness.TestCaseHarness) error {
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
		// indexSelect
		{"index_select_axis0", "[[7.0, 8.0, 9.0], [1.0, 2.0, 3.0], [10.0, 11.0, 12.0]]", "indexSelect axis=0 values"},
		{"index_select_axis0_shape", "[3, 3]", "indexSelect axis=0 shape"},
		{"index_select_axis1", "[[3.0, 1.0], [6.0, 4.0], [9.0, 7.0], [12.0, 10.0]]", "indexSelect axis=1 values"},
		{"index_select_axis1_shape", "[4, 2]", "indexSelect axis=1 shape"},
		{"index_select_1d", "[50.0, 20.0, 20.0, 40.0]", "indexSelect 1D"},
		{"index_select_dup", "[[1.0, 2.0, 3.0], [1.0, 2.0, 3.0], [7.0, 8.0, 9.0]]", "indexSelect duplicate indices"},
		// scatterAdd
		{"scatter_add_basic", "[[2.0, 2.0, 2.0], [0.0, 0.0, 0.0], [1.0, 1.0, 1.0], [0.0, 0.0, 0.0]]", "scatterAdd basic"},
		{"scatter_add_dup", "[[8.0, 10.0, 12.0], [4.0, 5.0, 6.0]]", "scatterAdd duplicate indices accumulation"},
		// maskedFill
		{"masked_fill_2d", "[[1.0, 2.0, -999.0], [-999.0, 5.0, 6.0]]", "maskedFill 2D"},
		{"masked_fill_1d", "[0.0, 20.0, 0.0]", "maskedFill 1D"},
		{"masked_fill_none", "[10.0, 20.0, 30.0]", "maskedFill no change"},
		// where
		{"where_2d", "[[10.0, -2.0], [-3.0, 40.0]]", "where 2D"},
		{"where_1d", "[-1.0, 2.0, -3.0, 4.0, 5.0]", "where 1D"},
		// errors
		{"masked_fill_shape_error", "ERROR", "maskedFill shape mismatch error"},
		{"where_shape_error", "ERROR", "where shape mismatch error"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All S14 tests passed!")
	return nil
}
