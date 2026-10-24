package classification

// Rule definitions and keyword dictionaries for difficulty estimation.

var BeginnerLabels = []string{
	"good first issue",
	"good-first-issue",
	"first-timers-only",
	"beginner",
	"easy",
	"starter",
	"up-for-grabs",
	"low-hanging-fruit",
	"documentation",
	"docs",
	"typo",
	"trivial",
	"junior",
	"level/beginner",
	"difficulty: easy",
	"difficulty: beginner",
}

var IntermediateLabels = []string{
	"help wanted",
	"help-wanted",
	"bug",
	"enhancement",
	"feature",
	"refactor",
	"testing",
	"test",
	"api",
	"ui",
	"cli",
	"cleanup",
	"intermediate",
	"level/intermediate",
	"difficulty: medium",
}

var AdvancedLabels = []string{
	"architecture",
	"performance",
	"security",
	"concurrency",
	"breaking-change",
	"breaking change",
	"core",
	"optimization",
	"distributed",
	"rfc",
	"expert",
	"advanced",
	"level/advanced",
	"difficulty: hard",
	"race-condition",
	"deadlock",
	"memory-leak",
}

var AdvancedKeywords = []string{
	"concurrency", "race condition", "memory leak", "deadlock", "goroutine leak",
	"breaking change", "under the hood", "distributed consensus", "zero-copy",
	"re-architect", "cryptographic", "assembly", "kernel",
}

var BeginnerKeywords = []string{
	"good first issue", "easy fix", "typo", "update readme", "fix doc",
	"add test case", "documentation", "spelling", "sample",
}
