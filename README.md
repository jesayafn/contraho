# Contraho
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/jesayafn/contraho?sort=semver)

CLI tool for gathering information of projects and applications within Sonarqube server.

## Without compile the code

### Search all projects or applications with password

#### Project

```shell
go run . project --host="http://localhost:9000" --username="admin" --password="admin" --filename="app-list.csv"
```

#### Application

```shell
go run . application --host="http://localhost:9000" --username="admin" --password="admin" --filename="app-list.csv"
```

### Search all projects or applications with token

#### Project

```shell
go run . project --host="http://localhost:9000" --token="squ_3131sa..." --filename="project-list.csv"
```

#### Application

```shell
go run . application --host="http://localhost:9000" --token="squ_3131sa..." --filename="app-list.csv"
```

## Options

### Global options
- `--filename`, filename of the csv output file
- `--token`, token credential will be used. Should be privileged user token.
- `--host`, server url of the Sonarqube server.
- `--username`, username will be used to access Sonarqube server. **Not recommended** and should be privileged user.
- `--password`, password will be used to access Sonarqube server. **Not recommended** and should be privileged user.

### `project`,`proj`, or `p`  argument's option
- `--unlisted-on-app`, list all projects that are not listed on any application. Can not be used with `--listed-on-app` option and onand on Sonarqube Community Edition.
- `--listed-on-app`, list all projects that are listed on any application. Can not be used with `--unlisted-on-app` option and on Sonarqube Community Edition.
- `--app [app1,app2]`, list all projects that are listed on any application defined in with this option. 
