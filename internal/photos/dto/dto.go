package dto

type CreatePhotoDto struct {
	Title  string `json:"title"`
	Base64 string `json:"base64"`
}
