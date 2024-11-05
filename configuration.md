# Configuration explanation

## Overview

The format of the configuration file used in this project is `json`. The JSON configuration contains two arrays, `input` and `output`, each of which contains a specific configuration for one or more input or output formats.

```json
{
  "input":  [],
  "output": []
}
```

## Supported formats

Supported `input` formats:

- **cutter**: Remove data from previous steps
- **maxmindGeoLite2CountryCSV**: Convert MaxMind GeoLite2 country CSV data to other formats
- **maxmindMMDB**: Convert MaxMind country mmdb database to other formats
- **private**: Convert LAN and private network CIDR to other formats
- **text**: Convert plaintext IP and CIDR to other formats
- **v2rayGeoIPDat**: Convert V2Ray GeoIP dat to other formats

Supported `output` formats:

- **text**: Convert data to plaintext CIDR format
- **v2rayGeoIPDat**: Convert data to V2Ray GeoIP dat format

## Configuration options for `input` formats

### **cutter**

- **type**: (required) the name of the input format
- **action**: (required) action type, the value must be `remove` (to remove IP / CIDR)
- **args**: (required)
  - **wantedList**: (required, array) specified wanted lists
  - **onlyIPType**: (optional) the IP address type to be processed, the value is `ipv4` or `ipv6`

```jsonc
{
  "type": "cutter",
  "action": "remove",                // remove IP or CIDR
  "args": {
    "wantedList": ["cn", "us", "jp"] // remove IPv4 and IPv6 addresses from lists called cn, us, jp, AKA remove the three entire lists
  }
}
```

```jsonc
{
  "type": "cutter",
  "action": "remove",                 // remove IP or CIDR
  "args": {
    "wantedList": ["cn", "us", "jp"],
    "onlyIPType": "ipv6"              // remove IPv6 addresses from lists called cn, us, jp
  }
}
```

### **maxmindGeoLite2CountryCSV**

- **type**: (required) the name of the input format
- **action**: (required) action type, the value could be `add`(to add IP / CIDR) or `remove`(to remove IP / CIDR)
- **args**: (optional)
  - **country**: (optional) the path to MaxMind GeoLite2 Country CSV location file (`GeoLite2-Country-Locations-en.csv`), can be local file path or remote `http` or `https` URL
  - **ipv4**: (optional) the path to MaxMind GeoLite2 Country IPv4 file (`GeoLite2-Country-Blocks-IPv4.csv`), can be local file path or remote `http` or `https` URL
  - **ipv6**: (optional) the path to MaxMind GeoLite2 Country IPv6 file (`GeoLite2-Country-Blocks-IPv6.csv`), can be local file path or remote `http` or `https` URL
  - **wantedList**: (optional, array) specified wanted lists
  - **onlyIPType**: (optional) the IP address type to be processed, the value is `ipv4` or `ipv6`

```jsonc
// Files to be used by default:
// ./geolite2/GeoLite2-Country-Locations-en.csv
// ./geolite2/GeoLite2-Country-Blocks-IPv4.csv
// ./geolite2/GeoLite2-Country-Blocks-IPv6.csv
{
  "type": "maxmindGeoLite2CountryCSV",
  "action": "add" // add IP or CIDR
}
```

```jsonc
{
  "type": "maxmindGeoLite2CountryCSV",
  "action": "add",                     // add IP or CIDR
  "args": {
    "country": "./geolite2/GeoLite2-Country-Locations-en.csv",
    "ipv4": "./geolite2/GeoLite2-Country-Blocks-IPv4.csv",
    "ipv6": "./geolite2/GeoLite2-Country-Blocks-IPv6.csv"
  }
}
```

```jsonc
{
  "type": "maxmindGeoLite2CountryCSV",
  "action": "add",                   // add IP or CIDR
  "args": {
    "wantedList": ["cn", "us", "jp"] // add IPv4 and IPv6 addresses to lists called cn, us, jp 
  }
}
```

```jsonc
{
  "type": "maxmindGeoLite2CountryCSV",
  "action": "remove",                 // remove IP or CIDR
  "args": {  
    "wantedList": ["cn", "us", "jp"], // only to remove IPv6 addresses from lists called cn, us, jp
    "onlyIPType": "ipv6"              // only to remove IPv6 addresses
  }
}
```

### **maxmindMMDB**

