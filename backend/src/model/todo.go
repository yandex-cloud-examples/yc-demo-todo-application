package model

type Todo struct {
    Id             int64     `json:"id" gorm:"autoIncrement;primaryKey"`
    Title          string    `json:"title"`
    Description    string    `json:"description"`
    Completed      bool      `json:"completed"`
}

