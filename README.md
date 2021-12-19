# protonutils

protonutils is a CLI tool that provides different utilities to make using the [Proton][0] compatibility tool more easily. It has the following features:

* List games by configured Proton version
* Download latest [Proton-GE][1] release
* Download older Proton-GE releases
* Clean-up unused runtimes automatically
* Search for app ID by name
* Print or open compatdata directory by game name (handy of you want to mess with savegames or mods, for example)

The commands `list`, `appid`, `compatdata`, and `installdir` do only work with (non-native) games that either have an explicit Proton/CompatTool mapping or have been started at least once with Proton.

## Download Binary

* [Latest version](https://github.com/nning/protonutils/releases/latest/download/protonutils) (Linux x86_64)  

Make sure, the downloaded binary is executable (e.g. by running `chmod +x protonutils`).

## Build on Arch Linux

There is a [package in the AUR][2].

    yay -S protonutils

## Manual Build

Minimal Go version is 1.17.

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
      -u, --user string    Steam user name (or SteamID3)

### Update Proton-GE

    $ protonutils ge update -h
    Download and extract the latest Proton-GE release

    Usage:
      protonutils ge update [flags]

    Flags:
      -f, --force   Force last version update
      -h, --help    help for update
      -k, --keep    Keep downloaded archive of last version

### Configure Default User

`uid` can be a Steam user name or an SteamID3.

    $ protonutils config user <uid>


[0]: https://github.com/ValveSoftware/Proton
[1]: https://github.com/GloriousEggroll/proton-ge-custom
[2]: https://aur.archlinux.org/packages/protonutils/
