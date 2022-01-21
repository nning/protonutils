package vdf2

type DeckCompatibilityCategory int

const (
	DeckCompatibilityUnknown DeckCompatibilityCategory = iota
	DeckCompatibilityUnsupported
	DeckCompatibilityPlayable
	DeckCompatibilityVerified
)

var DeckCompatibilityRegistry map[DeckCompatibilityCategory]string = map[DeckCompatibilityCategory]string{
	DeckCompatibilityUnknown:     "Unknown",
	DeckCompatibilityUnsupported: "Unsupported",
	DeckCompatibilityPlayable:    "Playable",
	DeckCompatibilityVerified:    "Verified",
}

func (c DeckCompatibilityCategory) String() string {
	return DeckCompatibilityRegistry[c]
}
