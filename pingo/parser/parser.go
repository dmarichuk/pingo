package parser

import (
	"fmt"
	"log"
	"os"
	j "pingo/job"
	"pingo/taskHandlers"
	"time"

	"gopkg.in/yaml.v2"
)

var YamlConfig Config

func init() {
	f, err := os.ReadFile("./pingo.example.yaml") // TODO pass config path as an argument in cli
	if err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}
	fmt.Printf("%v", string(f))

	err = yaml.Unmarshal(f, &YamlConfig)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err)
	}
}

type Config struct {
	Variables Variables                         `yaml:"variables"`
	Jobs      map[string]map[string]interface{} `yaml:"jobs,flow"`
}

func (c *Config) Parse() []j.Job {
	jobs := make([]j.Job, len(c.Jobs))
	idx := 0
	for k, v := range c.Jobs {
		jobs[idx] = ParseJob(k, v, &c.Variables)
		idx += 1
	}
	return jobs

}

func ParseJob(name string, jobMap map[string]interface{}, vars *Variables) j.Job {
	job := j.Job{Name: name}

	interval, hasInterval := jobMap["interval"]
	if !hasInterval {
		log.Fatalf("Job %s must include job interval", name)
	}
	parsedInterval, err := time.ParseDuration(interval.(string))
	if err != nil {
		log.Fatalf("Failed to parse job %s duration: %s", name, err)
	}
	job.Interval = parsedInterval

	jobType, hasType := jobMap["type"]
	if !hasType {
		log.Fatalf("Job %s must include job type!", name)
	}

	switch jobType {
	case j.SERVICE_PING:
		endpoint, hasEndpoint := jobMap["endpoint"]
		if !hasEndpoint {
			log.Fatalf("Job %s must include endpoint!", name)
		}
		job.Type = jobType.(string)
		job.Endpoint = endpoint.(string)
	default:
		log.Fatalf("Unknown job type %s", jobType)
	}

	rawTasks, hasTasks := jobMap["on_failure"]
	if !hasTasks {
		log.Fatalf("Job %s must include on failure tasks", name)
	}

	job.OnFailure = ParseTasks(rawTasks.([]interface{}), &job, vars)

	return job
}

func ParseTasks(rawTasks []interface{}, job *j.Job, vars *Variables) []j.Task {
	parsedTasks := TransformISliceToStrSlice(rawTasks)
	tasks := make([]j.Task, len(parsedTasks))
	for idx, t := range parsedTasks {
		switch t {
		case j.TELEGRAM_ALERT:
			if ok := vars.IsValidForTelegram(); !ok {
				log.Fatalf("Variables must include telegram_bot_token and telegram_chat_id!")
			}
			tasks[idx] = taskHandlers.NewTelegramTask(
				vars.TelegramBotToken,
				vars.TelegramChatID,
				fmt.Sprintf("Service %s is unavailable", job.Name),
			)
		default:
			log.Fatalf("Unknow task type in job %s: %s", job.Name, t)
		}
	}
	return tasks
}
