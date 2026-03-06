package service

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ImageService interface {
	UploadImage(ctx context.Context, file interface{}, folder string) (string, error)
	DeleteImage(ctx context.Context, publicID string) error
}

type cloudinaryService struct {
	clnd *cloudinary.Cloudinary
}

func NewImageService() (ImageService, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("cloudinary credentials not configured")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	return &cloudinaryService{clnd: cld}, nil
}

func (s *cloudinaryService) UploadImage(ctx context.Context, file interface{}, folder string) (string, error) {
	resp, err := s.clnd.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func (s *cloudinaryService) DeleteImage(ctx context.Context, publicID string) error {
	_, err := s.clnd.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
