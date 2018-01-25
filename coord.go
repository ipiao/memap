package memap

// redis 是采用的WGS84无误
// 坐标系之间不采用相同的半径，转换后的坐标算距离,误差率在 0.1%,采用相同的半径，误差率在 0.01%

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// 常量
const (
	GCJ02RADIUS = 6378137.0 //  高德坐标算法半径(单位：米)
	WGS84RADIUS = 6372797.560856

	EE = 0.00669342162296594323 // 椭球的偏心率。
	A  = 6378245.0              //  a: 卫星椭球坐标投影到平面地图坐标系的投影因子。
	// PI = 3.1415926535897932384626;
)

// Coord 是基础坐标
type Coord struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (c Coord) String() string {
	return fmt.Sprintf("%f,%f", c.Longitude, c.Latitude)
}

// StringParseToCoord 将字符串解析成坐标
func StringParseToCoord(c string) (*Coord, error) {
	s := strings.Split(c, ",")
	if len(s) != 2 {
		return nil, fmt.Errorf("Error value of coord '%s',format must be 'longitude,latitude'", c)
	}
	lon, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return nil, fmt.Errorf("Error value of longitude '%s'", s[0])
	}
	lat, err := strconv.ParseFloat(s[1], 64)
	if err != nil {
		return nil, fmt.Errorf("Error value of latitude '%s'", s[1])
	}
	return &Coord{Longitude: lon, Latitude: lat}, nil
}

// GeoHash 算法、墨卡托投影
// 高德\国测局：GCJ02坐标系(火星坐标系)
// Redis: WGS84坐标系(地心坐标系)
// 百度坐标:（BD09）

// GCJ02DistanceBetweenCoords 计算火星坐标系的两点间距离（单位：米）
func GCJ02DistanceBetweenCoords(c1, c2 *Coord) float64 {
	return DistanceBetweenGCJ02Coords(c1.Longitude, c1.Latitude, c2.Longitude, c2.Latitude)
}

// DistanceBetweenGCJ02Coords 计算高德坐标的两点间距离（单位：米）
func DistanceBetweenGCJ02Coords(lon1, lat1, lon2, lat2 float64) float64 {
	rad := math.Pi / 180.0
	dlat1 := lat1 * rad
	dlng1 := lon1 * rad
	dlat2 := lat2 * rad
	dlng2 := lon2 * rad
	theta := dlng2 - dlng1
	dist := math.Acos(math.Sin(dlat1)*math.Sin(dlat2) + math.Cos(dlat1)*math.Cos(dlat2)*math.Cos(theta))
	return dist * GCJ02RADIUS
}

// WGS84DistanceBetweenCoords 计算WGS84坐标系的两点间距离（单位：米）
func WGS84DistanceBetweenCoords(c1, c2 *Coord) float64 {
	return DistanceBetweenWGS84Coords(c1.Longitude, c1.Latitude, c2.Longitude, c2.Latitude)
}

// DistanceBetweenWGS84Coords 计算Redis坐标的两点间距离（单位：米）
func DistanceBetweenWGS84Coords(lon1, lat1, lon2, lat2 float64) float64 {
	rad := math.Pi / 180.0
	dlat1 := lat1 * rad
	dlng1 := lon1 * rad
	dlat2 := lat2 * rad
	dlng2 := lon2 * rad
	u := math.Sin((dlat2 - dlat1) / 2)
	v := math.Sin((dlng2 - dlng1) / 2)
	return 2.0 * WGS84RADIUS * math.Asin(math.Sqrt(u*u+math.Cos(dlat1)*math.Cos(dlat2)*v*v))
}

// GCJ02CoordToWGS84Coord 坐标转换
// 火星坐标系 转 地心坐标系
func GCJ02CoordToWGS84Coord(c *Coord) *Coord {
	lon, lat := GCJ02ToWGS84(c.Longitude, c.Latitude)
	return &Coord{Longitude: lon, Latitude: lat}
}

// GCJ02ToWGS84 坐标转换
// 高德坐标转Redis坐标
func GCJ02ToWGS84(lon, lat float64) (float64, float64) {
	mgLon, mgLat := WGS84ToGCJ02(lon, lat)
	return lon*2 - mgLon, lat*2 - mgLat
}

