package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nning/protonutils/internal/steamtest"
	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	steam.StoreAPIURL = "http://127.0.0.1:1/"
}

func resetFlags(c *cobra.Command) {
	c.Flags().VisitAll(func(f *pflag.Flag) {
		if err := f.Value.Set(f.DefValue); err != nil {
			panic(err)
		}
		f.Changed = false
	})

	for _, sub := range c.Commands() {
		resetFlags(sub)
	}
}

// run executes protonutils with args and returns its output
func run(t *testing.T, args ...string) (string, error) {
	t.Helper()

	resetFlags(rootCmd)
	s, cfg, root = nil, nil, ""

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetIn(strings.NewReader(""))
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	return out.String(), err
}

func mustRun(t *testing.T, args ...string) string {
	t.Helper()

	out, err := run(t, args...)
	require.NoError(t, err, out)
	return out
}

// setup creates a synthetic Steam root and configures it
func setup(t *testing.T, r steamtest.Root) string {
	t.Helper()

	dir := steamtest.Setup(t, r)
	mustRun(t, "config", "steam_root", dir)
	return dir
}

func Test_HelpWithoutSteam(t *testing.T) {
	steamtest.Setup(t, steamtest.Default())
	mustRun(t, "config", "steam_root", "/nonexistent")

	for _, args := range [][]string{{}, {"help"}, {"--help"}, {"list", "--help"}, {"version"}, {"config"}} {
		out, err := run(t, args...)
		assert.NoError(t, err, args)
		assert.NotEmpty(t, out, args)
	}

	// The completion command captures its output writer on first use, so only
	// the error is checked here
	_, err := run(t, "completion", "bash")
	assert.NoError(t, err)

	_, err = run(t, "list")
	assert.ErrorContains(t, err, "not a directory")
}

func Test_Config(t *testing.T) {
	steamtest.Setup(t, steamtest.Default())

	mustRun(t, "config", "user", "bob")
	assert.Equal(t, "bob\n", mustRun(t, "config", "user"))

	_, err := run(t, "config", "foo", "bar")
	assert.ErrorContains(t, err, "unknown config key")

	_, err = run(t, "config", "enable_viewed_steam_play", "maybe")
	assert.ErrorContains(t, err, "invalid boolean")

	mustRun(t, "config", "enable_viewed_steam_play", "true")
	assert.Contains(t, mustRun(t, "config"), "enable_viewed_steam_play: true")

	p := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "protonutils", "protonutils.yml")
	assert.FileExists(t, p)
}

func Test_UserFlag(t *testing.T) {
	setup(t, steamtest.Default())

	out := mustRun(t, "list")
	assert.Contains(t, out, "Half-Life")
	assert.Contains(t, out, "Cyberpunk 2077")

	// bob only has Portal 2 in localconfig.vdf
	out = mustRun(t, "list", "-u", "bob")
	assert.Contains(t, out, "Portal 2")
	assert.NotContains(t, out, "Cyberpunk 2077")

	out = mustRun(t, "list", "--user", "alice")
	assert.Contains(t, out, "Cyberpunk 2077")

	_, err := run(t, "list", "-u", "mallory")
	assert.ErrorContains(t, err, "Steam user not found")
}

func Test_IgnoreCacheFlag(t *testing.T) {
	setup(t, steamtest.Default())

	cacheFile := filepath.Join(os.Getenv("XDG_CACHE_HOME"), "protonutils", "appids.json")
	require.NoError(t, os.MkdirAll(filepath.Dir(cacheFile), 0700))
	require.NoError(t, os.WriteFile(cacheFile, []byte(`{"620":{"name":"Cached Name","valid":true,"updatedAt":1}}`), 0600))

	assert.Contains(t, mustRun(t, "list"), "Cached Name")

	out := mustRun(t, "list", "-c")
	assert.NotContains(t, out, "Cached Name")
	assert.Contains(t, out, "Portal 2")
}

func Test_ListJSON(t *testing.T) {
	setup(t, steamtest.Default())

	var tools map[string]struct {
		Name  string `json:"name"`
		Games map[string]struct {
			AppID string `json:"appID"`
			Name  string `json:"name"`
		} `json:"games"`
	}

	require.NoError(t, json.Unmarshal([]byte(mustRun(t, "list", "-j")), &tools))
	assert.Equal(t, "620", tools["proton_63"].Games["Portal 2"].AppID)
	assert.Equal(t, "Portal 2", tools["proton_63"].Games["Portal 2"].Name)
}

func Test_PathsWithoutPreviousList(t *testing.T) {
	dir := setup(t, steamtest.Default())

	out := mustRun(t, "compatdata", "path", "portal", "2")
	assert.Equal(t, filepath.Join(dir, "library", "steamapps", "compatdata", "620")+"\n", out)

	out = mustRun(t, "installdir", "path", "-v", "cyberpunk")
	assert.Equal(t, "Cyberpunk 2077\n"+filepath.Join(dir, "library", "steamapps", "common", "Cyberpunk 2077")+"\n", out)

	_, err := run(t, "installdir", "path", "p")
	var ambiguous *steam.AmbiguousNameError
	assert.True(t, errors.As(err, &ambiguous))

	var buf bytes.Buffer
	printError(&buf, err)
	assert.Contains(t, buf.String(), "620  Portal 2")
	assert.Contains(t, buf.String(), "Ambiguous name")

	out = mustRun(t, "appid", "portal", "2")
	assert.Equal(t, "       620  Portal 2\n", out)
}

