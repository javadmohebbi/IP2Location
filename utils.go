package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ip2location/ip2location-go"
	"github.com/sparrc/go-ping"
)

// IP2Location struct for fetched information
type IP2Location struct {
	CountryShort, CountryLong, Region, City, State, Timezone string
	Lat, Lon                                                 float32
}

const dbFileName string = "IP2LOCATION-LITE-DB11.IPV6.BIN"

var pathofExe string = ""
var tmpPath string = ""

// validateArgsAndCallFuncs - Validate Input Arguments and Call the needed Functions
func validateArgsAndCallFuncs(dir string, tmpDir string) {

	pathofExe = dir
	tmpPath = tmpDir

	version := "IP2Location | Version 0.0.1 | http://mjmohebbi.com \n\nLite Location Database by https://ip2location.com\n\n"

	argIP := flag.String("i", "", "(Could be Mandatory) IP Address (v4 or v6). FQND or Domain Name is not supported yet... ;-) ")
	argIPList := flag.String("list", "", "(Could be Mandatory) List of IP Address (v4 or v6). FQND or Domain Name is not supported yet... ;-) ")
	argCsv := flag.String("c", "", "(Optional) Print it in CSV format. Accepted values are comma and tab. This option can NOT be used with -t option.")
	argVer := flag.Bool("v", false, "Print Version & exit.")

	argTel := flag.String("t", "", "(Optional) Prepare it for InfluxDB Telegraf [[inputs.exec]]. It also will send 3 ICMP requets to the IP provided by the user to check the host availability. Accepted values are country_short, country_long, host or all  as a tag for InfluxDB. This option can NOT be used with -c option. Get more information about [[inputs.exec]] here: https://github.com/influxdata/telegraf/tree/master/plugins/inputs/exec")
	argMetric := flag.String("m", "ip2loc", "(Optional) Measurement name for InfluxDB. Default value is ip2loc. Use it with -t option")
	argTimeout := flag.String("timeout", "10", "(Optional) ICMP requet timeout in seconds.")
	argLocalDB := flag.String("local", "", "(Optional) Local CSV database for Private IPs. Only IPv4 is supported. It MUST be in this format: CountryShort,CountryLong,State,City,Timezone,Lat,Lon,StartIP,EndIP")
	argDownload := flag.Bool("dl", false, "Download the latest lite ip2location database and replace it with the current and exit")

	flag.Parse()

	if *argVer {
		fmt.Printf("\n%s\n", version)
		os.Exit(0)
	}

	if *argDownload {
		DownloadDatabase(dbFileName, pathofExe, tmpPath)
		os.Exit(0)
	}

	if *argIP == "" && *argIPList == "" {
		fmt.Println("***ERR*** IP address OR IP address list is not specified.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *argCsv != "" && *argTel != "" {
		fmt.Printf("\n***ERR*** -c and -t option can NOT be used at the same time. \n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *argCsv == "" && *argTel == "" {
		fmt.Printf("\n***ERR*** You should use one of -t or -c option. \n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *argCsv != "" && (*argCsv != "comma" && *argCsv != "tab") {
		fmt.Printf("\n***ERR*** Provided value (%s) for -c option is not acceptable!\n\n", *argCsv)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *argTel != "" && (*argTel != "country_long" && *argTel != "country_short" && *argTel != "city" && *argTel != "state" && *argTel != "timezone" && *argTel != "host" && *argTel != "all") {
		fmt.Printf("\n***ERR*** Provided value (%s) for -t option is not acceptable!\n\n", *argTel)
		flag.PrintDefaults()
		os.Exit(1)
	}

	var providedIPAddress []*string
	if *argIPList == "" {
		providedIPAddress = append(providedIPAddress, argIP)
	} else {
		providedIPAddress = ParsListAsArray(argIPList)
	}
	for _, ipaddress := range providedIPAddress {

		ip2loc, err := GetIPAddressInfo(ipaddress)
		// if err=false it means it might be a private IP address
		if err {
			if *argLocalDB == "" {
				fmt.Printf("\n***ERR*** For Private IP ranges you need to use -local option!\n\n")
				os.Exit(1)
			}
			ip2loc, err = GetPrivateIPAddressInfo(ipaddress, argLocalDB)
			if err {
				fmt.Printf("Provided IP (%s) is not provided in the local database (%s)\n", *ipaddress, *argLocalDB)
				// os.Exit(1)
				continue
			}
		}

		if *argCsv != "" {
			PrintCsv(ipaddress, argCsv, ip2loc)
		} else {
			PrintTelegraf(argTel, ipaddress, ip2loc, argMetric, argTimeout)
		}
	}

}

// ParsListAsArray - Pars list of string as []*string
func ParsListAsArray(argIPList *string) []*string {
	var returnArray []*string
	csvFile, err := os.Open(*argIPList)
	if err != nil {
		fmt.Println("Couldn't open the IP list file", err)
		os.Exit(1)
	}
	r := csv.NewReader(csvFile)
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			return returnArray
		}
		if err != nil {
			fmt.Printf("Can not read %s file", *argIPList)
			os.Exit(1)
		}
		// tmpIP := &record[0]
		// fmt.Println(*tmpIP)
		returnArray = append(returnArray, &record[0])
	}
}

// PrintCsv - Print Location Info in a CSV Format
func PrintCsv(ip *string, format *string, ip2loc IP2Location) {
	delimiter := "\t"
	if *format == "comma" {
		delimiter = ","
	}
	fmt.Printf("\"%s\"%s\"%s\"%s\"%s\"%s\"%s\"%s\"%s\"%s\"%s\"%s%f%s%f\n",
		*ip, delimiter,
		ip2loc.CountryShort, delimiter,
		ip2loc.CountryLong, delimiter,
		ip2loc.Region, delimiter,
		ip2loc.City, delimiter,
		ip2loc.Timezone, delimiter,
		ip2loc.Lat, delimiter,
		ip2loc.Lon)
}

// PrintTelegraf - Print Location Info for InfluxData Telegraf. Also it PINGs the requested IP address
func PrintTelegraf(tag *string, ip *string, ip2loc IP2Location, mes *string, tm *string) {
	tagFieldValue := ""

	pinger, err := ping.NewPinger(*ip)
	pinger.SetPrivileged(true)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.Timeout, _ = time.ParseDuration(*tm + "s")

	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats

	isHostOnline := 2 // 1 = Online, 2 = Offline
	if stats.PacketsRecv > 1 {
		isHostOnline = 1
	}

	if *tag == "countryLongTag" {
		tagFieldValue = ip2loc.CountryLong
	} else if *tag == "countryShortTag" {
		tagFieldValue = ip2loc.CountryShort
	} else if *tag == "CityTag" {
		tagFieldValue = ip2loc.City
	} else if *tag == "StateTag" {
		tagFieldValue = ip2loc.Region
	} else if *tag == "TimeZoneTag" {
		tagFieldValue = ip2loc.Timezone
	} else if *tag == "host" {
		tagFieldValue = *ip
	} else {
		tagFieldValue = ""
	}

	minRTT := float64(stats.MinRtt) / float64(time.Millisecond)
	avgRTT := float64(stats.AvgRtt) / float64(time.Millisecond)
	maxRTT := float64(stats.MaxRtt) / float64(time.Millisecond)

	// Escaping Space and Comma
	tgCL := strings.Replace(ip2loc.CountryLong, ",", " ", -1)
	tgCL = strings.Replace(tgCL, " ", "_", -1)
	tgCS := strings.Replace(ip2loc.CountryShort, " ", "-", -1)
	tgCT := strings.Replace(ip2loc.City, ",", " ", -1)
	tgCT = strings.Replace(tgCT, " ", "-", -1)
	tgRG := strings.Replace(ip2loc.Region, ",", " ", -1)
	tgRG = strings.Replace(tgRG, " ", "-", -1)
	tgFieldVal := strings.Replace(tagFieldValue, ",", " ", -1)
	tgFieldVal = strings.Replace(tgFieldVal, " ", "-", -1)
	if tagFieldValue != "" {
		fmt.Printf("%s,online=%d,%s=%s country_long=\"%s\",country_short=\"%s\",city=\"%s\",state=\"%s\",timezone=\"%s\",latitude=%f,longitude=%f,packetSent=%d,packetRecv=%d,packetLost=%f,minRtt=%v,avgRtt=%v,maxRtt=%v,online=%d", *mes, isHostOnline, *tag, tgFieldVal, ip2loc.CountryLong, ip2loc.CountryShort, ip2loc.City, ip2loc.Region, ip2loc.Timezone, ip2loc.Lat, ip2loc.Lon, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss, minRTT, avgRTT, maxRTT, isHostOnline)
	} else {
		fmt.Printf("%s,online=%d,host=%s,countryLongTag=%s,countryShortTag=%s,CityTag=%s,StateTag=%s,TimeZoneTag=%s country_long=\"%s\",country_short=\"%s\",city=\"%s\",state=\"%s\",timezone=\"%s\",latitude=%f,longitude=%f,packetSent=%d,packetRecv=%d,packetLost=%f,minRtt=%v,avgRtt=%v,maxRtt=%v,online=%d", *mes, isHostOnline, *ip, tgCL, tgCS, tgCT, tgRG, ip2loc.Timezone, ip2loc.CountryLong, ip2loc.CountryShort, ip2loc.City, ip2loc.Region, ip2loc.Timezone, ip2loc.Lat, ip2loc.Lon, stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss, minRTT, avgRTT, maxRTT, isHostOnline)
	}
	fmt.Printf("\n")

}

// GetIPAddressInfo - Fetch Address Info
func GetIPAddressInfo(ip *string) (IP2Location, bool) {
	ip2location.Open(pathofExe + "db/" + dbFileName)

	results := ip2location.Get_all(*ip)

	if results.Country_short == "Invalid IP address." {
		fmt.Println("Provided IP is not valid: ", *ip)
		// os.Exit(1)
		return IP2Location{}, true
	}

	if results.Country_short == "-" {
		// fmt.Println("Provided IP might be a private address: ", *ip)
		// os.Exit(1)
		return IP2Location{}, true
	}

	tmpIP2Loc := IP2Location{CountryLong: results.Country_long, CountryShort: results.Country_short, Region: results.Region, City: results.City, Timezone: results.Timezone, Lat: results.Latitude, Lon: results.Longitude}

	ip2location.Close()

	return tmpIP2Loc, false
}

// GetPrivateIPAddressInfo - If Private CSV files is provided, It will use this file show information
func GetPrivateIPAddressInfo(ip *string, localDB *string) (IP2Location, bool) {
	CountryShort, CountryLong, State, City, Timezone, Lati, Loni, StartIP, EndIP := 0, 1, 2, 3, 4, 5, 6, 7, 8

	csvFile, err := os.Open(*localDB)
	if err != nil {
		fmt.Println("Couldn't open the csv file", err)
		return IP2Location{}, true
		// os.Exit(1)
	}
	r := csv.NewReader(csvFile)
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			// fmt.Printf("Can not read %s file", *localDB)
			// os.Exit(1)
			return IP2Location{}, true
		}
		if err != nil {
			fmt.Printf("Can not read %s file", *localDB)
			os.Exit(1)
		}

		// check if the provided IP is in the range
		if IsItInTheRangeIPv4(ip, record[StartIP], record[EndIP]) {
			tmpLat, _ := strconv.ParseFloat(record[Lati], 32)
			tmpLon, _ := strconv.ParseFloat(record[Loni], 32)
			return IP2Location{
				CountryShort: record[CountryShort],
				CountryLong:  record[CountryLong],
				Region:       record[State],
				City:         record[City],
				State:        record[State],
				Timezone:     record[Timezone],
				Lat:          float32(tmpLat),
				Lon:          float32(tmpLon),
			}, false
		}
	}
}

// IsItInTheRangeIPv4 return true if its in the range
func IsItInTheRangeIPv4(ip *string, startIP string, endIP string) bool {
	trial := net.ParseIP(*ip)

	// ip is NOT an IPv4
	if trial.To4() == nil {
		fmt.Println("NOT v4")
		return false
	}

	// ip is in the range
	if bytes.Compare(trial, net.ParseIP(startIP)) >= 0 && bytes.Compare(trial, net.ParseIP(endIP)) <= 0 {
		return true
	}

	// ip is not in the range
	return false
}
