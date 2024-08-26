package portal

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/sethvargo/go-githubactions"
	"go.uber.org/zap"
)

// actionPkg manages business logic around action request
type actionPkg interface {
	Context() (*githubactions.GitHubContext, error)
	Infof(msg string, args ...any)
	Debugf(msg string, args ...any)
	Errorf(msg string, args ...any)
	Fatalf(msg string, args ...any)
	SetOutput(k string, v string)
}

// Handler manages portal requests
type Handler struct {

	// isRunningLocal is true when running locally
	isRunningLocal bool

	// actionPkg represents the githubactions package
	actionPkg actionPkg

	// embeddedContent embedded content of the web app
	embeddedContent fs.FS

	// embeddedContentFilePathPrefix path prefix of the embedded content
	embeddedContentFilePathPrefix string

	// githubToken is the github token used to make Api calls
	githubToken string
}

// NewHandler returns portal handler
func NewHandler(actionPkg actionPkg, isRunningLocal bool, embeddedContent fs.FS, embeddedContentFilePathPrefix, githubToken string) *Handler {
	return &Handler{
		isRunningLocal:                isRunningLocal,
		actionPkg:                     actionPkg,
		embeddedContent:               embeddedContent,
		embeddedContentFilePathPrefix: embeddedContentFilePathPrefix,
		githubToken:                   githubToken,
	}
}

// CancelPortal returns response for request to cancel the portal
func (h *Handler) CancelPortal(w http.ResponseWriter, r *http.Request) {

	additionalContext := map[string]string{
		"JobUrl": "",
	}

	actionContext, err := h.actionPkg.Context()
	if err != nil {
		h.actionPkg.Errorf("Unable to get action context: %v", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	repoOwner, repoName := actionContext.Repo()

	additionalContext["JobUrl"] = fmt.Sprintf(
		"%s/%s/actions/runs/%d",
		actionContext.ServerURL,
		repoOwner+"/"+repoName,
		actionContext.RunID,
	)

	// Parse template
	parsedTemplates, err := template.ParseFS(h.embeddedContent, fmt.Sprintf("%sweb/ui/html/partials/responses/cancel.tmpl.html", h.embeddedContentFilePathPrefix))
	if err != nil {
		h.actionPkg.Errorf("Unable to parse referenced template: %v", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Added templates needed for htmx replacement
	w.Header().Set("HX-Trigger", "template-executed")
	w.Header().Set("HX-Trigger-After-Swap", "template-swapped")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")

	// Write template to response
	err = parsedTemplates.Execute(w, additionalContext)
	if err != nil {
		h.actionPkg.Errorf("Unable to execute parsed template: %v", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.actionPkg.Infof("Cancel request received")

	go func(actionContext *githubactions.GitHubContext) {

		runId := actionContext.RunID
		h.actionPkg.Infof("Cancelling job within run %d", runId)
		time.Sleep(3 * time.Second)

		h.actionPkg.Fatalf("Job within run %d cancelled", runId)

	}(actionContext)

}

// SubmitPortal returns response for request to submit the portal
func (h *Handler) SubmitPortal(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	additionalContext := map[string]string{
		"JobUrl": "",
	}

	if h.isRunningLocal {
		h.actionPkg.Infof("Running locally, will only print the form data to stdout")
	}

	for key, value := range r.Form {
		h.actionPkg.Infof("%s: %s", key, value)

		if !h.isRunningLocal {
			// Can't use when running locally
			h.actionPkg.SetOutput(key, strings.Join(value, ","))
		}
	}

	actionContext, err := h.actionPkg.Context()
	if err != nil {
		h.actionPkg.Errorf("Unable to get action context: %v", zap.Error(err))
	}

	if err == nil {
		repoOwner, repoName := actionContext.Repo()

		additionalContext["JobUrl"] = fmt.Sprintf(
			"%s/%s/actions/runs/%d",
			actionContext.ServerURL,
			repoOwner+"/"+repoName,
			actionContext.RunID,
		)

		// TODO: Figure out if there is a way to get the the exact job Id of current
		// job.
		// if actionContext.Job != "" {
		// 	additionalContext["JobUrl"] = additionalContext["JobUrl"] + "/job/" + actionContext.Job
		// }
	}

	// Parse template
	parsedTemplates, err := template.ParseFS(h.embeddedContent, fmt.Sprintf("%sweb/ui/html/partials/responses/success.tmpl.html", h.embeddedContentFilePathPrefix))
	if err != nil {
		h.actionPkg.Errorf("Unable to parse referenced template: %v", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Added templates needed for htmx replacement
	w.Header().Set("HX-Trigger", "template-executed")
	w.Header().Set("HX-Trigger-After-Swap", "template-swapped")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")

	// Write template to response
	err = parsedTemplates.Execute(w, additionalContext)
	if err != nil {
		h.actionPkg.Errorf("Unable to execute parsed template: %v", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.actionPkg.Infof("Your inputs have successfully been received!")

	// put an exit command in background so that the action can finish
	go func() {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
}
