package vncalendar

import (
	"github.com/openhoangnc/vncalendar/lunar"
)

// SolarDate is a struct to store solar date
type SolarDate struct {
	Year, Month, Day int
}

// LunarDate is a struct to store lunar date
type LunarDate struct {
	Year, Month, Day int
	Leap             bool
}

// Solar2lunar convert solar date to lunar date
func Solar2lunar(yyyy, mm, dd, timeZoneOffset int) LunarDate {
	lunarDay, lunarMonth, lunarYear, lunarLeap := lunar.ConvertSolar2Lunar(
		float64(dd),
		float64(mm),
		float64(yyyy),
		float64(timeZoneOffset),
	)

	return LunarDate{
		Year:  int(lunarYear),
		Month: int(lunarMonth),
		Day:   int(lunarDay),
		Leap:  lunarLeap == 1,
	}
}

// Lunar2solar convert lunar date to solar date
func Lunar2solar(lunarYear, lunarMonth, lunarDay int, lunarLeap bool, timeZoneOffset int) SolarDate {
	dd, mm, yyyy := lunar.ConvertLunar2Solar(
		float64(lunarDay),
		float64(lunarMonth),
		float64(lunarYear),
		bool2float64(lunarLeap),
		float64(timeZoneOffset),
	)

	return SolarDate{
		Year:  int(yyyy),
		Month: int(mm),
		Day:   int(dd),
	}
}

func bool2float64(lunarLeap bool) float64 {
	if lunarLeap {
		return 1
	}
	return 0
}
