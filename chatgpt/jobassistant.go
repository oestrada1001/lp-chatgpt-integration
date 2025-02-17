package chatgpt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
	"log"
	"oestrada1001/lp-chatgpt-integration/models"
	"oestrada1001/lp-chatgpt-integration/services"
	"os"
	"sync"
)

type FunctionResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

func JobAssistant(jobOpportunity models.JobOpportunity) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	assistantId := "asst_EduFYodMbZXx2XM3MdWloGrc"
	client := openai.NewClient(
		option.WithHeader("OpenAI-Beta", "assistants=v2"),
		option.WithAPIKey(openaiApiKey),
	)
	ctx := context.Background()

	thread, err := client.Beta.Threads.New(ctx, openai.BetaThreadNewParams{
		Messages: openai.F([]openai.BetaThreadNewParamsMessage{
			{
				Content: openai.F([]openai.MessageContentPartParamUnion{
					openai.TextContentBlockParam{
						Text: openai.String(fmt.Sprintf("Job Oppotunity Id: %d, Job Title: %s, Job Description: %s, Company Name: %s", jobOpportunity.Id, jobOpportunity.Title, jobOpportunity.Description, jobOpportunity.CompanyName)),
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
		AdditionalInstructions: openai.String("Process the job opportunity and call functions as needed to complete hard skill processing. You will need to make multiple sequential calls to the functions because the 'hard_skillables' field is a list of object reference Ids. The first call should execute the 'create_or_get_hard_skill_types', 'create_or_get_proficiency_levels', and 'create_or_get_hard_skill_contexts' functions. The second call should execute the 'create_or_get_hard_skills' function and use the response from the first call associate the 'hard_skill_types' with the 'hard_skills'. The third call should execute the 'create_hard_skillables' function and use the response from the first and second call to provide the list of Ids to associate with the 'hard_skillables'."),
	}, 0)

	if err != nil {
		panic(err.Error())
	}

	if run.Status == openai.RunStatusRequiresAction {

		actions := run.RequiredAction.SubmitToolOutputs.ToolCalls
		fmt.Print(actions)

		var wg sync.WaitGroup
		wg.Add(len(actions))

		var toolOutputs []openai.BetaThreadRunSubmitToolOutputsParamsToolOutput
		for _, action := range actions {
			log.Printf("Processing required action: %s", action.Type)

			if action.Type == "function" {
				actionId := action.ID
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
					var wrapper struct {
						HardSkills []models.HardSkill `json:"hard_skills"`
					}
					err := json.Unmarshal([]byte(functionArgs), &wrapper)
					if err != nil {
						fmt.Println(err)
						continue
					}
					response, err = services.CreateOrGetHardSkill(wrapper.HardSkills)
					if err != nil {
						fmt.Println(err)
						continue
					}
				case "create_or_get_proficiency_levels":
					var wrapper struct {
						ProficiencyLevels []models.ProficiencyLevel `json:"proficiency_levels"`
					}
					err := json.Unmarshal([]byte(functionArgs), &wrapper)
					if err != nil {
						fmt.Println(err)
						continue
					}
					response, err = services.CreateOrGetProficiencyLevels(wrapper.ProficiencyLevels)
					if err != nil {
						fmt.Println(err)
						continue
					}
				case "create_or_get_hard_skill_contexts":
					var wrapper struct {
						HardSkillContexts []models.HardSkillContext `json:"hard_skill_contexts"`
					}
					err := json.Unmarshal([]byte(functionArgs), &wrapper)
					if err != nil {
						fmt.Println(err)
						continue
					}
					response, err = services.CreateOrGetHardSkillContexts(wrapper.HardSkillContexts)
					if err != nil {
						fmt.Println(err)
						continue
					}
				case "create_hard_skillables":
					var wrapper struct {
						HardSkillables []models.HardSkillable `json:"hard_skillables"`
					}
					err := json.Unmarshal([]byte(functionArgs), &wrapper)
					if err != nil {
						fmt.Println(err)
						continue
					}
					response, err = services.CreateOrGetHardSkillables(wrapper.HardSkillables)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
				fmt.Println(response)

				newToolOutputs := []openai.BetaThreadRunSubmitToolOutputsParamsToolOutput{
					{
						Output:     openai.F(response),
						ToolCallID: openai.F(actionId),
					},
				}
				toolOutputs = append(toolOutputs, newToolOutputs...)
			}
			defer wg.Done()
		}
		wg.Wait()
		params := openai.BetaThreadRunSubmitToolOutputsParams{
			ToolOutputs: openai.F(toolOutputs),
		}
		response, submitErr := client.Beta.Threads.Runs.SubmitToolOutputs(ctx, thread.ID, run.ID, params)
		fmt.Println(response)

		if submitErr != nil {
			log.Printf("Error submitting tool outputs.")
		} else {
			log.Printf("Successfully submitted tool outputs.")
			RunJobOpportunityAssociations(client, ctx, assistantId, run, thread)
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

func RunJobOpportunityAssociations(client *openai.Client, ctx context.Context, assistantId string, run *openai.Run, thread *openai.Thread) {
	println("Created thread with id", thread.ID)
	secondRun, err := client.Beta.Threads.Runs.NewAndPoll(ctx, thread.ID, openai.BetaThreadRunNewParams{
		AssistantID:            openai.F(assistantId),
		AdditionalInstructions: openai.String("Provide me with list of `hard_skills` with the associated `hard_skill_type_id` and the `hard_skillables` with the associated `hard_skill_id` and `proficiency_level_id` and `hard_skill_context_id`."),
	}, 0)

	if err != nil {
		panic(err.Error())
	}

	if secondRun.Status == openai.RunStatusRequiresAction {
		actions := secondRun.RequiredAction.SubmitToolOutputs.ToolCalls
		fmt.Print(actions)

		var toolOutputs []openai.BetaThreadRunSubmitToolOutputsParamsToolOutput
		for _, action := range actions {
			log.Printf("Processing required action: %s", action.Type)

			if action.Type == "function" {
				actionId := action.ID
				functionName := action.Function.Name
				functionArgs := action.Function.Arguments
				fmt.Printf("Function name: %s, function args: %s\n", functionName, functionArgs)

				var response string
				switch functionName {
				case "create_hard_skillables":
					var wrapper struct {
						HardSkillables []models.HardSkillable `json:"hard_skillables"`
					}
					err := json.Unmarshal([]byte(functionArgs), &wrapper)
					if err != nil {
						fmt.Println(err)
						continue
					}
					response, err = services.CreateOrGetHardSkillables(wrapper.HardSkillables)
					if err != nil {
						fmt.Println(err)
						continue
					}
				case "create_or_get_hard_skills":
					var wrapper struct {
						HardSkills []models.HardSkill `json:"hard_skills"`
					}
					err := json.Unmarshal([]byte(functionArgs), &wrapper)
					if err != nil {
						fmt.Println(err)
						continue
					}
					response, err = services.CreateOrGetHardSkill(wrapper.HardSkills)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
				fmt.Println(response)
				newToolOutputs := []openai.BetaThreadRunSubmitToolOutputsParamsToolOutput{
					{
						Output:     openai.F(response),
						ToolCallID: openai.F(actionId),
					},
				}
				toolOutputs = append(toolOutputs, newToolOutputs...)
			}
		}

		params := openai.BetaThreadRunSubmitToolOutputsParams{
			ToolOutputs: openai.F(toolOutputs),
		}

		response, submitErr := client.Beta.Threads.Runs.SubmitToolOutputs(ctx, thread.ID, run.ID, params)
		fmt.Println(response)

		if submitErr != nil {
			log.Printf("Error second submitting tool outputs.")
		} else {
			log.Printf("Successfully second submitted tool outputs.")
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
