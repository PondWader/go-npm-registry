package pkg

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/PondWader/go-npm-registry/pkg/database"
	"github.com/PondWader/go-npm-registry/pkg/response"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Regex for validating package names taken from https://github.com/dword-design/package-name-regex/blob/master/src/index.js
var PackageNameRegex, _ = regexp.Compile(`^(@[a-z0-9-~][a-z0-9-._~]*\/)?[a-z0-9-~][a-z0-9-._~]*$`)

type PublishPackageBody struct {
	Name     string            `json:"name"`
	DistTags map[string]string `json:"dist-tags"`
	Versions map[string]struct {
		Name        string         `json:"name"`
		Main        string         `json:"main"`
		Module      string         `json:"module"`
		Exports     map[string]any `json:"exports"`
		Description string         `json:"description"`
		Author      struct {
			Name string `json:"name"`
		} `json:"author"`
		Homepage             string            `json:"homepage"`
		Bugs                 string            `json:"bugs"`
		License              string            `json:"license"`
		Bin                  map[string]string `json:"bin"`
		Dependencies         map[string]string `json:"dependencies"`
		PeerDependencies     map[string]string `json:"peerDependencies"`
		OptionalDependencies map[string]string `json:"optionalDependencies"`
		Engines              map[string]string `json:"engines"`
		Types                string            `json:"types"`
		Dist                 struct {
			Integrity string `json:"integrity"`
			Shasum    string `json:"shasum"`
			Tarball   string `json:"tarball"`
		} `json:"dist"`
	} `json:"versions"`
	Attachments map[string]struct {
		ContentType string `json:"content_type"`
		Data        string `json:"data"`
		Length      int    `json:"length"`
	} `json:"_attachments"`
}

func PublishPackage(ctx RequestContext, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body PublishPackageBody
	decoder.Decode(&body)

	// Validate name
	if !PackageNameRegex.MatchString(body.Name) {
		response.Error(w, http.StatusBadRequest, "Invalid package name")
		return
	}

	for tag, version := range body.DistTags {
		versionData, ok := body.Versions[version]
		if !ok {
			continue
		}

		// Check that there is not already a package published with the same version
		queryResult := database.PackageVersion{}
		ctx.DB.Select("ID").Where("version = ?", version).First(&queryResult)
		if queryResult.Version != version {
			response.Error(w, http.StatusBadGateway, "Version already exists")
			return
		}

		distUrl, err := url.Parse(versionData.Dist.Tarball)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid tarball URL")
			return
		}
		// Assume that path begins with /<package name>/-/@passiveapp
		attachmentName := distUrl.Path[len(body.Name)+4:]
		attachment, ok := body.Attachments[attachmentName]
		if !ok {
			response.Error(w, http.StatusBadRequest, "Missing attachment")
			return
		}

		if attachment.ContentType != "application/octet-stream" {
			response.Error(w, http.StatusUnsupportedMediaType, "Attachment is not of expected type application/octet-stream")
			return
		}

		data, err := base64.StdEncoding.DecodeString(attachment.Data)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Bad attachment data")
			return
		}
		if len(data) != attachment.Length {
			response.Error(w, http.StatusBadRequest, "Data length does not match attachment length value")
			return
		}

		// TODO: Check .tgz file is valid?

		fileName := body.Name + "-" + url.QueryEscape(version) + ".tgz"
		if err := ctx.Storage.Write(fileName, data); err != nil {
			response.Error(w, http.StatusInternalServerError, "An internal error occured saving file")
			fmt.Println("Failed to save", fileName+":", err)
			return
		}

		err = ctx.DB.Transaction(func(tx *gorm.DB) error {
			// Make sure package exists in DB
			var packageRecord database.Package
			res := tx.Where("name = ?", body.Name).Find(&packageRecord)
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				res = tx.Create(&database.Package{
					Name: body.Name,
					DistTags: datatypes.JSONMap{
						tag: version,
					},
				})
				if res.Error != nil {
					response.Error(w, http.StatusInternalServerError, "An internal error occured")
					fmt.Println("Failed to create package record:", res.Error)
					return err
				}
			} else if res.Error != nil {
				response.Error(w, http.StatusInternalServerError, "An internal error occured")
				fmt.Println("Failed to query package record:", res.Error)
				return err
			} else {
				// Add dist tag
				res = tx.Model(&database.Package{}).Where("id = ?", packageRecord.Name).UpdateColumn("dist_tags", datatypes.JSONSet("dist_tags").Set(tag, version))
				if res.Error != nil {
					response.Error(w, http.StatusInternalServerError, "An internal error occured")
					fmt.Println("Failed to update package dist tags:", res.Error)
					return err
				}
			}

			uuid, err := uuid.NewRandom()
			if err != nil {
				response.Error(w, http.StatusInternalServerError, "An internal error occured")
				fmt.Println("Failed to generate UUID:", err)
				return err
			}

			res = tx.Create(&database.PackageVersion{
				ID:        uuid,
				PackageID: 5,
			})
			if res.Error != nil {
				response.Error(w, http.StatusInternalServerError, "An internal error occured")
				fmt.Println("Failed to generate UUID:", res.Error)
				return err
			}

			return nil
		})
		if err != nil {
			return
		}
	}

	w.WriteHeader(200)
}
