# ecs-cli
A kubectl-style command-line interface for AWS Elastic Container Service (ECS) that simplifies cluster management and container operations.

Important: This CLI is designed to interact with existing ECS clusters and services. It does NOT create or manage AWS infrastructure resources. For infrastructure provisioning, please use AWS CLI, AWS CloudFormation, or Terraform. While task deletion is supported, service deletion has many implications (load balancers, auto-scaling groups, task definitions) and is intentionally not supported. 

## Features
- Context-based configuration management similar to kubectl  
- Easy service and task management for existing ECS resources  
- Real-time container logs viewing  
- Service scaling capabilities  
- AWS profile and region support  
- Intuitive command structure

## Scope
### What this CLI does:

- Manage and switch between multiple ECS cluster contexts  
- List and describe existing services and tasks  
- View and follow container logs  
- Scale existing services  

### What this CLI doesn't do:
- Create or delete ECS clusters  
- Create or modify AWS infrastructure  
- Manage IAM roles or permissions  
- Handle service definitions or task definitions  
- Manage Auto Scaling configurations  
- Create or modify Load Balancers  

## Installation 
### Using Binary
```bash
VERSION=v0.0.9
wget -O ecs-cli.tar.gz https://github.com/yogendratamang48/ecs-cli/releases/download/$VERSION/ecs-cli_Linux_x86_64.tar.gz
tar -xzf ecs-cli.tar.gz
chmod +x ecs
mv ecs /usr/local/bin/
ecs version
```
### Building from Source
```bash
git clone https://github.com/yogendratamanga48/ecs-cli.git
cd ecs-cli
go build -o ecs
mv ecs /usr/local/bin/
```
## Configuration
The CLI uses a context-based configuration system similar to kubectl, powered by Viper. Configurations are stored in `$HOME/.ecs/config.yaml` 
## Context Management
setup new context:
```bash
ecs config set-context <my-context> \
    --cluster <my-cluster> \
    --profile <aws-profile> \
    --region <region>
```
Other context operations:
```bash
ecs config get contexts
ecs config use-context <context-name>
ecs config delete-context <context-name>
```
## Usage
```bash
# get service and tasks
ecs get tasks
ecs get services -o json
ecs get tasks -o wide

# delete task
ecs delete task <task-id>

# Describe Service and Tasks
ecs describe service <service-name>
ecs describe task <task-id>

# Show logs
ecs logs <task-id>
ecs logs <task-id> --follow

# scale service
ecs scale service-name --replicas=N
```

## Development
This CLI is built using:
- [Cobra](https://github.com/spf13/cobra) - CLI framework  
- [Viper](https://github.com/spf13/viper) - Configuration management  
- [pflag](https://github.com/spf13/pflag) - Flag parsing  

## Status
Note: This CLI is under active development. Breaking changes may occur.