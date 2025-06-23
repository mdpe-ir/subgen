package userlink

type UUID struct {
	ID   uint   `gorm:"primaryKey"`
	UUID string `gorm:"uniqueIndex;not null"`
}
