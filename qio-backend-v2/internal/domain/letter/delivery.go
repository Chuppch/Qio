package letter

import (
	"math"
	"time"
)

// 投递时长推算。
//
// 迁移自 v1 的 PositionUtil.getDistance 与散落在 LetterServiceImpl.sendLetterPre
// 里的两条业务规则（距离下限、投递速度）。这两条规则属于信件域的投递策略，
// 不是通用工具，因此收在本包内且不对外导出。

const (
	// earthRadiusMeters 地球平均半径，单位米
	earthRadiusMeters = 6371000.0

	// minDeliveryDistanceMeters 是投递距离下限。
	// 同城或近距离寄信也需要一定的等待时间，避免瞬间送达。
	minDeliveryDistanceMeters = 200000.0

	// deliverySpeedMetersPerHour 是信件的名义投递速度。
	deliverySpeedMetersPerHour = 40000.0
)

// Coordinate 是一个地理坐标。
type Coordinate struct {
	Longitude float64
	Latitude  float64
}

// EstimateDuration 返回从 from 寄往 to 的预计投递时长。
//
// 实际距离低于下限时按下限计算。
func EstimateDuration(from, to Coordinate) time.Duration {
	distance := math.Max(distanceMeters(from, to), minDeliveryDistanceMeters)
	hours := distance / deliverySpeedMetersPerHour
	return time.Duration(hours * float64(time.Hour))
}

// distanceMeters 用 Haversine 公式计算两点间的球面距离，单位米。
func distanceMeters(from, to Coordinate) float64 {
	radLat1 := radians(from.Latitude)
	radLat2 := radians(to.Latitude)
	deltaLat := radLat1 - radLat2
	deltaLng := radians(from.Longitude) - radians(to.Longitude)

	h := math.Pow(math.Sin(deltaLat/2), 2) +
		math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(deltaLng/2), 2)

	return 2 * math.Asin(math.Sqrt(h)) * earthRadiusMeters
}

func radians(degree float64) float64 {
	return degree * math.Pi / 180.0
}
