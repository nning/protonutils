package steam

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
)

// InnerOffsetShortcuts sets byte count before appid match in shortcuts
const InnerOffsetShortcuts = 1

func (s *Steam) initShortcuts() error {
	p := path.Join(s.Root, "userdata", s.UID, "config", "shortcuts.vdf")
	in, err := ioutil.ReadFile(p)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if log.GetLevel() == log.DebugLevel {
		log.Debug("Steam.initShortcuts(", p, "): ", humanize.Bytes(uint64(len(in))))
	}

	s.Shortcuts = &BinaryVdf{
		Bytes: in,
		Path:  p,
	}

	return nil
}
