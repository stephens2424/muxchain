package muxchainutil

import "testing"

func TestPathMatch(t *testing.T) {
	testPathMatch("/x/", "/x/", true, t)
	testPathMatch("/x/*", "/x/", true, t)
	testPathMatch("/x/*", "/x", true, t)
	testPathMatch("/x/*/*", "/x", true, t)
	testPathMatch("/x/*/*", "/x/y", true, t)
	testPathMatch("/x/*", "/x/y/z", true, t)
	testPathMatch("/x/*/z", "/x/y/z", true, t)
	testPathMatch("/x/*/z", "/x/z", false, t)
}

func testPathMatch(pattern, path string, expect bool, t *testing.T) {
	if pathMatch(pattern, path) != expect {
		if expect {
			t.Logf("%s should have matched pattern %s", path, pattern)
		} else {
			t.Logf("%s should not have matched pattern %s", path, pattern)
		}
		t.Fail()
	}
}
