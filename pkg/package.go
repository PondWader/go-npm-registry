package pkg

import (
	"errors"
	"net/http"

	"github.com/PondWader/go-npm-registry/pkg/database"
	"gorm.io/gorm"
)

func GetPackage(ctx RequestContext, w http.ResponseWriter, r *http.Request) {
	var pkg database.PackageModel
	if tx := ctx.DB.First(&pkg, "name = ?", "name here from req param"); errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		w.WriteHeader(404)
		w.Write([]byte(`{"error": "Not found"}`))
		return
	} else if tx.Error != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "An unkown error occured"}`))
		return
	}

}
