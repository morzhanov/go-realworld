package services

import (
	"github.com/jmoiron/sqlx"
	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	. "github.com/morzhanov/go-realworld/internal/pictures/models"
)

type PictureService struct {
	db *sqlx.DB
}

func (s *PictureService) GetUserPictures(userId string) (*prpc.PicturesMessage, error) {
	q := `SELECT * FROM pictures
		WHERE pictures.user_id = $1`
	rows, err := s.db.Query(q, userId)
	if err != nil {
		return nil, err
	}

	pictures := make([]*Picture, 0)
	for rows.Next() {
		pic := Picture{}
		err = rows.Scan(&pic.ID, &pic.Title, &pic.Base64, &pic.UserId)
		if err != nil {
			return nil, err
		}
		pictures = append(pictures, &pic)
	}
	res := prpc.PicturesMessage{}
	for _, picture := range pictures {
		res.Pictures = append(res.Pictures, &prpc.PictureMessage{
			Id:     picture.ID,
			Title:  picture.Title,
			Base64: picture.Base64,
			UserId: picture.UserId,
		})
	}
	return &res, nil
}

func (s *PictureService) GetUserPicture(userId string, pictureId string) (*prpc.PictureMessage, error) {
	q := `SELECT * FROM pictures
		WHERE id = $1 AND pictures.user_id = $2`
	row := s.db.QueryRow(q, pictureId, userId)

	res := &Picture{}
	err := row.Scan(&res.ID, &res.Title, &res.Base64, &res.UserId)
	if err != nil {
		return nil, err
	}
	return &prpc.PictureMessage{
		Id:     res.ID,
		Title:  res.Title,
		Base64: res.Base64,
		UserId: res.UserId,
	}, nil
}

func (s *PictureService) CreateUserPicture(data *prpc.CreateUserPictureRequest) (*prpc.PictureMessage, error) {
	q := `INSERT INTO pictures (title, base64, user_id)
		VALUES ($1, $2, $3)
		RETURNING *`
	row, err := s.db.Query(q, data.Title, data.Base64, data.UserId)
	if err != nil {
		return nil, err
	}

	res := &Picture{}
	if err = row.Scan(&res.ID, &res.Title, &res.Base64, &res.UserId); err != nil {
		return nil, err
	}
	return &prpc.PictureMessage{
		Id:     res.ID,
		Title:  res.Title,
		Base64: res.Base64,
		UserId: res.UserId,
	}, nil
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
