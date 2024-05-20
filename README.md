# Quick Usage Example

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

### Other options

#### Usage with `project`,`proj`,and `p`  argument
- `--unlisted-on-app`, list all projects that are not listed on any application. Can not be used with `--unlisted-on-app` option and onand on Sonarqube Community Edition.
- `--listed-on-app`, list all projects that are listed on any application. Can not be used with `--unlisted-on-app` option and on Sonarqube Community Edition.
- `--app [app1,app2]`, list all projects that are listed on any application defined in with this option. 
