package completion

import (
	"encoding/json"
	"time"

	"github.com/carapace-sh/carapace"
)

type gcloudProject struct {
	Labels         map[string]string `json:"labels"`
	LifecycleState string            `json:"lifecycleState"`
	Name           string            `json:"name"`
	ProjectID      string            `json:"projectId"`
}

func ActionGCloudProjects() carapace.Action {
	return carapace.ActionExecCommand("gcloud", "projects", "list", "--format=json")(func(output []byte) carapace.Action {
		var projects []gcloudProject
		err := json.Unmarshal(output, &projects)
		if err != nil {
			return carapace.ActionValues()
		}

		var values []string
		for _, project := range projects {
			if project.LifecycleState != "ACTIVE" {
				continue
			}

			firebaseLabel, ok := project.Labels["firebase"]
			if !ok {
				continue
			}
			if firebaseLabel != "enabled" {
				continue
			}

			values = append(values, project.ProjectID, project.Name)
		}

		return carapace.ActionValuesDescribed(values...)
	}).Cache(time.Second * 5)
}
