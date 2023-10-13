# gg - Go Github CLI


`gg` is a command-line tool that allows you to access information about GitHub Pull Requests and Repositories.

## Features

- List Latest GitHub Pull Requests by user.
- List GitHub Pull Requests by repository with optional status.
- List a user's GitHub Repositories (owned, followed, or both).
- Check if a Github Repository has workflows.

## Installation with Docker (Recommended)
To use gg, you can pull the Docker container from the official repository on Docker Hub. Here's how to get started:

1. **Install Docker:** If you don't have Docker installed, you can download it for your specific platform from the [Docker website.](https://www.docker.com/get-started/)
2. **Pull the Docker Image:** Run the following command to pull the gg Docker image from Docker Hub:
```bash
docker pull carolinafsilva/go-github-cli
```
3. ***Run gg Using Docker:*** Once you've pulled the Docker image, you can use the gg command by running it within a Docker container:
```bash
docker run --rm -e GITHUB_ACCESS_TOKEN=<token> carolinafsilva/go-github-cli gg [command]
```
Replace `[command]` with the specific gg command you want to use.

Notice you need to pass your GitHub access token as an environment variable. If you do not have one, you can generate one in [https://github.com/settings/tokens](https://github.com/settings/tokens).

## Examples

### List 10 latest PRs from `<user>`
```bash
gg pr author <user> --size 10
```

### List owned and followed repositories of `<user>`
```bash
gg repo list <user>
```

### List PRs from a repository (`<user>/<repo>`) with status
```bash
gg pr repo <user>/<repo> --status
```

### Check if a repository (`<user>/<repo>`) has workflows
```bash
gg repo workflow <user>/<repo>
```
