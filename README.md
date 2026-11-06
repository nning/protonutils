# protonutils

protonutils is a CLI tool that provides different utilities to make using the [Proton][0] compatibility tool more easily. It has the following features:

* List games by configured Proton version
* Download latest and older [Proton-GE][1] releases
* Download latest and older [Luxtorpeda][3] releases
* Clean-up unused compatibility tools automatically
* Search for app ID of game in your library by name
* Print or open compatdata and install directories by game name  
  (handy of you want to mess with savegames or mods, for example)
* Update assigned compatibility tool for one or more games
* Clean-up unused `compatdata` directories

## Download Binary

* [Latest version](https://github.com/nning/protonutils/releases/latest/download/protonutils) (Linux x86_64)  
* [Latest version](https://github.com/nning/protonutils/releases/latest/download/protonutils-arm64) (Linux arm64)  
* [SHA256SUMS](https://github.com/nning/protonutils/releases/latest/download/SHA256SUMS)

Make sure, the downloaded binary is executable (e.g. by running `chmod +x protonutils`).

## Build on Arch Linux

There is a [package in the AUR][2].

    yay -S protonutils

## Manual Build

Minimal Go version is 1.21.

This step is only necessary if aforementioned binary does not suit your needs for some reason. Install [Go](https://golang.org/), make sure `$GOPATH` is set correctly, then run...

    go install github.com/nning/protonutils/cmd/protonutils@latest
    protonutils

## Example Output

    $ protonutils list
    Proton-6.21-GE-2
            Cyberpunk 2077
            DEATHLOOP
            Dishonored 2
            Frostpunk
            Horizon Zero Dawn
            Iron Harvest
            Kena - Bridge of Spirits [SHORTCUT]

    Proton 6.3-8 (Default)
            Age of Empires II: Definitive Edition
            Fallout 4
            Grand Theft Auto V
            It Takes Two
            Metro Exodus
            Shadow Tactics: Blades of the Shogun
            The Witcher 3: Wild Hunt

    proton_experimental
            Mass Effect™ Legendary Edition
            Red Dead Redemption 2

## Usage

The two outputs are just examples for two prominent commands, see full usage
documentation by running `protonutils -h` on your machine (or `man protonutils`
on Arch Linux if you installed from the AUR).

### List Games by Version

    $ protonutils list -h
    List games by runtime
    
    Usage:
      protonutils list [flags]
    
    Flags:
      -a, --all            List both installed and non-installed games
      -h, --help           help for list
      -c, --ignore-cache   Ignore app ID/name cache
      -j, --json           Output JSON (implies -a and -i)
      -i, --show-id        Show app ID
      -s, --size           Show disk usage of compatdata and game directories (JSON: also paths)
      -u, --user string    Steam user name, steamID64 or account ID (SteamID3)

### Update Proton-GE

    $ protonutils ge update -h
    Download and extract the latest Proton-GE release

    Usage:
      protonutils ge update [flags]

    Flags:
      -f, --force     Force last version update
      -h, --help      help for update
      -k, --keep      Keep downloaded archive of last version
      -m, --migrate   Switch installed games (and the default) from older Proton-GE versions to this version
      -y, --yes       Do not ask

With `--migrate`, installed games and non-Steam shortcuts using an older
Proton-GE version (and the default compatibility tool, if it is an older
Proton-GE version) are switched to the new version after confirmation. Newer
versions and special builds (like `-rtsp`) are not changed.

Downloads are verified with the SHA-512 checksum published with the release
(if available) and only moved into `compatibilitytools.d` after successful
extraction.

### Bulk Update Compatibiliy Tool Versions

    $ protonutils compattool migrate -h
    Migrate compatibility tool version mappings of installed games from one version to another. Version parameters have to be version IDs. See `compattool list` for list of possible options. If "default" is used as toVersion, explicit mappings are removed and games use default compatibility tool.

    Usage:
      protonutils compattool migrate [flags] <fromVersion> <toVersion>

    Flags:
      -h, --help           help for migrate
      -c, --ignore-cache   Ignore app ID/name cache
      -r, --remove         Remove fromVersion after migration (only custom tools in compatibilitytools.d)
      -u, --user string    Steam user name, steamID64 or account ID (SteamID3)
      -y, --yes            Do not ask

    $ protonutils compattool migrate Proton-6.21-GE-2 Proton-7.0rc2-GE-1
    Proton-6.21-GE-2 -> Proton-7.0rc2-GE-1

      * Cyberpunk 2077
      * DEATHLOOP
      * Horizon Zero Dawn

    Really update? [y/N]

### Clean Unused Compatibility Tools

`compattool clean` deletes compatibility tools in `compatibilitytools.d` that
are neither used by an installed game nor configured as default compatibility
tool.

### Clean Unused `compatdata` Directories

    $ protonutils compatdata clean
    Calculating unused compatdata directory sizes...

          19 B  Fallout 4

    Total size: 19 B
    WARNING: Backup save game data for games without Steam Cloud support!
    Do you want to delete compatdata directories? [y/N]

## Shell Completion

Completion scripts can be generated with `protonutils completion
bash|zsh|fish`. Besides commands and flags, game names (for `compatdata`,
`installdir`, `appid`, `compattool set`) and compatibility tool versions are
completed. Game names are known after any command that reads games (like
`list`) has run once.

## Configuration

The configuration is stored in `$XDG_CONFIG_HOME/protonutils/protonutils.yml`
(default: `~/.config/protonutils/protonutils.yml`), caches in
`$XDG_CACHE_HOME/protonutils` (default: `~/.cache/protonutils`). Older versions
stored the configuration file directly in `$XDG_CONFIG_HOME` if it was set; it
is moved automatically.

### Default User

`uid` can be a Steam user name, a steamID64 or an account ID (SteamID3). If no
user is configured and several users exist, the most recently logged in user
is used.

    $ protonutils config user <uid>

### Steam Root Path

By default, the Steam root is detected automatically (`~/.local/share/Steam`,
`~/.steam/root`, `~/.steam/steam`, and Steam from Flatpak in
`~/.var/app/com.valvesoftware.Steam`). It can be set explicitly:

    $ protonutils config steam_root ~/.var/app/com.valvesoftware.Steam/data/Steam

### Reset Option

    $ protonutils config steam_root ""

### Backups

Before `compattool set` or `compattool migrate` change Steam's `config.vdf`, the
previous version is saved as `config.vdf.protonutils.bak`.

### Debug Output

    $ DEBUG=1 protonutils list


[0]: https://github.com/ValveSoftware/Proton
[1]: https://github.com/GloriousEggroll/proton-ge-custom
[2]: https://aur.archlinux.org/packages/protonutils/
[3]: https://github.com/luxtorpeda-dev/luxtorpeda
