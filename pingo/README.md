# Pingo
## Simple yaml-configured service to ping other services and notify if anything goes wrong

# TODO
- Container health checker and server ping command
- Email alert
- Job states dumps
- Graphs on html

# HOWTO
- Install
> Currently you can only build from source (works only on linux machines)

- Create your configuration

```yaml-configured

# common variables for tasks (more like slots)
# for now there are only variables for telegram, that are used within tasks
# in future is probably going to extend
variables: 
  telegram_bot_token: ${{TELEGRAM_BOT_TOKEN}} # variables in ${{...}} are parsed from environment 
  telegram_chat_id: 111111111

# list of jobs
jobs:

  an-important-service: # name of your job
    type: service-ping # type of your job (now only service-ping, ram-usage and disk-usage are available)
    endpoint: http://192.168.0.1:9000/ping # endpoint, that is going to be requested for healthcheck
    interval: 10s # interval for a job. Expect string as for https://pkg.go.dev/time#ParseDuration 
    on_failure: # tasks to execute on failures
      - telegram-alert # send telegram message with bot token and chat from Variables
    on_recovery: # tasks to execute on recovery (change job task from failed to success)
      - telegram-alert

  local-memory-check:
    type: ram-usage # check current ram usage
    threshold: 0.9 # threshold for failing - current ram usage / total ram usage
    interval: 10s
    on_failure:
      - telegram-alert
    on_recovery:
      - telegram-alert

  local-disk-check:
    type: disk-usage
    threshold: 0.8 # threshold for failing - current disk usage / total size
    interval: 10s
    path: / # path for disk usage
    on_failure:
      - telegram-alert
    on_recovery:
      - telegram-alert
```

- Launch
```bash
pingo -config <path_to_your_config (default - "./pingo.yaml")>
```
