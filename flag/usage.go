package flag

import (
	"bytes"
	"fmt"
	"path"
	"sort"
)

func (f *Flag) Usage() string {
	type info struct {
		shorthand   string
		name        string
		typeDefault string
		usage       string
	}

	var argumentInfos []*info
	var flagInfos []*info

	for _, name := range f.arguments {
		finfo := f.flagInfos[name]
		defaultValue := finfo.Type.String()
		if finfo.DefaultValue != "" {
			defaultValue = defaultValue + "=" + finfo.DefaultValue
		}
		argumentInfos = append(argumentInfos, &info{
			shorthand:   "",
			name:        finfo.Name,
			typeDefault: "[" + defaultValue + "]",
			usage:       finfo.Usage,
		})
	}

	var keys []string
	for key, info := range f.flagInfos {
		if info.IsArgument {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		finfo := f.flagInfos[key]
		defaultValue := finfo.Type.String()
		if finfo.DefaultValue != "" {
			defaultValue = defaultValue + "=" + finfo.DefaultValue
		}
		shorthand := ""
		if finfo.Shorthand != "" {
			shorthand = "-" + finfo.Shorthand
		}
		flagInfos = append(flagInfos, &info{
			shorthand:   shorthand,
			name:        "--" + finfo.Name,
			typeDefault: "[" + defaultValue + "]",
			usage:       finfo.Usage,
		})
	}

	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	var shorthandWidth, nameWidth, typeDefaultWidth int
	for _, i := range append(argumentInfos, flagInfos...) {
		shorthandWidth = max(len(i.shorthand), shorthandWidth)
		nameWidth = max(len(i.name), nameWidth)
		typeDefaultWidth = max(len(i.typeDefault), typeDefaultWidth)
	}

	var buffer bytes.Buffer

	buffer.WriteString("usage: ")
	buffer.WriteString(path.Base(f.name))
	for _, key := range f.arguments {
		buffer.WriteString(fmt.Sprintf(" [%v]", key))
	}

	for _, key := range keys {
		finfo := f.flagInfos[key]
		nameShorthand := finfo.Name
		if finfo.Shorthand != "" {
			nameShorthand = finfo.Shorthand + "," + finfo.Name
		}
		if finfo.DefaultValue != "" {
			buffer.WriteString(fmt.Sprintf(" [-%v %v=%v]", nameShorthand, finfo.Type, finfo.DefaultValue))
		} else if finfo.Required {
			buffer.WriteString(fmt.Sprintf(" <-%v %v>", nameShorthand, finfo.Type))
		} else {
			buffer.WriteString(fmt.Sprintf(" [-%v %v]", nameShorthand, finfo.Type))
		}
	}
	buffer.WriteString("\n")

	if len(argumentInfos) != 0 {
		buffer.WriteString("\narguments:\n")
		posFormat := fmt.Sprintf("  %%%dv  %%-%dv  %%-%dv  %%v\n", shorthandWidth, nameWidth, typeDefaultWidth)
		for _, i := range argumentInfos {
			buffer.WriteString(fmt.Sprintf(posFormat, i.shorthand, i.name, i.typeDefault, i.usage))
		}
	}
	buffer.WriteString("\noptions:\n")
	format := fmt.Sprintf("  %%%dv, %%-%dv  %%-%dv  %%v\n", shorthandWidth, nameWidth, typeDefaultWidth)
	for _, i := range flagInfos {
		buffer.WriteString(fmt.Sprintf(format, i.shorthand, i.name, i.typeDefault, i.usage))
	}

	return buffer.String()
}
