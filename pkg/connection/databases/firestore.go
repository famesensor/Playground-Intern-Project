package databases

// Commenter: sifer169966,
// These code should be into pkg/connection ...

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/HangoKub/Hango-service/config"
	"google.golang.org/api/option"
)

type FirebaseConn struct {
	AuthFirestore *firestore.Client
	UserFirestore *firestore.Client
	RestFirestore *firestore.Client
	FeedFirestore *firestore.Client
	CloudStorage  *storage.BucketHandle
}

var firebaseConn = &FirebaseConn{}

func NewFirestoreAuth(cfg *config.Config) (*FirebaseConn, error) {
	ctx := context.Background()

	var otps option.ClientOption
	switch cfg.Env {
	case "local":
		otps = option.WithCredentialsFile(cfg.FirestoreCert)
		// otps = option.WithEndpoint("localhost:4000")
		break
	case "dev":
		otps = option.WithCredentialsFile(cfg.FirestoreCert)
		break
	default:
		return nil, errors.New("unexpected run mode from config.NewFirestoreAuth")
	}
	conn, err := firebase.NewApp(ctx, nil, otps)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	firestore, err := conn.Firestore(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	firebaseConn.AuthFirestore = firestore
	return firebaseConn, nil
}

func NewFirestoreUser(cfg *config.Config) (*FirebaseConn, error) {
	ctx := context.Background()

	var otps option.ClientOption
	switch cfg.Env {
	case "local":
		otps = option.WithCredentialsFile(cfg.FirestoreCert)
		// otps = option.WithEndpoint("localhost:4000")
		break
	case "dev":
		otps = option.WithCredentialsFile(cfg.FirestoreCert)
		break
	default:
		return nil, errors.New("unexpected run mode from config.NewFirestoreUser")
	}
	conn, err := firebase.NewApp(ctx, nil, otps)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	firestore, err := conn.Firestore(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	firebaseConn.UserFirestore = firestore
	return firebaseConn, nil
}

func NewFirestoreRestaurant(cfg *config.Config) (*FirebaseConn, error) {
	ctx := context.Background()

	var otps option.ClientOption
	switch cfg.Env {
	case "local":
		otps = option.WithCredentialsFile(cfg.FirestoreCert)
		// otps = option.WithEndpoint("localhost:4000")
		break
	case "dev":
		otps = option.WithCredentialsFile(cfg.FirestoreCert)
		break
	default:
		return nil, errors.New("unexpected run mode from config.NewFirestoreRestaurant")
	}
	conn, err := firebase.NewApp(ctx, nil, otps)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	firestore, err := conn.Firestore(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	firebaseConn.RestFirestore = firestore
	return firebaseConn, nil
}

func NewFirestoreFeed(cfg *config.Config) (*FirebaseConn, error) {
	ctx := context.Background()

	var otps option.ClientOption
	switch cfg.Env {
	case "local":
		otps = option.WithCredentialsFile(cfg.FirestoreCert)
		// otps = option.WithEndpoint("localhost:4000")
		break
	case "dev":
		otps = option.WithCredentialsFile(cfg.FirestoreCert)
		break
	default:
		return nil, errors.New("unexpected run mode from config.NewFirestorePostFeed")
	}
	conn, err := firebase.NewApp(ctx, nil, otps)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	firestore, err := conn.Firestore(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	firebaseConn.FeedFirestore = firestore
	return firebaseConn, nil
}

func NewStorage(cfg *config.Config) (*FirebaseConn, error) {
	ctx := context.Background()

	var otps option.ClientOption
	switch cfg.Env {
	case "local":
		otps = option.WithCredentialsFile(cfg.FirestoreCert)
		// otps = option.WithEndpoint("localhost:4000")
		break
	case "dev":
		otps = option.WithCredentialsFile(cfg.FirestoreCert)
		break
	default:
		return nil, errors.New("unexpected run mode from config.NewFirestoreUser")
	}
	config := &firebase.Config{
		StorageBucket: cfg.StorageName,
	}
	conn, err := firebase.NewApp(ctx, config, otps)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	client, err := conn.Storage(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	firebaseConn.CloudStorage = bucket
	return firebaseConn, nil
}
