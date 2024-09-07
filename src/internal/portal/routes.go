package portal

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// portalEventHandler expected methods for valid portal event handler
type portalEventHandler interface {
	SubmitPortal(w http.ResponseWriter, r *http.Request)
	CancelPortal(w http.ResponseWriter, r *http.Request)
	UploadToPortal(w http.ResponseWriter, r *http.Request)
}

// uiHandler expected methods for valid ui handler
type uiHandler interface {
	Home(w http.ResponseWriter, r *http.Request)
}

// AttachRoutesRequest holds everything needed to attach portal
// routes to router
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

// AttachRoutes attaches portal handlers to corresponding
// routes on router
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
	request.Router.HandleFunc("/upload", request.PortalEventHandler.UploadToPortal).Methods("POST")

}
