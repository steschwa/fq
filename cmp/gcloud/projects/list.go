package projects

import (
	"encoding/json"

	"github.com/carapace-sh/carapace"
)

type (
	glcoudProject struct {
		Labels         map[string]string `json:"labels"`
		LifecycleState string            `json:"lifecycleState"`
		Name           string            `json:"name"`
		ProjectID      string            `json:"projectId"`
	}
)

func ActionProjects() carapace.Action {
	return carapace.ActionExecCommand("gcloud", "projects", "list", "--format=json")(func(output []byte) carapace.Action {
		var projects []glcoudProject
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
	})
}
