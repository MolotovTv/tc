package cmd

import "testing"

func TestRenameReleaseBranchForProdBase(t *testing.T) {
	for _, testCase := range []struct {
		in  string
		out string
	}{
		{"release-1.2.3", "1.2.3"},
		{"release-1.2.3.4", "1.2.3.4"},
		{"release-0.23.3", "0.23.3"},
		{"release-1.2.3-blabla", ""}, // theres a bug in ci with -xxx releases, better fail than build wrong version
		{"1.2.3", ""},
		{"release", ""},
	} {
		out := renameBranchForProd(testCase.in)
		if out != testCase.out {
			t.Fatalf("testcase %s fail : %s instead of %s", testCase.in, out, testCase.out)
		}
	}

}
