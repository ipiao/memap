package memap

// PNPoly算法

// Polygon 多边形
type Polygon struct {
	points []Coord
	box
}

type box struct {
	minLon float64
	maxLon float64

	minLat float64
	maxLat float64
}

// NewPolygon 创建一个多边形
func NewPolygon(points []Coord) *Polygon {
	polygon := &Polygon{points: points}
	box := box{}
	for _, p := range points {
		if p.Longitude > box.maxLon {
			box.maxLon = p.Longitude
		}
		if p.Longitude < box.minLon {
			box.minLon = p.Longitude
		}
		if p.Latitude > box.maxLat {
			box.maxLat = p.Latitude
		}
		if p.Latitude < box.minLat {
			box.minLat = p.Latitude
		}
	}
	polygon.box = box
	return polygon
}

// IsInside 是否在多边形内
// PNPoly算法
func (p *Polygon) IsInside(c *Coord) bool {
	if c.Longitude > p.box.maxLon || c.Longitude < p.box.minLat ||
		c.Latitude > p.box.maxLat || c.Latitude < p.box.minLat {
		return false
	}
	isIn := false
	n := len(p.points)

	for i, j := 0, n-1; i < n; i++ {
		if (p.points[i].Latitude < c.Latitude && p.points[j].Latitude >= c.Latitude ||
			p.points[j].Latitude < c.Latitude && p.points[i].Latitude >= c.Latitude) &&
			(p.points[i].Longitude <= c.Longitude || p.points[j].Longitude <= c.Longitude) { // 从待测点引出一条水平向左的射线
			// 有交点,必定应 : 交点.Longitude<c.Longitude
			if p.points[i].Longitude+(c.Latitude-p.points[i].Latitude)/(p.points[j].Latitude-p.points[i].Latitude)*
				(p.points[j].Longitude-p.points[i].Longitude) < c.Longitude {
				isIn = !isIn
			}
		}
		j = i
	}
	return isIn
}
