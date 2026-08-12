package mysql

// countryPO 映射 country 表。
//
// 地理字典表，被信件投递距离计算与用户地址共用，因此单独成文件而不归入某个域。
//
// 注意：country.id 在表中是 bigint unsigned，这里用 int64 承载。
// 国家数量远小于 int64 上限，实际不会溢出。
type countryPO struct {
	ID                 int64   `gorm:"column:id;primaryKey;autoIncrement"`
	CountryName        string  `gorm:"column:country_name"`
	CountryNameEnglish string  `gorm:"column:country_name_english"`
	CapitalName        string  `gorm:"column:capital_name"`
	CapitalLongitude   float64 `gorm:"column:capital_longitude"`
	CapitalLatitude    float64 `gorm:"column:capital_latitude"`
}

func (countryPO) TableName() string { return "country" }
