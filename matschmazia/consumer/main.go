// Consumer is a REST service that waits for sensors data as flat JSON and store them in an InfluxDB database.
package main

// More information on: https://browser.lter.eurac.edu/p/info.md
const (
	time              = iota // Date/time of measurement (UTC +1).
	station                  // Station code.
	landuse                  // me = meadows, pa = pasture, bs = bare soil, fo = forest
	altitude                 // Altitude of the station in meters.
	latitude                 // Latitude, coordinates in decimal degrees.
	longitude                // Longitude, coordinates in decimal degrees.
	airRelHumidityAvg        // Relative humidity in percent (15 min average).
	airTempAvg               // Air temperature in degree celsius (15 min average).
	nrUpSwAvg                // Undocumented.
	precipRtNrtTot           // Precipitation in mm (15 min cumulative sum).
	snowHeight               // Snow height in meter.
	srAvg                    // Global solar radiation in Watt square meter (15 min average).
	windDir                  // Wind direction in degrees (15 min average).
	windSpeed                // Undocumented.
	windSpeedAvg             // Wind speed in m/s.
	windSpeedMax             // Wind gust in m/s.
)

func main() {
}
