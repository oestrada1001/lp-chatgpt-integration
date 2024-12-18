package chatgpt

import (
	"context"
	_ "context"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
	"log"
	"oestrada1001/lp-chatgpt-integration/models"
	"oestrada1001/lp-chatgpt-integration/services"
)

type FunctionResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

func JobAssistant(jobOpportunity models.JobOpportunity) {
	assistantId := "asst_eHRybtmHVjcrz4Keik1IGqA1"
	client := openai.NewClient(
		option.WithHeader("OpenAI-Beta", "assistants=v2"),
		option.WithAPIKey(""),
	)
	ctx := context.Background()

	thread, err := client.Beta.Threads.New(ctx, openai.BetaThreadNewParams{
		Messages: openai.F([]openai.BetaThreadNewParamsMessage{
			{
				Content: openai.F([]openai.MessageContentPartParamUnion{
					openai.TextContentBlockParam{
						Text: openai.String("Job Oppotunity Id: 3, Job Title: " + jobOpportunity.Title + ", Job Description: " + jobOpportunity.Description + ", Company Name: " + jobOpportunity.CompanyName),
						Type: openai.F(openai.TextContentBlockParamTypeText),
					},
				}),
				Role: openai.F(openai.BetaThreadNewParamsMessagesRoleUser),
			},
		}),
	})
	if err != nil {
		panic(err)
	}

	println("Created thread with id", thread.ID)
	run, err := client.Beta.Threads.Runs.NewAndPoll(ctx, thread.ID, openai.BetaThreadRunNewParams{
		AssistantID:            openai.F(assistantId),
		AdditionalInstructions: openai.String("Process the job opportunity and call functions as needed to complete hard skill processing."),
	}, 0)

	if err != nil {
		panic(err.Error())
	}

	if run.Status == openai.RunStatusRequiresAction {
		actions := run.RequiredAction.SubmitToolOutputs.ToolCalls
		fmt.Print(actions)

		for _, action := range actions {
			log.Printf("Processing required action: %s", action.Type)

			if action.Type == "function" {
				functionName := action.Function.Name
				functionArgs := action.Function.Arguments // {hard_skill_types: {{label: '', value: '', description: ''}}}
				fmt.Printf("Function name: %s, function args: %s\n", functionName, functionArgs)

				var response string
				switch functionName {
				case "create_or_get_hard_skill_types":
					var wrapper struct {
						HardSkillTypes []models.HardSkillType `json:"hard_skill_types"`
					}
					err := json.Unmarshal([]byte(functionArgs), &wrapper)
					if err != nil {
						fmt.Println(err)
						continue
					}
					response, err = services.CreateOrGetHardSkillTypes(wrapper.HardSkillTypes)
					if err != nil {
						fmt.Println(err)
						continue
					}
				case "create_or_get_hard_skills":
				case "create_or_get_proficiency_levels":
				case "create_or_get_hard_skill_contexts":
				case "create_hard_skillables":
				}
				fmt.Println(response)
			}
		}
	}
	if run.Status == openai.RunStatusCompleted {
		messages, err := client.Beta.Threads.Messages.List(ctx, thread.ID, openai.BetaThreadMessageListParams{})

		if err != nil {
			panic(err.Error())
		}

		for _, data := range messages.Data {
			for _, content := range data.Content {
				println(content.Text.Value)
			}
		}
	}
}
