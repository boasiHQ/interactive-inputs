package webui

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/boasihq/interactive-inputs/internal/config"
	"github.com/boasihq/interactive-inputs/internal/toolbox"
	githubactions "github.com/sethvargo/go-githubactions"
	"go.uber.org/zap"
)

// NewWebAppHandlerRequest is the request needed to create an ui handler
type NewWebAppHandlerRequest struct {
	EmbeddedContent fs.FS
	// EmbeddedContentFilePathPrefix the prefix used to access the embedded files
	EmbeddedContentFilePathPrefix string
	// Config is the configuration of the action
	Config *config.Config
}

// NewWebAppHandler creates a new instance of an ui handler
func NewWebAppHandler(r *NewWebAppHandlerRequest) *Handler {
	return &Handler{
		embeddedFileSystem:            r.EmbeddedContent,
		embeddedContentFilePathPrefix: r.EmbeddedContentFilePathPrefix,
		action:                        r.Config.Action,
		config:                        r.Config,
	}
}

// Handler manages request for webapp
type Handler struct {
	embeddedFileSystem            fs.FS
	embeddedContentFilePathPrefix string
	action                        *githubactions.Action
	config                        *config.Config
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {

	var response *CreateInteractiveInputsPortalRequest

	actionContext, err := h.action.Context()
	if err != nil {
		h.action.Errorf("Unable to get action context: %v", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// shape the request
	repoOwner, _ := actionContext.Repo()
	response = &CreateInteractiveInputsPortalRequest{
		RepoOwner: repoOwner,
		Title:     h.config.Title,
		Fields:    h.config.Fields,
		Timeout:   toolbox.SecondsToMinutes(h.config.Timeout),
	}

	// Calculate the base path once for the template
    basePath := strings.Trim(h.config.RunnerEndpointKey, "/ ")
    if basePath == "" {
        basePath = "runner"
    }
    response.BasePath = "/" + basePath

	// Build balloon suggestion data from field properties and environment
    balloonData := make(map[string][]string)
    preOutput := make(map[string]struct{ Title, Value string })
    if h.config.Fields != nil {
        for _, f := range h.config.Fields.Fields {
            var suggestions []string

            // Static suggestions
            if len(f.Properties.BalloonValues) > 0 {
                suggestions = append(suggestions, f.Properties.BalloonValues...)
            }

            // Env-derived suggestions
            if len(f.Properties.BalloonValueEnvKeys) > 0 {
                for _, key := range f.Properties.BalloonValueEnvKeys {
                    val := strings.TrimSpace(os.Getenv(key))
                    if val == "" {
                        continue
                    }
                    suggestions = append(suggestions, val)
                }
            }

            if len(suggestions) > 0 {
                balloonData[f.Label] = suggestions
            }

            // PreOutput from a single env var
            if k := strings.TrimSpace(f.Properties.OutputFromEnvKey); k != "" {
                if v := strings.TrimSpace(os.Getenv(k)); v != "" {
                    t := strings.TrimSpace(f.Properties.OutputTitle)
                    preOutput[f.Label] = struct{ Title, Value string }{Title: t, Value: v}
                }
            }
        }
    }
    response.BalloonData = balloonData
    response.PreOutput = preOutput

	// list of template files to parse, must be in order of inheritence
	templateFilesToParse := []string{
		fmt.Sprintf("%sweb/ui/html/index.tmpl.html", h.embeddedContentFilePathPrefix),
		fmt.Sprintf("%sweb/ui/html/partials/shared/head-meta.tmpl.html", h.embeddedContentFilePathPrefix),
		fmt.Sprintf("%sweb/ui/html/pages/@landing.tmpl.html", h.embeddedContentFilePathPrefix),
		fmt.Sprintf("%sweb/ui/html/partials/shared/tailwind-dash-script.tmpl.html", h.embeddedContentFilePathPrefix),
	}

	// Parse template
	parsedTemplates, err := template.ParseFS(h.embeddedFileSystem, templateFilesToParse...)
	if err != nil {
		h.action.Errorf("Unable to parse referenced template: %v", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write template to response
	err = parsedTemplates.ExecuteTemplate(w, "base", response)
	if err != nil {
		h.action.Errorf("Unable to execute parsed template: %v", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
