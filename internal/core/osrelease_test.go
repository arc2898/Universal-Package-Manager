package core

import "testing"

func TestNativeManagerPreferenceFromOSRelease(t *testing.T) {
	cases := []struct {
		name string
		data string
		want string
	}{
		{"debian", "ID=debian\n", "apt"},
		{"ubuntu", "ID=ubuntu\nID_LIKE=debian\n", "apt"},
		{"arch", "ID=arch\n", "pacman"},
		{"fedora", "ID=fedora\n", "dnf"},
		{"opensuse", "ID=opensuse-tumbleweed\n", "zypper"},
		{"void", "ID=void\n", "xbps"},
		{"gentoo", "ID=gentoo\n", "portage"},
		{"suse-like", `ID=my-suse\nID_LIKE="suse"\n`, "zypper"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := nativeManagerPreferenceFromOSRelease(tc.data)
			if len(got) == 0 || got[0] != tc.want {
				t.Fatalf("preference = %v, want first manager %q", got, tc.want)
			}
		})
	}
}
