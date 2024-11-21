# ecs-cli
Simple kubectl like CLI tool for AWS Elastic Container Service.

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
```