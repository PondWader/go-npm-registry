package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PondWader/go-npm-registry/pkg/response"
)

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

	for tag, version := range body.DistTags {
		versionData, ok := body.Versions[version]
		if !ok {
			continue
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
		fmt.Println(tag, data)
	}

	response.Error(w, http.StatusBadRequest, "Version already exists")
}
