# Ltser

This module is part of my **goex project** containing several exercises and experimental code that I wrote while learning Golang.

Main purposes are:
- Gathering data from http://lter.eurac.edu sensor network.
- Creating a server capable of receiving those data and store them in an appropriate data storage.
- Make some hands-on data analysis.

Following some assumptions:
- Data will be downloaded in **.CSV file** format from https://browser.lter.eurac.edu/de.
- A tool developed in Go (**pusher**) will load data rows from the .CSV file, parse them, convert them to JSON and post data to the server.
- A server developed in Go (**ingestor**) will expose a REST API to receive the raw data JSON. After some validation it will store the data in the database.
- The database will be an **InfluxDB** instance (InfluxDB Cloud at first).
- In a second phase, everything will be containerized with docker (so, also InfluxDB Cloud will be replaced by a dockerized version).
