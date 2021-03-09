package repositories

import (
	"bytes"
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type UploadStorageRepo struct {
	cloudStorage *storage.BucketHandle
	bucketName   string
}

func NewUploadStorage(cloudstorage *storage.BucketHandle, bucketName string) *UploadStorageRepo {
	return &UploadStorageRepo{
		cloudstorage,
		bucketName,
	}
}

func (r *UploadStorageRepo) UploadFiletoStorage(ctx context.Context, files []bytes.Buffer, collection, typeFile, id string) ([]string, error) {
	var urls []string

	for _, file := range files {
		newFilename := "hango-picture-" + uuid.Must(uuid.NewRandom()).String() + typeFile
		imagePath := collection + "/" + id + "/" + newFilename

		wc := r.cloudStorage.Object(imagePath).NewWriter(ctx)
		_, err := io.Copy(wc, &file)
		if err != nil {
			return nil, err
		}
		if err := wc.Close(); err != nil {
			return nil, err
		}
		url := "https://firebasestorage.googleapis.com/v0/b/" + r.bucketName + "/o/" + collection + "%2F" + id + "%2F" + newFilename + "?alt=media"
		// url := "https://storage.cloud.google.com/" + r.bucketName + "/" + imagePath
		urls = append(urls, url)
	}

	return urls, nil
}

// TODO: function delete file image in storage
