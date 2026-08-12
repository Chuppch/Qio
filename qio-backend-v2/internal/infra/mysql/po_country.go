package mysql

import "github.com/Chuppch/Qio/qio-backend-v2/internal/domain/dict"

// 字典表的持久化对象。
//
// avatar 与 country 都是全站共用的静态数据，对应 internal/domain/dict。

// avatarPO 映射 avatar 表，无审计字段。
type avatarPO struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`
	URL  string `gorm:"column:url"`
}

func (avatarPO) TableName() string { return "avatar" }

func (p avatarPO) toDomain() dict.Avatar {
	return dict.Avatar{ID: p.ID, Name: p.Name, URL: p.URL}
}

// countryPO 映射 country 表。
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

func (p countryPO) toDomain() dict.Country {
	return dict.Country{
		ID:               p.ID,
		Name:             p.CountryName,
		NameEnglish:      p.CountryNameEnglish,
		CapitalName:      p.CapitalName,
		CapitalLongitude: p.CapitalLongitude,
		CapitalLatitude:  p.CapitalLatitude,
	}
}
