package grafana

import (
	"fmt"
	"path/filepath"

	"github.com/grafana/grizzly/pkg/grizzly"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// AlertRuleHandler is a Grizzly Handler for Grafana dashboard AlertRules
type AlertRuleHandler struct {
	Provider Provider
}

// NewAlertRuleHandler returns configuration defining a new Grafana AlertRule Handler
func NewAlertRuleHandler(provider Provider) *AlertRuleHandler {
	return &AlertRuleHandler{
		Provider: provider,
	}
}

// Kind returns the name for this handler
func (h *AlertRuleHandler) Kind() string {
	return "AlertRule"
}

// Validate returns the uid of resource
func (h *AlertRuleHandler) Validate(resource grizzly.Resource) error {
	uid, exist := resource.GetSpecString("uid")
	if exist {
		if uid != resource.Name() {
			return fmt.Errorf("uid '%s' and name '%s', don't match", uid, resource.Name())
		}
	}

	return nil
}

// APIVersion returns the group and version for the provider of which this handler is a part
func (h *AlertRuleHandler) APIVersion() string {
	return h.Provider.APIVersion()
}

// GetExtension returns the file name extension for an alert rule
func (h *AlertRuleHandler) GetExtension() string {
	return "json"
}

const (
	alertruleGlob    = "alertrules/*/alertrule-*"
	alertrulePattern = "alertrules/%s/alertrule-%s.%s"
)

// FindResourceFiles identifies files within a directory that this handler can process
func (h *AlertRuleHandler) FindResourceFiles(dir string) ([]string, error) {
	path := filepath.Join(dir, alertruleGlob)
	return filepath.Glob(path)
}

// ResourceFilePath returns the location on disk where a resource should be updated
func (h *AlertRuleHandler) ResourceFilePath(resource grizzly.Resource, filetype string) string {
	return fmt.Sprintf(alertrulePattern, resource.GetMetadata("folder"), resource.Name(), filetype)
}

// Parse parses a manifest object into a struct for this resource type
func (h *AlertRuleHandler) Parse(m manifest.Manifest) (grizzly.Resources, error) {
	resource := grizzly.Resource(m)
	resource.SetSpecString("uid", resource.Name())
	if !resource.HasMetadata("folder") {
		resource.SetMetadata("folder", generalFolderUID)
	}
	return grizzly.Resources{resource}, nil
}

// Unprepare removes unnecessary elements from a remote resource ready for presentation/comparison
func (h *AlertRuleHandler) Unprepare(resource grizzly.Resource) *grizzly.Resource {
	return &resource
}

// Prepare gets a resource ready for dispatch to the remote endpoint
func (h *AlertRuleHandler) Prepare(existing, resource grizzly.Resource) *grizzly.Resource {
	return &resource
}

// GetUID returns the UID for a resource
func (h *AlertRuleHandler) GetUID(resource grizzly.Resource) (string, error) {
	return resource.Name(), nil
}

// GetByUID retrieves JSON for a resource from an endpoint, by UID
func (h *AlertRuleHandler) GetByUID(UID string) (*grizzly.Resource, error) {
	resource, err := getRemoteAlertRule(h.Provider.client, UID)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving alert rule %s: %v", UID, err)
	}

	return resource, nil
}

// GetRemote retrieves a alert rule as a resource
func (h *AlertRuleHandler) GetRemote(resource grizzly.Resource) (*grizzly.Resource, error) {
	return getRemoteAlertRule(h.Provider.client, resource.Name())
}

// ListRemote retrieves as list of UIDs of all remote resources
func (h *AlertRuleHandler) ListRemote() ([]string, error) {
	return getRemoteAlertRuleList(h.Provider.client)
}

// Add pushes a new AlertRule to Grafana via the API
func (h *AlertRuleHandler) Add(resource grizzly.Resource) error {
	return postAlertRule(h.Provider.client, resource)
}

// Update pushes a AlertRule to Grafana via the API
func (h *AlertRuleHandler) Update(existing, resource grizzly.Resource) error {
	return putAlertRule(h.Provider.client, resource)
}
