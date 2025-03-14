package kong

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kong/semver/v4"
)

// Version represents a three or four digit version.
type Version struct {
	// version represents a three of four digit version using the forked kong/semver
	// library.
	version semver.Version
	// str represents the original version string when creating Version.
	str string
}

// Range represents a range of versions which can be used to validate if a Version is valid
// for a Range.
type Range func(Version) bool

// NewVersion creates a new instance of a Version.
func NewVersion(versionStr string) (Version, error) {
	version, err := semver.Parse(versionStr)
	if err != nil {
		return Version{}, fmt.Errorf("unable to create version: %w", err)
	}

	// Remove pre-release and build metadata versioning for finalized version comparisons
	version.Pre = []semver.PRVersion{}
	version.Build = []string{}

	return Version{
		version: version,
		str:     versionStr,
	}, nil
}

// MustNewVersion creates a new instance of a Version; however it will panic if it cannot
// be created.
func MustNewVersion(versionStr string) Version {
	version, err := NewVersion(versionStr)
	if err != nil {
		panic(err)
	}
	return version
}

// Major returns the major digit of Version.
func (v Version) Major() uint64 {
	return v.version.Major
}

// Minor returns the minor digit of Version.
func (v Version) Minor() uint64 {
	return v.version.Minor
}

// Patch returns the patch digit of Version.
func (v Version) Patch() uint64 {
	return v.version.Patch
}

// Revision returns the revision digit of Version; if revision has not been set then an error
// will be returned.
func (v Version) Revision() (uint64, error) {
	if v.version.Revision < 0 {
		return 0, errors.New("revision is unavailable for version")
	}
	return uint64(v.version.Revision), nil //nolint:gosec
}

// PreRelease returns the pre-release string of Version.
func (v Version) PreRelease() string {
	var preReleaseStr strings.Builder
	if len(v.version.Pre) > 0 {
		fmt.Fprintf(&preReleaseStr, "-%s", v.version.Pre[0].String())
		for _, preRelease := range v.version.Pre[1:] {
			preReleaseStr.WriteString(".")
			preReleaseStr.WriteString(preRelease.String())
		}
	}
	return preReleaseStr.String()
}

// Build returns the build metadata string of Version.
func (v Version) Build() string {
	var buildStr strings.Builder
	if len(v.version.Build) > 0 {
		fmt.Fprintf(&buildStr, "+%s", v.version.Build[0])
		for _, build := range v.version.Build[1:] {
			buildStr.WriteString(".")
			buildStr.WriteString(build)
		}
	}
	return buildStr.String()
}

// String returns the textual or display value of the Version.
func (v Version) String() string {
	return v.version.String()
}

// IsKongGatewayEnterprise determines if a Version represents a Kong Gateway enterprise edition.
func (v Version) IsKongGatewayEnterprise() bool {
	return v.version.Revision >= 0 || strings.Contains(v.str, "enterprise")
}

// NewRange creates an instance of a Range.
// Valid ranges can consist of multiple comparisons and three/four digit versions:
//   - "<1.0.0" || "<v1.0.0.0"
//   - "<=1.0.0" || "<=1.0.0.0"
//   - ">1.0.0" || ">1.0.0.0"
//   - ">=1.0.0" || >= 1.0.0.0
//   - "1.0.0", "=1.0.0", "==1.0.0" || "1.0.0.0", "=1.0.0.0", "==1.0.0.0"
//   - "!1.0.0", "!=1.0.0" || "!1.0.0.0", "!=1.0.0.0"
//
// A Range can consist of multiple ranges separated by space:
// Ranges can be linked by logical AND:
//   - ">1.0.0 <2.0.0" would match between both ranges, so "1.1.1" and "1.8.7" but not "1.0.0" or "2.0.0"
//   - ">1.0.0 <3.0.0 !2.0.3-beta.2" would match every version between 1.0.0 and 3.0.0 except 2.0.3-beta.2
//
// Four digit versions can be used in ranges with three digit version and linked by logical AND:
//   - ">1.0.0 <2.0.0.0" would match between both ranges, so "1.0.0.1" and "1.8.7" but not "1.0.0", "2.0.0"
//   - ">1.0.0 <3.0.0 !2.0.3.0-beta.2" would match every version between 1.0.0 and 3.0.0 except 2.0.3-beta.2 and
//     2.0.3.0-beta2
//
// Ranges can also be linked by logical OR:
//   - "<2.0.0 || >=3.0.0" would match "1.x.x" and "3.x.x" but not "2.x.x"
//
// Four digit versions can be used in ranges with three digit version and linked by logical OR:
//
// AND has a higher precedence than OR. It's not possible to use brackets.
//
// Ranges can be combined by both AND and OR:
//   - ">1.0.0 <2.0.0.0 || >3.0.0 !4.2.1" would match "1.2.3", "1.0.0.1", "1.9.9", "3.1.1", but not "4.2.1", "2.1.1"
func NewRange(rangeStr string) (Range, error) {
	rng, err := semver.ParseRange(rangeStr)
	if err != nil {
		return nil, fmt.Errorf("unable to create range: %w", err)
	}
	return Range(func(version Version) bool {
		return rng(version.version)
	}), nil
}

// MustNewRange creates a new instance of a Version; however it will panic if it cannot
// be created.
func MustNewRange(rangeStr string) Range {
	rng, err := NewRange(rangeStr)
	if err != nil {
		panic(err)
	}
	return rng
}
