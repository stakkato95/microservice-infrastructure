package domain

type Tweet struct {
	// gorm.Model
	Id     int
	UserId int
	Text   string
}
