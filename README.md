# KWoC Backend v2.0
KWoC backend server v2.0 (also) written in Go (but better).

## Table of Contents
- [Development](#development)
  - [Setting Up Locally](#setting-up-locally)
    - [Setting Up Pre-Commit Hooks](#setting-up-pre-commit-hooks)
  - [Building](#building)
  - [File Naming Convention](#file-naming-convention)
  - [Testing](#testing)
- [Project Structure](#project-structure)
  - [Libraries Used](#libraries-used)
  - [File Structure](#file-structure)
  - [Endpoints](#endpoints)
  - [Middleware](#middleware)
  - [Utils](#utils)
  - [Database Models](#database-models)
  - [Command-Line Arguments](#command-line-arguments)
  - [Environment Variables](#environment-variables)
  - [Github OAuth](#github-oauth)

## Development
### Setting Up Locally
- Install [Go](https://go.dev/).
- Clone this repository.
- Run `go mod download` in the cloned repository to download all the dependencies.
- Create a `.env` file to store all the [environment variables](#environment-variables). You can use the `.env.template` file for this.
- Set the `DEV` environment variable to `true`.
- Optionally set up [Github OAuth](#github-oauth) to test the endpoints which require login. (See also: [Endpoints](#endpoints))
- Run `go run cmd/backend.go` to start the server.
- Optionally install [SQLiteStudio](https://sqlitestudio.pl/) or a similar tool to help manage the local database `devDB.db` (sqlite3).
- Optionally install [Postman](https://www.postman.com/) or a similar tool to test the API endpoints.
- Optionally (but **recommended**) [set up pre-commit hooks](#setting-up-pre-commit-hooks).

#### Setting Up Pre-Commit Hooks
- Check if `golangci-lint` is installed, if not, install from [golangci-lint](https://golangci-lint.run/usage/install/).
- Run `git config core.hooksPath .githooks`, see [core.hooksPath](https://git-scm.com/docs/git-config#Documentation/git-config.txt-corehooksPath).

### Building
- Install all the dependencies using `go mod tidy` or `go mod download`.
- Run `go build cmd/backend.go` to build the server executable. The executable file will be named `backend`.

### File Naming Convention
See also [File Structure](#file-structure).

1. Test Files: Tests for a particular `file.go` are placed in a file named `file_test.go` and are placed in the same directory.
2. Model Files: Database models are placed in the `models` directory. Each file corresponds to one table in the database and the name of the file corresponds with the name of the database table.

### Testing
[WIP]

## Project Structure
[WIP]

### Libraries Used
- [gorilla/mux](https://github.com/gorilla/mux): Used for routing.
- [gorm.io/gorm](https://gorm.io): Used for database modelling.
- [joho/godotenv](https://github.com/joho/godotenv): Used for loading environment variables from a `.env` file.
- [rs/zerolog](https://github.com/rs/zerolog): Used for logging information, errors, and warnings.

### File Structure
```
├── cmd
│   ├── backend.go
│   └── ...
├── controllers
│   ├── index.go
│   └── ...
├── server
│   ├── router.go
│   ├── routes.go
│   └── ...
├── models
│   ├── mentors.go
│   └── ...
├── utils
│   ├── database.go
│   └── ...
└── middleware
    ├── logger.go
    └── ...
```

- `cmd` : Contains executable files including the entrypoing of the server (`backend.go`).
- `controllers` : Contains controllers for the endpoints. (See also [Endpoints](#endpoints))
- `server` : Contains the router logic and routes.
- `models` : Contains KWoC database models.
- `utils` : Contains misc functions like database utils.
- `middleware`: Contains all middleware.

### Endpoints
TODO: Find a good, automated way to generate endpoint documentation. (See also: [File Structure](#file-structure))

### Middleware
The `middleware/` directory contains all the middleware used in the server. The middleware are used in `server/routes.go` and `server/router.go` files. The following middelware are exported under the `middelware` package.

All middleware take an `http.Handler` function as an argument and return the wrapper `http.Handler` function.

#### Logger
File: `middleware/logger.go`

Logs information regarding the incoming request and the time taken to handle the request.

#### Login
File: `middleware/login.go`

Handles login/authentication for requests. Requests must included the `Bearer` key in the header with the JWT string for login.

The middleware responds to invalid/unauthenticated requests and only passes valid requests to the inner function. The middleware also adds the logged in user's username to the request's context.

A constant `LOGIN_CTX_USERNAME_KEY` is exported by the middleware. This is the key used to access the login username.
```go
login_username := r.Context().Value(middleware.LOGIN_CTX_USERNAME_KEY).(string)
```

#### Wrap
File `middleware/wrap.go`

Adds an instance of the `App` struct (also defined in the same file) to the requests's context. This struct contains the database (`*gorm.DB`) used by the server.

A constant `APP_CTX_KEY` is exported by the middleware. This is the key used to access the `App`.
```go
app := r.Context().Value(middleware.APP_CTX_KEY).(*middleware.App)
db := app.Db
```

### Utils
The `utils/` directory contains utility functions reused in multiple controllers. (See also: [File Structure](#file-structure))

#### Database
File: `utils/database.go`

Contains utilities related to database handling.
- `GetDB()`: Connects to the database and returns a `*gorm.DB`.
- `MigrateModels()`: Automigrates database models.

#### JSON
File: `utils/json.go`

Contains utilities for handling of JSON body of request and response.
- `DecodeJSONBody()`: Decodes the JSON body of an HTTP request.
- `RespondWithJson()`: Takes a response struct and responds to with the JSON string of the response with the appropriate headers set.

#### JWT
File: `utils/jwt.go`

Contains utilities for handling [JSON Web Tokens (JWTs)](https://jwt.io/).
- `ParseLoginJwtString()`: Parses a JWT string used for login and returns the claims.
- `GenerateLoginJwtString()`: Generates a JWT string used for login using the given claims.

#### Logging
File `utils/log.go`

Contains functions for logging information, warnings, and errors encountered during the handling of an HTTP request in a consistent manner.
- `LogInfo()`: Logs an information message with information regarding the HTTP request that triggered the log.
- `LogWarn()`: Logs a warning with information regarding the HTTP request that triggered the warning.
- `LogWarnAndRespond()`: Logs a warning, same as the `LogWarn()` function and responds to the HTTP request with the warning message.
- `LogErr()`: Logs an error with an error message and information regarding the HTTP request that triggered the error.
- `LogErrAndRespond()`: Logs an error, same as the `LogErr()` function and responds to the HTTP request with the error message.

#### OAuth
File: `utils/oauth.go`

Contains functions for authenticating a user via [Github OAuth](#github-oauth). (See also: [Web application flow - Github OAuth Docs](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#web-application-flow))
- `GetOauthAccessToken()`: Gets an access token from the Github API using the given code generated during the authentication process. (See [this](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#2-users-are-redirected-back-to-your-site-by-github) for more information)
- `GetOauthUserInfo()`: Gets basic information about the user from the Github API using an access token. (See: [this](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#3-use-the-access-token-to-access-the-api) for more information)

### Database Models
The `models/` directory contains database models for the KWoC datatabase tables. (See also: [File Structure](#file-structure) and [File Naming Convention](#file-naming-convention))

#### Mentors
- Table: `mentors`
- Structure:
  - `name` (string): Name of the mentor.
  - `email` (string): Email of the mentor.
  - `username` (string): Username of the mentor.

#### Projects
- Table: `projects`
- Structure:
  - `name` (string): Name of the project.
  - `description` (string): Description for the project.
  - `tags` (string): A list of tags for the project.
  - `comm_channel` (string): A link to the project's official communication channel.
  - `readme_link` (string): A link to the project's README file.
  - `project_status` (bool): Whether the project is approved.
  - `last_pull_time` (int64): The timestamp of merging of the last tracked (for statistics) pull request.
  - `commit_count` (uint): The number of commits contributed to this project during KWoC.
  - `pull_count` (uint): The number of pull requests contributed to this project during KWoC.
  - `lines_added` (uint): The number of lines added to this project during KWoC.
  - `lines_removed` (uint): The number of lines removed from this project during KWoC.
  - `contributors` (string): A list of usernames of students who contributed to the project during KWoC, separated by comma(,).
  - `pulls` (string): A list of links to pull requests contributed to the project during KWoC, separated by comma(,).
  - `mentor_id` (int32): The ID of the project's primary mentor.
  - `secondary_mentor_id` (int32):The ID of the project's secondary mentor.

#### Stats
- Table: `stats`
- Structure:
  - `total_commit_count` (uint): The total number of commits contributed during KWoC.
  - `total_pull_count` (uint): The total number of pull requests contributed during KWoC.
  - `total_lines_added` (uint): The total number of lines added during KWoC.
  - `total_lines_removed` (uint): The total number of lines removed during KWoC.

#### Students
- Table: `students`
- Structure:
  - `name` (string): The name of the KWoC student.
  - `email` (string): The email of the KWoC student.
  - `college` (string): The college in which the KWoC student is enrolled.
  - `username` (string): The username of the KWoC student.
  - `passed_mid_evals` (bool): Whether the student has passed the mid evals.
  - `passed_end_evals` (bool): Whether the student has passed the end evals.
  - `blog_link` (string): A link to the student's final KWoC blog.
  - `commit_count` (uint): The number of commits contributed by the student during KWoC.
  - `pull_count` (uint): The number of pull requests contributed by the student during KWoC.
  - `lines_added` (uint): The number of lines added by the student during KWoC.
  - `lines_removed` (uint): The number of lines removed by the student during KWoC.
  - `languages_used` (string): A list of languages used by the student in KWoC contributions, separated by comma(,).
  - `projects_worked` (string): A list of IDs of projects the student contributed to during KWoC.
  - `pulls` (string): A list of links to pull requests contributed by the student during KWoC, separated by comma(,).


### Command-Line Arguments
The following command-line arguments are accepted by `cmd/backend.go`. `--argument=value`, `--argument value`, `-argument=value`, and `-argument value` are all acceptable formats to pass a value to the command-line argument.
- `envFile`: A file to load environment variables from. (Default: `.env`)

### Environment Variables
Environment variables can be set using a `.env` (See [Command-Line Arguments](#command-line-arguments) to use a different file) file. The following variables are used. (See the `.env.template` file for an example)
- `DEV`: When set to `true`, uses a local sqlite3 database from a file `devDB.db` (or `DEV_DB_PATH` if set).
- `DEV_DB_PATH`: The path to a local sqlite3 database file to be used in development. (Valid when `DEV` is set to `true`) (Default: `devDB.db`) (NOTE: `testDB.db` is used for testing)
- `BACKEND_PORT`: The port on which the backend server runs. (Default: `8080`)
- `DATABASE_USERNAME`: The username used to log into the database. (Valid when `DEV` is not set to `true`)
- `DATABASE_PASSWORD`: The password used to log into the database. (valid when `DEV` is not set to `true`)
- `DATABASE_NAME`: The name of the database to log into. (Valid when `DEV` is not set to `true`)
- `DATABASE_HOST`: The host/url used to log into the database. (Valid when `DEV` is not set to `true`)
- `DATABASE_PORT`: The port used to log into the database. (Valid when `DEV` is not set to `true`)
- `GITHUB_OAUTH_CLIENT_ID`: The client id used for Github OAuth. (See [Github OAuth](#github-oauth))
- `GITHUB_OAUTH_CLIENT_SECRET`: The client secret used for Github OAuth. (See [Github OAuth](#github-oauth))
- `JWT_SECRET_KEY`: The secret key used to create a JWT token. (It can be a randomly generated string)
- `JWT_VALIDITY_TIME`: The amount of time (in hours) for which the generated JWT tokens should be valid.

### Github OAuth
KWoC uses Github [OAuth](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/differences-between-github-apps-and-oauth-apps#about-github-apps-and-oauth-apps) for log in authentication instead of passwords.

To set up the KWoC server, a Github OAuth application has to be created and the client id and secret has to be set in the [environment variables](#environment-variables).

- Follow [this](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app) documentation to create an OAuth app. In the production server, use the `koss-service` account to create the application.
- Set the Homepage URL to `https://kwoc.kossiitkgp.org` and the authorization callback URL to `https://kwoc.kossiitkgp.org/oauth/` in the production application.
- Copy the client ID and the client secret (this should NEVER be made public) and set the `GITHUB_OAUTH_CLIENT_ID` and `GITHUB_OAUTH_CLIENT_SECRET` [environment variables](#environment-variables) to these values.

****
> Please update this documentation if you make changes to the KWoC backend or any other part of KWoC which affects the backend. Future humans will praise you.