func toolDirs(t *testing.T, dir string) []string {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join(dir, "compatibilitytools.d"))
	require.NoError(t, err)

	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}

	return names
}

func Test_CompatToolClean(t *testing.T) {
	dir := setup(t, steamtest.Default())

	out := mustRun(t, "compattool", "clean", "-y")
	assert.Contains(t, out, "GE-Proton6-1")
	assert.Contains(t, out, "GE-Proton7-1")

	// GE-Proton9-1 is the default and must be kept
	assert.Equal(t, []string{"GE-Proton8-1", "GE-Proton9-1"}, toolDirs(t, dir))

	assert.Contains(t, mustRun(t, "compattool", "clean"), "No unused compatibility tool found")
}

func Test_CompatToolCleanAborted(t *testing.T) {
	dir := setup(t, steamtest.Default())

	out := mustRun(t, "compattool", "clean")
	assert.Contains(t, out, "Aborted")
	assert.Len(t, toolDirs(t, dir), 4)
}

func Test_CompatToolMigrateRemove(t *testing.T) {
	dir := setup(t, steamtest.Default())

	out := mustRun(t, "compattool", "migrate", "-y", "-r", "GE-Proton8-1", "proton_63")
	assert.Contains(t, out, "Half-Life")

	// Only fromVersion is removed, other unused tools are kept
	assert.Equal(t, []string{"GE-Proton6-1", "GE-Proton7-1", "GE-Proton9-1"}, toolDirs(t, dir))

	out = mustRun(t, "list", "-i")
	assert.Contains(t, out, "Proton 6.3-8\n\tHalf-Life (70)")

	_, err := run(t, "compattool", "migrate", "-y", "-r", "proton_63", "GE-Proton9-1")
	assert.ErrorContains(t, err, "not a custom compatibility tool")
}

func Test_CompatToolMigrateToDefault(t *testing.T) {
	setup(t, steamtest.Default())

	mustRun(t, "compattool", "migrate", "-y", "GE-Proton8-1", "default")
	assert.Contains(t, mustRun(t, "list"), "GE-Proton9-1 (Default)\n\tCyberpunk 2077\n\tHalf-Life\n")
}

func Test_CompatToolMigrateRemoveDefault(t *testing.T) {
	r := steamtest.Default()
	r.Mapping["620"] = "GE-Proton9-1"
	dir := setup(t, r)

	_, err := run(t, "compattool", "migrate", "-y", "-r", "GE-Proton9-1", "proton_63")
	assert.ErrorContains(t, err, "configured as default")
	assert.Len(t, toolDirs(t, dir), 4)
}

func Test_CompatToolSet(t *testing.T) {
	dir := setup(t, steamtest.Default())

	mustRun(t, "compattool", "set", "-y", "cyberpunk", "GE-Proton6-1")
	assert.Contains(t, mustRun(t, "list"), "GE-Proton6-1\n\tCyberpunk 2077\n")
	assert.FileExists(t, filepath.Join(dir, "config", "config.vdf"+steam.BackupSuffix))

	assert.Contains(t, mustRun(t, "compattool", "set", "cyberpunk", "GE-Proton6-1"), "already using")

	mustRun(t, "compattool", "set", "-y", "cyberpunk", "default")
	assert.Contains(t, mustRun(t, "list"), "GE-Proton9-1 (Default)\n\tCyberpunk 2077\n")

	_, err := run(t, "compattool", "set", "-y", "cyberpunk", "invalid")
	assert.ErrorContains(t, err, "Invalid version")

	out := mustRun(t, "compattool", "list")
	assert.Contains(t, out, "Proton 6.3-8 [proton_63]\n")
	assert.Contains(t, out, "GE-Proton6-1\n")
}

type archiveFile struct {
	name    string
	content string
}

func tarGz(t *testing.T, files ...archiveFile) []byte {
	t.Helper()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	for _, f := range files {
		require.NoError(t, tw.WriteHeader(&tar.Header{Name: f.name, Mode: 0644, Size: int64(len(f.content)), Typeflag: tar.TypeReg}))
		_, err := tw.Write([]byte(f.content))
		require.NoError(t, err)
	}

	require.NoError(t, tw.Close())
	require.NoError(t, gz.Close())

	return buf.Bytes()
}

func serve(t *testing.T, files map[string][]byte) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, ok := files[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Write(content)
	}))
	t.Cleanup(srv.Close)

	old := githubURL
	githubURL = srv.URL
	t.Cleanup(func() { githubURL = old })
}

