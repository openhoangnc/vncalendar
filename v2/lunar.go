package lunar

// Ported from https://www.informatik.uni-leipzig.de/~duc/amlich/amlich-aa98.js

import (
	"math"
)

/* Discard the fractional part of a number, e.g., Floor(3.2) = 3 */
func Floor(d float64) float64 {
	return math.Floor(d)
}

/* Compute the (integral) Julian day number of day dd/mm/yyyy, i.e., the number
 * of days between 1/1/4713 BC (Julian calendar) and dd/mm/yyyy.
 * Formula from http://www.tondering.dk/claus/calendar.html
 */
func jdFromDate(dd, mm, yy float64) float64 {
	var a, y, m, jd float64
	a = Floor((14 - mm) / 12)
	y = yy + 4800 - a
	m = mm + 12*a - 3
	jd = dd + Floor((153*m+2)/5) + 365*y + Floor(y/4) - Floor(y/100) + Floor(y/400) - 32045
	if jd < 2299161 {
		jd = dd + Floor((153*m+2)/5) + 365*y + Floor(y/4) - 32083
	}
	return jd
}

/* Convert a Julian day number to day/month/year. Parameter jd is an integer */
func jdToDate(jd float64) (float64, float64, float64) {
	var a, b, c, d, e, m, day, month, year float64
	if jd > 2299160 { // After 5/10/1582, Gregorian calendar
		a = jd + 32044
		b = Floor((4*a + 3) / 146097)
		c = a - Floor((b*146097)/4)
	} else {
		b = 0
		c = jd + 32082
	}
	d = Floor((4*c + 3) / 1461)
	e = c - Floor((1461*d)/4)
	m = Floor((5*e + 2) / 153)
	day = e - Floor((153*m+2)/5) + 1
	month = m + 3 - 12*Floor(m/10)
	year = b*100 + d - 4800 + Floor(m/10)
	return day, month, year
}

/* Compute the time of the k-th new moon after the new moon of 1/1/1900 13:52 UCT
 * (measured as the number of days since 1/1/4713 BC noon UCT, e.g., 2451545.125 is 1/1/2000 15:00 UTC).
 * Returns a floating number, e.g., 2415079.9758617813 for k=2 or 2414961.935157746 for k=-2
 * Algorithm from: "Astronomical Algorithms" by Jean Meeus, 1998
 */
func NewMoon(k float64) float64 {
	var T, T2, T3, dr, Jd1, M, Mpr, F, C1, deltat, JdNew float64
	T = k / 1236.85 // Time in Julian centuries from 1900 January 0.5
	T2 = T * T
	T3 = T2 * T
	dr = math.Pi / 180
	Jd1 = 2415020.75933 + 29.53058868*k + 0.0001178*T2 - 0.000000155*T3
	Jd1 = Jd1 + 0.00033*math.Sin((166.56+132.87*T-0.009173*T2)*dr) // Mean new moon
	M = 359.2242 + 29.10535608*k - 0.0000333*T2 - 0.00000347*T3    // Sun's mean anomaly
	Mpr = 306.0253 + 385.81691806*k + 0.0107306*T2 + 0.00001236*T3 // Moon's mean anomaly
	F = 21.2964 + 390.67050646*k - 0.0016528*T2 - 0.00000239*T3    // Moon's argument of latitude
	C1 = (0.1734-0.000393*T)*math.Sin(M*dr) + 0.0021*math.Sin(2*dr*M)
	C1 = C1 - 0.4068*math.Sin(Mpr*dr) + 0.0161*math.Sin(dr*2*Mpr)
	C1 = C1 - 0.0004*math.Sin(dr*3*Mpr)
	C1 = C1 + 0.0104*math.Sin(dr*2*F) - 0.0051*math.Sin(dr*(M+Mpr))
	C1 = C1 - 0.0074*math.Sin(dr*(M-Mpr)) + 0.0004*math.Sin(dr*(2*F+M))
	C1 = C1 - 0.0004*math.Sin(dr*(2*F-M)) - 0.0006*math.Sin(dr*(2*F+Mpr))
	C1 = C1 + 0.0010*math.Sin(dr*(2*F-Mpr)) + 0.0005*math.Sin(dr*(2*Mpr+M))
	if T < -11 {
		deltat = 0.001 + 0.000839*T + 0.0002261*T2 - 0.00000845*T3 - 0.000000081*T*T3
	} else {
		deltat = -0.000278 + 0.000265*T + 0.000262*T2
	}
	JdNew = Jd1 + C1 - deltat
	return JdNew
}

/* Compute the longitude of the sun at any time.
 * Parameter: floating number jdn, the number of days since 1/1/4713 BC noon
 * Algorithm from: "Astronomical Algorithms" by Jean Meeus, 1998
 */
func SunLongitude(jdn float64) float64 {
	var T, T2, dr, M, L0, DL, L float64
	T = (jdn - 2451545.0) / 36525 // Time in Julian centuries from 2000-01-01 12:00:00 GMT
	T2 = T * T
	dr = math.Pi / 180                                             // degree to radian
	M = 357.52910 + 35999.05030*T - 0.0001559*T2 - 0.00000048*T*T2 // mean anomaly, degree
	L0 = 280.46645 + 36000.76983*T + 0.0003032*T2                  // mean longitude, degree
	DL = (1.914600 - 0.004817*T - 0.000014*T2) * math.Sin(dr*M)
	DL = DL + (0.019993-0.000101*T)*math.Sin(dr*2*M) + 0.000290*math.Sin(dr*3*M)
	L = L0 + DL // true longitude, degree
	L = L * dr
	L = L - math.Pi*2*(Floor(L/(math.Pi*2))) // Normalize to (0, 2*math.Pi)
	return L
}

