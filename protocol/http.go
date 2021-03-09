package protocol

import (
	"log"

	"github.com/HangoKub/Hango-service/config"
	"github.com/HangoKub/Hango-service/internal/core/usecases"
	"github.com/HangoKub/Hango-service/internal/handlers"
	"github.com/HangoKub/Hango-service/internal/repositories"
	"github.com/HangoKub/Hango-service/pkg/connection/databases"
	"github.com/HangoKub/Hango-service/pkg/genkey"
	"github.com/HangoKub/Hango-service/pkg/validators"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ServeHttp() error {
	app := fiber.New()

	app.Use(cors.New())

	cfg := config.ParseConfig()
	FirestoreConn, err := databases.NewFirestoreAuth(cfg)
	if err != nil {
		log.Fatal(err)
	}
	FirestoreConn, err = databases.NewFirestoreUser(cfg)
	if err != nil {
		log.Fatal(err)
	}
	FirestoreConn, err = databases.NewFirestoreRestaurant(cfg)
	if err != nil {
		log.Fatal(err)
	}
	FirestoreConn, err = databases.NewFirestoreFeed(cfg)
	if err != nil {
		log.Fatal(err)
	}
	FirestoreConn, err = databases.NewStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	validators.InitializeTranslator()
	privKey, publiKey := genkey.GenerateRsaKey(cfg.PrivKey, cfg.PublicKey)
	uploadFileRepo := repositories.NewUploadStorage(FirestoreConn.CloudStorage, cfg.StorageName)
	authFirestoreRepo := repositories.NewAuthFirestore(FirestoreConn.AuthFirestore)
	authUc := usecases.NewAuthUsecase(authFirestoreRepo, privKey, publiKey, cfg.MessageBridKey)
	authMiddleware := handlers.NewMiddlewareAuth(privKey)
	userFirestoreRepo := repositories.NewUserFirestore(FirestoreConn.UserFirestore)
	userUc := usecases.NewUserUsecase(userFirestoreRepo)
	authHandler := handlers.NewAuthHandler(authUc, userUc)
	restFirestoreRepo := repositories.NewRestaurantFirestore(FirestoreConn.RestFirestore)
	restUc := usecases.NewRestaurantUsecase(restFirestoreRepo)
	restHandler := handlers.NewRestaurantHandler(restUc)
	postFirestoreRepo := repositories.NewPostFeedFirestore(FirestoreConn.FeedFirestore)
	postUc := usecases.NewPostUsecase(restFirestoreRepo, postFirestoreRepo, uploadFileRepo)
	postHandler := handlers.NewPostFeedHandler(postUc)
	commentFirestoreRepo := repositories.NewCommentPostFirestore(FirestoreConn.FeedFirestore)
	commentUc := usecases.NewCommentUsecase(commentFirestoreRepo, uploadFileRepo)
	commentHandler := handlers.NewCommentHandler(postUc, commentUc)

	setupRoute(app, authHandler, restHandler, postHandler, commentHandler, authMiddleware)

	defer func() {
		FirestoreConn.AuthFirestore.Close()
		FirestoreConn.RestFirestore.Close()
		FirestoreConn.UserFirestore.Close()
		FirestoreConn.FeedFirestore.Close()
	}()

	log.Fatal(app.Listen(":" + cfg.ServerPort))
	return nil
}
