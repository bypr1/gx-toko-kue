## CREATE NEW PROJECT

If you haven't downloaded ["go-create-project.sh"](https://storage.globalxtreme-gateway.net/link/installations/go-create-project.sh) yet, please download it from the following link ["go-create-project.sh"](https://storage.globalxtreme-gateway.net/link/installations/go-create-project.sh) and place it in the directory where you want to install your project.

Navigate to the file where you saved **go-create-project.sh** earlier, and execute the following command.
```shell
cd /path/your-go-create-project-dir

chmod +x go-create-project.sh
```

After that, you can create a new project using **go-create-project.sh** by running the following command.
```shell
./go-create-project.sh your-project-path-and-name
```

####
## Setup Executable
If you want to use go-create-project like commands such as php, composer, redis-server, and others in a generic way, you can follow the steps below:

### Install shc
This tool will allow you to convert shell scripts into binary executable files.
```shell
brew install shc
```

### Convert the shell script "go-create-project.sh" to a binary file
```shell
shc -f go-create-project.sh
```

### Rename and move
Rename the resulting binary file (if necessary) and move it to a location in your system's PATH, so it can be executed from anywhere.
```shell
mv go-create-project.sh.x /usr/local/bin/go-create-project
```

By following these steps, you can execute go-create-project from any location in the terminal, just like other commands.
```shell
go-create-project your-project-path-and-name
```

For documentation, you can read it on [Go-Lang Backend Service](https://www.notion.so/globalxtreme/Go-Lang-Backend-Service-527f335297b8465f838fc2598538dae7?pvs=4), which we will create later!


#### how to install and run the application
```shell
# Install Application
go build -o application main.go

# Root (API)
./application

# Migration
./application xtreme:migration

# Seeder
./application xtreme:seeder

# gRPC
./application xtreme:grpc

# Queue
./application xtreme:queue

# RabbitMQ
./application xtreme:rabbitmq

# Schedule
./application xtreme:schedule

# Generator

# Migration file
./application gen:migration <Name>

# Handler File
./application gen:handler <Name> --type=<web/mobile> --resource

# Model File
./application gen:model <Name> --migration //--migration for autocreate migration 

# Parser File
./application gen:parser <Name> --model

# Custom Commands (Example)
./application dev-test
```
Add **"--dev"** for development mode.
