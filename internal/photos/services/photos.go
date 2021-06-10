package services

import (
	"github.com/jmoiron/sqlx"
	. "github.com/morzhanov/go-realworld/internal/photos/dto"
	. "github.com/morzhanov/go-realworld/internal/photos/models"
)

type PhotosService struct {
	db *sqlx.DB
}

func (s *PhotosService) GetUserPhotos(userId string) (photos []*Photo, err error) {
	q := `SELECT * FROM photos
		WHERE photos.user_id = $1`
	rows, err := s.db.Query(q, userId)
	if err != nil {
		return nil, err
	}

	photos = make([]*Photo, 0)
	for rows.Next() {
		photo := Photo{}
		err = rows.Scan(photo.ID, photo.Title, photo.Base64, photo.UserId)
		if err != nil {
			return nil, err
		}
		photos = append(photos, &photo)
	}

	return photos, nil
}

func (s *PhotosService) GetUserPhoto(userId string, photoId string) (res *Photo, err error) {
	q := `SELECT * FROM photos
		WHERE id = $1 AND photos.user_id = $2`
	row := s.db.QueryRow(q, photoId, userId)

	res = &Photo{}
	err = row.Scan(res.ID, res.Title, res.Base64, res.UserId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *PhotosService) CreateUserPhoto(userId string, data *CreatePhotoDto) error {
	q := `INSERT INTO photos (title, base64, user_id)
		VALUES ($1, $2, $3)`
	_, err := s.db.Query(q, data.Title, data.Base64, userId)
	return err
}

func (s *PhotosService) DeleteUserPhoto(userId string, photoId string) error {
	q := `DELETE FROM photos
		WHERE id = $1 AND user_id = $2`
	_, err := s.db.Query(q, photoId, userId)
	return err
}

func NewPhotosService(db *sqlx.DB) *PhotosService {
	return &PhotosService{db}
}
