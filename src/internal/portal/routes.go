package portal

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// portalEventHandler describes the portal event endpoints mounted by AttachRoutes.
type portalEventHandler interface {
	SubmitPortal(w http.ResponseWriter, r *http.Request)
	CancelPortal(w http.ResponseWriter, r *http.Request)
	UploadToPortal(w http.ResponseWriter, r *http.Request)
	ResetUpload(w http.ResponseWriter, r *http.Request)
}

// uiHandler describes the web UI endpoint mounted by AttachRoutes.
type uiHandler interface {
	Home(w http.ResponseWriter, r *http.Request)
}

// AttachRoutesRequest holds everything needed to mount the portal routes.
type AttachRoutesRequest struct {

	// Router main router being served by API
	Router *mux.Router

	// PortalEventHandler valid portal event handler
	PortalEventHandler portalEventHandler

	// UiHandler valid ui handler
	UiHandler uiHandler

	// EmbeddedContent embedded content of the web app
	EmbeddedContent fs.FS

	// EmbeddedContentFilePathPrefix path prefix of the embedded content
	EmbeddedContentFilePathPrefix string

	// ActionPkg represents the githubactions package
	ActionPkg actionPkg
}

// AttachRoutes mounts static assets, page routes, and upload API routes on the
// provided router.
func AttachRoutes(request *AttachRoutesRequest) {

	// Create filesystem only holding static assets
	staticSubFS, err := fs.Sub(request.EmbeddedContent, fmt.Sprintf("%sweb/ui/static", request.EmbeddedContentFilePathPrefix))
	if err != nil {
		request.ActionPkg.Errorf("unable-to-create-file-system-for-static-assets: %v", err)
		os.Exit(1)
	}

	// Create path for handling static assets
	request.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticSubFS))))

	request.Router.HandleFunc("/", request.UiHandler.Home).Methods("GET")
	request.Router.HandleFunc("/submit", request.PortalEventHandler.SubmitPortal).Methods("POST")
	request.Router.HandleFunc("/cancel", request.PortalEventHandler.CancelPortal).Methods("POST")

	apiRouter := request.Router.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/upload", request.PortalEventHandler.UploadToPortal).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc(fmt.Sprintf("/reset/{%s}", InputFieldLabelUriVariableId), request.PortalEventHandler.ResetUpload).Methods("DELETE", "OPTIONS")

}
