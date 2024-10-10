---
sidebar_position: 4
---

# Developer Guide
The following guide will help you in setting up the development environment to work on Vitta.


## Prerequisites
Make sure, before you start any development, following things are installed and available on your system of choice:

- [Git](https://git-scm.com)
- [Golang](https://go.dev/doc/install)
- [Node.js](https://nodejs.org/en)
- [Docker](https://www.docker.com)
- [Golang Migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md)
- [Common Sense](https://en.wikipedia.org/wiki/Common_sense)

## Getting Codebase
Git clone the repository using:
```bash
git clone git@github.com:prabhuomkar/vitta.git
```

## Setup

### Database
Run the following command to spin up database:
```bash
docker compose up -d database # requires golang-migrate to be installed
```
Run the following command to setup database migrations:
```bash
cd api
make migrate-up # requires golang-migrate to be installed
```

### API
Run the following command to run the API server:
```bash
cd api
go mod tidy # requires golang to be installed
make run
```

### UI
TODO(akshay): Add steps to run UI

### Docs
Run the following command to setup docs and run it:
```bash
cd docs
npm install # requires node.js to be installed
npm run start
```
