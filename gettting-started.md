# Getting Started

This is the GitHub Action for Interactive Inputs (IAIP), a Go-base GitHub Action project that empowers users to bring interactivity to their GitHub Action workflows.

> Note, unless specified it's assumed that the any reference codeblocks will be ran from the Go actions' source code's root directory (./src)

## Usage

To test this action locally, you can run the following command:

```sh
env \
  'GITHUB_API_URL=https://api.github.com' \
  'GITHUB_REPOSITORY=blend/repo-that-uses-an-action' \
  "GITHUB_WORKSPACE=$(pwd)" \
  "INPUT_TITLE=Start something exciting, dynamically..." \
  'INPUT_INTERACTIVE=fields:
  - label: overview
    properties:
      type: textarea
      description: Information on what this action does
      defaultValue: "This example is a powerful demonstration of how you can utilize the boasiHQ/interactive-inputs action to tailor the dynamic portal to your specific needs and desired output."
      readOnly: true
  - label: name
    properties:
      display: What is your name?
      type: text
      description: Name of the user
      maxLength: 20
      required: true
  - label: age
    properties:
      display: How old are you?
      type: number
      description: Age of the user
      placeholder: 18
      maxNumber: 20
      minNumber: 1
      required: false
  - label: city
    properties:
      display: What city do you live in?
      type: text
      description: City of the user
      maxLength: 20
      required: false 
  - label: car
    properties:
      display: Favourite Car
      type: select
      description: The name of your favourite car
      choices:
        - Ford
        - Toyota
        - Honda
        - Volvo
        - BMW
        - Mercedes
        - Audi
        - Lexus
        - Tesla
      required: true
  - label: colour
    properties:
      display: What are your favourite colours
      type: multiselect
      choices: 
        ["Red", "Green", "Blue", "Orange", "Purple", "Pink", "Yellow"]
      required: true
  - label: verify
    properties:
      display: Are you sure you want to continue?
      defaultValue: 'false'
      type: boolean
      required: true' \
  'INPUT_NOTIFIER-SLACK-ENABLED=false' \
  'INPUT_NOTIFIER-SLACK-TOKEN=xoxb-secret-token' \
  'INPUT_NOTIFIER-SLACK-CHANNEL=#random' \
  'INPUT_NOTIFIER-SLACK-BOT=' \
  'INPUT_NOTIFIER-DISCORD-ENABLED=false' \
  'INPUT_NOTIFIER-DISCORD-WEBHOOK=secret-webhook' \
  'INPUT_NOTIFIER-DISCORD-USERNAME=' \
  'INPUT_GITHUB-TOKEN=github-secret-token' \
  'INPUT_NGROK-AUTHTOKEN=1234567890' \
  'IAIP_LOCAL_RUN=true' \
  'IAIP_SKIP_CONFIG_PARSE=1' \
  go run main.go
```

> To skip the config parse, set the `IAIP_SKIP_CONFIG_PARSE` environment variable to `1`

### Hot reloading

Install reflex

`go install github.com/cespare/reflex@latest`

> You can find more information in the repo https://github.com/cespare/reflex

Once installed, run the server

```sh
reflex -r '\.(html|go|css|png|svg|ico|js|woff2|woff|ttf|eot)$' -s -- env \
  'GITHUB_API_URL=https://api.github.com' \
  'GITHUB_REPOSITORY=blend/repo-that-uses-an-action' \
  "GITHUB_WORKSPACE=$(pwd)" \
  "INPUT_TITLE=Start something exciting, dynamically..." \
  'INPUT_INTERACTIVE=fields:
  - label: overview
    properties:
      type: textarea
      description: Information on what this action does
      defaultValue: "This example is a powerful demonstration of how you can utilize the boasiHQ/interactive-inputs action to tailor the dynamic portal to your specific needs and desired output."
      readOnly: true
  - label: name
    properties:
      display: What is your name?
      type: text
      description: Name of the user
      maxLength: 20
      required: true
  - label: age
    properties:
      display: How old are you?
      type: number
      description: Age of the user
      placeholder: 18
      maxNumber: 20
      minNumber: 1
      required: false
  - label: city
    properties:
      display: What city do you live in?
      type: text
      description: City of the user
      maxLength: 20
      required: false 
  - label: car
    properties:
      display: Favourite Car
      type: select
      description: The name of your favourite car
      choices:
        - Ford
        - Toyota
        - Honda
        - Volvo
        - BMW
        - Mercedes
        - Audi
        - Lexus
        - Tesla
      required: true
  - label: colour
    properties:
      display: What are your favourite colours
      type: multiselect
      choices: 
        ["Red", "Green", "Blue", "Orange", "Purple", "Pink", "Yellow"]
      required: true
  - label: verify
    properties:
      display: Are you sure you want to continue?
      defaultValue: 'false'
      type: boolean
      required: true' \
  'INPUT_NOTIFIER-SLACK-ENABLED=false' \
  'INPUT_NOTIFIER-SLACK-TOKEN=xoxb-secret-token' \
  'INPUT_NOTIFIER-SLACK-CHANNEL=#random' \
  'INPUT_NOTIFIER-SLACK-BOT=' \
  'INPUT_NOTIFIER-DISCORD-ENABLED=false' \
  'INPUT_NOTIFIER-DISCORD-WEBHOOK=secret-webhook' \
  'INPUT_NOTIFIER-DISCORD-USERNAME=' \
  'INPUT_GITHUB-TOKEN=github-secret-token' \
  'INPUT_NGROK-AUTHTOKEN=1234567890' \
  'IAIP_LOCAL_RUN=true' \
  'IAIP_SKIP_CONFIG_PARSE=1' \
  go run main.go
```


### Formatting code

You can use the following command to format the code:

```sh
go fmt  ./...
```

### Testing code

You can use the following command to format the code:

```sh
go test -v  ./...
```

### Building binary

There are two methods of building the binary. It can be done with native go build command.


```sh
# For Linux AMD64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags=\"-w -s\" -o dist/action-amd64 main.go

# For Linux ARM64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -ldflags=\"-w -s\" -o dist/action-arm64 main.go
```

**or**

It can be done using the package.json script (this is done automatically on push)

```sh
npm run package
```