- **type**: (required) the name of the input format
- **action**: (required) action type, the value could be `add`(to add IP / CIDR) or `remove`(to remove IP / CIDR)
- **args**: (optional)
  - **uri**: (optional) the path to MaxMind GeoLite2 Country mmdb file(`GeoLite2-Country.mmdb`), can be local file path or remote `http` or `https` URL
  - **wantedList**: (optional, array) specified wanted lists
  - **onlyIPType**: (optional) the IP address type to be processed, the value is `ipv4` or `ipv6`

```jsonc
// The file to be used by default:
// ./geolite2/GeoLite2-Country.mmdb
{
  "type": "maxmindMMDB",
  "action": "add"       // add IP or CIDR
}
```

```jsonc
{
  "type": "maxmindMMDB",
  "action": "add",       // add IP or CIDR
  "args": {
    "uri": "./geolite2/GeoLite2-Country.mmdb"
  }
}
```

```jsonc
{
  "type": "maxmindMMDB",
  "action": "add",                        // add IP or CIDR
  "args": {
    "uri": "https://example.com/my.mmdb",
    "wantedList": ["cn", "us", "jp"],    // add IPv4 addresses to lists called cn, us, jp
    "onlyIPType": "ipv4"                 // only to add IPv4 addresses
  }
}
```

```jsonc
{
  "type": "maxmindMMDB",
  "action": "remove",                    // add IP or CIDR
  "args": {
    "uri": "https://example.com/my.mmdb",
    "wantedList": ["cn", "us", "jp"],    // only to remove IPv4 addresses from lists called cn, us, jp
    "onlyIPType": "ipv4"                 // only to remove IPv4 addresses
  }
}
```

### **private**

- **type**: (required) the name of the input format
- **action**: (required) action type, the value could be `add`(to add IP / CIDR) or `remove`(to remove IP / CIDR)
- **args**: (optional)
  - **onlyIPType**: (optional) the IP address type to be processed, the value is `ipv4` or `ipv6`

