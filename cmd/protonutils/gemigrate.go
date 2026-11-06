package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
)

var (
	geVersionNew = regexp.MustCompile(`^GE-Proton([0-9]+)-([0-9]+)(-.*)?$`)
	geVersionOld = regexp.MustCompile(`^Proton-([0-9]+)\.([0-9]+)(rc([0-9]+))?-GE-([0-9]+)(-.*)?$`)
)

// geVersion is a parsed Proton-GE tool ID that can be compared
type geVersion struct {
	parts [4]int
	// variant is the suffix of special builds (like "-rtsp"), without
	// architecture suffix
	variant string
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// parseGEVersion parses Proton-GE tool IDs like "GE-Proton9-27",
// "GE-Proton11-12-x86_64" or "Proton-7.0rc6-GE-1"
func parseGEVersion(id string) (geVersion, bool) {
	var v geVersion
	var suffix string

	if m := geVersionNew.FindStringSubmatch(id); m != nil {
		v.parts = [4]int{1, atoi(m[1]), atoi(m[2]), 0}
		suffix = m[3]
	} else if m := geVersionOld.FindStringSubmatch(id); m != nil {
		// Release candidates are older than the final release
		rc := 1 << 20
		if m[3] != "" {
			rc = atoi(m[4])
		}

		v.parts = [4]int{0, atoi(m[1])*1000 + atoi(m[2]), rc, atoi(m[5])}
		suffix = m[6]
	} else {
		return v, false
	}

	for _, arch := range []string{"-x86_64", "-aarch64"} {
		suffix = strings.TrimSuffix(suffix, arch)
	}
	v.variant = suffix

	return v, true
}

// olderThan returns whether v is an older version of the same variant as o
func (v geVersion) olderThan(o geVersion) bool {
	if v.variant != o.variant {
		return false
	}

	for i := range v.parts {
		if v.parts[i] != o.parts[i] {
			return v.parts[i] < o.parts[i]
		}
	}

	return false
}

// geMigrate switches the default and all installed games (and non-Steam
// shortcuts) that use an older Proton-GE version of the same variant to id
func geMigrate(cmd *cobra.Command, id string) error {
	out := cmd.OutOrStdout()

	newVersion, ok := parseGEVersion(id)
	if !ok {
		return fmt.Errorf("cannot migrate: %v is no Proton-GE version", id)
	}

	err := s.ReloadCustomTools()
	if err != nil {
		return err
	}

	if s.CustomTools.Get(id) == nil {
		return fmt.Errorf("cannot migrate: %v is not installed in %v", id, s.GetCompatibilityToolsDir())
	}

	type change struct {
		appID string
		name  string
		from  string
	}
	var changes []change

	ctm := s.CompatToolMapping
	for x := ctm.Node.FirstSubTree(); x != nil; x = x.NextSubTree() {
		appID := x.Name()
		from := x.FirstByName("name").String()

		v, ok := parseGEVersion(from)
		if !ok || !v.olderThan(newVersion) {
			continue
		}

		if appID == "0" {
			changes = append(changes, change{appID, "Default compatibility tool", from})
			continue
		}

		if !steam.IsShortcut(appID) && !s.LibraryConfig.IsInstalled(appID) {
			continue
		}

		name := appID
		if game, valid, err := s.GetGameData(appID); err == nil && valid {
			name = game.Name
		}

		changes = append(changes, change{appID, name, from})
	}

	if len(changes) == 0 {
		fmt.Fprintln(out, "No installed games use an older Proton-GE version")
		return nil
	}

	sort.Slice(changes, func(i, j int) bool {
		if (changes[i].appID == "0") != (changes[j].appID == "0") {
			return changes[i].appID == "0"
		}

		return changes[i].name < changes[j].name
	})

	warnSteamRunning(out)
	fmt.Fprintf(out, "Switch to %v:\n\n", id)
	for _, c := range changes {
		fmt.Fprintf(out, "  * %v (%v)\n", c.name, c.from)
	}
	fmt.Fprintln(out)

	ok, err = confirm(cmd, "Really update?")
	if !ok || err != nil {
		return err
	}

	for _, c := range changes {
		err = ctm.Update(c.appID, id)
		if err != nil {
			return err
		}
	}

	err = ctm.Save()
	if err != nil {
		return err
	}

	fmt.Fprintln(out, "Done")
	return nil
}