func Test_GeDownload(t *testing.T) {
	dir := setup(t, steamtest.Default())

	archive := tarGz(t, archiveFile{"GE-Proton9-20/proton", "#!/bin/sh\n"})
	sum := sha512.Sum512(archive)
	base := "/" + geRepo + "/releases/download/GE-Proton9-20/"

	serve(t, map[string][]byte{
		base + "GE-Proton9-20.tar.gz":    archive,
		base + "GE-Proton9-20.sha512sum": []byte(hex.EncodeToString(sum[:]) + "  GE-Proton9-20.tar.gz\n"),
		"/" + geRepo + "/releases.atom":  []byte(`<feed><entry><link href="https://github.com/x/y/releases/tag/GE-Proton9-20"/></entry></feed>`),
	})

	out := mustRun(t, "ge", "update")
	assert.Contains(t, out, "Checksum OK")
	assert.FileExists(t, filepath.Join(dir, "compatibilitytools.d", "GE-Proton9-20", "proton"))
	assert.NotContains(t, strings.Join(toolDirs(t, dir), ","), tmpDirPrefix)

	assert.Contains(t, mustRun(t, "ge", "download", "GE-Proton9-20"), "already available")

	mustRun(t, "ge", "download", "-f", "-k", "GE-Proton9-20")
	assert.FileExists(t, filepath.Join(dir, "compatibilitytools.d", "GE-Proton9-20.tar.gz"))
	assert.FileExists(t, filepath.Join(dir, "compatibilitytools.d", "GE-Proton9-20", "proton"))
}

func Test_GeDownloadChecksumMismatch(t *testing.T) {
	dir := setup(t, steamtest.Default())

	base := "/" + geRepo + "/releases/download/GE-Proton9-20/"
	serve(t, map[string][]byte{
		base + "GE-Proton9-20.tar.gz":    tarGz(t, archiveFile{"GE-Proton9-20/proton", "x"}),
		base + "GE-Proton9-20.sha512sum": []byte(strings.Repeat("00", sha512.Size) + "  GE-Proton9-20.tar.gz\n"),
	})

	_, err := run(t, "ge", "download", "GE-Proton9-20")
	assert.ErrorContains(t, err, "checksum mismatch")
	assert.Len(t, toolDirs(t, dir), 4)
}

func Test_GeDownloadInvalidTag(t *testing.T) {
	setup(t, steamtest.Default())

	for _, tag := range []string{"GE-Proton9-1; touch /tmp/x", "../GE-Proton9-1", "GE-Proton9-1/x", "foo"} {
		_, err := run(t, "ge", "download", tag)
		assert.ErrorContains(t, err, "No valid GE version tag", tag)
	}

	_, err := run(t, "luxtorpeda", "download", "v1;id")
	assert.ErrorContains(t, err, "No valid Luxtorpeda version tag")
}

func Test_LuxtorpedaDownload(t *testing.T) {
	dir := setup(t, steamtest.Default())

	vdf := "\"compatibilitytools\"\n{\n\t\"compat_tools\"\n\t{\n\t\t\"luxtorpeda\"\n\t\t{\n\t\t\t\"install_path\"\t\".\"\n\t\t\t\"display_name\"\t\"Luxtorpeda\"\n\t\t}\n\t}\n}\n"
	archive := tarGz(t,
		archiveFile{"luxtorpeda/compatibilitytool.vdf", vdf},
		archiveFile{"luxtorpeda/luxtorpeda", "bin"},
	)

	// Only the old archive name (without "v") exists
	serve(t, map[string][]byte{
		"/" + luxtorpedaRepo + "/releases/download/v68.0/luxtorpeda-68.0.tar.xz": archive,
	})

	mustRun(t, "luxtorpeda", "download", "v68.0")

	toolDir := filepath.Join(dir, "compatibilitytools.d", "luxtorpeda-68.0")
	assert.FileExists(t, filepath.Join(toolDir, "luxtorpeda"))

	tools, err := steam.ReadCustomTools(filepath.Join(dir, "compatibilitytools.d"))
	require.NoError(t, err)
	require.NotNil(t, tools["luxtorpeda-68.0"])
	assert.Equal(t, "Luxtorpeda 68.0", tools["luxtorpeda-68.0"].DisplayName)
}

func Test_CompatToolMigrateRemoveStillUsed(t *testing.T) {
	r := steamtest.Default()
	// Shortcut of another user also uses GE-Proton8-1
	r.Mapping["3000000000"] = "GE-Proton8-1"
	dir := setup(t, r)

	out := mustRun(t, "compattool", "migrate", "-y", "-r", "GE-Proton8-1", "proton_63")
	assert.Contains(t, out, "is not deleted")
	assert.Contains(t, toolDirs(t, dir), "GE-Proton8-1")
	assert.Contains(t, mustRun(t, "list"), "Proton 6.3-8\n\tHalf-Life\n")
}

func Test_CompatToolCleanOffline(t *testing.T) {
	r := steamtest.Default()
	// Installed game without appinfo.vdf entry, store API is unreachable
	r.Mapping["999999"] = "GE-Proton7-1"
	r.Installed = append(r.Installed, "999999")
	dir := setup(t, r)

	mustRun(t, "compattool", "clean", "-y")
	assert.Equal(t, []string{"GE-Proton7-1", "GE-Proton8-1", "GE-Proton9-1"}, toolDirs(t, dir))
}
