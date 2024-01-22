package users

import (
	"github.com/mileusna/useragent"
	"manga-explorer/internal/util"
	"math"
	"strings"
)

func NewRole(val string) (Role, error) {
	switch val {
	case "admin":
		return RoleAdmin, nil
	case "user":
		return RoleUser, nil
	default:
		return Role(math.MaxUint8), ErrUnknownRole
	}
}

type Role uint8

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleUser:
		return "user"
	default:
		return "unknown"
	}
}

func (r Role) Underlying() uint8 {
	return (uint8)(r)
}

func (r Role) Validate() error {
	if r.Underlying() > 1 {
		return ErrUnknownRole
	}
	return nil
}

var RoleUser = Role(0)
var RoleAdmin = Role(1)

type Device struct {
	Name string `json:"name"`
}

func ParseDeviceName(agent *useragent.UserAgent) string {
	deviceName := agent.Device + " " + agent.OS + " " + agent.OSVersion
	if agent.Mobile || agent.Tablet {
		deviceName += " (Mobile)"
	} else if agent.Desktop {
		deviceName += " (Desktop)"
	} else if agent.Bot {
		deviceName += " (Bot)"
	}

	deviceName = strings.TrimSpace(deviceName)
	util.SetDefaultString(&deviceName, agent.Name)
	return deviceName
}
