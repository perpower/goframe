// 时区设置，转换
package ptime

import (
	"os"
	"sync"
	"time"

	"github.com/perpower/goframe/funcs/normal"
	"github.com/perpower/goframe/utils/errors"
)

var (
	setTimeZoneMu   sync.Mutex
	setTimeZoneName string
	zoneMap         = make(map[string]*time.Location)
	zoneMu          sync.RWMutex
)

// SetTimeZone sets the time zone for current whole process.
// The parameter `zone` is an area string specifying corresponding time zone,
// eg: Asia/Shanghai.
//
// PLEASE VERY NOTE THAT:
// 1. This should be called before package "time" import.
// 2. This function should be called once.
// 3. Please refer to issue: https://github.com/golang/go/issues/34814
func SetTimeZone(zone string) (err error) {
	setTimeZoneMu.Lock()
	defer setTimeZoneMu.Unlock()
	if setTimeZoneName != "" && !normal.Equal(zone, setTimeZoneName) {
		return errors.Newf(errors.ERROR_CODE.Code, `process timezone already set using "%s"`, nil, setTimeZoneName)
	}
	defer func() {
		if err == nil {
			setTimeZoneName = zone
		}
	}()

	// It is already set to time.Local.
	if normal.Equal(zone, time.Local.String()) {
		return
	}

	// Load zone info from specified name.
	location, err := time.LoadLocation(zone)
	if err != nil {
		err = errors.Newf(errors.ERROR_CODE.Code, `time.LoadLocation failed for zone "%s"`, err, zone)
		return err
	}

	// Update the time.Local for once.
	time.Local = location

	// Update the timezone environment for *nix systems.
	var (
		envKey   = "TZ"
		envValue = location.String()
	)
	if err = os.Setenv(envKey, envValue); err != nil {
		err = errors.Newf(errors.ERROR_CODE.Code, `set environment failed with key "%s", value "%s"`, err, envKey, envValue)
	}
	return
}

// ToLocation converts current time to specified location.
func (t *Time) ToLocation(location *time.Location) *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.In(location)
	return newTime
}

// ToZone converts current time to specified zone like: Asia/Shanghai.
func (t *Time) ToZone(zone string) (*Time, error) {
	if location, err := t.getLocationByZoneName(zone); err == nil {
		return t.ToLocation(location), nil
	} else {
		return nil, err
	}
}

func (t *Time) getLocationByZoneName(name string) (location *time.Location, err error) {
	zoneMu.RLock()
	location = zoneMap[name]
	zoneMu.RUnlock()
	if location == nil {
		location, err = time.LoadLocation(name)
		if err != nil {
			err = errors.Newf(errors.ERROR_CODE.Code, `time.LoadLocation failed for name "%s"`, err, name)
		}
		if location != nil {
			zoneMu.Lock()
			zoneMap[name] = location
			zoneMu.Unlock()
		}
	}
	return
}

// Local converts the time to local timezone.
func (t *Time) Local() *Time {
	newTime := t.Clone()
	newTime.Time = newTime.Time.Local()
	return newTime
}

// Clone returns a new Time object which is a clone of current time object.
func (t *Time) Clone() *Time {
	return New(t.Time)
}
