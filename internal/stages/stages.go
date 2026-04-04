package stages

import (
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

// GetDefinition returns the TesterDefinition for the tinynum course.
func GetDefinition() tester_definition.TesterDefinition {
	return tester_definition.TesterDefinition{
		TestCases: []tester_definition.TestCase{
			// Phase 1: 数据模型
			e01StorageTestCase(),
			e02StridesTestCase(),
			e03ReshapeTestCase(),
			e04TransposeTestCase(),
			// Phase 2: 逐元素运算
			e05UnaryMathTestCase(),
			e06BinaryOpsTestCase(),
			e07BroadcastingTestCase(),
			e08SumMeanTestCase(),
			e09MaxVarTestCase(),
			// Phase 3: 线性代数 & 视图
			e10MatmulTestCase(),
			e11SlicingTestCase(),
			e12CreationTestCase(),
			// Phase 4: 工具箱
			e13JoinTestCase(),
			e14FancyTestCase(),
			e15CapstoneTestCase(),
		},
	}
}

// javaRule creates a LanguageRule for Java auto-detection.
// testDriver is the class name (e.g. "TestE01").
func javaRule(testDriver string) tester_definition.LanguageRule {
	return tester_definition.LanguageRule{
		DetectFile: "src/main/java/dev/tensorhero/tinynum/NDArray.java",
		Language:   "java",
		Source:     "src/main/java/dev/tensorhero/tinynum/NDArray.java",
		Flags: []string{
			"-encoding", "UTF-8",
			"src/main/java/dev/tensorhero/tinynum/Slice.java",
			"src/main/java/dev/tensorhero/tinynum/DType.java",
			"tests/" + testDriver + ".java",
		},
		RunCmd:  "java",
		RunArgs: []string{"-cp", ".", testDriver},
	}
}

// pythonRule creates a LanguageRule for Python auto-detection.
// testDriver is the module name without extension (e.g. "test_e01").
func pythonRule(testDriver string) tester_definition.LanguageRule {
	return tester_definition.LanguageRule{
		DetectFile: "tinynum/ndarray.py",
		Language:   "python",
		Source:     "tinynum/ndarray.py",
		RunCmd:     "python3",
		RunArgs:    []string{"tests/" + testDriver + ".py"},
	}
}

// autoCompileStep returns a CompileStep with auto-detection for Java/Python.
func autoCompileStep(javaDriver, pythonDriver string) *tester_definition.CompileStep {
	return &tester_definition.CompileStep{
		Language: "auto",
		AutoDetect: []tester_definition.LanguageRule{
			javaRule(javaDriver),
			pythonRule(pythonDriver),
		},
	}
}
