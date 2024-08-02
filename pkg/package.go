package pkg

import (
	"errors"
	"net/http"

	"github.com/PondWader/go-npm-registry/pkg/database"
	"github.com/PondWader/go-npm-registry/pkg/response"
	"gorm.io/gorm"
)

func GetPackage(ctx RequestContext, w http.ResponseWriter, r *http.Request) {
	var pkg database.Package
	if tx := ctx.DB.First(&pkg, "name = ?", "name here from req param"); errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		response.Error(w, http.StatusNotFound, "Not found")
		return
	} else if tx.Error != nil {
		response.Error(w, http.StatusInternalServerError, "An unkown error occured")
		return
	}

}
