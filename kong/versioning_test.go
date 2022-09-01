package kong

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersioning_Version(t *testing.T) {
	t.Run("valid version string with three or four digits creates a version", func(t *testing.T) {
		tests := []struct {
			versionStr         string
			expectedVersionStr string
			expectedMajor      uint64
			expectedMinor      uint64
			expectedPatch      uint64
			expectedRevision   uint64
			isEnterprise       bool
			hasRevision        bool
		}{
			// Three digit versions OSS
			{
				versionStr:         "0.33.3",
				expectedVersionStr: "0.33.3",
				expectedMajor:      0,
				expectedMinor:      33,
				expectedPatch:      3,
			},
			{
				versionStr:         "1.2.0-alpha",
				expectedVersionStr: "1.2.0",
				expectedMajor:      1,
				expectedMinor:      2,
			},
			{
				versionStr:         "2.1.1+build.metadata",
				expectedVersionStr: "2.1.1",
				expectedMajor:      2,
				expectedMinor:      1,
				expectedPatch:      1,
			},
			{
				versionStr:         "2.6.10-alpha+build.metadata",
				expectedVersionStr: "2.6.10",
				expectedMajor:      2,
				expectedMinor:      6,
				expectedPatch:      10,
			},
			{
				versionStr:         "3.0.0+alpha-build.metadata",
				expectedVersionStr: "3.0.0",
				expectedMajor:      3,
				expectedMinor:      0,
				expectedPatch:      0,
			},
			// Three digit versions enterprise
			{
				versionStr:         "0.33.3-enterprise",
				expectedVersionStr: "0.33.3",
				expectedMajor:      0,
				expectedMinor:      33,
				expectedPatch:      3,
				isEnterprise:       true,
			},
			{
				versionStr:         "1.2.0-alpha-enterprise",
				expectedVersionStr: "1.2.0",
				expectedMajor:      1,
				expectedMinor:      2,
				isEnterprise:       true,
			},
			{
				versionStr:         "2.1.1+build.metadata-enterprise-edition",
				expectedVersionStr: "2.1.1",
				expectedMajor:      2,
				expectedMinor:      1,
				expectedPatch:      1,
				isEnterprise:       true,
			},
			{
				versionStr:         "2.6.10-alpha+build.metadataenterprise",
				expectedVersionStr: "2.6.10",
				expectedMajor:      2,
				expectedMinor:      6,
				expectedPatch:      10,
				isEnterprise:       true,
			},
			{
				versionStr:         "3.0.0+alphaenterpriseedition-build.metadata",
				expectedVersionStr: "3.0.0",
				expectedMajor:      3,
				expectedMinor:      0,
				expectedPatch:      0,
				isEnterprise:       true,
			},

			// Four digit versions enterprise
			{
				versionStr:         "0.33.3.1",
				expectedVersionStr: "0.33.3.1",
				expectedMajor:      0,
				expectedMinor:      33,
				expectedPatch:      3,
				expectedRevision:   1,
				isEnterprise:       true,
				hasRevision:        true,
			},
			{
				versionStr:         "2.8.1.3-1",
				expectedVersionStr: "2.8.1.3",
				expectedMajor:      2,
				expectedMinor:      8,
				expectedPatch:      1,
				expectedRevision:   3,
				isEnterprise:       true,
				hasRevision:        true,
			},
			{
				versionStr:         "3.0.0.0",
				expectedVersionStr: "3.0.0.0",
				expectedMajor:      3,
				expectedMinor:      0,
				expectedPatch:      0,
				expectedRevision:   0,
				isEnterprise:       true,
				hasRevision:        true,
			},
			{
				versionStr:         "1.2.0.0-alpha",
				expectedVersionStr: "1.2.0.0",
				expectedMajor:      1,
				expectedMinor:      2,
				expectedRevision:   0,
				isEnterprise:       true,
				hasRevision:        true,
			},
			{
				versionStr:         "2.1.1.3+build.metadata",
				expectedVersionStr: "2.1.1.3",
				expectedMajor:      2,
				expectedMinor:      1,
				expectedPatch:      1,
				expectedRevision:   3,
				isEnterprise:       true,
				hasRevision:        true,
			},
			{
				versionStr:         "2.6.10.1234-alpha+build.metadata",
				expectedVersionStr: "2.6.10.1234",
				expectedMajor:      2,
				expectedMinor:      6,
				expectedPatch:      10,
				expectedRevision:   1234,
				isEnterprise:       true,
				hasRevision:        true,
			},
			{
				versionStr:         "3.0.0.56+alpha-build.metadata",
				expectedVersionStr: "3.0.0.56",
				expectedMajor:      3,
				expectedMinor:      0,
				expectedPatch:      0,
				expectedRevision:   56,
				isEnterprise:       true,
				hasRevision:        true,
			},
		}

		for _, test := range tests {
			version, err := NewVersion(test.versionStr)
			require.NoError(t, err)
			require.Equal(t, test.versionStr, version.str)
			require.Equal(t, test.expectedVersionStr, version.String())
			require.Equal(t, test.expectedMajor, version.Major())
			require.Equal(t, test.expectedMinor, version.Minor())
			require.Equal(t, test.expectedPatch, version.Patch())
			revision, err := version.Revision()
			if test.hasRevision {
				require.Equal(t, test.expectedRevision, revision)
			} else {
				require.EqualError(t, err, "revision is unavailable for version")
			}
			require.Equal(t, test.isEnterprise, version.IsKongGatewayEnterprise())
			require.Empty(t, version.PreRelease()) // Pre-release has been removed after parsing
			require.Empty(t, version.Build())      // Build has been removed after parsing
		}
	})

	t.Run("invalid version string returns error when creating a new version", func(t *testing.T) {
		tests := []string{
			"3.0",
			"one.two.three",
			"one.two.three.four",
			"1.000.1",
			"1.001.1",
			"1.01.1.1",
			"1.1.1alpha",
			"1.1.1.1alpha",
		}

		for _, invalidVersion := range tests {
			version, err := NewVersion(invalidVersion)
			require.ErrorContains(t, err, invalidVersion)
			require.Equal(t, Version{}, version)
		}
	})
}

