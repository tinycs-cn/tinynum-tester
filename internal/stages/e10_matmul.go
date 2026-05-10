package stages

import (
	"fmt"
	"time"

	"github.com/tinycs-cn/tester-utils/runner"
	"github.com/tinycs-cn/tester-utils/test_case_harness"
	"github.com/tinycs-cn/tester-utils/tester_definition"
	"github.com/tinycs-cn/tinynum-tester/internal/helpers"
)

func e10MatmulTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "matmul",
		Timeout:     30 * time.Second,
		TestFunc:    testE10Matmul,
		CompileStep: autoCompileStep("TestS10", "test_s10"),
	}
}

func testE10Matmul(harness *test_case_harness.TestCaseHarness) error {
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
		// dot
		{"dot_1d", "32.0", "dot product [1,2,3]·[4,5,6]"},
		{"dot_shape", "[]", "dot result is 0-D"},
		// matmul 2D
		{"matmul_2d", "[[19.0, 22.0], [43.0, 50.0]]", "2D matmul 2×2"},
		{"matmul_2d_shape", "[2, 2]", "2D matmul shape"},
		// matmul rectangular
		{"matmul_rect", "[[22.0, 28.0], [49.0, 64.0]]", "rectangular matmul (2,3)×(3,2)"},
		{"matmul_rect_shape", "[2, 2]", "rectangular matmul shape"},
		// identity
		{"matmul_identity", "[[2.0, 3.0, 4.0], [5.0, 6.0, 7.0], [8.0, 9.0, 10.0]]", "eye(3) @ X == X"},
		// batch
		{"matmul_batch", "[[[4.0, 2.0], [10.0, 5.0]], [[8.0, 16.0], [11.0, 22.0]]]", "batch matmul [2,2,3]@[2,3,2]"},
		{"matmul_batch_shape", "[2, 2, 2]", "batch matmul shape"},
		// batch broadcast
		{"matmul_batch_broadcast", "[[[4.0, 2.0], [10.0, 5.0]], [[16.0, 8.0], [22.0, 11.0]]]", "batch broadcast [2,2,3]@[1,3,2]"},
		{"matmul_batch_broadcast_shape", "[2, 2, 2]", "batch broadcast shape"},
		// errors
		{"dot_length_error", "ERROR", "dot mismatched lengths"},
		{"matmul_dim_error", "ERROR", "matmul inner dim mismatch"},
		{"dot_not_1d_error", "ERROR", "dot with 2D array"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E10 tests passed!")
	return nil
}
