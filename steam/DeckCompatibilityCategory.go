package steam

// DeckCompatibilityCategory represents deck compatibility category
type DeckCompatibilityCategory int

const (
	// DeckCompatibilityUnknown represents unknown deck compatibility
	DeckCompatibilityUnknown DeckCompatibilityCategory = iota
	// DeckCompatibilityUnsupported represents unsupported deck compatibility
	DeckCompatibilityUnsupported
	// DeckCompatibilityPlayable represents playable deck compatibility
	DeckCompatibilityPlayable
	// DeckCompatibilityVerified represents verified deck compatibility
	DeckCompatibilityVerified
)

// DeckCompatibilityRegistry translates DeckCompatibilityCategory values to
// human-readable strings
var DeckCompatibilityRegistry map[DeckCompatibilityCategory]string = map[DeckCompatibilityCategory]string{
	DeckCompatibilityUnknown:     "Unknown",
	DeckCompatibilityUnsupported: "Unsupported",
	DeckCompatibilityPlayable:    "Playable",
	DeckCompatibilityVerified:    "Verified",
}

func (c DeckCompatibilityCategory) String() string {
	return DeckCompatibilityRegistry[c]
}
