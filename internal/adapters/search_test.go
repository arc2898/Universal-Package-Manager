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

func TestParseSearchLinesRpm(t *testing.T) {
	results := parseSearchLines("bash - GNU Bourne Again Shell\ncoreutils - basic file utilities", "rpm")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "bash" || results[0].Source != "rpm" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Name != "coreutils" || results[1].Source != "rpm" {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestParseSearchLinesHomebrew(t *testing.T) {
	results := parseSearchLines("git - distributed version control\nvim - editor", "homebrew")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "git" || results[0].Source != "homebrew" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Name != "vim" || results[1].Source != "homebrew" {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestParseSearchLinesWinget(t *testing.T) {
	results := parseSearchLines("git - distributed version control\ncode - editor", "winget")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "git" || results[0].Source != "winget" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Name != "code" || results[1].Source != "winget" {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestParseSearchLinesChocolatey(t *testing.T) {
	results := parseSearchLines("git - distributed version control\ncode - editor", "chocolatey")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "git" || results[0].Source != "chocolatey" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Name != "code" || results[1].Source != "chocolatey" {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestParseSearchLinesScoop(t *testing.T) {
	results := parseSearchLines("git - distributed version control\ncode - editor", "scoop")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "git" || results[0].Source != "scoop" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Name != "code" || results[1].Source != "scoop" {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestParseSearchLinesCargo(t *testing.T) {
	results := parseSearchLines("git - distributed version control\ncode - editor", "cargo")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "git" || results[0].Source != "cargo" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Name != "code" || results[1].Source != "cargo" {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestParseSearchLinesNpm(t *testing.T) {
	results := parseSearchLines("git - distributed version control\ncode - editor", "npm")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "git" || results[0].Source != "npm" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Name != "code" || results[1].Source != "npm" {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestParseSearchLinesPip(t *testing.T) {
	results := parseSearchLines("git - distributed version control\ncode - editor", "pip")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "git" || results[0].Source != "pip" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Name != "code" || results[1].Source != "pip" {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestParseSearchLinesGem(t *testing.T) {
	results := parseSearchLines("git - distributed version control\ncode - editor", "gem")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "git" || results[0].Source != "gem" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Name != "code" || results[1].Source != "gem" {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestParseHomebrewSearchSections(t *testing.T) {
	results := parseHomebrewSearch("==> Formulae\nvlc\n==> Casks\nvlc")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2: %#v", len(results), results)
	}
	if results[0].Name != "vlc" || results[0].Description != "formula" {
		t.Fatalf("unexpected formula result: %#v", results[0])
	}
	if results[1].Name != "vlc" || results[1].Description != "cask" {
		t.Fatalf("unexpected cask result: %#v", results[1])
	}
}

func TestParseFlatpakSearchUsesApplicationID(t *testing.T) {
	results := parseFlatpakSearch("org.videolan.VLC VLC Multimedia player")
	if len(results) != 1 || results[0].Name != "org.videolan.VLC" || results[0].Source != "flatpak" || len(results[0].Aliases) != 1 || results[0].Aliases[0] != "VLC" {
		t.Fatalf("unexpected Flatpak results: %#v", results)
	}
}

func TestParsePortageSearchReturnsCategoryAtom(t *testing.T) {
	output := `*  app-misc/foo
      Latest version available: 1.2
      Description: A useful utility
*  media-video/vlc
      Description: VLC media player`
	results := parsePortageSearch(output)
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2: %#v", len(results), results)
	}
	if results[1].Name != "media-video/vlc" || results[1].Description != "VLC media player" {
		t.Fatalf("unexpected Portage result: %#v", results[1])
	}
}
