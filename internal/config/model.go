package config

type Config struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"not null"`
	Content string `gorm:"not null"`
}
