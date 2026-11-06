package steam

import (
	"os"
	"testing"

	"github.com/nning/protonutils/config"
)

// testSteamRoot contains a copy of real Steam data (decrypted from
// test/root.tar.xz.age, see test/decrypt.sh). Tests using it are skipped if it
// is not available.
const testSteamRoot = "../test/root"

var testConfig = &config.Config{
	SteamRoot: testSteamRoot,
}

func requireTestRoot(t *testing.T) {
	t.Helper()

	if _, err := os.Stat(testSteamRoot); err != nil {
		t.Skip("encrypted test data not available (see test/decrypt.sh)")
	}
}
