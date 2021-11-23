package REST_API

type Model struct {
	ID int `json:"id"`
}

type Book struct {
	Model
	Name   string  `json:"name" binding:"min=1,max=100" gorm:"unique"`
	Price  float64 `json:"price" binding:"min=0"`
	Genre  int     `json:"genre" binding:"min=1,max=3"`
	Amount int     `json:"amount" binding:"min=0"`
}

type Genre struct {
	Model
	Name string
}
