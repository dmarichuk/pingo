# Pingo
## Simple yaml-configured service to ping other services and notify if anything goes wrong

# TODO
- Job states dumps to SQLite
- Graphs on html
- ICMP check 
- Change Launch interface to create messages with current job data
- Add expected values to service-ping
- Container health checker ???

# HOWTO
- Install
> Currently you can only build from source (works only on linux machines)

- Create your configuration
Example in pingo.example.yaml

- Launch
```bash
pingo -config <path_to_your_config> # default - "./pingo.yaml"
```
