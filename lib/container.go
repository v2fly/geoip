package lib

import (
	"fmt"
	"strings"

	"go4.org/netipx"
)

type Container interface {
	GetEntry(name string) (*Entry, bool)
	Add(entry *Entry, opts ...IgnoreIPOption) error
	Remove(entry *Entry, rCase CaseRemove, opts ...IgnoreIPOption) error
	Len() int
	Loop() <-chan *Entry
}

type container struct {
	entries map[string]*Entry
}

func NewContainer() Container {
	return &container{
		entries: make(map[string]*Entry),
	}
}

func (c *container) isValid() bool {
	return c.entries != nil
}

func (c *container) GetEntry(name string) (*Entry, bool) {
	if !c.isValid() {
		return nil, false
	}
	val, ok := c.entries[strings.ToUpper(strings.TrimSpace(name))]
	if !ok {
		return nil, false
	}
	return val, true
}

func (c *container) Len() int {
	if !c.isValid() {
		return 0
	}
	return len(c.entries)
}

func (c *container) Loop() <-chan *Entry {
	ch := make(chan *Entry, c.Len())
	go func() {
		for _, val := range c.entries {
			ch <- val
		}
		close(ch)
	}()
	return ch
}

func (c *container) Add(entry *Entry, opts ...IgnoreIPOption) error {
	var ignoreIPType IPType
	for _, opt := range opts {
		if opt != nil {
			ignoreIPType = opt()
		}
	}

	name := entry.GetName()
	val, found := c.GetEntry(name)

	switch found {
	case true:
		var ipv4set, ipv6set *netipx.IPSet
		var err4, err6 error
		if entry.hasIPv4Builder() {
			ipv4set, err4 = entry.ipv4Builder.IPSet()
			if err4 != nil {
				return err4
			}
		}
		if entry.hasIPv6Builder() {
			ipv6set, err6 = entry.ipv6Builder.IPSet()
			if err6 != nil {
				return err6
			}
		}
		switch ignoreIPType {
		case IPv4:
			if !val.hasIPv6Builder() {
				val.ipv6Builder = new(netipx.IPSetBuilder)
			}
			val.ipv6Builder.AddSet(ipv6set)
		case IPv6:
			if !val.hasIPv4Builder() {
				val.ipv4Builder = new(netipx.IPSetBuilder)
			}
			val.ipv4Builder.AddSet(ipv4set)
		default:
			if !val.hasIPv4Builder() {
				val.ipv4Builder = new(netipx.IPSetBuilder)
			}
			if !val.hasIPv6Builder() {
				val.ipv6Builder = new(netipx.IPSetBuilder)
			}
			val.ipv4Builder.AddSet(ipv4set)
			val.ipv6Builder.AddSet(ipv6set)
		}

	case false:
		switch ignoreIPType {
		case IPv4:
			entry.ipv4Builder = nil
		case IPv6:
			entry.ipv6Builder = nil
		}
		c.entries[name] = entry
	}

	return nil
}

func (c *container) Remove(entry *Entry, rCase CaseRemove, opts ...IgnoreIPOption) error {
	name := entry.GetName()
	val, found := c.GetEntry(name)
	if !found {
		return fmt.Errorf("entry %s not found", name)
	}

	var ignoreIPType IPType
	for _, opt := range opts {
		if opt != nil {
			ignoreIPType = opt()
		}
	}

	switch rCase {
	case CaseRemovePrefix:
		var ipv4set, ipv6set *netipx.IPSet
		var err4, err6 error
		if entry.hasIPv4Builder() {
			ipv4set, err4 = entry.ipv4Builder.IPSet()
			if err4 != nil {
				return err4
			}
		}
		if entry.hasIPv6Builder() {
			ipv6set, err6 = entry.ipv6Builder.IPSet()
			if err6 != nil {
				return err6
			}
		}

		switch ignoreIPType {
		case IPv4:
			if !val.hasIPv6Builder() {
				val.ipv6Builder = new(netipx.IPSetBuilder)
			}
			val.ipv6Builder.RemoveSet(ipv6set)
		case IPv6:
			if !val.hasIPv4Builder() {
				val.ipv4Builder = new(netipx.IPSetBuilder)
			}
			val.ipv4Builder.RemoveSet(ipv4set)
		default:
			if !val.hasIPv4Builder() {
				val.ipv4Builder = new(netipx.IPSetBuilder)
			}
			if !val.hasIPv6Builder() {
				val.ipv6Builder = new(netipx.IPSetBuilder)
			}
			val.ipv4Builder.RemoveSet(ipv4set)
			val.ipv6Builder.RemoveSet(ipv6set)
		}

	case CaseRemoveEntry:
		switch ignoreIPType {
		case IPv4:
			val.ipv6Builder = nil
		case IPv6:
			val.ipv4Builder = nil
		default:
			delete(c.entries, name)
		}

	default:
		return fmt.Errorf("unknown remove case %d", rCase)
	}

	return nil
}
