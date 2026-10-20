package core

import "testing"

func TestValidatePackageName(t *testing.T) {
	for _, name := range []string{"git", "@scope/package", "python3-pip", "foo.bar"} {
		if err := validatePackageName(name); err != nil {
			t.Errorf("validatePackageName(%q) unexpectedly failed: %v", name, err)
		}
	}
	for _, name := range []string{"", "  ", "--help", "foo bar", "foo\nbar", "foo\x00bar"} {
		if err := validatePackageName(name); err == nil {
			t.Errorf("validatePackageName(%q) expected an error", name)
		}
	}
}
