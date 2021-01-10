# About this article!

I myself needed to monitor some hosts availability on a Map and by using some sort of Telegraf plugins & shellscript I was able to do that. Since some of our offices & branches are connected to HQ using a tunnel and we have private IP addresses which we could not find their location coordination, I developed a simple IP2Location terminal application which can do this for Public & Private IPs using IP2Location Lite database (for Public IPs) and your local IP2Location DB (for private IPs). So, Let's continue reading this blog to know how to use this tool.
What do we need?

This tutorial is tested on Ubuntu 18.04 and with some little change can be run on any Windows, Linux or Mac systems.

   InfluxDB & Grafana: If you don't know how to install them, you can read this blog post: How to install Grafana. Also if you need to Install InfluxDB & Telegraf you can read their documentation here: InfluxDB Installation and Telegraf Installation
   IP2Location: Is a terminal application which get information about IP addresses based on IP2Location.com Lite Database. This terminal application will use IP2Location Lite Package to fetch the information about provided IP addresses. If the IPs are private addresses, you can provide a CSV file to fetch location information from your local database. You can read it's documentation on it's Github repository: IP2Location Github Repository

# Download & Prepare IP2Location

- To download IP2Location executable package, run this command:
  - Linux AMD64: ```wget https://github.com/javadmohebbi/IP2Location/raw/master/dist/linux/amd64/ip2location```
  - Linux I386: ```wget https://github.com/javadmohebbi/IP2Location/raw/master/dist/linux/386/ip2location```

  - Copy this file to /usr/bin/ directory and makes it executable

  - ```sudo cp ip2location /usr/bin/ -v```
  - ```sudo chmod +x /usr/bin/ip2location -v```


- Get the latest version of IP2Location Lite DB: (Every month you might need to run this command to update DB) ```sudo ip2location -dl```
  - You must see this output:
  - ```
$ sudo ip2location -dl
*** Outout must be like this ***

Download Started
Downloading... 24 MB complete
Unzipped:
/etc/ip2location/db/IP2LOCATION-LITE-DB11.IPV6.BIN
Download Finished
```

- Since we need to send ICMP request without root access, we need to run these commands:
```
sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
sudo setcap cap_net_raw=+ep /usr/bin/ip2location
```

- In order to test your configuration run this command and you must get one line output in your terminal:
```
ip2location -i 8.8.8.8 -t all

*** Output ***
ip2loc,host=8.8.8.8,countryLongTag=United_States,countryShortTag=US,CityTag=Mountain-View,StateTag=California,TimeZoneTag=-07:00 country_long="United States",country_short="US",city="Mountain View",state="California",timezone="-07:00",latitude=37.405991,longitude=-122.078514,packetSent=3,packetRecv=3,packetLost=0.000000,minRtt=79.244016,avgRtt=79.890058,maxRtt=80.809445,online=1
```

# How to use Telegraf [[inputs.exec]]

Telegraf is a plugin-deriven agent for collecting and sending metrics and events from databases, systems and IoT sensors (InfluxData Website). You can read it's manual on this page: Telegraf Manual

I want to use Telegraf [[inputs.exec]] & run ip2location to get information from the hosts I want to monitor and store them in InfluxDB database. You can read it's documentation here: Telegraf [[inputs.exec]]

To make a telegraf configuration do the following steps:
```
   Create a configuration file: touch /etc/telegraf/telegraf.d/ip2location-telegraf.conf

   Edit this file like this:

   [outputs.influxdb]]
     urls = [ "http://127.0.0.1:8086"]
     database = "ip2location"


   [[inputs.exec]]
     commands = ["/usr/bin/ip2location -list /etc/ip2location/list.csv -local /etc/ip2location/local.csv -t all" ]
     interval = "2m"
     timeout = "Number of Host * 10"
     data_format = "influx"				
```         

