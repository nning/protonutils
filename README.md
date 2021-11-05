# list-proton-versions

Lists configured Proton version per game. Can be useful to clean up old runtimes, for example. (For now, only works if game was launched with Proton at least once.)

## Download Binary

* [Latest version](https://github.com/nning/list_proton_versions/releases/latest/download/list-proton-versions) (Linux x86_64)  

Make sure, the downloaded binary is executable (e.g. by running `chmod +x list-proton-versions`).

## Build & Run

This step is only necessary if aforementioned binary does not suit your needs for some reason. Install [Go](https://golang.org/), make sure `$GOPATH` is set correctly, then run...

    go install github.com/nning/list_proton_versions/cmd/list-proton-versions@latest
    list-proton-versions

## Usage

    $ ./list-proton-versions -h
    Usage of ./list-proton-versions:
      -a    List both installed and non-installed games
      -c    Ignore app ID/name cache
      -i    Show app ID
      -j    Output JSON (implies -a and -i)
      -u string
            Steam user name (or SteamID3)

## Example Output

    $ ./list-proton-versions
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
