package config

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// @see https://en.wikipedia.org/wiki/INI_file
type IniDecoder struct{}

func unescape(str string) string {
	var buf bytes.Buffer
	for i := 0; i < len(str); {
		if str[i] == '\\' && i+1 < len(str) {
			switch str[i+1] {
			case ';':
				buf.WriteByte(';')
			case '#':
				buf.WriteByte('#')
			case '=':
				buf.WriteByte('=')
			case ':':
				buf.WriteByte(':')
			case '\'':
				buf.WriteByte('\'')
			case '"':
				buf.WriteByte('"')
			case 'n':
				buf.WriteByte('\n')
			case 'r':
				buf.WriteByte('\r')
			case '\\':
				buf.WriteByte('\\')
			case 't':
				buf.WriteByte('\t')
			case 'b':
				buf.WriteByte('\b')
			case 'a':
				buf.WriteByte('\a')
			case '0':
				buf.WriteByte('\000')
			default:
				buf.WriteByte('\\')
				buf.WriteByte(str[i+1])
			}
			i += 2
		} else {
			buf.WriteByte(str[i])
			i++
		}
	}
	return buf.String()
}

func escape(str string) string {
	var buf bytes.Buffer
	for i := 0; i < len(str); i++ {
		switch str[i] {
		case ';':
			buf.WriteString(`\;`)
		case '#':
			buf.WriteString(`\#`)
		case '=':
			buf.WriteString(`\=`)
		case ':':
			buf.WriteString(`\:`)
		case '\'':
			buf.WriteString(`\'`)
		case '"':
			buf.WriteString(`\"`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\\':
			buf.WriteString(`\\`)
		case '\t':
			buf.WriteString(`\t`)
		case '\b':
			buf.WriteString(`\b`)
		case '\a':
			buf.WriteString(`\a`)
		case '\000':
			buf.WriteString(`\0`)
		default:
			buf.WriteByte(str[i])
		}
	}
	return buf.String()
}

func (d *IniDecoder) Decode(buf []byte) (*Storage, error) {
	storage, _ := NewStorage(nil)

	lines := strings.FieldsFunc(string(buf), func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	prefix := ""
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if line[0] == ';' {
			continue
		}
		if line[0] == '[' && line[len(line)-1] == ']' {
			prefix = unescape(line[1 : len(line)-1])
			continue
		}
		idx := strings.IndexAny(line, "=")
		key := strings.Trim(line[:idx], " ")
		val := strings.Trim(line[idx+1:], " ")
		if val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		key = unescape(key)
		val = unescape(val)
		if prefix != "" {
			key = prefix + "." + key
		}
		if err := storage.Set(key, val); err != nil {
			return nil, err
		}
	}

	return storage, nil
}

func (d *IniDecoder) Encode(storage *Storage) ([]byte, error) {
	var buf bytes.Buffer

	var globalKeys []string
	globalKeySet := map[string]bool{}
	if err := storage.Travel(func(key string, val interface{}) error {
		info, next, err := getToken(key)
		if err != nil {
			return err
		}
		if next != "" {
			return nil
		}
		if info.mod != MapMod {
			return fmt.Errorf("info is not a map")
		}
		globalKeys = append(globalKeys, info.key)
		globalKeySet[info.key] = true
		return nil
	}); err != nil {
		_ = storage.Travel(func(key string, val interface{}) error {
			str, _ := ToStringE(val)
			buf.WriteString(fmt.Sprintf("%v = %v\n", escape(key), escape(str)))
			return nil
		})
		return nil, err
	}

	sort.Strings(globalKeys)
	for _, key := range globalKeys {
		val, _ := storage.Get(key)
		str, _ := ToStringE(val)
		buf.WriteString(fmt.Sprintf("%v = %v\n", escape(key), escape(str)))
	}

	subs, _ := storage.SubMap("")
	var sortedKeys []string
	for key := range subs {
		if globalKeySet[key] {
			continue
		}
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	for _, key := range sortedKeys {
		sub := subs[key]
		buf.WriteString(fmt.Sprintf("\n[%v]\n", escape(key)))

		var lines []string
		_ = sub.Travel(func(key string, val interface{}) error {
			str, _ := ToStringE(val)
			lines = append(lines, fmt.Sprintf("%v = %v\n", escape(key), escape(str)))
			return nil
		})
		sort.Strings(lines)
		for _, line := range lines {
			buf.WriteString(line)
		}
	}

	return buf.Bytes(), nil
}