- Since the default timeout in ip2location is 10 seconds, you can configure the telegraf timeout for this command as (Number of Host in list.csv * 10s). For instance if you have 10 hosts, It must be 10 * 10 = 100s

- Interval defines the time that Telegraf must run this command. It must be more than timeout. For instance if you have 100s timeout you'd better to set it as 120s or 2m.


# Before we get start!

We must create two file list.csv (List of IP addresses which you want to monitor) & local.csv (Private IP2Location Database for your Private IP Addresses).

**list.csv** File:

Create a file inside ```/etc/ip2location/list.csv``` and put the IP addreses you want to monitor. First field of this comma delimited file must be IP address and you can add as many as fields you want for your information.
```
8.8.8.8,"Google"
172.16.50.16,"Tehran HQ"
10.0.0.6,"Arak Office"
10.5.3.5,"Isfahana Office #1"
172.16.2.3,"Isfahana Office #2"
```

**local.csv** File:

Create a file inside ```/etc/ip2location/local.csv```. This file must have this format: **CountryShort,CountryLong,State,City,Timezone,Lat,Lon,StartIP,EndIP**.

It's only support IPv4 for now!

My sample data file is:
```
IR,"Iran, Islamic Republic of",Tehran,Tehran,+03:30,35.6892,51.3890,172.16.50.1,172.16.50.254
IR,"Iran, Islamic Republic of",Markazi,Arak,+03:30,34.0954,49.7013,10.0.0.4,10.7.0.30
IR,"Iran, Islamic Republic of",Isfahan,Isfahan,+03:30,32.6539,51.6660,172.16.2.1,172.16.2.100
```

## Test your list.csv and local.csv

