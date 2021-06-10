package services

import (
	"github.com/jmoiron/sqlx"
	. "github.com/morzhanov/go-realworld/internal/pictures/dto"
	. "github.com/morzhanov/go-realworld/internal/pictures/models"
)

type PictureService struct {
	db *sqlx.DB
}

func (s *PictureService) GetUserPictures(userId string) (pictures []*Picture, err error) {
	q := `SELECT * FROM pictures
		WHERE pictures.user_id = $1`
	rows, err := s.db.Query(q, userId)
	if err != nil {
		return nil, err
	}

	pictures = make([]*Picture, 0)
	for rows.Next() {
		picture := Picture{}
		err = rows.Scan(picture.ID, picture.Title, picture.Base64, picture.UserId)
		if err != nil {
			return nil, err
		}
		pictures = append(pictures, &picture)
	}

	return pictures, nil
}

func (s *PictureService) GetUserPicture(userId string, pictureId string) (res *Picture, err error) {
	q := `SELECT * FROM pictures
		WHERE id = $1 AND pictures.user_id = $2`
	row := s.db.QueryRow(q, pictureId, userId)

	res = &Picture{}
	err = row.Scan(res.ID, res.Title, res.Base64, res.UserId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *PictureService) CreateUserPicture(userId string, data *CreatePicturesDto) error {
	q := `INSERT INTO pictures (title, base64, user_id)
		VALUES ($1, $2, $3)`
	_, err := s.db.Query(q, data.Title, data.Base64, userId)
	return err
}

func (s *PictureService) DeleteUserPicture(userId string, pictureId string) error {
	q := `DELETE FROM pictures
		WHERE id = $1 AND user_id = $2`
	_, err := s.db.Query(q, pictureId, userId)
	return err
}

func NewPicturesService(db *sqlx.DB) *PictureService {
	return &PictureService{db}
}
