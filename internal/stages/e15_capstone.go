package stages

import (
	"fmt"
	"time"

	"github.com/tinycs-cn/tester-utils/runner"
	"github.com/tinycs-cn/tester-utils/test_case_harness"
	"github.com/tinycs-cn/tester-utils/tester_definition"
	"github.com/tinycs-cn/tinynum-tester/internal/helpers"
)

func e15CapstoneTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "capstone-toolkit",
		Timeout:     30 * time.Second,
		TestFunc:    testE15Capstone,
		CompileStep: autoCompileStep("TestS15", "test_s15"),
	}
}

func testE15Capstone(harness *test_case_harness.TestCaseHarness) error {
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
		// tril
		{"tril_default", "[[1.0, 0.0, 0.0], [1.0, 1.0, 0.0], [1.0, 1.0, 1.0]]", "tril default diagonal"},
		{"tril_diag1", "[[1.0, 1.0, 0.0], [1.0, 1.0, 1.0], [1.0, 1.0, 1.0]]", "tril diagonal=1"},
		// triu
		{"triu_default", "[[1.0, 1.0, 1.0], [0.0, 1.0, 1.0], [0.0, 0.0, 1.0]]", "triu default diagonal"},
		{"triu_diag_neg1", "[[1.0, 1.0, 1.0], [1.0, 1.0, 1.0], [0.0, 1.0, 1.0]]", "triu diagonal=-1"},
		// norm
		{"norm_axis1", "[5.0, 10.0]", "norm axis=1"},
		{"norm_axis0", "[5.0, 13.0]", "norm axis=0"},
		// diff
		{"diff_1d", "[2.0, 3.0, 4.0]", "diff 1D"},
		{"diff_axis1", "[[2.0, 3.0, 4.0], [3.0, 4.0, 5.0]]", "diff axis=1 values"},
		{"diff_axis1_shape", "[2, 3]", "diff axis=1 shape"},
		{"diff_axis0", "[[1.0, 2.0, 3.0, 4.0]]", "diff axis=0"},
		// percentile
		{"percentile_50", "[3.0]", "percentile 50th"},
		{"percentile_0", "[1.0]", "percentile 0th"},
		{"percentile_100", "[5.0]", "percentile 100th"},
		{"percentile_25", "[2.0]", "percentile 25th"},
		// argsort
		{"argsort_1d", "[1.0, 2.0, 0.0]", "argsort 1D"},
		{"argsort_2d_axis1", "[[1.0, 2.0, 0.0], [2.0, 0.0, 1.0]]", "argsort 2D axis=1"},
		// unique
		{"unique", "[1.0, 2.0, 3.0]", "unique values"},
		{"unique_shape", "[3]", "unique shape"},
		// allClose
		{"allclose_true", "true", "allClose true"},
		{"allclose_false", "false", "allClose false"},
		{"allclose_shape_diff", "false", "allClose shape mismatch"},
		// astype
		{"astype_int8", "[127.0, -128.0, 50.0, 0.0]", "astype FLOAT32→INT8 clamp"},
		{"astype_float32", "[127.0, -128.0, 50.0]", "astype FLOAT32→FLOAT32"},
		// softmax integration
		{"softmax_row_sum", "true", "softmax row sums ≈ 1.0"},
		{"softmax_order", "true", "softmax ordering correct"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E15 tests passed!")
	return nil
}
