# Pingo
## Simple yaml-configured service to ping other services and notify if anything goes wrong

# TODO
- ICMP check 
- CPU check
- Asserts for the jobs
- Tests
- Dockerhub image
- Github Actions - linter, tests, push image, build binary, load to release

# HOWTO
- Install
> Currently you can only build from source (works only on linux machines)

- Create your configuration
Example in pingo.example.yaml

- Launch
```bash
pingo -config <path_to_your_config> # default - "./pingo.yaml"
```

- Dashboard
Dashboard is available on localhost:9080. You can change port by providing _-port_ flag
```bash
pingo -port 8080
```

- Docker
You can launch pingo from Docker. For now you have to build image first
```bash
docker build -t pingo:dev . \
&& docker run -ti --rm -v /path/to/config:/pingo/config -p 9080:9080 -e TELEGRAM_BOT_TOKEN=1 -e SMTP_PASSWORD=1 pingo:dev
```
