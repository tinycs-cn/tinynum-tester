package stages

import (
	"fmt"
	"time"

	"github.com/tinycs-cn/tester-utils/runner"
	"github.com/tinycs-cn/tester-utils/test_case_harness"
	"github.com/tinycs-cn/tester-utils/tester_definition"
	"github.com/tinycs-cn/tinynum-tester/internal/helpers"
)

func e03ReshapeTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "reshape",
		Timeout:     30 * time.Second,
		TestFunc:    testE03Reshape,
		CompileStep: autoCompileStep("TestE03", "test_e03"),
	}
}

func testE03Reshape(harness *test_case_harness.TestCaseHarness) error {
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
		// reshape
		{"reshape_shape", "[3, 2]", "reshape(3,2) shape == [3,2]"},
		{"reshape_toString", "[[1.0, 2.0], [3.0, 4.0], [5.0, 6.0]]", "reshape(3,2) toString"},
		{"reshape_neg1_shape", "[3, 2]", "reshape(3,-1) infers shape [3,2]"},
		{"reshape_neg1_3d_shape", "[2, 2, 3]", "reshape(2,-1,3) infers shape [2,2,3]"},
		{"reshape_zerocopy", "99.0", "reshape view shares data (zero-copy)"},

		// flatten
		{"flatten_shape", "[6]", "flatten().shape == [6]"},
		{"flatten_toString", "[1.0, 2.0, 3.0, 4.0, 5.0, 6.0]", "flatten() toString"},

		// duplicate
		{"duplicate_toString", "[[1.0, 2.0], [3.0, 4.0]]", "duplicate() toString matches original"},
		{"duplicate_independent", "1.0", "duplicate is independent (deep copy)"},

		// error
		{"error_reshape_size", "ERROR", "reshape size mismatch throws"},
		{"error_reshape_double_neg1", "ERROR", "reshape(-1,-1) two -1 dims throws"},

		// non-contiguous reshape: transpose then flatten
		{"reshape_noncontiguous_toString", "[1.0, 4.0, 2.0, 5.0, 3.0, 6.0]", "reshape non-contiguous: flatten of transposed [2,3] → [6]"},
		{"reshape_noncontiguous_copy", "1.0", "reshape non-contiguous: result is deep copy (original unchanged)"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All E03 tests passed!")
	return nil
}
