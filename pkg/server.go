package pkg

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PondWader/go-npm-registry/pkg/config"
	"github.com/PondWader/go-npm-registry/pkg/database"
	"github.com/PondWader/go-npm-registry/pkg/storage"
	"gorm.io/gorm"
)

type RequestContext struct {
	DB      *gorm.DB
	Storage storage.StorageDriver
	Config  config.Config
	UserKey string
}

type RequestHandler func(ctx RequestContext, w http.ResponseWriter, r *http.Request)

func StartServer(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	db, err := database.Open(cfg.DbPath)
	if err != nil {
		return err
	}

	storage, err := storage.New(cfg.StorageDriver, cfg.StorageDriverOpts)
	if err != nil {
		return err
	}

	baseReqCtx := RequestContext{
		DB:      db,
		Storage: storage,
		Config:  cfg,
	}

	http.HandleFunc("GET /{package...}", ContextMiddleware(baseReqCtx, GetPackage))
	http.HandleFunc("PUT /{package}", AuthMiddleware(baseReqCtx, PublishPackage))

	fmt.Println("Starting HTTP server on port", cfg.Port)

	return http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil)
}
