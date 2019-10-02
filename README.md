# IP2Location
A terminal application which use to get information about IP addresses based on IP2Location.com Lite Database

# About
This terminal application will use the IP2Location Lite Package to fetch the information about provided IP addresses. If the provided IPs are private addresses, you can provide a CSV file to fetch information from your local database. 

# Output
- It can output data as a CSV (-c csv) or TSV (-c tab)
- If you want to store the information in InfluxDB, you can use -t option and define it as a [[input.exec]] in telegraf.conf. Read more about Telegraf [[input.exec]] here: https://github.com/influxdata/telegraf/tree/master/plugins/inputs/exec
- When you use -t option, It will ping the provided IP and let you know the result based on each host


# Command line input paramentes:
```
  -c string
        (Optional) Print it in CSV format. Accepted values are comma and tab. This option can NOT be used with -t option.
  -i string
        (Could be Mandatory) IP Address (v4 or v6). FQND or Domain Name is not supported yet... ;-)
  -list string
        (Could be Mandatory) List of IP Address (v4 or v6). FQND or Domain Name is not supported yet... ;-)
  -local string
        (Optional) Local CSV database for Private IPs. Only IPv4 is supported. It MUST be in this format: CountryShort,CountryLong,State,City,Timezone,Lat,Lon,StartIP,EndIP
  -m string
        (Optional) Measurement name for InfluxDB. Default value is ip2loc. Use it with -t option (default "ip2loc")
  -t string
        (Optional) Prepare it for InfluxDB Telegraf [[inputs.exec]]. It also will send 3 ICMP requets to the IP provided by the user to check the host availability. Accepted values are country_short, country_long, host or all  as a tag for InfluxDB. This option can NOT be used with -c option. Get more information about [[inputs.exec]] here: https://github.com/influxdata/telegraf/tree/master/plugins/inputs/exec
  -timeout string
        (Optional) ICMP requet timeout in seconds. (default "10s")
  -v    Print Version & exit.
```

# Local Database Format
The local database file format must be like this: CountryShort,CountryLong,State,City,Timezone,Lat,Lon,StartIP,EndIP
**For now, IP2Location Local database is not supports IPv6.**

**A sample for local database**
```
IR,"Iran, Islamic Republic of",Tehran,Tehran,+03:30,35.6892,51.3890,172.16.50.1,172.16.50.254
IR,"Iran, Islamic Republic of",Markazi,Arak,+03:30,34.0954,49.7013,10.0.0.20,10.0.0.30
IR,"Iran, Islamic Republic of",Isfahan,Isfahan,+03:30,32.6539,51.6660,172.16.2.1,172.16.2.100
```


# Exmaples
There are some few examples which you can use it

**Example 1: Get Location Info in CSV format**
```
ip2location -i 8.8.8.8 -c comma

Output is:
"8.8.8.8","US","United States","California","Mountain View","-07:00",37.405991,-122.078514
```

**Example 2: Get Location Info in TSV format**
```
ip2location -i 8.8.8.8 -c tab

Output is:
"8.8.8.8"       "US"    "United States" "California"    "Mountain View" "-07:00"        37.405991       -122.078514
```

**Example 3: Get Location Info and prepare it for Telefraf**

After -t option you can use country_short, country_long, host or all as a tag for InfluxDB. In this example we will use host
```
ip2location -i 8.8.8.8 -t host

Output is:
ip2loc,host="8.8.8.8" latitude=37.405991,longitude=-122.078514,packetSent=3,packetRecv=3,packetLost=0.000000,minRtt=85.7711,avgRtt=90.0008,maxRtt=98.4593,online=1
```


**Example 4: Using LocalDB and List of IP addresses**

The only hint about List of IP addresses is, the first field of CSV file must be IP Address & You can add as much as field you want in another fields. Here is the example:
```
8.8.8.8,"Google"
172.16.50.16,"Tehran HQ"
10.0.0.6,"Arak Office"
10.5.3.5,"Isfahana Office #1"
172.16.2.3,"Isfahana Office #2"
```

And Here is the example and It's Out put in CSV format:
```
"8.8.8.8","US","United States","California","Mountain View","-07:00",37.405991,-122.078514
"172.16.50.16","IR","Iran, Islamic Republic of","Tehran","Tehran","+03:30",35.689201,51.389000
Provided IP (10.0.0.6) is not provided in the local database (local.csv)
Provided IP (10.5.3.5) is not provided in the local database (local.csv)
"172.16.2.3","IR","Iran, Islamic Republic of","Isfahan","Isfahan","+03:30",32.653900,51.666000
```



# How to build the package:
It's written in Go and I used https://github.com/ip2location/ip2location-go package to extract Information from IP2Location Lite Package. So to build executable from this source code download the repo and run the further command:

```
go build .\ip2location.go .\utils.go
```



