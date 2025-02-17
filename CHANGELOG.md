# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-02-17

### Added
- add dagger ci
- add dead letter queue
- move schedules to registered workers
- prune old workers
- deregister worker on close
- add worker distributor and model registry
- enable worker process
- add migration
- add executor (#3)
  Adds an executor which can process and dispatch events to a set of workers.
  Co-authored-by: kjuulh <contact@kjuulh.io>
  Co-committed-by: kjuulh <contact@kjuulh.io>
- add basic main.go
- add default

### Fixed
- *(deps)* update module github.com/spf13/cobra to v1.9.1
- *(deps)* update module github.com/spf13/cobra to v1.9.0
- *(deps)* update module golang.org/x/sync to v0.11.0
- *(deps)* update module github.com/golang-migrate/migrate/v4 to v4.18.2
- ci 
- *(deps)* update all dependencies
- use dedicated connection for scheduler process
- *(deps)* update module gitlab.com/greyxor/slogor to v1.6.1
- orbis demo

### Other
- add more specific log for when leader is acquired
- add basic leader election on top of postgres

- add orbis demo

- add basic scheduler

- add utility scripts
- add utility scripts
- add logger
- add cmd
