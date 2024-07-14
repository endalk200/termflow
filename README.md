# Termflow

Termflow is a versatile CLI application with an optional cloud storage feature. It allows users to save and manage terminal commands and workflows efficiently. Termflow offers the following major functionalities:

1. **Save Terminal Commands Locally**: Easily save terminal commands with a title and description for quick access and reference. This feature helps you organize and document your commonly used commands.

2. **Save and Execute Terminal Workflows Locally**: Create and store sequences of terminal commands as workflows. You can execute these workflows to automate repetitive tasks, saving time and reducing errors.

3. **Sync Commands and Workflows with Cloud**: Seamlessly sync your locally saved commands and workflows with a cloud-hosted API server. This ensures your commands and workflows are backed up and accessible from any device.

## Directory Structure

```
.
├── .github/
│   └── workflows
├── infrastructure/
│   └── azure/
│       └── main.tf
├── web/
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── cli/
│   ├── go.mod
│   └── go.sum
└── README.md
```

- **.github/workflows/**: Contains GitHub Actions workflows for CI/CD.
- **infrastructure/azure/**: Infrastructure as Code (IaC) using Terraform files for deploying the cloud components on Azure.
- **web/**: Source code for the API server, including dependencies.
- **cli/**: Source code for the CLI application, including dependencies.
- **README.md**: This file, providing an overview of the project.
