package version

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/params"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags).
	GitCommit    = ""
	GitDate      = ""
	VersionMajor = "0"
	VersionMinor = "0"
	VersionPatch = "0"
	VersionMeta  = ""
)

func init() {
	vmaj, _ := strconv.Atoi(VersionMajor)
	vmin, _ := strconv.Atoi(VersionMinor)
	vpatch, _ := strconv.Atoi(VersionPatch)
	params.VersionMajor = vmaj       // Major version component of the current release
	params.VersionMinor = vmin       // Minor version component of the current release
	params.VersionPatch = vpatch     // Patch version component of the current release
	params.VersionMeta = VersionMeta // Version metadata to append to the version string
}

func VersionWithCommit() string {
	vsn := params.VersionWithMeta()
	if len(GitCommit) >= 0 {
		vsn += "-" + GitDate
	}
	if (params.VersionMeta != "stable") && (GitDate != "") {
		vsn += "-" + GitDate
	}
	return vsn
}

func BigToString(b *big.Int) string {
	if len(b.Bytes()) > 8 {
		return "_malformed_version_"
	}
	return U64ToString(b.Uint64())
}

func AsString() string {
	return ToString(uint16(params.VersionMajor), uint16(params.VersionMinor), uint16(params.VersionPatch))
}

func AsU64() uint64 {
	return ToU64(uint16(params.VersionMajor), uint16(params.VersionMinor), uint16(params.VersionPatch))
}

func AsBigInt() *big.Int {
	return new(big.Int).SetUint64(AsU64())
}

func ToU64(vMajor, vMinor, vPatch uint16) uint64 {
	return uint64(vMajor)*1e12 + uint64(vMinor)*1e6 + uint64(vPatch)
}

func ToString(major, minor, patch uint16) string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

func U64ToString(v uint64) string {
	return ToString(uint16((v/1e12)%1e6), uint16((v/1e6)%1e6), uint16(v%1e6))
}
