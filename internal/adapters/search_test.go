package adapters

import "testing"

func TestParseSearchLines(t *testing.T) {
	results := parseSearchLines("git - distributed version control\n\nvim - editor", "apt")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "git" || results[0].Description != "distributed version control" || results[0].Source != "apt" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
}

func TestParseSearchLinesSkipsMetadataNotice(t *testing.T) {
	results := parseSearchLines("Last metadata expiration check: 0:01:00 ago\nfoo - package", "dnf")
	if len(results) != 1 || results[0].Name != "foo" {
		t.Fatalf("unexpected results: %#v", results)
	}
}
