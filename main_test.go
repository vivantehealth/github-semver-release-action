package main

import "testing"

func TestCalculateVersion(t *testing.T) {
	tests := []struct {
		name        string
		latestTag   string
		commitMsg   string
		defaultBump string
		vPrefix     bool
		wantTag     string
		wantMajor   string
		wantMinor   string
		wantPart    string
	}{
		{
			name:        "first release with no tags, default minor",
			latestTag:   "",
			commitMsg:   "initial commit",
			defaultBump: "minor",
			vPrefix:     true,
			wantTag:     "v0.1.0",
			wantMajor:   "v0",
			wantMinor:   "v0.1",
			wantPart:    "minor",
		},
		{
			name:        "first release with default patch",
			latestTag:   "",
			commitMsg:   "initial commit",
			defaultBump: "patch",
			vPrefix:     true,
			wantTag:     "v0.0.1",
			wantMajor:   "v0",
			wantMinor:   "v0.0",
			wantPart:    "patch",
		},
		{
			name:        "minor bump from existing tag",
			latestTag:   "v1.2.3",
			commitMsg:   "new feature",
			defaultBump: "minor",
			vPrefix:     true,
			wantTag:     "v1.3.0",
			wantMajor:   "v1",
			wantMinor:   "v1.3",
			wantPart:    "minor",
		},
		{
			name:        "patch bump from existing tag",
			latestTag:   "v1.2.3",
			commitMsg:   "bugfix",
			defaultBump: "patch",
			vPrefix:     true,
			wantTag:     "v1.2.4",
			wantMajor:   "v1",
			wantMinor:   "v1.2",
			wantPart:    "patch",
		},
		{
			name:        "#major overrides default",
			latestTag:   "v1.2.3",
			commitMsg:   "breaking change #major",
			defaultBump: "patch",
			vPrefix:     true,
			wantTag:     "v2.0.0",
			wantMajor:   "v2",
			wantMinor:   "v2.0",
			wantPart:    "major",
		},
		{
			name:        "#minor overrides default",
			latestTag:   "v1.2.3",
			commitMsg:   "new feature #minor",
			defaultBump: "patch",
			vPrefix:     true,
			wantTag:     "v1.3.0",
			wantMajor:   "v1",
			wantMinor:   "v1.3",
			wantPart:    "minor",
		},
		{
			name:        "#patch overrides default",
			latestTag:   "v1.2.3",
			commitMsg:   "fix #patch",
			defaultBump: "major",
			vPrefix:     true,
			wantTag:     "v1.2.4",
			wantMajor:   "v1",
			wantMinor:   "v1.2",
			wantPart:    "patch",
		},
		{
			name:        "#none skips release",
			latestTag:   "v1.2.3",
			commitMsg:   "docs only #none",
			defaultBump: "minor",
			vPrefix:     true,
			wantTag:     "",
			wantPart:    "none",
		},
		{
			name:        "no v-prefix",
			latestTag:   "1.0.0",
			commitMsg:   "next",
			defaultBump: "minor",
			vPrefix:     false,
			wantTag:     "1.1.0",
			wantMajor:   "1",
			wantMinor:   "1.1",
			wantPart:    "minor",
		},
		{
			name:        "keyword in commit body",
			latestTag:   "v1.2.3",
			commitMsg:   "some feature\n\nThis is a detailed description.\n#major",
			defaultBump: "patch",
			vPrefix:     true,
			wantTag:     "v2.0.0",
			wantMajor:   "v2",
			wantMinor:   "v2.0",
			wantPart:    "major",
		},
		{
			name:        "case insensitive keyword",
			latestTag:   "v1.0.0",
			commitMsg:   "change #MAJOR",
			defaultBump: "patch",
			vPrefix:     true,
			wantTag:     "v2.0.0",
			wantMajor:   "v2",
			wantMinor:   "v2.0",
			wantPart:    "major",
		},
		{
			name:        "multiple keywords takes highest priority",
			latestTag:   "v1.2.3",
			commitMsg:   "change #patch #major",
			defaultBump: "minor",
			vPrefix:     true,
			wantTag:     "v2.0.0",
			wantMajor:   "v2",
			wantMinor:   "v2.0",
			wantPart:    "major",
		},
		{
			name:        "large version numbers",
			latestTag:   "v99.99.99",
			commitMsg:   "bump",
			defaultBump: "patch",
			vPrefix:     true,
			wantTag:     "v99.99.100",
			wantMajor:   "v99",
			wantMinor:   "v99.99",
			wantPart:    "patch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateVersion(tt.latestTag, tt.commitMsg, tt.defaultBump, tt.vPrefix)

			if result.NewTag != tt.wantTag {
				t.Errorf("NewTag = %q, want %q", result.NewTag, tt.wantTag)
			}
			if result.Part != tt.wantPart {
				t.Errorf("Part = %q, want %q", result.Part, tt.wantPart)
			}
			if tt.wantPart != "none" {
				if result.Major != tt.wantMajor {
					t.Errorf("Major = %q, want %q", result.Major, tt.wantMajor)
				}
				if result.Minor != tt.wantMinor {
					t.Errorf("Minor = %q, want %q", result.Minor, tt.wantMinor)
				}
			}
		})
	}
}

func TestDetectBump(t *testing.T) {
	tests := []struct {
		name        string
		commitMsg   string
		defaultBump string
		want        string
	}{
		{"major keyword", "breaking #major change", "minor", "major"},
		{"minor keyword", "feature #minor", "patch", "minor"},
		{"patch keyword", "fix #patch", "major", "patch"},
		{"none keyword", "docs #none", "minor", "none"},
		{"no keyword uses default", "just a commit", "patch", "patch"},
		{"case insensitive", "#MINOR update", "patch", "minor"},
		{"major wins over patch", "#patch and #major", "minor", "major"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectBump(tt.commitMsg, tt.defaultBump); got != tt.want {
				t.Errorf("detectBump() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name                       string
		tag                        string
		wantMajor, wantMinor, wantPatch uint64
	}{
		{"with v prefix", "v1.2.3", 1, 2, 3},
		{"without v prefix", "1.2.3", 1, 2, 3},
		{"empty string", "", 0, 0, 0},
		{"invalid string", "notaversion", 0, 0, 0},
		{"zeros", "v0.0.0", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := parseVersion(tt.tag)
			if v.Major() != tt.wantMajor || v.Minor() != tt.wantMinor || v.Patch() != tt.wantPatch {
				t.Errorf("parseVersion(%q) = %d.%d.%d, want %d.%d.%d",
					tt.tag, v.Major(), v.Minor(), v.Patch(),
					tt.wantMajor, tt.wantMinor, tt.wantPatch)
			}
		})
	}
}
