package pkg

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PondWader/go-npm-registry/pkg/database"
	"github.com/PondWader/go-npm-registry/pkg/response"
	"gorm.io/gorm"
)

type VersionDist struct {
	Integrity    string `json:"integrity"`
	Shasum       string `json:"shasum"`
	Tarball      string `json:"tarball"`
	FileCount    uint   `json:"fileCount"`
	UnpackedSize uint   `json:"unpackedSize"`
}

type VersionData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Author  struct {
		Name string `json:"name"`
	} `json:"author"`
	Description          *string        `json:"description"`
	Dependencies         map[string]any `json:"dependencies"`
	PeerDependencies     map[string]any `json:"peerDependencies"`
	OptionalDependencies map[string]any `json:"optionalDependencies"`
	Dist                 VersionDist    `json:"dist"`
}

type PackageData struct {
	Time        map[string]string      `json:"time"`
	Name        string                 `json:"name"`
	DistTags    map[string]any         `json:"dist-tags"`
	Versions    map[string]VersionData `json:"versions"`
	Description *string
}

func GetPackage(ctx RequestContext, w http.ResponseWriter, r *http.Request) {
	pkgName := r.PathValue("package")
	if strings.HasSuffix(pkgName, ".tgz") {
		paramSplit := strings.Split(pkgName, "/-/")
		if len(paramSplit) == 2 {
			DownloadPackage(ctx, w, paramSplit[0], paramSplit[1])
			return
		}
	}

	var pkg database.Package
	if tx := ctx.DB.Where("name = ?", pkgName).First(&pkg); errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		response.Error(w, http.StatusNotFound, "Not found")
		return
	} else if tx.Error != nil {
		response.Error(w, http.StatusInternalServerError, "An unkown error occured")
		return
	}

	versionRecords := make([]database.PackageVersion, 1)
	if tx := ctx.DB.Find(&versionRecords); tx.Error != nil {
		response.Error(w, http.StatusInternalServerError, "An unkown error occured")
		return
	}

	pkgNameSplit := strings.Split(pkgName, "/")
	pkgScopedName := pkgNameSplit[len(pkgNameSplit)-1]

	versions := make(map[string]VersionData, len(versionRecords))
	for _, version := range versionRecords {

		data := VersionData{
			Name:    pkgName,
			Version: version.Version,
			Dist: VersionDist{
				Integrity:    version.DistIntegrity,
				Shasum:       version.DistShasum,
				Tarball:      ctx.Config.Url + "/" + pkgName + "/-/" + pkgScopedName + "-" + version.Version + ".tgz",
				FileCount:    version.DistFileCount,
				UnpackedSize: version.DistUnpackedSize,
			},
		}
		if version.Author != nil {
			data.Author = struct {
				Name string `json:"name"`
			}{*version.Author}
		}
		if version.Description != nil {
			data.Description = version.Description
		}
		data.Dependencies = version.Dependencies
		data.PeerDependencies = version.PeerDependencies
		data.OptionalDependencies = version.OptionalDependencies

		versions[version.Version] = data
	}

	data := PackageData{
		Time: map[string]string{
			"created":  pkg.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
			"modified": pkg.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
		},
		Name:        pkg.Name,
		DistTags:    pkg.DistTags,
		Versions:    versions,
		Description: versionRecords[len(versionRecords)-1].Description,
	}

	response.Json(w, data)
}

func DownloadPackage(ctx RequestContext, w http.ResponseWriter, pkgName string, fileName string) {
	var pkg database.Package
	if tx := ctx.DB.Where("name = ?", pkgName).First(&pkg); errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		response.Error(w, http.StatusNotFound, "Not found")
		return
	} else if tx.Error != nil {
		response.Error(w, http.StatusInternalServerError, "An unkown error occured")
		return
	}

	pkgNameSplit := strings.Split(pkgName, "/")
	pkgScopedName := pkgNameSplit[len(pkgNameSplit)-1]
	version := fileName[len(pkgScopedName)+1 : len(fileName)-4]

	reader, err := ctx.Storage.NewReader(pkgName + "-" + url.QueryEscape(version) + ".tgz")
	if err != nil {
		fmt.Println(err)
		response.Error(w, http.StatusInternalServerError, "An internal error occured")
		return
	}

	io.Copy(w, reader)
}
