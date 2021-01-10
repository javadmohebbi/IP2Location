# IP2Location
A terminal application which use to get information about IP addresses based on IP2Location.com Lite Database

# About
This terminal application will use the IP2Location Lite Package to fetch the information about provided IP addresses. If the provided IPs are private addresses, you can provide a CSV file to fetch information from your local database.

You can find a how to guide in [this article](https://github.com/javadmohebbi/IP2Location/blob/master/GRAFANA-MONITOR-HOST.md) to Monitor your hosts using InfluxDB, Telegraf & Granfana. 


# Output
- It can output data as a CSV (-c csv) or TSV (-c tab)
- If you want to store the information in InfluxDB, you can use -t option and define it as a [[input.exec]] in telegraf.conf. Read more about Telegraf [[input.exec]] here: https://github.com/influxdata/telegraf/tree/master/plugins/inputs/exec
- When you use -t option, It will ping the provided IP and let you know the result based on each host


# Command line input parameters:
```
  -c string
        (Optional) Print it in CSV format. Accepted values are comma and tab. This option can NOT be used with -t option.
  -dl
        Download the latest lite ip2location database and replace it with the current and exit
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
        (Optional) ICMP requet timeout in seconds. (default "10")
  -v    Print Version & exit.
```


# How to build?
It's written in Go and I used https://github.com/ip2location/ip2location-go package to extract Information from IP2Location Lite Package. So to build executable from this source code download the repo and run the further command:

```
go get github.com/ip2location/ip2location-go
go get github.com/sparrc/go-ping
go get github.com/dustin/go-humanize
go build .\ip2location.go .\utils.go .\update.go
```

# Download binaries
Currently compiled executables are:
1. **Windows**:
      - [ip2location.exe (i386)](https://github.com/javadmohebbi/IP2Location/raw/master/dist/windows/386/ip2location.exe)
      - [ip2location.exe (amd64)](https://github.com/javadmohebbi/IP2Location/raw/master/dist/windows/amd64/ip2location.exe)
1. **Linux**:
      - [ip2location (i386)](https://github.com/javadmohebbi/IP2Location/raw/master/dist/linux/386/ip2location)
      - [ip2location (amd64)](https://github.com/javadmohebbi/IP2Location/raw/master/dist/linux/amd64/ip2location)
1. **Mac**:
      - [ip2location (i386)](https://github.com/javadmohebbi/IP2Location/raw/master/dist/darwin/386/ip2location)
      - [ip2location (amd64)](https://github.com/javadmohebbi/IP2Location/raw/master/dist/darwin/amd64/ip2location)

### NOTE for Linux Users
1. After download the executable you need to let it to be executable by running ```chmod +x ip2location```.
2. One of the other thing is if you use ```-t``` option, ip2location will try to send 3 icmp request to provided IP address. So you need to run it as **root**.
      - An alternate solution is: [Read Go-Ping Doc](https://github.com/sparrc/go-ping#note-on-linux-support)
      ```
            sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
            sudo setcap cap_net_raw=+ep /path/to/ip2location
      ```



# How to Use?
You need to download the Compiled executables from the ** build ** folder or build it from the source. After that, You need to download the latest Lite Package using ```-dl``` option.

## Get the latest Lite DB
Run the further command and it will download the latest Lite DB and place it in ```%ProgramFiles%\IP2Location\db\%``` (windows) or ```\etc\ip2location\db``` (Linux).
You need to run it as **root** (linux) or system **administrator** (windows)
```
# ip2location -dl

*** Outout ***

Download Started
Downloading... 24 MB complete
Unzipped:
db\IP2LOCATION-LITE-DB11.IPV6.BIN
Download Finished
```




# Local Database Format
The local database file format must be like this: CountryShort,CountryLong,State,City,Timezone,Lat,Lon,StartIP,EndIP
**For now, IP2Location Local database is not supports IPv6.**

**A sample for local CSV database**
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

*** Output ***
ip2loc,onlineTag=1,host=8.8.8.8 country_long="United States",country_short="US",city="Mountain View",state="California",timezone="-07:00",latitude=37.405991,longitude=-122.078514,packetSent=3,packetRecv=3,packetLost=0.000000,minRtt=82.852251,avgRtt=83.647769,maxRtt=84.680822,online=1
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



# Grafana Dashboard
I've made a simple dashboard for using this tool to monitor your hosts. It might give you a chance to get the idea behind IP2Location Tool. You might want to make more effective dashboards and share it with others.

You can download this dashboard from Grafana website. Here is the link to the dashboard: [Grafana IP2Location Dashbaord](https://grafana.com/grafana/dashboards/10964)

 After you import this dashbaord to your Grafana, You might see something like the further Screenshot:

![IP2Location + IP Monitor Simple Dashboard](http://mjmohebbi.com/public/img/uploads/upload__10-06-2019-054447.png)
![IP2Location + IP Monitor Simple Dashboard](http://mjmohebbi.com/public/img/uploads/upload__10-06-2019-054514.png)
