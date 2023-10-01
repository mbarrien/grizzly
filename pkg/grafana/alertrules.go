package grafana

import (
	"encoding/json"
	"errors"
	"net/http"

	gerrors "github.com/go-openapi/errors"
	gclient "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/provisioning"
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/grafana/grizzly/pkg/grizzly"
)

// getRemoteAlertRule retrieves a alertrule object from Grafana
func getRemoteAlertRule(client *gclient.GrafanaHTTPAPI, uid string) (*grizzly.Resource, error) {
	h := AlertRuleHandler{}

	params := provisioning.NewRouteGetAlertRuleParams().WithUID(uid)
	alertruleOk, err := client.Provisioning.RouteGetAlertRule(params, nil)
	if err != nil {
		var gErr *gerrors.Error
		if errors.As(err, gErr) && (*gErr).Code() == http.StatusNotFound {
			return nil, grizzly.ErrNotFound
		}
		return nil, err
	}
	alertrule := alertruleOk.GetPayload()

	// TODO: Turn spec into a real models.ProvisionedAlertRule object
	spec, err := structToMap(alertrule)
	if err != nil {
		return nil, err
	}

	resource := grizzly.NewResource(h.APIVersion(), h.Kind(), uid, spec)
	resource.SetMetadata("folder", *alertrule.FolderUID)
	return &resource, nil
}

func getRemoteAlertRuleList(client *gclient.GrafanaHTTPAPI) ([]string, error) {
	params := provisioning.NewRouteGetAlertRulesParams()
	alertrulesOk, err := client.Provisioning.RouteGetAlertRules(params, nil)
	if err != nil {
		return nil, err
	}
	alertrules := alertrulesOk.GetPayload()

	uids := make([]string, len(alertrules))
	for i, alertrule := range alertrules {
		uids[i] = alertrule.UID
	}
	return uids, nil
}

func postAlertRule(client *gclient.GrafanaHTTPAPI, resource grizzly.Resource) error {
	// TODO: Turn spec into a real models.ProvisionedAlertRule object
	data, err := json.Marshal(resource.Spec())
	if err != nil {
		return err
	}

	var alertrule models.ProvisionedAlertRule
	err = json.Unmarshal(data, &alertrule)
	if err != nil {
		return err
	}
	disableProvenance := "true"
	params := provisioning.NewRoutePostAlertRuleParams().WithBody(&alertrule).WithXDisableProvenance(&disableProvenance)
	_, err = client.Provisioning.RoutePostAlertRule(params, nil)
	return err
}

func putAlertRule(client *gclient.GrafanaHTTPAPI, resource grizzly.Resource) error {
	// TODO: Turn spec into a real models.ProvisionedAlertRule object
	data, err := json.Marshal(resource.Spec())
	if err != nil {
		return err
	}

	var alertrule models.ProvisionedAlertRule
	err = json.Unmarshal(data, &alertrule)
	if err != nil {
		return err
	}
	disableProvenance := "true"
	params := provisioning.NewRoutePutAlertRuleParams().WithBody(&alertrule).WithUID(alertrule.UID).WithXDisableProvenance(&disableProvenance)
	_, err = client.Provisioning.RoutePutAlertRule(params, nil)
	return err
}
