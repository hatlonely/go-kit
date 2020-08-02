// Code generated by "enumer -type Level -trimprefix Level -text"; DO NOT EDIT.

//
package logger

import (
	"fmt"
)

const _LevelName = "DebugInfoWarnError"

var _LevelIndex = [...]uint8{0, 5, 9, 13, 18}

func (i Level) String() string {
	i -= 1
	if i < 0 || i >= Level(len(_LevelIndex)-1) {
		return fmt.Sprintf("Level(%d)", i+1)
	}
	return _LevelName[_LevelIndex[i]:_LevelIndex[i+1]]
}

var _LevelValues = []Level{1, 2, 3, 4}

var _LevelNameToValueMap = map[string]Level{
	_LevelName[0:5]:   1,
	_LevelName[5:9]:   2,
	_LevelName[9:13]:  3,
	_LevelName[13:18]: 4,
}

// LevelString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func LevelString(s string) (Level, error) {
	if val, ok := _LevelNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Level values", s)
}

// LevelValues returns all values of the enum
func LevelValues() []Level {
	return _LevelValues
}

// IsALevel returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Level) IsALevel() bool {
	for _, v := range _LevelValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalText implements the encoding.TextMarshaler interface for Level
func (i Level) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Level
func (i *Level) UnmarshalText(text []byte) error {
	var err error
	*i, err = LevelString(string(text))
	return err
}