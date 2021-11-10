# protonutils

protonutils is a CLI tool that provides different utilities to make using the Proton compatibility tool more easy. Currently, it lists games by configured Proton version (which can be useful to clean up old runtimes) and it updates [Proton-GE][0].

## Download Binary

* [Latest version](https://github.com/nning/protonutils/releases/latest/download/protonutils) (Linux x86_64)  

Make sure, the downloaded binary is executable (e.g. by running `chmod +x protonutils`).

## Build & Run

This step is only necessary if aforementioned binary does not suit your needs for some reason. Install [Go](https://golang.org/), make sure `$GOPATH` is set correctly, then run...

    go install github.com/nning/protonutils/cmd/protonutils@latest
    protonutils

## Usage

### List

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

## Example Output

    $ protonutils list
    Proton-6.20-GE-1
            Cyberpunk 2077
            DEATHLOOP
            Dishonored 2
            Frostpunk
            Horizon Zero Dawn
            Iron Harvest
            Kena - Bridge of Spirits [SHORTCUT]

    proton_63 (Default)
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


[0]: https://github.com/GloriousEggroll/proton-ge-custom
