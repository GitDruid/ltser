// Package models provide common data structures for matschmazia tools.
package models // import "goex/ltser/matschmazia/models"

// RawData contains the raw data coming from the sensors (all in string format).
// More information on: https://browser.lter.eurac.edu/p/info.md
type RawData struct {
	Time              string `json:"time"`              // Date/time of measurement (UTC +1).
	Station           string `json:"station"`           // Station code.
	Landuse           string `json:"landuse"`           // me = meadows, pa = pasture, bs = bare soil, fo = forest
	Altitude          string `json:"altitude"`          // Altitude of the station in meters.
	Latitude          string `json:"latitude"`          // Latitude, coordinates in decimal degrees.
	Longitude         string `json:"longitude"`         // Longitude, coordinates in decimal degrees.
	AirRelHumidityAvg string `json:"air_rh_avg"`        // Relative humidity in percent (15 min average).
	AirTempAvg        string `json:"air_t_avg"`         // Air temperature in degree celsius (15 min average).
	NrUpSwAvg         string `json:"nr_up_sw_avg"`      // Undocumented.
	PrecipRtNrtTot    string `json:"precip_rt_nrt_tot"` // Precipitation in mm (15 min cumulative sum).
	SnowHeight        string `json:"snow_height"`       // Snow height in meter.
	SrAvg             string `json:"sr_avg"`            // Global solar radiation in Watt square meter (15 min average).
	WindDir           string `json:"wind_dir"`          // Wind direction in degrees (15 min average).
	WindSpeed         string `json:"wind_speed"`        // Undocumented.
	WindSpeedAvg      string `json:"wind_speed_avg"`    // Wind speed in m/s.
	WindSpeedMax      string `json:"wind_speed_max"`    // Wind gust in m/s.
}
