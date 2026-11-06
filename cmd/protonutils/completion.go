package main

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/nning/protonutils/cache"
	"github.com/nning/protonutils/config"
	"github.com/nning/protonutils/steam"
	"github.com/spf13/cobra"
)

// Shell completion only reads caches and compatibilitytools.d, so that it is
// fast and works without full Steam initialization. Game names are known after
// any command that reads games (like "list") has run once.

// completeGame returns a completion function for game names. If joined is true,
// the game search string consists of all arguments joined by spaces (and the
// remaining words of a name are completed), otherwise only the first argument
// is completed.
func completeGame(joined bool) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		directive := cobra.ShellCompDirectiveNoFileComp

		if !joined && len(args) > 0 {
			return nil, directive
		}

		c, err := cache.New("appids", -1)
		if err != nil {
			return nil, directive
		}

		prefix := toComplete
		if joined && len(args) > 0 {
			prefix = strings.Join(args, " ") + " " + toComplete
		}

		// Words already typed are not part of the completed word
		skip := len(prefix) - len(toComplete)

		var completions []cobra.Completion
		seen := make(map[string]bool)

		for _, m := range steam.SearchAppIDAndNames(c, prefix) {
			if len(m.Name) < len(prefix) || !strings.EqualFold(m.Name[:len(prefix)], prefix) {
				continue
			}

			name := m.Name[skip:]
			if seen[name] {
				continue
			}
			seen[name] = true

			completions = append(completions, name)
		}

		return completions, directive
	}
}

// versionCompletions returns compatibility tool IDs (with names as
// description) known from caches and compatibilitytools.d
func versionCompletions(withDefault bool) []cobra.Completion {
	names := make(map[string]string)

	if c, err := cache.New("proton-names", -1); err == nil {
		for id, v := range c.Dump() {
			if v.Valid && strings.HasPrefix(id, "proton_") {
				names[id] = v.Name
			}
		}
	}

	if cfg, err := config.New(); err == nil {
		if root, err := steam.FindRoot(cfg.SteamRoot); err == nil {
			if tools, err := steam.ReadCustomTools(filepath.Join(root, "compatibilitytools.d")); err == nil {
				for id, tool := range tools {
					names[id] = tool.DisplayName
				}
			}
		}
	}

	ids := make([]string, 0, len(names))
	for id := range names {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var completions []cobra.Completion
	if withDefault {
		completions = append(completions, cobra.CompletionWithDesc("default", "Default compatibility tool"))
	}

	for _, id := range ids {
		if names[id] != "" && names[id] != id {
			completions = append(completions, cobra.CompletionWithDesc(id, names[id]))
		} else {
			completions = append(completions, id)
		}
	}

	return completions
}

func init() {
	for _, c := range []*cobra.Command{appidCmd, compatdataCleanCmd} {
		c.ValidArgsFunction = completeGame(true)
	}

	compatToolSetCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		switch len(args) {
		case 0:
			return completeGame(false)(cmd, args, toComplete)
		case 1:
			return versionCompletions(true), cobra.ShellCompDirectiveNoFileComp
		}

		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	compatToolMigrateCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) < 2 {
			return versionCompletions(len(args) == 1), cobra.ShellCompDirectiveNoFileComp
		}

		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}