> The default CIDRs to be added to or removed from `private`, see [private.go](https://github.com/v2fly/geoip/blob/HEAD/plugin/special/private.go#L16-L36).

```jsonc
{
  "type": "private",
  "action": "add"   // add IP or CIDR
}
```

```jsonc
{
  "type": "private",
  "action": "remove" // remove IP or CIDR
}
```

```jsonc
{
  "type": "private",
  "action": "add",       // add IP or CIDR
  "args": {
    "onlyIPType": "ipv4" // add IPv4 addresses only
  }
}
```

```jsonc
{
  "type": "private",
  "action": "remove",    // remove IP or CIDR
  "args": {
    "onlyIPType": "ipv6" // remove IPv6 addresses only
  }
}
```

### **text**

- **type**: (required) the name of the input format
- **action**: (required) action type, the value could be `add`(to add IP / CIDR) or `remove`(to remove IP / CIDR)
- **args**: (required)
  - **name**: (optional) the list name (cannot be used with `inputDir`; must be used with `uri` or `ipOrCIDR`)
  - **uri**: (optional) the path to plaintext txt file, can be local file path or remote `http` or `https` URL (cannot be used with `inputDir`; must be used with `name`; can be used with `ipOrCIDR`)
  - **ipOrCIDR**: (optional, array) an array of plaintext IP addresses or CIDRs (cannot be used with `inputDir`; must be used with `name`; can be used with `uri`)
  - **inputDir**: (optional) the directory of the files to walk through (excluded children directories). (the filename will be the list name; cannot be used with `name` or `uri` or `ipOrCIDR`)
  - **wantedList**: (optional, array) specified wanted files. (used with `inputDir`)
  - **onlyIPType**: (optional) the IP address type to be processed, the value is `ipv4` or `ipv6`
  - **removePrefixesInLine**: (optional, array) the array of string prefixes to be removed in each line
  - **removeSuffixesInLine**: (optional, array) the array of string suffixes to be removed in each line

```jsonc
{
  "type": "text",
  "action": "add",                                // add IP or CIDR
  "args": {
    "name": "cn",
    "uri": "./cn.txt",                            // get IPv4 and IPv6 addresses from local file cn.txt, and add to list cn
    "removePrefixesInLine": ["Host,", "IP-CIDR"], // remove all prefixes from each line of the file
    "removeSuffixesInLine": [",no-resolve"]       // remove all suffixes from each line of the file
  }
}
```

```jsonc
{
  "type": "text",
  "action": "add",                        // add IP or CIDR
  "args": {
    "name": "cn",
    "ipOrCIDR": ["1.0.0.1", "1.0.0.1/24"] // add IP or CIDR to cn list
  }
}
```

```jsonc
{
  "type": "text",
  "action": "remove",                     // remove IP or CIDR
  "args": {
    "name": "cn",
    "ipOrCIDR": ["1.0.0.1", "1.0.0.1/24"] // remove IP or CIDR from cn list
  }
}
```

```jsonc
{
  "type": "text",
  "action": "add",                        // add IP or CIDR
  "args": {
    "name": "cn",
    "uri": "./cn.txt",                    // get IPv4 and IPv6 addresses from local file cn.txt, and add to list cn
    "ipOrCIDR": ["1.0.0.1", "1.0.0.1/24"] // add IP or CIDR to cn list
  }
}
```

```jsonc
{
  "type": "text",
  "action": "add", // add IP or CIDR
  "args": {
    "inputDir": "./text",                         // walk through all files in directory ./text (excluded children directories by default)
    "wantedList": ["cn", "us", "jp"],             // wanted lists called cn, us, jp without file extension
    "onlyIPType": "ipv6",                         // add IPv6 addresses only
    "removePrefixesInLine": ["Host,", "IP-CIDR"], // remove all prefixes from each line of each file
    "removeSuffixesInLine": [",no-resolve"]       // remove all suffixes from each line of each file
  }
}
```

```jsonc
{
  "type": "text",
  "action": "remove",                             // remove IP or CIDR
  "args": {
    "name": "cn",
    "uri": "https://example.com/cn.txt",          // read the content of the remote file
    "onlyIPType": "ipv6",                         // remove only IPv6 addresses from list called cn
    "removePrefixesInLine": ["Host,", "IP-CIDR"], // remove all prefixes from each line of the file
  }
}
```

```jsonc
{
  "type": "text",
  "action": "remove",                       // remove IP or CIDR
  "args": {
    "name": "cn",
    "uri": "https://example.com/cn.txt",    // read the content of the remote file
    "onlyIPType": "ipv6",                   // remove only IPv6 addresses from list called cn
    "removeSuffixesInLine": [",no-resolve"] // remove all suffixes from each line of the file
  }
}
```

### **v2rayGeoIPDat**

- **type**: (required) the name of the input format
- **action**: (required) action type, the value could be `add`(to add IP / CIDR) or `remove`(to remove IP / CIDR)
- **args**: (required)
  - **uri**: (required) the path to V2Ray dat format geoip file, can be local file path or remote `http` or `https` URL
  - **wantedList**: (optional, array) specified wanted lists
  - **onlyIPType**: (optional) the IP address type to be processed, the value is `ipv4` or `ipv6`

```jsonc
{
  "type": "v2rayGeoIPDat",
  "action": "add",         // add IP or CIDR
  "args": {
    "uri": "./cn.dat"      // add IPv4 and IPv6 addresses of local file cn.dat to lists
  }
}
```

```jsonc
{
  "type": "v2rayGeoIPDat",
  "action": "add",                    // add IP or CIDR
  "args": {
    "uri": "./geoip.dat",             // read from local file geoip.dat
    "wantedList": ["cn", "us", "jp"], // wanted lists called cn, us, jp
    "onlyIPType": "ipv6"              // only to add IPv6 addresses
  }
}
```

```jsonc
{
  "type": "v2rayGeoIPDat",
  "action": "remove",                       // remove IP or CIDR
  "args": {
    "uri": "https://example.com/geoip.dat", // read the content of the remote file
    "onlyIPType": "ipv6"                    // remove only IPv6 addresses from all lists
  }
}
```

## Configuration options for `output` formats

### **text**

- **type**: (required) the name of the output format
- **action**: (required) action type, the value must be `output`
- **args**: (optional)
  - **outputDir**: (optional) path to the output directory
  - **outputExtension**: (optional) the extension of the output file
  - **wantedList**: (optional, array) specified wanted lists
  - **excludedList**: (optional, array) specified lists to be excluded when output
  - **onlyIPType**: (optional) the IP address type to output, the value is `ipv4` or `ipv6`
  - **addPrefixInLine**: (optional) the prefix to be added in each line
  - **addSuffixInLine**: (optional) the suffix to be added in each line

```jsonc
// The output directory by default:
// ./output/text
{
  "type": "text",
  "action": "output",
  "args": {
    "outputDir": "./text",           // output files to directory ./text
    "outputExtension": ".conf",      // the extension of the output files are .conf
    "addPrefixInLine": "IP-CIDR,",   // add prefix to each line
    "addSuffixInLine": ",no-resolve" // add suffix to each line
  }
}
```

```jsonc
{
  "type": "text",
  "action": "output",
  "args": {
    "outputDir": "./text",           // output files to directory ./text
    "outputExtension": ".conf",      // the extension of the output files are .conf
    "addPrefixInLine": "IP-CIDR,",
    "addSuffixInLine": ",no-resolve"
  }
}
```

```jsonc
{
  "type": "text",
  "action": "output",
  "args": {
    "outputDir": "./text",            // output files to directory ./text
    "outputExtension": ".conf",       // the extension of the output files are .conf
    "wantedList": ["cn", "us", "jp"], // output IPv4 and IPv6 addresses of lists called cn, us, jp
    "addPrefixInLine": "HOST,"        // add prefix to each line of each file
  }
}
```

```jsonc
{
  "type": "text",
  "action": "output",
  "args": {
    "outputDir": "./text",            // output files to directory ./text
    "outputExtension": ".conf",       // the extension of the output files are .conf
    "wantedList": ["cn", "us", "jp"], // output only IPv4 addresses of lists called cn, us, jp
    "onlyIPType": "ipv4",             // output IPv4 addresses only
    "addSuffixInLine": ";"            // add suffix to each line of each file
  }
}
```

```jsonc
{
  "type": "text",
  "action": "output",
  "args": {
    "outputDir": "./text",              // output files to directory ./text
    "outputExtension": ".conf",         // the extension of the output files are .conf
    "excludedList": ["cn", "us", "jp"], //  exclude lists called cn, us, jp when output
    "addPrefixInLine": "HOST,"
  }
}
```

### **v2rayGeoIPDat**

- **type**: (required) the name of the output format
- **action**: (required) action type, the value must be `output`
- **args**: (optional)
  - **outputName**: (optional) the output filename
  - **outputDir**: (optional) path to the output directory
  - **wantedList**: (optional, array) specified wanted lists or files
  - **excludedList**: (optional, array) specified lists to be excluded when output
  - **onlyIPType**: (optional) the IP address type to output, the value is `ipv4` or `ipv6`
  - **oneFilePerList**: (optional) output every single list to a new file, the value is `true` or `false`(default value)

```jsonc
// The output directory by default:
// ./output/dat
{
  "type": "v2rayGeoIPDat",
  "action": "output"      // output all lists to one file
}
```

```jsonc
{
  "type": "v2rayGeoIPDat",
  "action": "output",
  "args": {
    "oneFilePerList": true // output every single list to a new file
  }
}
```

```jsonc
{
  "type": "v2rayGeoIPDat",
  "action": "output",
  "args": {
    "outputDir": "./output",                   // output to ./output directory
    "outputName": "geoip-only-cn-private.dat", // output file called geoip-only-cn-private.dat
    "wantedList": ["cn", "private"]            // only output lists called cn, private
  }
}
```

```jsonc
{
  "type": "v2rayGeoIPDat",
  "action": "output",
  "args": {
    "outputDir": "./output",                      // output to ./output directory
    "outputName": "geoip-without-cn-private.dat", // output file called geoip-without-cn-private.dat
    "excludedList": ["cn", "private"]             // exclude lists called cn, private when output
  }
}
```

```jsonc
{
  "type": "v2rayGeoIPDat",
  "action": "output",
  "args": {
    "outputName": "geoip-asn.dat",        // output file called geoip-asn.dat
    "wantedList": ["telegram", "google"], // only output lists called telegram, google
    "onlyIPType": "ipv4"                  // output only IPv4 addresses of lists called telegram, google
  }
}
```

```jsonc
{
  "type": "v2rayGeoIPDat",
  "action": "output",
  "args": {
    "wantedList": ["telegram", "google"], // only output lists called telegram, google
    "onlyIPType": "ipv4",                 // output only IPv4 addresses of lists called telegram, google
    "oneFilePerList": true                // output every single list to a new file
  }
}
```
