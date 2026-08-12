package dict

// Avatar 是一个可选头像。
type Avatar struct {
	ID   int64
	Name string
	URL  string
}

// Country 是一个国家。
//
// 首都坐标在寄往境外时作为投递距离的计算依据——v1 对境外地址不记录精确坐标，
// 一律按首都位置估算。
type Country struct {
	ID               int64
	Name             string
	NameEnglish      string
	CapitalName      string
	CapitalLongitude float64
	CapitalLatitude  float64
}
