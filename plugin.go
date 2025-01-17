package main

import (
	"encoding/json"
	"time"

	"github.com/AlertFlow/runner/pkg/executions"
	"github.com/AlertFlow/runner/pkg/flows"
	"github.com/AlertFlow/runner/pkg/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CollectFlowDataPlugin struct{}

func (p *CollectFlowDataPlugin) Init() models.Plugin {
	return models.Plugin{
		Name:    "Collect Flow Data",
		Type:    "action",
		Version: "1.0.5",
		Creator: "JustNZ",
	}
}

func (p *CollectFlowDataPlugin) Details() models.PluginDetails {
	params := []models.Param{
		{
			Key:         "FlowID",
			Type:        "text",
			Default:     "00000000-0000-0000-0000-00000000",
			Required:    true,
			Description: "The Flow ID to collect data from",
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
	}

	return models.PluginDetails{
		Action: models.ActionDetails{
			Name:        "Collect Flow Data",
			Description: "Collects Flow data from AlertFlow",
			Icon:        "solar:book-bookmark-broken",
			Type:        "collect_flow_data",
			Category:    "Data",
			Function:    p.Execute,
			IsHidden:    true,
			Params:      json.RawMessage(paramsJSON),
		},
	}
}

func (p *CollectFlowDataPlugin) Execute(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
	flowID := ""

	if action.Params == nil {
		flowID = execution.FlowID
	} else {
		for _, param := range action.Params {
			if param.Key == "FlowID" {
				flowID = param.Value
			}
		}
	}

	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{"Collecting flow data from AlertFlow"},
		Pending:        false,
		Running:        true,
		StartedAt:      time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	flow, err = flows.GetFlowData(flowID)
	if err != nil {
		err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
			ID:             step.ID,
			ActionMessages: []string{"Failed to get Flow Data"},
			Error:          true,
			Finished:       true,
			Running:        false,
			FinishedAt:     time.Now(),
		})
		if err != nil {
			log.Error("Error updating step: ", err)
		}

		return nil, false, false, false, true
	}

	err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Flow Data collected"},
		Running:        false,
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		log.Error("Error updating step: ", err)
	}

	return map[string]interface{}{"flow": flow}, true, false, false, false
}

func (p *CollectFlowDataPlugin) Handle(context *gin.Context) {}

var Plugin CollectFlowDataPlugin