// WGS84CoordToGCJ02Coord 坐标转换
// 地心坐标系 转 火星坐标系
func WGS84CoordToGCJ02Coord(c *Coord) *Coord {
	lon, lat := WGS84ToGCJ02(c.Longitude, c.Latitude)
	return &Coord{Longitude: lon, Latitude: lat}
}

// WGS84ToGCJ02 坐标转换
// Redis坐标转高德坐标
func WGS84ToGCJ02(lon, lat float64) (float64, float64) {
	dLat := transformLat(lon-105.0, lat-35.0)
	dLon := transformLon(lon-105.0, lat-35.0)
	radLat := lat / 180.0 * math.Pi
	magic := math.Sin(radLat)
	magic = 1 - EE*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((A * (1 - EE)) / (magic * sqrtMagic) * math.Pi)
	dLon = (dLon * 180.0) / (A / sqrtMagic * math.Cos(radLat) * math.Pi)
	mgLat := lat + dLat
	mgLon := lon + dLon
	return mgLon, mgLat
}

func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*math.Pi) + 40.0*math.Sin(y/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*math.Pi) + 320*math.Sin(y*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

func transformLon(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*math.Pi) + 40.0*math.Sin(x/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*math.Pi) + 300.0*math.Sin(x/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}

// CoordOutOfChina 判断坐标是否超出中国
func CoordOutOfChina(c *Coord) bool {
	return OutOfChina(c.Longitude, c.Latitude)
}

// OutOfChina 坐标是否超出中国
func OutOfChina(lon, lat float64) bool {
	if lon < 72.004 || lon > 137.8347 {
		return true
	}
	if lat < 0.8293 || lat > 55.8271 {
		return true
	}
	return false
}

// WGS84RouteDistance WGS84路线距离
func WGS84RouteDistance(route []Coord) float64 {
	return routedistance(route, WGS84DistanceBetweenCoords)
}

// GCJ02RouteDistance WGS84路线距离
func GCJ02RouteDistance(route []Coord) float64 {
	return routedistance(route, GCJ02DistanceBetweenCoords)
}

// 计算路线距离，分段计算
func routedistance(route []Coord, fn func(c1, c2 *Coord) float64) float64 {
	var dist float64
	if len(route) < 2 {
		return 0
	}
	for i := 1; i < len(route); i++ {
		dist += fn(&route[i-1], &route[i])
	}
	return dist
}

// OutOfFence 检验坐标是否在围栏之外
func OutOfFence(c *Coord, fence []Coord) bool {
	if len(fence) < 3 {
		return false
	}
	polygon := NewPolygon(fence)
	return !polygon.IsInside(c)
}

// ExtractNCoordsFromRoute 从路线中提取特定数目的点，包含首、末两点
// 采用均分法 n=1 区 1/2,n=2 取 1/3,2/3 ...
func ExtractNCoordsFromRoute(route []Coord, n int) []Coord {
	if n <= 2 {
		return ExtractNInnerCoordsFromRoute(route, n)
	}
	length := len(route)
	ret := []Coord{route[0]}
	ret = append(ret, ExtractNInnerCoordsFromRoute(route[1:length-1], n-2)...)
	ret = append(ret, route[length-1])
	return ret
}

// ExtractNInnerCoordsFromRoute 从路线向内提取特定数目的点
// 采用均分法 n=1 区 1/2,n=2 取 1/3,2/3 ...
func ExtractNInnerCoordsFromRoute(route []Coord, n int) []Coord {
	length := len(route)
	if length <= n {
		return route
	}
	ret := []Coord{}
	if n >= 1 {
		for i := 0; i < n; i++ {
			ind := float64(i) * float64(length) / float64(n+1)
			ceilInd := math.Ceil(ind)
			if ind-ceilInd > 0.5 {
				ret = append(ret, route[int(ceilInd+1)])
			} else {
				ret = append(ret, route[int(ceilInd)])
			}
		}
	}
	return ret
}
