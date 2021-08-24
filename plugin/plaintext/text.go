package plaintext

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/v2fly/geoip/lib"
)

const (
	typeText = "text"
	descText = "Convert data to text format"
)

var (
	defaultOutputDir = filepath.Join("./", "output", "text")
)

func init() {
	lib.RegisterOutputConfigCreator(typeText, func(action lib.Action, data json.RawMessage) (lib.OutputConverter, error) {
		return newText(action, data)
	})
	lib.RegisterOutputConverter(typeText, &text{})
}

func newText(action lib.Action, data json.RawMessage) (lib.OutputConverter, error) {
	var tmp struct {
		OutputDir  string     `json:"outputDir"`
		Want       []string   `json:"wantedList"`
		OnlyIPType lib.IPType `json:"onlyIPType"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &tmp); err != nil {
			return nil, err
		}
	}

	if tmp.OutputDir == "" {
		tmp.OutputDir = defaultOutputDir
	}

	return &text{
		Type:        typeText,
		Action:      action,
		Description: descText,
		OutputDir:   tmp.OutputDir,
		Want:        tmp.Want,
		OnlyIPType:  tmp.OnlyIPType,
	}, nil
}

type text struct {
	Type        string
	Action      lib.Action
	Description string
	OutputDir   string
	Want        []string
	OnlyIPType  lib.IPType
}

func (t *text) GetType() string {
	return t.Type
}

func (t *text) GetAction() lib.Action {
	return t.Action
}

func (t *text) GetDescription() string {
	return t.Description
}

func (t *text) Output(container lib.Container) error {
	// Filter want list
	wantList := make(map[string]bool)
	for _, want := range t.Want {
		if want = strings.ToUpper(strings.TrimSpace(want)); want != "" {
			wantList[want] = true
		}
	}

	switch len(wantList) {
	case 0:
		for entry := range container.Loop() {
			cidrList, err := t.marshalText(entry)
			if err != nil {
				return err
			}
			filename := strings.ToLower(entry.GetName()) + ".txt"
			if err := t.writeFile(filename, cidrList); err != nil {
				return err
			}
		}

	default:
		for name := range wantList {
			entry, found := container.GetEntry(name)
			if !found {
				log.Printf("entry %s not found", name)
				continue
			}
			cidrList, err := t.marshalText(entry)
			if err != nil {
				return err
			}
			filename := strings.ToLower(entry.GetName()) + ".txt"
			if err := t.writeFile(filename, cidrList); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *text) marshalText(entry *lib.Entry) ([]string, error) {
	var entryCidr []string
	var err error
	switch t.OnlyIPType {
	case lib.IPv4:
		entryCidr, err = entry.MarshalText(lib.IgnoreIPv6)
		if err != nil {
			return nil, err
		}
	case lib.IPv6:
		entryCidr, err = entry.MarshalText(lib.IgnoreIPv4)
		if err != nil {
			return nil, err
		}
	default:
		entryCidr, err = entry.MarshalText()
		if err != nil {
			return nil, err
		}
	}

	return entryCidr, nil
}

func (t *text) writeFile(filename string, cidrList []string) error {
	var buf bytes.Buffer
	for _, cidr := range cidrList {
		buf.WriteString(cidr)
		buf.WriteString("\n")
	}
	cidrBytes := buf.Bytes()

	if err := os.MkdirAll(t.OutputDir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(t.OutputDir, filename), cidrBytes, 0644); err != nil {
		return err
	}

	log.Printf("✅ [%s] %s --> %s", t.Type, filename, t.OutputDir)

	return nil
}
