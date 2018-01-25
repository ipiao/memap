package memap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试坐标转换,差距在0.5m之内
func TestCoord(t *testing.T) {
	var origins = [][2]float64{
		{127.0, 35.0},        // det~2
		{127.2, 35.4},        // det~1
		{127.23, 35.43},      // det~0.5
		{127.234, 35.433},    // det~0.4
		{127.2342, 35.4333},  // det~0.4
		{126.0, 33.0},        // det~2
		{126.2, 33.4},        // det~1.5
		{126.23, 33.43},      // det~0.8
		{126.234, 33.433},    // det~0.7
		{126.2342, 33.43335}, // det~0.7
		{136.0, 37.0},        // det~2.7
		{136.2, 37.4},        // det~1.7
		{136.23, 37.43},      // det~0.9
		{136.234, 37.433},    // det~0.8
		{136.2342, 37.43335}, // det~0.8
	}
	for i, origin := range origins {
		lon, lat := GCJ02ToWGS84(origin[0], origin[1])
		t.Log(lon, lat)
		rlon, rlat := WGS84ToGCJ02(lon, lat)
		t.Log(rlon, rlat)
		dist := DistanceBetweenGCJ02Coords(origin[0], origin[1], rlon, rlat)
		t.Logf("第%d组转换后差距:%f\n", i+1, dist)

		lon1, lat1 := WGS84ToGCJ02(origin[0], origin[1])
		t.Log(lon1, lat1)
		rlon1, rlat1 := GCJ02ToWGS84(lon1, lat1)
		t.Log(rlon1, rlat1)
		dist1 := DistanceBetweenGCJ02Coords(origin[0], origin[1], rlon1, rlat1)
		t.Logf("第%d组转换后差距:%f\n", i+1, dist1)
	}
}

func TestDist(t *testing.T) {
	var origins = [][][2]float64{
		{{127.234, 35.433}, {127.434, 35.345}},
		{{123.232, 35.433}, {122.434, 34.345}},
		{{126.234, 37.433}, {124.434, 37.345}},
		{{122.274, 33.733}, {135.434, 36.445}},

		{{127.2342, 35.4333}, {127.4343, 35.3453}},
		{{123.232, 35.4333}, {122.4344, 34.3452}},
		{{126.2345, 37.4336}, {124.4346, 37.3455}},
		{{122.2747, 33.7337}, {135.4345, 36.4455}},
	}
	for i, origin := range origins {
		lon1, lat1 := GCJ02ToWGS84(origin[0][0], origin[0][1])
		lon2, lat2 := GCJ02ToWGS84(origin[1][0], origin[1][1])
		dist1 := DistanceBetweenGCJ02Coords(origin[0][0], origin[0][1], origin[1][0], origin[1][1])
		dist2 := DistanceBetweenWGS84Coords(lon1, lat1, lon2, lat2)
		t.Logf("GCJ02ToWGS84 第%d组转换后距离:%f-%f=%f\n,差距比率:%f", i+1, dist1, dist2, dist1-dist2, (dist1-dist2)/dist1)

		lon3, lat3 := WGS84ToGCJ02(origin[0][0], origin[0][1])
		lon4, lat4 := WGS84ToGCJ02(origin[1][0], origin[1][1])
		dist3 := DistanceBetweenWGS84Coords(origin[0][0], origin[0][1], origin[1][0], origin[1][1])
		dist4 := DistanceBetweenGCJ02Coords(lon3, lat3, lon4, lat4)
		t.Logf("WGS84ToGCJ02 第%d组转换后距离:%f-%f=%f\n,差距比率:%f%%", i+1, dist3, dist4, dist3-dist4, (dist3-dist4)/dist3*100)
	}
}

func TestWGS84DistWithRedis(t *testing.T) {
	lon1, lat1 := 15.08726745843887329, 37.50266842333162032
	lon2, lat2 := 15.23300260305404663, 67.32299876601062749
	dist := DistanceBetweenWGS84Coords(lon1, lat1, lon2, lat2)
	assert.Equal(t, dist-3316817.6881 < 0.001, true)
}
