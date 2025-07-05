package database

type User struct {
	Discord_id   string `gorm:"primaryKey"`
	Osu_id       string
	Osu_username string
}

type Map struct {
	Channel_id string `gorm:"primaryKey"`
	Map_id     string
}
