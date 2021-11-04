# list-proton-versions

Lists configured Proton version per game. Can be useful to clean up old runtimes, for example. (For now, only works if game was launched with Proton at least once.)

## Build & Run

Install [Go](https://golang.org/), make sure `$GOPATH` is set correctly, then run...

    go install github.com/nning/list_proton_versions/cmd/list-proton-versions@latest
    list-proton-versions

## Manually Download Binary

[Download Linux x86_64 binary from CI](https://github.com/nning/list_proton_versions/suites/4259569597/artifacts/110638934)  
(Sorry, for now only works with GitHub account.)

You can download the most recent version [from continuous builds](https://github.com/nning/list_proton_versions/actions/workflows/build.yml). Choose the last build and find the download link under "Artifacts".

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
