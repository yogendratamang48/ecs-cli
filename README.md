# ecs-cli
**Note: This cli is still under development**  
Simple kubectl like CLI tool for AWS Elastic Container Service.

## Installation 
### Using Binary
```bash
VERSION=v0.0.7
wget -O ecs-cli.tar.gz https://github.com/yogendratamang48/ecs-cli/releases/download/$VERSION/ecs-cli_Linux_x86_64.tar.gz
tar -xzf ecs-cli.tar.gz
chmod +x ecs
mv ecs /usr/local/bin/
ecs version
```

## Supported commands
```
# get help
ecs --help
# Context setup 
ecs config set-context default --cluster default --profile airflow --region us-east-1  

# view contexts
ecs config get-contexts

# get service and tasks
ecs get services
ecs get tasks

# Describe Service and Tasks
ecs describe service <service-name>
ecs describe task <task-id>

# Show logs
ecs logs <task-id>
ecs logs <task-id> --follow

# scale service
ecs scale service-name --replicas=N
```