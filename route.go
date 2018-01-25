package memap

import (
	"time"
)

// TimeCoord 是时序坐标
type TimeCoord struct {
	Coord
	Time time.Time `json:"time"`
}

//  常量
const (
	DefaultMinRoutePointsNum = 3    // 默认清理路线的最少点数
	DefaultMinSpeed          = 3.0  // m/s
	DefaultMaxSpeed          = 50.0 // m/s
)

var (
	minRoutePointsNum = DefaultMinRoutePointsNum
	minSpeed          = DefaultMinSpeed
	maxSpeed          = DefaultMaxSpeed
)

// SetMinRoutePointsNum 设置路线的最少点数，至少为3
func SetMinRoutePointsNum(n int) {
	if n > DefaultMinRoutePointsNum {
		minRoutePointsNum = n
	}
}

// SetRouteSpeedLimit 设置路线速度限制
func SetRouteSpeedLimit(min, max float64) {
	if min > max || min <= 0 || max <= 0 {
		return
	}
	minSpeed = min
	maxSpeed = max
}

// CleanTimeRoute 对轨迹作出一定的处理
// 清理一些轨迹坐标
func CleanTimeRoute(route []TimeCoord, fn func(c1, c2 *Coord) float64) []TimeCoord {
	n := len(route)
	if n < minRoutePointsNum {
		return route
	}
	ret := []TimeCoord{route[0]}

	for i := 1; i < n-1; i++ {
		next := i + 1
		// 去重点
		interval := route[next].Time.Sub(route[i].Time).Seconds()
		if interval == 0 ||
			(route[next].Latitude == route[i].Latitude && route[next].Longitude == route[i].Longitude) {
			continue
		}
		// 速度验证
		dist := fn(&route[next].Coord, &route[i].Coord)
		speed := dist / interval
		if speed < minSpeed || speed > minSpeed {
			continue
		}
		ret = append(ret, route[i])
	}
	ret = append(ret, route[n-1])
	return ret
}
