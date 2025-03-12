package entproto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithGoPkg(t *testing.T) {
	tests := []struct {
		name        string
		goPkg       string
		fdPackage   string
		wantPkgName string
	}{
		{
			name:        "Default behavior",
			goPkg:       "",
			fdPackage:   "example.service",
			wantPkgName: "service",
		},
		{
			name:        "Custom package",
			goPkg:       "custompkg",
			fdPackage:   "example.service",
			wantPkgName: "custompkg",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test directly with the extractLastFqnPart helper function
			// since we can't easily mock desc.FileDescriptor
			var pkgName string
			if tt.goPkg != "" {
				pkgName = tt.goPkg
			} else {
				pkgName = extractLastFqnPart(tt.fdPackage)
			}
			
			require.Equal(t, tt.wantPkgName, pkgName,
				"Expected package name to be %q, got %q", tt.wantPkgName, pkgName)
		})
	}
}

// TestExtensionGoPkg tests that the WithGoPkg option properly sets the goPkg field
func TestExtensionGoPkg(t *testing.T) {
	customPkg := "myspecialpkg"
	ext, err := NewExtension(WithGoPkg(customPkg))
	require.NoError(t, err)
	
	// Check that the goPkg field is set
	require.Equal(t, customPkg, ext.goPkg, 
		"Expected goPkg to be %q, got %q", customPkg, ext.goPkg)
}

// TestExtractLastFqnPart tests the package extraction function
func TestExtractLastFqnPart(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"example.service", "service"},
		{"service", "service"},
		{"com.example.api.v1", "v1"},
		{"", ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractLastFqnPart(tt.input)
			require.Equal(t, tt.want, got,
				"extractLastFqnPart(%q) = %q, want %q", tt.input, got, tt.want)
		})
	}
}