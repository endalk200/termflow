[![Lint](https://github.com/endalk200/termflow/actions/workflows/lint.yaml/badge.svg)](https://github.com/endalk200/termflow/actions/workflows/lint.yaml)

# Termflow

Termflow is a CLI-first application designed to help developers and command-line users bookmark, organize, and efficiently use terminal commands.
With Termflow, you can save and categorize your most-used commands by tags, making them easy to find and use when needed.
You no longer need to remember complex terminal commands—you simply save them in Termflow and retrieve them instantly.

## Motivation

Working with the terminal is powerful, but remembering all those commands can become overwhelming.
Whether you're jumping between projects or using multiple machines, finding the right command can break your flow.
That's where Termflow comes in. It's built to save you time and make your terminal experience smoother.
By allowing you to organize your commands with tags and access them from anywhere with cloud sync, Termflow keeps your workflow uninterrupted.
You can focus on getting things done, rather than searching for the right command.

## Repository Structure

This repository contains all the components that make up Termflow, divided into distinct directories for easy navigation and contribution.

### cli/ (Command-Line Interface)

- Description: This directory houses the Termflow CLI application, written in Go. It allows users to bookmark, organize, and retrieve terminal commands locally. The CLI can also sync with the Termflow webapp for cloud-based management of commands.
- Tech: Go, Viper (for configuration), SQLite3 (for local storage), SQLC (for Go structs and queries), Goose (for database migrations)

### api/ (Backend API)

- Description: The API, also written in Go, powers Termflow’s webapp and enables syncing between multiple devices. It handles user authentication, storage, and retrieval of commands from the cloud.
- Tech: Go, Gin (HTTP web framework), SQLC, PostgreSQL, RabbitMQ (for messaging)

### infra/ (Infrastructure)

- Description: The infrastructure code for Termflow, including Ansible playbooks for configuration management and Terraform scripts to provision cloud resources for the Termflow API and webapp.
- Tech: Ansible, Terraform, AWS (cloud provider)

### web/ (Webapp)

- Description: The web interface for Termflow, where users can manage, edit, and sync their terminal commands across devices. The webapp integrates with the API to provide a smooth experience for command organization and synchronization.
- Tech: Next.js, React, TypeScript, Tailwind CSS

## Tech Stack

- Go: Backend and CLI development
- Viper: Configuration management for the CLI
- SQLite3: Local command storage for the CLI
- SQLC: Database access layer for both the CLI and API
- Goose: Database migrations
- PostgreSQL: Cloud-based command storage
- Ansible: Configuration management
- Terraform: Infrastructure provisioning
- AWS: Cloud platform
- Next.js: Webapp frontend framework
- React: Frontend library
- TypeScript: Typed JavaScript for web development
- Tailwind CSS: Utility-first CSS framework
