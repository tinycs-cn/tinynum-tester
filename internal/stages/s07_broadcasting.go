package stages

import (
	"fmt"
	"time"

	"github.com/tinycs-cn/tester-utils/runner"
	"github.com/tinycs-cn/tester-utils/test_case_harness"
	"github.com/tinycs-cn/tester-utils/tester_definition"
	"github.com/tinycs-cn/tinynum-tester/internal/helpers"
)

func e07BroadcastingTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "broadcasting",
		Timeout:     30 * time.Second,
		TestFunc:    testS07Broadcasting,
		CompileStep: autoCompileStep("TestS07", "test_s07"),
	}
}

func testS07Broadcasting(harness *test_case_harness.TestCaseHarness) error {
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
		// broadcastShapes
		{"bc_same", "[2, 3]", "broadcastShapes same shape"},
		{"bc_right_align", "[2, 3]", "broadcastShapes right-align [2,3]+[3]"},
		{"bc_both_expand", "[3, 4]", "broadcastShapes both expand [3,1]+[1,4]"},
		{"bc_3d", "[2, 3, 4]", "broadcastShapes 3D [2,1,4]+[3,1]"},
		{"bc_scalar_like", "[2, 3]", "broadcastShapes scalar-like [2,3]+[1]"},
		{"bc_error", "ERROR", "broadcastShapes incompatible throws error"},
		// broadcastTo
		{"bt_shape", "[4, 3]", "broadcastTo shape [1,3]→[4,3]"},
		{"bt_get_0_1", "2.0", "broadcastTo get(0,1)"},
		{"bt_get_3_2", "3.0", "broadcastTo get(3,2)"},
		{"bt_toString", "[[1.0, 2.0, 3.0], [1.0, 2.0, 3.0], [1.0, 2.0, 3.0], [1.0, 2.0, 3.0]]", "broadcastTo toString"},
		{"bt_error", "ERROR", "broadcastTo incompatible throws error"},
		// Auto-broadcast binary ops
		{"add_mat_row", "[[11.0, 22.0, 33.0], [14.0, 25.0, 36.0]]", "add matrix+row auto-broadcast"},
		{"mul_outer", "[[10.0, 20.0, 30.0, 40.0], [20.0, 40.0, 60.0, 80.0], [30.0, 60.0, 90.0, 120.0]]", "mul col*row outer product"},
		{"sub_3d", "[[[-9.0, -18.0, -27.0]], [[-6.0, -15.0, -24.0]]]", "sub 3D auto-broadcast"},
		{"gt_broadcast", "[[0.0, 0.0, 0.0], [1.0, 1.0, 1.0]]", "gt comparison with broadcast"},
		// Zero-copy & identity
		{"bt_zerocopy", "99.0", "broadcastTo is zero-copy (shares data)"},
		{"bt_identity", "[1.0, 2.0]", "broadcastTo same shape is identity"},
		{"bt_scalar", "[5.0, 5.0, 5.0, 5.0]", "broadcastTo scalar-like [1]→[4]"},
	}

	for _, tc := range tests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All S07 tests passed!")
	return nil
}
