# Git Deploy
Git Deploy is a lightweight HTTP micro-service that allows you to securely deploy code by taking advantage of git and your own hardware infrastructure.

Fundementally different to pushing code, Git Deploy pulls the code from existing git repositories on the target file systems and runs build scripts locally via a simple RESTful API that has support of many webhook providers

## Planned
- Providers:
    - GitHub Webhooks
    - Azure Webhooks
    - GitLab Webhooks
    - API POST

- Deploy Types
    - Script 
    - Composer
    - NPM
    - DOTNET
    - Docker

- Remote Deploy via SSH
    - Deploy to different servers via SSH, and keep your sensitive SSH Keys on your own hardware
- Enviroment Variables
    - Setup enviroment variables before building
- Completed Webhooks
    - Get a notification in Discord or your favourite app when a site deploys
- Deploy History
    - A list of deployment times and reasons (git commits) accessible via the API
- Branch Changing
    - Deployments will be able to switch branches (configurable) based on the API/Webhook provider
- Basic Run Scripts
    - While this project is intended to **only manage deployments and builds**, some simple scripts should be allowed to restart certain applications (npm, dotnet, docker)

## Not Planned

- API to add projects
    - It is all configured in a yaml file and should stay that way
- Project Port Mapping
    - Use apache, k8', or anything else, to handle your projects. This is just a simple deploy script
- Build Actions
    - All build scripts are just that, bash scripts. Advance features such as "actions" can be provided by a local [GitHub Action Runner](https://github.com/actions/runner)
- Pushing builds to servers
    - This is specifically for _git_ projects. Pushing code is a security concern and requires more effort into securing the API with things such as checksums and encryption. The idea of this project is to take a HTTP request and do a `git pull` from that, effectively making it so it never downloads arbitary code.