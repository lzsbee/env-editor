package update

import "testing"

func TestIsNewer(t *testing.T) {
	tests := []struct {
		latest  string
		current string
		want    bool
	}{
		{"1.2.0", "1.1.0", true},
		{"v1.1.1", "1.1.0", true},
		{"1.1.0", "1.1.0", false},
		{"1.0.9", "1.1.0", false},
		{"2.0.0", "1.9.9", true},
	}

	for _, tt := range tests {
		if got := isNewer(tt.latest, tt.current); got != tt.want {
			t.Fatalf("isNewer(%q, %q) = %v, want %v", tt.latest, tt.current, got, tt.want)
		}
	}
}
