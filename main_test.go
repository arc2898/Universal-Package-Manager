package main

import (
	"reflect"
	"testing"
)

func TestParseBulkArgs(t *testing.T) {
	got, err := parseBulkArgs([]string{"install", "apt", "git, vim", "snap", "code"})
	if err != nil {
		t.Fatalf("parseBulkArgs returned error: %v", err)
	}
	want := map[string][]string{
		"apt":  {"git", "vim"},
		"snap": {"code"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseBulkArgs() = %#v, want %#v", got, want)
	}
}

func TestParseBulkArgsRejectsIncompletePairs(t *testing.T) {
	if _, err := parseBulkArgs([]string{"install", "apt"}); err == nil {
		t.Fatal("expected incomplete pair to fail")
	}
	if _, err := parseBulkArgs([]string{"install", "apt", "git", "snap"}); err == nil {
		t.Fatal("expected trailing manager to fail")
	}
}

func TestParseBulkArgsRejectsEmptyPackage(t *testing.T) {
	if _, err := parseBulkArgs([]string{"remove", "apt", "git,,vim"}); err == nil {
		t.Fatal("expected empty package to fail")
	}
}