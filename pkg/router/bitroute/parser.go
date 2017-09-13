// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bitroute

import (
	"sort"

	"github.com/takama/k8sapp/pkg/router"
)

const (
	maxLevel = 255
	asterisk = "*"
)

type handle func(router.Control)

type parser struct {
	fields   map[uint8]records
	static   map[string]handle
	wildcard records
}

type record struct {
	key    uint16
	handle handle
	parts  []string
}

type param struct {
	key   string
	value string
}
type records []*record

func (n records) Len() int           { return len(n) }
func (n records) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n records) Less(i, j int) bool { return n[i].key < n[j].key }

func newParser() *parser {
	return &parser{
		fields:   make(map[uint8]records),
		static:   make(map[string]handle),
		wildcard: records{},
	}
}

func (p *parser) register(path string, h handle) bool {
	if trim(path, " ") == asterisk {
		p.static[asterisk] = h

		return true
	}
	if parts, ok := split(path); ok {
		var static, dynamic, wildcard uint16
		for _, value := range parts {
			if len(value) >= 1 && value[0:1] == ":" {
				dynamic++
			} else if len(value) == 1 && value == "*" {
				wildcard++
			} else {
				static++
			}
		}
		if wildcard > 0 {
			p.wildcard = append(p.wildcard, &record{key: dynamic<<8 + static, handle: h, parts: parts})
		} else if dynamic == 0 {
			p.static["/"+join(parts)] = h
		} else {
			level := uint8(len(parts))
			p.fields[level] = append(p.fields[level], &record{key: dynamic<<8 + static, handle: h, parts: parts})
			sort.Sort(records(p.fields[level]))
		}
		return true
	}

	return false
}

func (p *parser) get(path string) (h handle, result []param, ok bool) {
	if h, ok := p.static[asterisk]; ok {
		return h, nil, true
	}
	if h, ok := p.static[path]; ok {
		return h, nil, true
	}
	if parts, ok := split(path); ok {
		if h, ok := p.static["/"+join(parts)]; ok {
			return h, nil, true
		}
		if data := p.fields[uint8(len(parts))]; data != nil {
			if h, result, ok := parseParams(data, parts); ok {
				return h, result, ok
			}
		}
		// try to match wildcard route
		if h, result, ok := parseParams(p.wildcard, parts); ok {
			return h, result, ok
		}
	}

	return nil, nil, false
}

func split(path string) ([]string, bool) {
	sdata := explode(trim(path, "/"))
	if len(sdata) == 0 {
		return sdata, true
	}
	var result []string
	ind := 0
	if len(sdata) < maxLevel {
		result = make([]string, len(sdata))
		for _, value := range sdata {
			if v := trim(value, " "); v == "" {
				continue
			} else {
				result[ind] = v
				ind++
			}
		}
		return result[0:ind], true
	}

	return nil, false
}

func trim(str, sep string) string {
	result := str
	for {
		if len(result) >= 1 && result[0:1] == sep {
			result = result[1:]
		} else {
			break
		}
	}
	for {
		if len(result) >= 1 && result[len(result)-1:] == sep {
			result = result[:len(result)-1]
		} else {
			break
		}
	}
	return result
}

func join(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	n := len(parts) - 1
	for i := 0; i < len(parts); i++ {
		n += len(parts[i])
	}

	b := make([]byte, n)
	bp := copy(b, parts[0])
	for _, s := range parts[1:] {
		bp += copy(b[bp:], "/")
		bp += copy(b[bp:], s)
	}
	return string(b)
}

func explode(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	n := 1
	sep := "/"
	c := sep[0]
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			n++
		}
	}
	start := 0
	a := make([]string, n)
	na := 0
	for i := 0; i+1 <= len(s) && na+1 < n; i++ {
		if s[i] == c {
			a[na] = s[start:i]
			na++
			start = i + 1
		}
	}
	a[na] = s[start:]
	return a[0 : na+1]
}

func parseParams(data records, parts []string) (h handle, result []param, ok bool) {
	for _, nds := range data {
		values := nds.parts
		result = nil
		found := true
		for idx, value := range values {
			if len(value) == 1 && value == "*" {
				break
			} else if value != parts[idx] && !(len(value) >= 1 && value[0:1] == ":") {
				found = false
				break
			} else {
				if len(value) >= 1 && value[0:1] == ":" {
					result = append(result, param{key: value, value: parts[idx]})
				}
			}
		}
		if found {
			return nds.handle, result, true
		}
	}

	return nil, nil, false
}

func (p *parser) routes() []string {
	var rs []string
	for path := range p.static {
		rs = append(rs, path)
	}
	for _, records := range p.fields {
		for _, record := range records {
			rs = append(rs, "/"+join(record.parts))
		}
	}
	for _, record := range p.wildcard {
		rs = append(rs, "/"+join(record.parts))
	}

	return rs
}
