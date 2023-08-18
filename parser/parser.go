package parser

import (
	"log"
	"os"
	"pingo/alert"
	j "pingo/job"
	"time"

	"gopkg.in/yaml.v2"
)

var YamlConfig Config

type Config struct {
	Settings Settings                          `yaml:"settings"`
	Jobs     map[string]map[string]interface{} `yaml:"jobs,flow"`
}

func (c *Config) ReadFromFile(path string) {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err)
	}

	err = yaml.Unmarshal(f, &YamlConfig)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err)
	}
}

func (c *Config) Parse() []j.Job {
	c.Settings.PostParse()

	jobs := make([]j.Job, len(c.Jobs))
	idx := 0
	for k, v := range c.Jobs {
		jobs[idx] = c.ParseJob(k, v)
		idx += 1
	}
	return jobs
}

func (c *Config) ParseJob(name string, jobMap map[string]interface{}) j.Job {
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
	job.Type = jobType.(string)

	switch job.Type {
	case j.ENDPOING_HEALTH:
		endpoint, hasEndpoint := jobMap["endpoint"]
		if !hasEndpoint {
			log.Fatalf("Job %s must include endpoint!", name)
		}
		job.Endpoint = endpoint.(string)
	case j.RAM_USAGE:
		ramThreshold, hasRamThreshold := jobMap["threshold"]
		if !hasRamThreshold {
			log.Fatalf("Job %s must include threshold!", name)
		}
		job.RamThreshold = ramThreshold.(float64)
		if ok := IsFloatInLimits(job.RamThreshold, 0, 1); !ok {
			log.Fatalf("Job %s threshold must be in range of 0 an 1!", name)
		}
	case j.DISK_USAGE:
		diskThreshold, hasDiskThreshold := jobMap["threshold"]
		if !hasDiskThreshold {
			log.Fatalf("Job %s must include threshold!", name)
		}
		job.DiskThreshold = diskThreshold.(float64)
		diskPath, hasDiskPath := jobMap["path"]
		if !hasDiskPath {
			log.Fatalf("Job %s must include path!", name)
		}
		job.DiskPath = diskPath.(string)
	default:
		log.Fatalf("Unknown job type %s", jobType)
	}

	rawOnFailure, hasOnFailure := jobMap["on_failure"]
	if !hasOnFailure {
		log.Fatalf("Job %s must include on failure tasks", name)
	}

	rawOnRecovery, hasOnRecovery := jobMap["on_recovery"]
	if !hasOnRecovery {
		rawOnRecovery = []interface{}{}
	}
	job.OnFailure = c.ParseAlerts(rawOnFailure.([]interface{}))
	job.OnRecovery = c.ParseAlerts(rawOnRecovery.([]interface{}))
	return job
}

func (c *Config) ParseAlerts(rawAlerts []interface{}) []j.Alert {
	parsedAlerts := TransformISliceToStrSlice(rawAlerts)
	var alerts []j.Alert
	for _, t := range parsedAlerts {
		switch t {
		case j.TELEGRAM_ALERT:
			if ok := c.Settings.IsValidForTelegram(); !ok {
				log.Fatalf("Settings must include telegram_bot_token and telegram_chats!")
			}
			for _, chat := range c.Settings.TelegramChats {
				alerts = append(
					alerts,
					alert.NewTelegramAlert(
						c.Settings.TelegramBotToken,
						chat,
					),
				)
			}
		case j.EMAIL_ALERT:
			if ok := c.Settings.IsValidForSMTP(); !ok {
				log.Fatalf("Variables must include smtp_host and recipients!")
			}
			alerts = append(
				alerts,
				alert.NewEmailAlert(
					c.Settings.SmtpIdentity,
					c.Settings.SmtpUsername,
					c.Settings.SmtpPassword,
					c.Settings.SmtpHost,
					c.Settings.SmtpPort,
					c.Settings.SmtpFrom,
					c.Settings.SmtpTo,
				),
			)
		default:
			log.Fatalf("Unknow task type: %s", t)
		}
	}
	return alerts
}
