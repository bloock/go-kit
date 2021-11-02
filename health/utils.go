package health

import (
	"runtime/debug"
	"strings"

	"github.com/icza/bitio"
)

func DepVersion(d string) string {
	_ = bitio.NewReader
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "Unknown"
	}

	for _, dep := range bi.Deps {
		if strings.Contains(dep.Path, d) {
			return dep.Version
		}
	}
	return "Unknown"
}
