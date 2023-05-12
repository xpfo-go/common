package gis

import "math"

// WGS84坐标系：即地球坐标系，国际上通用的坐标系。
// GCJ02坐标系：即火星坐标系，WGS84坐标系经加密后的坐标系。Google Maps，高德在用。
// BD09坐标系：即百度坐标系，GCJ02坐标系经加密后的坐标系。

const (
	XPi    = math.Pi * 3000.0 / 180.0
	OFFSET = 0.00669342162296594323
	AXIS   = 6378245.0
)

// BD09toGCJ02 百度坐标系->火星坐标系
func BD09toGCJ02(lng, lat float64) (float64, float64) {
	x := lng - 0.0065
	y := lat - 0.006

	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*XPi)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*XPi)

	gLon := z * math.Cos(theta)
	gLat := z * math.Sin(theta)

	return gLon, gLat
}

// GCJ02toBD09 火星坐标系->百度坐标系
func GCJ02toBD09(lng, lat float64) (float64, float64) {
	z := math.Sqrt(lng*lng+lat*lat) + 0.00002*math.Sin(lat*XPi)
	theta := math.Atan2(lat, lng) + 0.000003*math.Cos(lng*XPi)

	bdLon := z*math.Cos(theta) + 0.0065
	bdLat := z*math.Sin(theta) + 0.006

	return bdLon, bdLat
}

// WGS84toGCJ02 WGS84坐标系->火星坐标系
func WGS84toGCJ02(lng, lat float64) (float64, float64) {
	if isOutOFChina(lng, lat) {
		return lng, lat
	}

	mgLon, mgLat := delta(lng, lat)

	return mgLon, mgLat
}

// GCJ02toWGS84 火星坐标系->WGS84坐标系
func GCJ02toWGS84(lng, lat float64) (float64, float64) {
	if isOutOFChina(lng, lat) {
		return lng, lat
	}

	mgLon, mgLat := delta(lng, lat)

	return lng*2 - mgLon, lat*2 - mgLat
}

// BD09toWGS84 百度坐标系->WGS84坐标系
func BD09toWGS84(lng, lat float64) (float64, float64) {
	lng, lat = BD09toGCJ02(lng, lat)
	return GCJ02toWGS84(lng, lat)
}

// WGS84toBD09 WGS84坐标系->百度坐标系
func WGS84toBD09(lng, lat float64) (float64, float64) {
	lng, lat = WGS84toGCJ02(lng, lat)
	return GCJ02toBD09(lng, lat)
}

func delta(lng, lat float64) (float64, float64) {
	dLat := transformLat(lng-105.0, lat-35.0)
	dLon := transformLng(lng-105.0, lat-35.0)

	radLat := lat / 180.0 * math.Pi
	magic := math.Sin(radLat)
	magic = 1 - OFFSET*magic*magic
	sqrtMagic := math.Sqrt(magic)

	dLat = (dLat * 180.0) / ((AXIS * (1 - OFFSET)) / (magic * sqrtMagic) * math.Pi)
	dLon = (dLon * 180.0) / (AXIS / sqrtMagic * math.Cos(radLat) * math.Pi)

	mgLat := lat + dLat
	mgLon := lng + dLon

	return mgLon, mgLat
}

func transformLat(lng, lat float64) float64 {
	var ret = -100.0 + 2.0*lng + 3.0*lat + 0.2*lat*lat + 0.1*lng*lat + 0.2*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*math.Pi) + 20.0*math.Sin(2.0*lng*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*math.Pi) + 40.0*math.Sin(lat/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*math.Pi) + 320*math.Sin(lat*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

func transformLng(lng, lat float64) float64 {
	var ret = 300.0 + lng + 2.0*lat + 0.1*lng*lng + 0.1*lng*lat + 0.1*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*math.Pi) + 20.0*math.Sin(2.0*lng*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lng*math.Pi) + 40.0*math.Sin(lng/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lng/12.0*math.Pi) + 300.0*math.Sin(lng/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}

func isOutOFChina(lng, lat float64) bool {
	return !(lng > 73.66 && lng < 135.05 && lat > 3.86 && lat < 53.55)
}
