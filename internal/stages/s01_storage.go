package stages

import (
	"fmt"
	"time"

	"github.com/tinycs-cn/tester-utils/runner"
	"github.com/tinycs-cn/tester-utils/test_case_harness"
	"github.com/tinycs-cn/tester-utils/tester_definition"
	"github.com/tinycs-cn/tinynum-tester/internal/helpers"
)

func e01StorageTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "storage-and-shape",
		Timeout:     30 * time.Second,
		TestFunc:    testS01Storage,
		CompileStep: autoCompileStep("TestS01", "test_s01"),
	}
}

func testS01Storage(harness *test_case_harness.TestCaseHarness) error {
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

	type tc struct {
		name     string
		expected string
		label    string
	}
	runGroup := func(groupName string, tests []tc) error {
		logger.Infof("--- %s ---", groupName)
		for _, t := range tests {
			if err := helpers.AssertEqual(results, t.name, t.expected); err != nil {
				return err
			}
			logger.Successf("✓ %s", t.label)
		}
		return nil
	}

	zerosTests := []tc{
		{"zeros_size", "6", "zeros(2,3).size() == 6"},
		{"zeros_ndim", "2", "zeros(2,3).ndim() == 2"},
		{"zeros_shape", "[2, 3]", "zeros(2,3).shape() == [2,3]"},
		{"zeros_toString", "[[0.0, 0.0, 0.0], [0.0, 0.0, 0.0]]", "zeros(2,3) toString"},
	}
	onesTests := []tc{
		{"ones_size", "12", "ones(3,4).size() == 12"},
		{"ones_toString", "[[1.0, 1.0, 1.0], [1.0, 1.0, 1.0]]", "ones(2,3) toString values are 1.0"},
	}
	fromArrayTests := []tc{
		{"fromArray_2d_toString", "[[1.0, 2.0, 3.0], [4.0, 5.0, 6.0]]", "fromArray 2D toString"},
		{"1d_toString", "[1.0, 2.0, 3.0]", "1D fromArray toString"},
		{"3d_ndim", "3", "3D array ndim() == 3"},
		{"3d_size", "12", "3D array size() == 12"},
		{"3d_toString", "[[[1.0, 2.0], [3.0, 4.0], [5.0, 6.0]], [[7.0, 8.0], [9.0, 10.0], [11.0, 12.0]]]", "3D toString recursive nesting"},
	}
	factoryTests := []tc{
		{"full_toString", "[[7.0, 7.0], [7.0, 7.0]]", "full(7.0, 2,2) toString"},
		{"zerosLike_toString", "[[0.0, 0.0], [0.0, 0.0]]", "zerosLike toString"},
		{"onesLike_toString", "[[1.0, 1.0], [1.0, 1.0]]", "onesLike toString"},
	}
	errorTests := []tc{
		{"error_mismatch", "ERROR", "fromArray shape mismatch throws"},
		{"data_isolation", "[[1.0, 2.0], [3.0, 4.0]]", "fromArray 数据隔离（修改原数组不影响 NDArray）"},
		{"shape_copy", "OK", "shape() 返回副本"},
	}

	for _, g := range []struct {
		name  string
		tests []tc
	}{
		{"zeros", zerosTests},
		{"ones", onesTests},
		{"fromArray", fromArrayTests},
		{"full / zerosLike / onesLike", factoryTests},
		{"错误处理", errorTests},
	} {
		if err := runGroup(g.name, g.tests); err != nil {
			return err
		}
	}

	total := len(zerosTests) + len(onesTests) + len(fromArrayTests) + len(factoryTests) + len(errorTests)
	logger.Successf("All %d S01 tests passed!", total)
	return nil
}
