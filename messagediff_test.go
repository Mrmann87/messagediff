package messagediff

import (
	"testing"
	"time"

	"github.com/Mrmann87/messagediff/testdata"
)

type testStruct struct {
	A, b int
	C    []int
	D    [3]int
}

type RecursiveStruct struct {
	Key   int
	Child *RecursiveStruct
}

func newRecursiveStruct(key int) *RecursiveStruct {
	a := &RecursiveStruct{
		Key: key,
	}
	b := &RecursiveStruct{
		Key:   key,
		Child: a,
	}
	a.Child = b
	return a
}

type testCase struct {
	a, b  interface{}
	diff  string
	equal bool
}

func checkTestCases(t *testing.T, testData []testCase, opt ...Option) {
	for i, td := range testData {
		diff, equal := PrettyDiff(td.a, td.b, opt...)
		if diff != td.diff {
			t.Errorf("%d. PrettyDiff(%#v, %#v) diff = %#v; not %#v", i, td.a, td.b, diff, td.diff)
		}
		if equal != td.equal {
			t.Errorf("%d. PrettyDiff(%#v, %#v) equal = %#v; not %#v", i, td.a, td.b, equal, td.equal)
		}
	}
}

func TestPrettyDiff(t *testing.T) {
	testData := []testCase{
		{
			true,
			false,
			"\x1b[1m---a\x1b[0m\n\x1b[1m+++b\x1b[0m\n\x1b[22;31m-true\x1b[0m\n\x1b[22;32m+false\x1b[0m\n",
			false,
		},
		{
			true,
			0,
			"\x1b[1m---a\x1b[0m\n\x1b[1m+++b\x1b[0m\n\x1b[22;31m-true\x1b[0m\n\x1b[22;32m+0\x1b[0m\n",
			false,
		},
		{
			[]int{0, 1, 2},
			[]int{0, 1, 2, 3},
			"\x1b[1m---a[3]\x1b[0m\n\x1b[1m+++b[3]\x1b[0m\n\x1b[22;32m+3\x1b[0m\n",
			false,
		},
		{
			[]int{0, 1, 2, 3},
			[]int{0, 1, 2},
			"\x1b[1m---a[3]\x1b[0m\n\x1b[1m+++b[3]\x1b[0m\n\x1b[22;31m-3\x1b[0m\n",
			false,
		},
		{
			[]int{0},
			[]int{1},
			"\x1b[1m---a[0]\x1b[0m\n\x1b[1m+++b[0]\x1b[0m\n\x1b[22;31m-0\x1b[0m\n\x1b[22;32m+1\x1b[0m\n",
			false,
		},
		{
			&[]int{0},
			&[]int{1},
			"\x1b[1m---a[0]\x1b[0m\n\x1b[1m+++b[0]\x1b[0m\n\x1b[22;31m-0\x1b[0m\n\x1b[22;32m+1\x1b[0m\n",
			false,
		},
		{
			map[string]int{"a": 1, "b": 2},
			map[string]int{"b": 4, "c": 3},
			"\x1b[1m---a[\"a\"]\x1b[0m\n\x1b[1m+++b[\"a\"]\x1b[0m\n\x1b[22;31m-1\x1b[0m\n\x1b[1m---a[\"b\"]\x1b[0m\n\x1b[1m+++b[\"b\"]\x1b[0m\n\x1b[22;31m-2\x1b[0m\n\x1b[22;32m+4\x1b[0m\n\x1b[1m---a[\"c\"]\x1b[0m\n\x1b[1m+++b[\"c\"]\x1b[0m\n\x1b[22;32m+3\x1b[0m\n",
			false,
		},
		{
			testStruct{1, 2, []int{1}, [3]int{4, 5, 6}},
			testStruct{1, 3, []int{1, 2}, [3]int{4, 5, 6}},
			"\x1b[1m---a.C[1]\x1b[0m\n\x1b[1m+++b.C[1]\x1b[0m\n\x1b[22;32m+2\x1b[0m\n\x1b[1m---a.b\x1b[0m\n\x1b[1m+++b.b\x1b[0m\n\x1b[22;31m-2\x1b[0m\n\x1b[22;32m+3\x1b[0m\n",
			false,
		},
		{
			nil,
			nil,
			"",
			true,
		},
		{
			&struct{}{},
			nil,
			"\x1b[1m---a\x1b[0m\n\x1b[1m+++b\x1b[0m\n\x1b[22;31m-&struct {}{}\x1b[0m\n",
			false,
		},
		{
			nil,
			&struct{}{},
			"\x1b[1m---a\x1b[0m\n\x1b[1m+++b\x1b[0m\n\x1b[22;32m+&struct {}{}\x1b[0m\n",
			false,
		},
		{
			time.Time{},
			time.Time{},
			"",
			true,
		},
		{
			testdata.MakeTest(10, "duck"),
			testdata.MakeTest(20, "foo"),
			"\x1b[1m---a.a\x1b[0m\n\x1b[1m+++b.a\x1b[0m\n\x1b[22;31m-10\x1b[0m\n\x1b[22;32m+20\x1b[0m\n\x1b[1m---a.b\x1b[0m\n\x1b[1m+++b.b\x1b[0m\n\x1b[22;31m-\"duck\"\x1b[0m\n\x1b[22;32m+\"foo\"\x1b[0m\n",
			false,
		},
		{
			time.Date(2018, 7, 24, 14, 06, 59, 0, &time.Location{}),
			time.Date(2018, 7, 24, 14, 06, 59, 0, time.UTC),
			"",
			true,
		},
		{
			time.Date(2017, 1, 1, 0, 0, 0, 0, &time.Location{}),
			time.Date(2018, 7, 24, 14, 06, 59, 0, time.UTC),
			"\x1b[1m---a\x1b[0m\n\x1b[1m+++b\x1b[0m\n\x1b[22;31m-\"2017-01-01 00:00:00 +0000 UTC\"\x1b[0m\n\x1b[22;32m+\"2018-07-24 14:06:59 +0000 UTC\"\x1b[0m\n",
			false,
		},
	}
	checkTestCases(t, testData)
}

