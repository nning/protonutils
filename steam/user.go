package steam

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// steamID64Base is the steamID64 of the individual account with account ID 0
// in the public universe
const steamID64Base = 76561197960265728

// UserNotFoundError is returned if a Steam user cannot be found
type UserNotFoundError struct {
	User string
}

func (e *UserNotFoundError) Error() string {
	return "Steam user not found: " + e.User
}

// id64ToID32 converts a steamID64 to the account ID (steamID3 without
// prefixes) as used for userdata directories
func id64ToID32(id64 string) (string, error) {
	n, err := strconv.ParseUint(id64, 10, 64)
	if err != nil {
		return "", err
	}

	if n <= steamID64Base {
		return "", fmt.Errorf("invalid steamID64: %v", id64)
	}

	return strconv.FormatUint(n&0xffffffff, 10), nil
}

// userToID32 converts a user name, steamID64 or account ID into the account ID
func (s *Steam) userToID32(u string) (string, error) {
	if id64 := s.LoginUsers.GetID64(u); id64 != "" {
		return id64ToID32(id64)
	}

	n, err := strconv.ParseUint(u, 10, 64)
	if err != nil {
		return "", &UserNotFoundError{u}
	}

	if n > steamID64Base {
		return id64ToID32(u)
	}

	if n == 0 || n > 0xffffffff {
		return "", &UserNotFoundError{u}
	}

	return u, nil
}

// getUserIDs returns sorted list of account IDs that have a userdata directory
func (s *Steam) getUserIDs() ([]string, error) {
	dir := filepath.Join(s.Root, "userdata")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	users := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() || name == "0" {
			continue
		}

		if _, err := strconv.ParseUint(name, 10, 32); err != nil {
			continue
		}

		users = append(users, name)
	}

	sort.Strings(users)
	return users, nil
}

// resolveUID returns the account ID of the user whose data is used. If user is
// empty, the only available user, the most recently logged in user, or the
// first user (in that order) is chosen. Warnings are written to w.
func (s *Steam) resolveUID(user string, w io.Writer) (string, error) {
	users, err := s.getUserIDs()
	if err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "", errors.New("No Steam user found in " + filepath.Join(s.Root, "userdata"))
	}

	if user != "" {
		uid, err := s.userToID32(user)
		if err != nil {
			return "", err
		}

		for _, u := range users {
			if u == uid {
				return uid, nil
			}
		}

		return "", fmt.Errorf("No userdata found for Steam user %v (%v)", user, uid)
	}

	if len(users) == 1 {
		return users[0], nil
	}

	if id64 := s.LoginUsers.GetMostRecentID64(); id64 != "" {
		uid, err := id64ToID32(id64)
		if err == nil {
			for _, u := range users {
				if u == uid {
					slog.Debug("Steam.resolveUID: using most recent user", "uid", uid)
					return uid, nil
				}
			}
		}
	}

	fmt.Fprintln(w,
		"Warning: Several Steam users available, using "+users[0]+"\n"+
			"All available users: "+strings.Join(users, ", ")+"\n"+
			"Option \"-u\" can be used to specify user")
	fmt.Fprintln(w)

	return users[0], nil
}
