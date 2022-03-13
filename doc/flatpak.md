# protonutils Flatpak

## Add repo

    curl https://nning.io/0xae5fc712.asc > ae5fc712.asc
    flatpak remote-add protonutils https://nning.io/protonutils --gpg-import=ae5fc712.asc

# Install

    flatpak install protonutils io.nning.protonutils

# Run

    flatpak run io.nning.protonutils list
