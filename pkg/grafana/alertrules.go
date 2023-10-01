package grafana

import (
	"encoding/json"
	"strings"

	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/grafana/grizzly/pkg/grizzly"
)

// getRemoteAlertRule retrieves a alertrule object from Grafana
func getRemoteAlertRule(uid string) (*grizzly.Resource, error) {
	h := AlertRuleHandler{}
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	alertrule, err := client.AlertRule(uid)
	if err != nil {
		if strings.HasPrefix(err.Error(), "status: 404") {
			return nil, grizzly.ErrNotFound
		}
	}

	// TODO: Turn spec into a real gapi.AlertRule object
	var spec map[string]interface{}
	data, err := json.Marshal(alertrule)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &spec)
	if err != nil {
		return nil, err
	}

	resource := grizzly.NewResource(h.APIVersion(), h.Kind(), uid, spec)
	resource.SetMetadata("folder", alertrule.FolderUID)
	return &resource, nil
}

func getRemoteAlertRuleList() ([]string, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	alertrules, err := client.AlertRules()
	if err != nil {
		return nil, err
	}

	uids := make([]string, len(alertrules))
	for i, alertrule := range alertrules {
		uids[i] = alertrule.UID
	}
	return uids, nil
}

func postAlertRule(resource grizzly.Resource) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	// TODO: Turn spec into a real gapi.AlertRule object
	data, err := json.Marshal(resource.Spec())
	if err != nil {
		return err
	}

	var alertrule gapi.AlertRule
	err = json.Unmarshal(data, &alertrule)
	if err != nil {
		return err
	}
	_, err = client.NewAlertRule(&alertrule)
	return err
}

func putAlertRule(resource grizzly.Resource) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	// TODO: Turn spec into a real gapi.AlertRule object
	data, err := json.Marshal(resource.Spec())
	if err != nil {
		return err
	}

	var alertrule gapi.AlertRule
	err = json.Unmarshal(data, &alertrule)
	if err != nil {
		return err
	}
	return client.UpdateAlertRule(&alertrule)
}