func TestVersioning_ForceNewVersion(t *testing.T) {
	t.Run("ensure valid version does not panic with three and four digit version", func(t *testing.T) {
		tests := []string{
			"1.2.3",
			"1.2.3.4",
		}
		for _, test := range tests {
			version := MustNewVersion(test)
			require.Equal(t, test, version.String())
		}
	})

	t.Run("ensure panic occurs for invalid version", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("panic test did not panic")
			}
		}()
		_ = MustNewVersion("invalid.version")
	})
}

func TestVersioning_Range(t *testing.T) {
	// Create two versions, three and four digits, for comparison checks with ranges
	v123, err := NewVersion("1.2.3")
	require.NoError(t, err)
	v1234, err := NewVersion("1.2.3.4")
	require.NoError(t, err)

	t.Run("valid range string with three or four digits creates a range comparison function", func(t *testing.T) {
		// Note: When comparing three digit versions to four digit versions the revision is ignored
		tests := []struct {
			rangeStr                 string
			expectedThreeDigitResult bool
			expectedFourDigitResult  bool
		}{
			// Three digit version ranges
			{
				rangeStr:                 "<= 2.0.0",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr:                 "<= 1.2.0",
				expectedThreeDigitResult: false,
				expectedFourDigitResult:  false,
			},
			{
				rangeStr:                 ">= 1.2.0",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr:                 "<= 1.2.3",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr:                 "== 1.2.3",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr: "!= 1.2.3",
			},
			{
				rangeStr:                 ">= 1.2.3",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr:                 ">= 1.2.3+build.metadata",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			// Note: Pre-release versions are concerned less than normal versions
			{
				rangeStr: "== 1.2.3-alpha",
			},
			{
				rangeStr:                 "> 1.2.3-alpha",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			// Four digit version ranges
			{
				rangeStr:                 "<= 2.0.0.0",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr:                 "<= 1.2.3.0",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  false,
			},
			{
				rangeStr:                 ">= 1.2.3.0",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr:                 "<= 1.2.3.4",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr:                 "== 1.2.3.4",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr:                 "!= 1.2.3.2",
				expectedThreeDigitResult: false,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr:                 ">= 1.2.3.4",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			{
				rangeStr:                 ">= 1.2.3.4+build.metadata",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			// Note: Pre-release versions are concerned less than normal versions
			{
				rangeStr: "== 1.2.3.4-alpha",
			},
			{
				rangeStr:                 "> 1.2.3.4-beta",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  true,
			},
			// Ensure v2 is within range, but three digit v1 is not satisfied
			{
				rangeStr:                 ">= 1.2.3 < 1.2.3.5",
				expectedThreeDigitResult: false,
				expectedFourDigitResult:  true,
			},
			// Ensure v1 is within range, but four digit v2 is not satisfied
			{
				rangeStr:                 ">= 1.2.3.5 < 1.2.5",
				expectedThreeDigitResult: true,
				expectedFourDigitResult:  false,
			},
			// Ensure v1 is excluded from range and four digit v2 is not satisfied
			{
				rangeStr:                 ">= 1.2.3.5 < 1.2.5 != 1.2.3",
				expectedThreeDigitResult: false,
				expectedFourDigitResult:  false,
			},
			// Ensure v1 is excluded by build meta range and four digit v2 is satisfied
			// since as range is valid for 1.2.3.4 using pre-release and build metadata
			// in ranges; pre-releases are considered less than normal/finalized versions
			{
				rangeStr:                 "> 1.2.3.4-alpha < 1.2.5 != 1.2.3.5+build.metadata",
				expectedThreeDigitResult: false,
				expectedFourDigitResult:  true,
			},
		}

		for _, test := range tests {
			rng, err := NewRange(test.rangeStr)
			require.NoError(t, err)
			require.Equal(t, test.expectedThreeDigitResult, rng(v123))
			require.Equal(t, test.expectedFourDigitResult, rng(v1234))
		}
	})

	t.Run("invalid range string returns error when creating a new range", func(t *testing.T) {
		tests := []string{
			"<= one.two.three",
			">= one.two.three.four",
			"!= 1.000.1 < 1.0.0",
			"~@ 1.001.1",
			"< 1.01.1.1",
			"> 1.1.1alpha",
			"== 1.1.1.1alpha",
			"$1.1.1",
		}

		for _, invalidRange := range tests {
			version, err := NewRange(invalidRange)
			require.Error(t, err)
			require.Nil(t, version)
		}
	})
}

func TestVersioning_ForceNewRange(t *testing.T) {
	t.Run("ensure valid range does not panic with three and four digit version", func(t *testing.T) {
		tests := []string{
			"<= 1.2.3",
			"<= 1.2.3.4",
		}
		for _, test := range tests {
			rng := MustNewRange(test)
			require.NotNil(t, rng)
		}
	})

	t.Run("ensure panic occurs for invalid range", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("panic test did not panic")
			}
		}()
		_ = MustNewRange("<= invalid.range")
	})
}
