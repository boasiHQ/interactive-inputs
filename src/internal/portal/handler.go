package portal

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/ooaklee/reply"
	"github.com/sethvargo/go-githubactions"
	"go.uber.org/zap"
)

// actionPkg manages business logic around action request
type actionPkg interface {
	Context() (*githubactions.GitHubContext, error)
	Infof(msg string, args ...any)
	Warningf(msg string, args ...any)
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

	// inputFieldLabelToCacheDirMapping mapping of input field label to its cache directory
	inputFieldLabelToCacheDirMapping map[string]string
}

// NewHandler returns portal handler
func NewHandler(actionPkg actionPkg, isRunningLocal bool, embeddedContent fs.FS, embeddedContentFilePathPrefix, githubToken string, inputFieldLabelToCacheDirMapping map[string]string) *Handler {
	return &Handler{
		isRunningLocal:                   isRunningLocal,
		actionPkg:                        actionPkg,
		embeddedContent:                  embeddedContent,
		embeddedContentFilePathPrefix:    embeddedContentFilePathPrefix,
		githubToken:                      githubToken,
		inputFieldLabelToCacheDirMapping: inputFieldLabelToCacheDirMapping,
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

		// handle file/multifile inputs
		if cacheDir := h.getInputFieldCacheDir(key); cacheDir != "" {

			h.actionPkg.Infof("%s: %s", key, cacheDir)

			if !h.isRunningLocal {
				// Can't use when running locally
				h.actionPkg.SetOutput(key, cacheDir)
			}

			continue
		}

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

// UploadToPortal returns response for request to upload file(s) to portal
// for later use
func (h *Handler) UploadToPortal(w http.ResponseWriter, r *http.Request) {

	const indexKeySplitter string = "__index__"
	var totalFiles int
	var fileCount int = 0
	var successFileUploads []string = []string{}
	var failedFileUploads []string = []string{}

	h.actionPkg.Infof("Uploading File(s)")

	r.ParseMultipartForm(10 << 20)

	// If no files are uploaded, return an error
	if r.MultipartForm == nil {
		h.actionPkg.Errorf("No files detected in upload request")

		//nolint will set up default fallback later
		getBaseResponseHandler().NewHTTPErrorResponse(w, errors.New(ErrNoFilesProvidedWithUploadRequest))
		return
	}

	// Get a reference to the parsed multipart form
	form := r.MultipartForm
	files := form.File

	totalFiles = len(files)
	h.actionPkg.Infof("Total pushed files: %d", totalFiles)

	// Get a reference to the parsed file
	for k, _ := range files {

		fileCount++

		h.actionPkg.Infof("[%d of %d] Initiating file upload flow", fileCount, totalFiles)

		file, handler, err := r.FormFile(k)
		if err != nil {
			h.actionPkg.Errorf("[%d of %d] Error Retrieving the file: %v", fileCount, totalFiles, err)
			failedFileUploads = append(failedFileUploads, k)
			continue
		}

		defer file.Close()

		// split index from file name to get the input name
		indexArray := strings.Split(k, indexKeySplitter)
		inputFieldLabel := indexArray[0]

		h.actionPkg.Debugf("  • Input Field: %+v", inputFieldLabel)
		h.actionPkg.Debugf("  • Uploaded File: %+v", handler.Filename)
		h.actionPkg.Debugf("  • File Size: %+v", handler.Size)
		h.actionPkg.Debugf("  • MIME Header: %+v", handler.Header)
		h.actionPkg.Debugf("")

		// Read file into byte array
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			h.actionPkg.Errorf("[%d of %d] Unable to read file: %v", fileCount, totalFiles, err)
			failedFileUploads = append(failedFileUploads, handler.Filename)
			continue
		}

		// create placeholder file in temp directory to hold uploaded file
		inputCacheDir := h.getInputFieldCacheDir(inputFieldLabel)
		err = os.WriteFile(fmt.Sprintf("%s/%s", inputCacheDir, handler.Filename), fileBytes, 0644)
		if err != nil {
			h.actionPkg.Errorf("[%d of %d] Unable to write file to input field cache dir: %s", fileCount, totalFiles, inputCacheDir)
			failedFileUploads = append(failedFileUploads, handler.Filename)
			continue
		}

		// add file to successful uploads
		successFileUploads = append(successFileUploads, handler.Filename)

	}

	h.actionPkg.Infof("Successfully uploaded %d of %d files", len(successFileUploads), totalFiles)

	response := UploadToPortalResponse{
		UploadedFiles: successFileUploads,
		FailedFiles:   failedFileUploads,
	}

	// TODO: better handle these failure/ partial failure situation
	// might make sense to return a 4XX and allow the frontend to
	// retrospectively remove the reference of the failed files from
	// the input fiels's value/state (user will have to be notified
	// that this is the case)
	if len(failedFileUploads) > 0 && len(successFileUploads) > 0 {
		response.Status = "partial success"
	}
	if len(successFileUploads) == 0 && len(failedFileUploads) > 0 {
		response.Status = "failed"
	}

	if len(failedFileUploads) == 0 && len(successFileUploads) == totalFiles {
		response.Status = "success"
	}

	//nolint will set up default fallback later
	getBaseResponseHandler().NewHTTPDataResponse(w, http.StatusOK, &response)
	return
}

// ResetUpload returns response for request to reset upload,
// which removes all files from the cache directory for the given input field name.
func (h *Handler) ResetUpload(w http.ResponseWriter, r *http.Request) {
	var inputFieldLabel string

	// Get the input field name from the request
	if inputFieldLabel = mux.Vars(r)[InputFieldLabelUriVariableId]; inputFieldLabel == "" {
		h.actionPkg.Errorf("Input field label not found in request")

		//nolint will set up default fallback later
		getBaseResponseHandler().NewHTTPErrorResponse(w, errors.New(ErrKeyInvalidInputFieldId))
		return
	}

	status, totalFilesToDelete, totalFilesDeleted, deletedFiles, failedFiles, err := h.cleanUpCacheDir(inputFieldLabel)
	if err != nil && err.Error() == ErrKeyUnableToRemoveCacheDirContents {

		//nolint will set up default fallback later
		getBaseResponseHandler().NewHTTPErrorResponse(w, errors.New(ErrKeyUnableToRemoveCacheDirContents),
			reply.WithMeta(map[string]interface{}{"data": ResetUploadResponse{
				Status:             status,
				DeletedFiles:       deletedFiles,
				FailedFiles:        failedFiles,
				TotalFilesToDelete: totalFilesToDelete,
				TotalFilesDeleted:  totalFilesDeleted,
			}}))
		return
	}

	if err != nil {
		//nolint will set up default fallback later
		getBaseResponseHandler().NewHTTPErrorResponse(w, err)
		return
	}

	h.actionPkg.Infof("Cache directory contents reseted for input field label: %s", inputFieldLabel)
	//nolint will set up default fallback later
	getBaseResponseHandler().NewHTTPDataResponse(w, http.StatusOK, &ResetUploadResponse{
		Status:             status,
		DeletedFiles:       deletedFiles,
		FailedFiles:        failedFiles,
		TotalFilesToDelete: totalFilesToDelete,
		TotalFilesDeleted:  totalFilesDeleted,
	})
}

// cleanUpCacheDir removes all files from the cache directory for the given input field name
func (h *Handler) cleanUpCacheDir(inputFieldLabel string) (string, int, int, []string, []string, error) {

	var (
		status             string
		totalFilesToDelete int      = 0
		totalFilesDeleted  int      = 0
		deletedFiles       []string = []string{}
		failedFiles        []string = []string{}
	)

	// Remove all files from the cache directory for the given input field name
	cacheDir := h.getInputFieldCacheDir(inputFieldLabel)
	if cacheDir == "" {
		h.actionPkg.Errorf("No cache directory found for input field label: %s", inputFieldLabel)

		return status, totalFilesToDelete, totalFilesDeleted, deletedFiles, failedFiles, errors.New(ErrKeyNoInputFieldCacheDirFound)
	}

	h.actionPkg.Infof("Initiating the reseting of the cache directory contents for the input field label: %s (%s)", inputFieldLabel, cacheDir)

	// Remove the cache directory contents for the given input field name
	readCacheDir, err := os.ReadDir(cacheDir)
	if err != nil {
		h.actionPkg.Errorf("Unable to read cache directory: %s", cacheDir)

		return status, totalFilesToDelete, totalFilesDeleted, deletedFiles, failedFiles, errors.New(ErrKeyUnableToReadCacheDir)
	}

	if len(readCacheDir) == 0 {
		h.actionPkg.Infof("No cache directory contents found for input field label: %s", inputFieldLabel)

		return status, totalFilesToDelete, totalFilesDeleted, deletedFiles, failedFiles, nil
	}

	totalFilesToDelete = len(readCacheDir)

	h.actionPkg.Infof("Cache directory contents (%d) found for input field label: %s", totalFilesToDelete, inputFieldLabel)

	for _, content := range readCacheDir {
		contentFullPath := path.Join([]string{cacheDir, content.Name()}...)

		h.actionPkg.Debugf("  • Removing file: %s (%s)", content.Name(), contentFullPath)
		err = os.RemoveAll(contentFullPath)
		if err != nil {
			h.actionPkg.Errorf("Unable to remove file: %s", contentFullPath)
			failedFiles = append(failedFiles, content.Name())

			if len(failedFiles) > 0 && len(deletedFiles) > 0 {
				status = "partial success"
			}

			if len(failedFiles) == 0 && len(deletedFiles) == totalFilesToDelete {
				status = "success"
			}

			if len(deletedFiles) == 0 && len(failedFiles) > 0 {
				status = "failed"
			}

			return status, totalFilesToDelete, totalFilesDeleted, deletedFiles, failedFiles, errors.New(ErrKeyUnableToRemoveCacheDirContents)
		}

		totalFilesDeleted++
		deletedFiles = append(deletedFiles, content.Name())
	}

	status = "success"

	return status, totalFilesToDelete, totalFilesDeleted, deletedFiles, failedFiles, nil

}

// getInputFieldCacheDir returns the cache directory path for the given input field name.
func (h *Handler) getInputFieldCacheDir(inputFieldName string) string {
	return h.inputFieldLabelToCacheDirMapping[inputFieldName]
}

// getBaseResponseHandler returns response handler configured with respective error map
func getBaseResponseHandler() *reply.Replier {
	return reply.NewReplier(append([]reply.ErrorManifest{}, portalErrorMap))
}