/* Compute sun position at midnight of the day with the given Julian day number.
 * The time zone if the time difference between local time and UTC: 7.0 for UTC+7:00.
 * The func returns a number between 0 and 11.
 * From the day after March equinox and the 1st major term after March equinox, 0 is returned.
 * After that, return 1, 2, 3 ...
 */
func getSunLongitude(dayNumber, timeZone float64) float64 {
	return Floor(SunLongitude(dayNumber-0.5-timeZone/24) / math.Pi * 6)
}

/* Compute the day of the k-th new moon in the given time zone.
 * The time zone if the time difference between local time and UTC: 7.0 for UTC+7:00
 */
func getNewMoonDay(k, timeZone float64) float64 {
	return Floor(NewMoon(k) + 0.5 + timeZone/24)
}

/* Find the day that starts the luner month 11 of the given year for the given time zone */
func getLunarMonth11(yy, timeZone float64) float64 {
	var k, off, nm, sunLong float64
	//off = jdFromDate(31, 12, yy) - 2415021.076998695;
	off = jdFromDate(31, 12, yy) - 2415021
	k = Floor(off / 29.530588853)
	nm = getNewMoonDay(k, timeZone)
	sunLong = getSunLongitude(nm, timeZone) // sun longitude at local midnight
	if sunLong >= 9 {
		nm = getNewMoonDay(k-1, timeZone)
	}
	return nm
}

/* Find the index of the leap month after the month starting on the day a11. */
func getLeapMonthOffset(a11, timeZone float64) float64 {
	var k, last, arc, i float64
	k = Floor((a11-2415021.076998695)/29.530588853 + 0.5)
	last = 0
	i = 1 // We start with the month following lunar month 11
	arc = getSunLongitude(getNewMoonDay(k+i, timeZone), timeZone)
	for {
		last = arc
		i = i + 1
		arc = getSunLongitude(getNewMoonDay(k+i, timeZone), timeZone)
		if !(arc != last && i < 14) {
			break
		}
	}
	return i - 1
}

/* Comvert solar date dd/mm/yyyy to the corresponding lunar date */
func convertSolar2Lunar(dd, mm, yy, timeZone float64) (float64, float64, float64, float64) {
	var k, dayNumber, monthStart, a11, b11, lunarDay, lunarMonth, lunarYear, lunarLeap float64
	dayNumber = jdFromDate(dd, mm, yy)
	k = Floor((dayNumber - 2415021.076998695) / 29.530588853)
	monthStart = getNewMoonDay(k+1, timeZone)
	if monthStart > dayNumber {
		monthStart = getNewMoonDay(k, timeZone)
	}
	//alert(dayNumber+" -> "+monthStart);
	a11 = getLunarMonth11(yy, timeZone)
	b11 = a11
	if a11 >= monthStart {
		lunarYear = yy
		a11 = getLunarMonth11(yy-1, timeZone)
	} else {
		lunarYear = yy + 1
		b11 = getLunarMonth11(yy+1, timeZone)
	}
	lunarDay = dayNumber - monthStart + 1
	var diff = Floor((monthStart - a11) / 29)
	lunarLeap = 0
	lunarMonth = diff + 11
	if b11-a11 > 365 {
		var leapMonthDiff = getLeapMonthOffset(a11, timeZone)
		if diff >= leapMonthDiff {
			lunarMonth = diff + 10
			if diff == leapMonthDiff {
				lunarLeap = 1
			}
		}
	}
	if lunarMonth > 12 {
		lunarMonth = lunarMonth - 12
	}
	if lunarMonth >= 11 && diff < 4 {
		lunarYear -= 1
	}
	return lunarDay, lunarMonth, lunarYear, lunarLeap
}

/* Convert a lunar date to the corresponding solar date */
func convertLunar2Solar(lunarDay, lunarMonth, lunarYear, lunarLeap, timeZone float64) (float64, float64, float64) {
	var k, a11, b11, off, leapOff, leapMonth, monthStart float64
	if lunarMonth < 11 {
		a11 = getLunarMonth11(lunarYear-1, timeZone)
		b11 = getLunarMonth11(lunarYear, timeZone)
	} else {
		a11 = getLunarMonth11(lunarYear, timeZone)
		b11 = getLunarMonth11(lunarYear+1, timeZone)
	}
	k = Floor(0.5 + (a11-2415021.076998695)/29.530588853)
	off = lunarMonth - 11
	if off < 0 {
		off += 12
	}
	if b11-a11 > 365 {
		leapOff = getLeapMonthOffset(a11, timeZone)
		leapMonth = leapOff - 2
		if leapMonth < 0 {
			leapMonth += 12
		}
		if lunarLeap != 0 && lunarMonth != leapMonth {
			return 0, 0, 0
		} else if lunarLeap != 0 || off >= leapOff {
			off += 1
		}
	}
	monthStart = getNewMoonDay(k+off, timeZone)
	return jdToDate(monthStart + lunarDay - 1)
}
