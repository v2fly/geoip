package maxmind

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/v2fly/geoip/lib"
)

const (
	typeCountryCSV = "maxmindGeoLite2CountryCSV"
	descCountryCSV = "Convert MaxMind GeoLite2 country CSV data to other formats"
)

var (
	defaultCCFile   = filepath.Join("./", "geolite2", "GeoLite2-Country-Locations-en.csv")
	defaultIPv4File = filepath.Join("./", "geolite2", "GeoLite2-Country-Blocks-IPv4.csv")
	defaultIPv6File = filepath.Join("./", "geolite2", "GeoLite2-Country-Blocks-IPv6.csv")
)

func init() {
	lib.RegisterInputConfigCreator(typeCountryCSV, func(action lib.Action, data json.RawMessage) (lib.InputConverter, error) {
		return newGeoLite2CountryCSV(action, data)
	})
	lib.RegisterInputConverter(typeCountryCSV, &geoLite2CountryCSV{})
}

func newGeoLite2CountryCSV(action lib.Action, data json.RawMessage) (lib.InputConverter, error) {
	var tmp struct {
		CountryCodeFile string   `json:"country"`
		IPv4File        string   `json:"ipv4"`
		IPv6File        string   `json:"ipv6"`
		Want            []string `json:"wantedList"`
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &tmp); err != nil {
			return nil, err
		}
	}

	if tmp.CountryCodeFile == "" {
		tmp.CountryCodeFile = defaultCCFile
	}

	if tmp.IPv4File == "" && tmp.IPv6File == "" {
		tmp.IPv4File = defaultIPv4File
		tmp.IPv6File = defaultIPv6File
	}

	return &geoLite2CountryCSV{
		Type:            typeCountryCSV,
		Action:          action,
		Description:     descCountryCSV,
		CountryCodeFile: tmp.CountryCodeFile,
		IPv4File:        tmp.IPv4File,
		IPv6File:        tmp.IPv6File,
		Want:            tmp.Want,
	}, nil
}

type geoLite2CountryCSV struct {
	Type            string
	Action          lib.Action
	Description     string
	CountryCodeFile string
	IPv4File        string
	IPv6File        string
	Want            []string
}

func (g *geoLite2CountryCSV) GetType() string {
	return g.Type
}

func (g *geoLite2CountryCSV) GetAction() lib.Action {
	return g.Action
}

func (g *geoLite2CountryCSV) GetDescription() string {
	return g.Description
}

func (g *geoLite2CountryCSV) Input(container lib.Container) (lib.Container, error) {
	ccMap, err := g.getCountryCode()
	if err != nil {
		return nil, err
	}

	entries := make(map[string]*lib.Entry)

	if g.IPv4File != "" {
		if err := g.process(g.IPv4File, ccMap, entries); err != nil {
			return nil, err
		}
	} else {
		log.Printf("type %s | action %s skip IPv4", g.Type, g.Action)
	}

	if g.IPv6File != "" {
		if err := g.process(g.IPv6File, ccMap, entries); err != nil {
			return nil, err
		}
	} else {
		log.Printf("type %s | action %s skip IPv6", g.Type, g.Action)
	}

	for name, entry := range entries {
		switch g.Action {
		case lib.ActionAdd:
			if err := container.Add(entry); err != nil {
				return nil, err
			}
		case lib.ActionRemove:
			container.Remove(name)
		case lib.ActionReplace:
			container.Replace(entry)
		default:
			return nil, lib.ErrUnknownAction
		}
	}

	return container, nil
}

func (g *geoLite2CountryCSV) getCountryCode() (map[string]string, error) {
	ccReader, err := os.Open(g.CountryCodeFile)
	if err != nil {
		return nil, err
	}
	defer ccReader.Close()

	reader := csv.NewReader(ccReader)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	ccMap := make(map[string]string)
	for _, line := range lines[1:] {
		id := strings.TrimSpace(line[0])
		countryCode := strings.TrimSpace(line[4])
		if id == "" || countryCode == "" {
			continue
		}
		ccMap[id] = strings.ToUpper(countryCode)
	}
	return ccMap, nil
}

func (g *geoLite2CountryCSV) process(file string, ccMap map[string]string, entries map[string]*lib.Entry) error {
	if len(ccMap) == 0 {
		return errors.New("country code list must be specified")
	}
	if entries == nil {
		entries = make(map[string]*lib.Entry)
	}

	fReader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fReader.Close()

	reader := csv.NewReader(fReader)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Filter want list
	wantList := make(map[string]bool)
	for _, want := range g.Want {
		if want = strings.ToUpper(strings.TrimSpace(want)); want != "" {
			wantList[want] = true
		}
	}

	for _, line := range lines[1:] {
		ccID := strings.TrimSpace(line[1])
		if countryCode, found := ccMap[ccID]; found {
			if len(wantList) > 0 {
				if _, found := wantList[countryCode]; !found {
					continue
				}
			}
			cidrStr := strings.ToLower(strings.TrimSpace(line[0]))
			entry, found := entries[countryCode]
			if !found {
				entry = lib.NewEntry(countryCode)
			}
			if err := entry.AddPrefix(cidrStr); err != nil {
				return err
			}
			entries[countryCode] = entry
		}
	}

	return nil
}