func TestPrettyDiffRecursive(t *testing.T) {
	testData := []testCase{
		{
			newRecursiveStruct(1),
			newRecursiveStruct(1),
			"",
			true,
		},
		{
			newRecursiveStruct(1),
			newRecursiveStruct(2),
			"\x1b[1m---a.Child.Key\x1b[0m\n\x1b[1m+++b.Child.Key\x1b[0m\n\x1b[22;31m-1\x1b[0m\n\x1b[22;32m+2\x1b[0m\n\x1b[1m---a.Key\x1b[0m\n\x1b[1m+++b.Key\x1b[0m\n\x1b[22;31m-1\x1b[0m\n\x1b[22;32m+2\x1b[0m\n",
			false,
		},
	}
	checkTestCases(t, testData)
}

func TestPathString(t *testing.T) {
	testData := []struct {
		in   Path
		want string
	}{{
		Path{StructField("test"), SliceIndex(1), MapKey{"blue"}, MapKey{12.3}},
		".test[1][\"blue\"][12.3]",
	}}
	for i, td := range testData {
		if out := td.in.String(); out != td.want {
			t.Errorf("%d. %#v.String() = %#v; not %#v", i, td.in, out, td.want)
		}
	}
}

type ignoreStruct struct {
	A int `testdiff:"ignore"`
	a int
	B [3]int `testdiff:"ignore"`
	b [3]int
}

func TestIgnoreTag(t *testing.T) {
	testData := []testCase{
		{
			ignoreStruct{1, 1, [3]int{1, 2, 3}, [3]int{4, 5, 6}},
			ignoreStruct{2, 1, [3]int{1, 8, 3}, [3]int{4, 5, 6}},
			"",
			true,
		},
		{
			ignoreStruct{1, 1, [3]int{1, 2, 3}, [3]int{4, 5, 6}},
			ignoreStruct{2, 2, [3]int{1, 8, 3}, [3]int{4, 9, 6}},
			"\x1b[1m---a.a\x1b[0m\n\x1b[1m+++b.a\x1b[0m\n\x1b[22;31m-1\x1b[0m\n\x1b[22;32m+2\x1b[0m\n\x1b[1m---a.b[1]\x1b[0m\n\x1b[1m+++b.b[1]\x1b[0m\n\x1b[22;31m-5\x1b[0m\n\x1b[22;32m+9\x1b[0m\n",
			false,
		},
	}
	checkTestCases(t, testData)
}

func TestIgnoreStructFieldOption(t *testing.T) {
	testData := []testCase{
		{
			struct {
				X string
				Y string
			}{
				"x",
				"y",
			},
			struct {
				X string
				Y string
			}{
				"xx",
				"y",
			},
			"",
			false,
		},
	}

	testData[0].equal = true
	checkTestCases(t, testData, IgnoreStructField("X"))

	testData[0].diff = "\x1b[1m---a.X\x1b[0m\n\x1b[1m+++b.X\x1b[0m\n\x1b[22;31m-\"x\"\x1b[0m\n\x1b[22;32m+\"xx\"\x1b[0m\n"
	testData[0].equal = false
	checkTestCases(t, testData, IgnoreStructField("Y"))
}
