package dto

type GenreResponse struct {
	Name string `json:"name"`
}

type GenreCreateInput struct {
	Name string `json:"name"`
}
