package portal

import (
	"net/http"

	"github.com/ooaklee/reply"
)

// portalErrorMap holds Error keys, their corresponding human-friendly message, and response status code
var portalErrorMap = map[string]reply.ErrorManifestItem{
	ErrNoFilesProvidedWithUploadRequest:  {Title: "Bad Request", Detail: "No files detected. Verify file(s) submitted with upload request", StatusCode: http.StatusBadRequest},
	ErrKeyInvalidInputFieldId:            {Title: "Bad Request", Detail: "Target input field id (label) missing or malformatted", StatusCode: http.StatusBadRequest},
	ErrKeyNoInputFieldCacheDirFound:      {Title: "Bad Request", Detail: "No cache directory found for input field label", StatusCode: http.StatusBadRequest},
	ErrKeyUnableToReadCacheDir:           {Title: "Internal Server Error", Detail: "Unable to read cache directory", StatusCode: http.StatusInternalServerError},
	ErrKeyUnableToRemoveCacheDirContents: {Title: "Internal Server Error", Detail: "Unable to remove cache directory content(s)", StatusCode: http.StatusInternalServerError},
}
