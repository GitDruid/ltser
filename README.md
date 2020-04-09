Ltser

This module is part of my goex project containing several exercises and experimental code that I wrote while learning Golang.

Briefly its purposes are:
- Gathering data from http://lter.eurac.edu sensor network.
- Creating a server capable of receiving those data and store them in an appropriate data storage.
- Make some hands-on data analysis.

Following some choises taken:
- Data will be downloaded in .CSV file format from https://browser.lter.eurac.edu/de.
- A tool developed in Go (pusher) will load from .CSV, parse, convert in JSON and post data to the server.
- A server developed in Go will expose a REST API to receive the raw data. After some validation it will store the data.
- The database will be an InfluxDB instance.
