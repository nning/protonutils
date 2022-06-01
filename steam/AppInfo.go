package steam

import (
	"errors"
	"io/ioutil"
	"path"

	"github.com/BenLubar/vdf"
	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
)

// InnerOffsetAppInfo sets byte count before appid match in appinfo
const InnerOffsetAppInfo = 10

// GetCompatToolName returns human-readable name of compatibility tool,
// for example: "proton_63" -> "Proton 6.3-8"
func (s *Steam) GetCompatToolName(id string) (string, error) {
	if id == "" {
		return "", nil
	}

	str, _ := s.VersionNameCache.Get(id)
	if str != "" {
		return str, nil
	}

	var name string
	var n *vdf.Node

	// TODO Extract name from tool's own compatibilitytool.vdf

	// Search for app ID 891390 ("SteamPlay 2.0 Manifests")
	i, err := s.AppInfo.GetNextEntryStartByID(0, InnerOffsetAppInfo, "891390")
	if i < 0 || err != nil {
		goto cache_and_return_error
	}

	n, err = ParseBinaryVdf(s.AppInfo.Bytes[i:])
	if err != nil {
		goto cache_and_return_error
	}

	n = n.FirstByName("extended").FirstByName("compat_tools")
	n = n.FirstByName(id).FirstByName("display_name")

	name = n.String()
	if name == "" {
		goto cache_and_return_error
	}

	s.VersionNameCache.Add(id, name, true)
	return name, nil

cache_and_return_error:
	s.VersionNameCache.Add(id, id, false)
	return id, err
}

// GetInstalldir returns installation directory of game by app ID
func (s *Steam) GetInstalldir(id string) (string, error) {
	if id == "" {
		return "", nil
	}

	var n *vdf.Node

	i, err := s.AppInfo.GetNextEntryStartByID(0, InnerOffsetAppInfo, id)
	if i < 0 || err != nil {
		return "", err
	}

	n, err = ParseBinaryVdf(s.AppInfo.Bytes[i:])
	if err != nil {
		return "", err
	}

	n = n.FirstByName("config").FirstByName("installdir")

	dir := n.String()
	if dir == "" {
		return "", errors.New("installdir not found")
	}

	return dir, nil
}

func (s *Steam) initAppInfo() error {
	p := path.Join(s.Root, "appcache", "appinfo.vdf")
	in, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	if log.GetLevel() == log.DebugLevel {
		log.Debug("Steam.initAppInfo(", p, "): ", humanize.Bytes(uint64(len(in))))
	}

	s.AppInfo = &BinaryVdf{
		Bytes: in,
		Path:  p,
	}

	return nil
}
