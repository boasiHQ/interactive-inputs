package portal

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// portalEventHandler expected methods for valid portal event handler
type portalEventHandler interface {
	SubmitPortal(w http.ResponseWriter, r *http.Request)
	CancelPortal(w http.ResponseWriter, r *http.Request)
	UploadToPortal(w http.ResponseWriter, r *http.Request)
	ResetUpload(w http.ResponseWriter, r *http.Request)
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

	// Base path Todo
	// e.g. "run-12345" resulting in endpoints like /run-12345/submit
	BasePath string
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

    // Namespace all app routes under the provided base path
    trimmedBase := strings.Trim(request.BasePath, "/ ")
    if trimmedBase == "" {
        trimmedBase = "runner"
    }
    // Redirect root to the namespaced base path and normalize trailing slash
    request.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        prefix := strings.TrimRight(r.Header.Get("X-Forwarded-Prefix"), "/ ")
        target := "/" + trimmedBase + "/"
        if prefix != "" {
            target = prefix + target
        }
        http.Redirect(w, r, target, http.StatusPermanentRedirect)
    }).Methods("GET")
    request.Router.HandleFunc("/"+trimmedBase, func(w http.ResponseWriter, r *http.Request) {
        prefix := strings.TrimRight(r.Header.Get("X-Forwarded-Prefix"), "/ ")
        target := "/" + trimmedBase + "/"
        if prefix != "" {
            target = prefix + target
        }
        http.Redirect(w, r, target, http.StatusPermanentRedirect)
    }).Methods("GET")
    baseRouter := request.Router.PathPrefix("/" + trimmedBase).Subrouter()

    baseRouter.HandleFunc("/", request.UiHandler.Home).Methods("GET")
    baseRouter.HandleFunc("/submit", request.PortalEventHandler.SubmitPortal).Methods("POST")
    baseRouter.HandleFunc("/cancel", request.PortalEventHandler.CancelPortal).Methods("POST")

    apiRouter := baseRouter.PathPrefix("/api/v1").Subrouter()
    apiRouter.HandleFunc("/upload", request.PortalEventHandler.UploadToPortal).Methods("POST", "OPTIONS")
    apiRouter.HandleFunc(fmt.Sprintf("/reset/{%s}", InputFieldLabelUriVariableId), request.PortalEventHandler.ResetUpload).Methods("DELETE", "OPTIONS")

}
