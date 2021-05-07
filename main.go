// GeoIP generator
//
// Before running this file, the GeoIP database must be downloaded and present.
// To download GeoIP database: https://dev.maxmind.com/geoip/geoip2/geolite2/
// Inside you will find block files for IPv4 and IPv6 and country code mapping.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/v2fly/v2ray-core/v4/app/router"
	"github.com/v2fly/v2ray-core/v4/common"
	"github.com/v2fly/v2ray-core/v4/infra/conf/rule"
	"google.golang.org/protobuf/proto"
)

var (
	countryCodeFile = flag.String("country", "GeoLite2-Country-Locations-en.csv", "Path to the country code file")
	ipv4File        = flag.String("ipv4", "GeoLite2-Country-Blocks-IPv4.csv", "Path to the IPv4 block file")
	ipv6File        = flag.String("ipv6", "GeoLite2-Country-Blocks-IPv6.csv", "Path to the IPv6 block file")
	outputName      = flag.String("outputname", "geoip.dat", "Name of the generated file")
	outputDir       = flag.String("outputdir", "./", "Path to the output directory")
)

var privateIPs = []string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.0.0.0/24",
	"192.0.2.0/24",
	"192.88.99.0/24",
	"192.168.0.0/16",
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"224.0.0.0/4",
	"240.0.0.0/4",
	"255.255.255.255/32",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
}

var testIPs = []string{
	"127.0.0.0/8",
}

func getCountryCodeMap() (map[string]string, error) {
	countryCodeReader, err := os.Open(*countryCodeFile)
	if err != nil {
		return nil, err
	}
	defer countryCodeReader.Close()

	m := make(map[string]string)
	reader := csv.NewReader(countryCodeReader)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, line := range lines[1:] {
		id := line[0]
		countryCode := line[4]
		if len(countryCode) == 0 {
			continue
		}
		m[id] = strings.ToUpper(countryCode)
	}
	return m, nil
}

func getCidrPerCountry(file string, m map[string]string, list map[string][]*router.CIDR) error {
	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	reader := csv.NewReader(fileReader)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for _, line := range lines[1:] {
		cidrStr := line[0]
		countryID := line[1]
		if countryCode, found := m[countryID]; found {
			cidr, err := rule.ParseIP(cidrStr)
			if err != nil {
				return err
			}
			cidrs := append(list[countryCode], cidr)
			list[countryCode] = cidrs
		}
	}
	return nil
}

func getPrivateIPs() *router.GeoIP {
	cidr := make([]*router.CIDR, 0, len(privateIPs))
	for _, ip := range privateIPs {
		c, err := rule.ParseIP(ip)
		common.Must(err)
		cidr = append(cidr, c)
	}
	return &router.GeoIP{
		CountryCode: "PRIVATE",
		Cidr:        cidr,
	}
}

func getTestIPs() *router.GeoIP {
	cidr := make([]*router.CIDR, 0, len(testIPs))
	for _, ip := range testIPs {
		c, err := rule.ParseIP(ip)
		common.Must(err)
		cidr = append(cidr, c)
	}
	return &router.GeoIP{
		CountryCode: "TEST",
		Cidr:        cidr,
	}
}

func main() {
	flag.Parse()

	ccMap, err := getCountryCodeMap()
	if err != nil {
		fmt.Println("Error reading country code map:", err)
		os.Exit(1)
	}

	cidrList := make(map[string][]*router.CIDR)
	if err := getCidrPerCountry(*ipv4File, ccMap, cidrList); err != nil {
		fmt.Println("Error loading IPv4 file:", err)
		os.Exit(1)
	}
	if err := getCidrPerCountry(*ipv6File, ccMap, cidrList); err != nil {
		fmt.Println("Error loading IPv6 file:", err)
		os.Exit(1)
	}

	geoIPList := new(router.GeoIPList)
	for cc, cidr := range cidrList {
		geoIPList.Entry = append(geoIPList.Entry, &router.GeoIP{
			CountryCode: cc,
			Cidr:        cidr,
		})
	}
	geoIPList.Entry = append(geoIPList.Entry, getPrivateIPs())
	geoIPList.Entry = append(geoIPList.Entry, getTestIPs())

	geoIPBytes, err := proto.Marshal(geoIPList)
	if err != nil {
		fmt.Println("Error marshalling geoip list:", err)
		os.Exit(1)
	}

	// Create output directory if not exist
	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(*outputDir, 0755); mkErr != nil {
			fmt.Println("Failed: ", mkErr)
			os.Exit(1)
		}
	}

	if err := ioutil.WriteFile(filepath.Join(*outputDir, *outputName), geoIPBytes, 0644); err != nil {
		fmt.Println("Error writing geoip to file:", err)
		os.Exit(1)
	} else {
		fmt.Println(*outputName, "has been generated successfully.")
	}
}