To test your configuration, you can run this command:
```
sudo -u telegraf telegraf --test --config /etc/telegraf/telegraf.d/ip2location-telegraf.conf

Output must be like this:


2019-10-06T06:57:34Z I! Starting Telegraf 1.11.1
> ip2loc,CityTag=Mountain-View,StateTag=California,TimeZoneTag=-07:00,countryLongTag=United_States,countryShortTag=US,host=8.8.8.8,onlineTag=1 avgRtt=84.102891,city="Mountain View",country_long="United States",country_short="US",latitude=37.405991,longitude=-122.078514,maxRtt=85.904206,minRtt=82.71992,online=1,packetLost=0,packetRecv=3,packetSent=3,state="California",timezone="-07:00" 1570345087000000000
> ip2loc,CityTag=Tehran,StateTag=Tehran,TimeZoneTag=+03:30,countryLongTag=Iran__Islamic_Republic_of,countryShortTag=IR,host=172.16.50.143,onlineTag=1 avgRtt=0.527725,city="Tehran",country_long="Iran, Islamic Republic of",country_short="IR",latitude=35.689201,longitude=51.389,maxRtt=0.699684,minRtt=0.404185,online=1,packetLost=0,packetRecv=3,packetSent=3,state="Tehran",timezone="+03:30" 1570345087000000000
> ip2loc,CityTag=Tehran,StateTag=Tehran,TimeZoneTag=+03:30,countryLongTag=Iran__Islamic_Republic_of,countryShortTag=IR,host=172.16.50.16,onlineTag=1 avgRtt=0.564527,city="Tehran",country_long="Iran, Islamic Republic of",country_short="IR",latitude=35.689201,longitude=51.389,maxRtt=0.645868,minRtt=0.487628,online=1,packetLost=0,packetRecv=3,packetSent=3,state="Tehran",timezone="+03:30" 1570345087000000000
> ip2loc,CityTag=Arak,StateTag=Markazi,TimeZoneTag=+03:30,countryLongTag=Iran__Islamic_Republic_of,countryShortTag=IR,host=10.0.0.6,onlineTag=2 avgRtt=0,city="Arak",country_long="Iran, Islamic Republic of",country_short="IR",latitude=34.095402,longitude=49.701302,maxRtt=0,minRtt=0,online=2,packetLost=100,packetRecv=0,packetSent=3,state="Markazi",timezone="+03:30" 1570345087000000000
> ip2loc,CityTag=Arak,StateTag=Markazi,TimeZoneTag=+03:30,countryLongTag=Iran__Islamic_Republic_of,countryShortTag=IR,host=10.5.3.5,onlineTag=2 avgRtt=0,city="Arak",country_long="Iran, Islamic Republic of",country_short="IR",latitude=34.095402,longitude=49.701302,maxRtt=0,minRtt=0,online=2,packetLost=100,packetRecv=0,packetSent=3,state="Markazi",timezone="+03:30" 1570345087000000000
> ip2loc,CityTag=Isfahan,StateTag=Isfahan,TimeZoneTag=+03:30,countryLongTag=Iran__Islamic_Republic_of,countryShortTag=IR,host=172.16.2.3,onlineTag=2 avgRtt=0,city="Isfahan",country_long="Iran, Islamic Republic of",country_short="IR",latitude=32.6539,longitude=51.666,maxRtt=0,minRtt=0,online=2,packetLost=100,packetRecv=0,packetSent=3,state="Isfahan",timezone="+03:30" 1570345087000000000
> ip2loc,CityTag=Arak,StateTag=Markazi,TimeZoneTag=+03:30,countryLongTag=Iran__Islamic_Republic_of,countryShortTag=IR,host=10.0.2.3,onlineTag=2 avgRtt=0,city="Arak",country_long="Iran, Islamic Republic of",country_short="IR",latitude=34.095402,longitude=49.701302,maxRtt=0,minRtt=0,online=2,packetLost=100,packetRecv=0,packetSent=3,state="Markazi",timezone="+03:30" 1570345087000000000
> ip2loc,CityTag=Tehran,StateTag=Tehran,TimeZoneTag=+03:30,countryLongTag=Iran__Islamic_Republic_of,countryShortTag=IR,host=172.16.50.144,onlineTag=1 avgRtt=0.218722,city="Tehran",country_long="Iran, Islamic Republic of",country_short="IR",latitude=35.689201,longitude=51.389,maxRtt=0.287862,minRtt=0.089485,online=1,packetLost=0,packetRecv=3,packetSent=3,state="Tehran",timezone="+03:30" 1570345087000000000
> ip2loc,CityTag=Dallas,StateTag=Texas,TimeZoneTag=-05:00,countryLongTag=United_States,countryShortTag=US,host=4.2.2.4,onlineTag=1 avgRtt=83.457814,city="Dallas",country_long="United States",country_short="US",latitude=32.783058,longitude=-96.806671,maxRtt=88.532515,minRtt=80.184843,online=1,packetLost=0,packetRecv=3,packetSent=3,state="Texas",timezone="-05:00" 1570345087000000000
> ip2loc,CityTag=Tehran,StateTag=Tehran,TimeZoneTag=+03:30,countryLongTag=Iran__Islamic_Republic_of,countryShortTag=IR,host=172.16.50.143,onlineTag=1 avgRtt=0.817406,city="Tehran",country_long="Iran, Islamic Republic of",country_short="IR",latitude=35.689201,longitude=51.389,maxRtt=1.214626,minRtt=0.578192,online=1,packetLost=0,packetRecv=3,packetSent=3,state="Tehran",timezone="+03:30" 1570345087000000000
```

As you can see, Your public IP (8.8.8.8) location information fetched from the IP2Location Lite DB and Private IP Adderesses (172.16.50.16, 10.0.0.6, 10.5.3.5, 172.16.2.3) fetched from your local db (local.csv).
The online=1 or online=2 means that provided IP address is alive(1) or unreachable(2). The other metrics are clear but if you had any questions, I would be happy to answer ;-)
Start using IP2Location Telegraf Input

To start using it, You need to restart Telegraf service.
```
sudo service telegraf restart
```
