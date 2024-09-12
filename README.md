# Interactive Inputs

Interactive Inputs now enables GitHub Actions to support dynamic runtime user inputs for workflows and composite actions. This action allows you to leverage various [input field types](#input-fields-types), including but not limited to [text](#text-input---text), [multiple choice](#multi-select-input---multiselect), and even [uploading files](#multifile-input---multifile), creating dynamic workflows that adapt to your CI/CD needs.

Use the reference **`boasihq/interactive-inputs@v2`** in your workflows to ensure you are using the latest features of this action.

## Motivation and Context

The action was developed to address the issue of GitHub Actions not having a core feature found in other CI tools such as Jenkins - dynamic runtime inputs. With this action, you can now create an dynamic runtime input portal that will prompt the user for input(s) during runtime. Then, you can use that given input value in the workflow via a deterministic output id, some examples can be found in the [examples](#examples) section.

Since the initial release, have continued to add features, i.e., [**Slack**/ **Discord** integration](#sending-notifications-to-slack-discord) for portal notifications, in-runtime file uploads via the [`files`](#file-input---file)/[`multifile`](#multifile-input---multifile) field type and more. We are eager to develop an excellent GitHub Action that allows you to deliver exceptional value with your CI/CD pipelines. If you have any ideas for additional features that could enhance the value of this action, please [**create an issue**](https://github.com/boasiHQ/interactive-inputs/issues/new/choose) with an explanation of your thoughts and use case(s). This will allow us to review and discuss it and then work on implementing it as soon as possible.


</details>


<details>
<summary><h2 id="view-action-inputs"> •• View Action Inputs •• </h2></summary><br>

To see the complete list of input field types and their respective properties supported by the `interactive` input mentioned below, go to the [Input Fields Types](#input-fields-types) section of this README.<br>


<!-- action-docs-inputs source="action.yml" -->
## Inputs

| name | description | required | default |
| --- | --- | --- | --- |
| `title` | <p>The title of the interactive inputs form</p> | `false` | `""` |
| `interactive` | <p>The representation (in yaml) of fields to be displayed</p> | `true` | `fields:   - label: requested-files     properties:       display: Upload desired files       type: multifile       required: true       description: Upload desired files that are to be uploaded to the runner for processing   - label: random-string     properties:       display: Enter a random string       type: text       description: A random string up to 20 characters long       maxLength: 20       required: false   - label: choice     properties:       display: Select a monitoring tool       type: select       description: Available options to chose from       choices: ["datadog", "sentry", "grafana"]       required: true ` |
| `timeout` | <p>The timeout in seconds for the interactive inputs form</p> | `false` | `300` |
| `ngrok-authtoken` | <p>The authtoken for ngrok used to expose the interactive inputs form</p> | `true` | `""` |
| `github-token` | <p>The token used to authenticate with GitHub API</p> | `true` | `${{ github.token }}` |
| `notifier-slack-enabled` | <p>Whether to send a notification to Slack about the status of the interative inputs form</p> | `true` | `false` |
| `notifier-slack-thread-ts` | <p>The timestamp of the message to reply to in the thread</p> | `false` | `""` |
| `notifier-slack-token` | <p>The token used to authenticate with Slack API</p> | `true` | `xoxb-secret-token` |
| `notifier-slack-channel` | <p>The channel to send the notification to</p> | `true` | `#notificaitons` |
| `notifier-slack-bot` | <p>The name of the bot to send the notification as</p> | `false` | `""` |
| `notifier-discord-enabled` | <p>Whether to send a notification to Discord about the status of the interative inputs form</p> | `true` | `false` |
| `notifier-discord-thread-id` | <p>The ID of the Discord thread the message should be sent to</p> | `false` | `""` |
| `notifier-discord-webhook` | <p>The webhook URL used to send the notification(s) to Discord</p> | `true` | `secret-webhook` |
| `notifier-discord-username` | <p>The username to send the notification(s) as</p> | `false` | `""` |
<!-- action-docs-inputs source="action.yml" -->
</details>

### See In Action

Here are some examples of the Interactive Input action... in action 👀😔:

.<img src="./assets/InteractiveInputs_LiveFileInputTypeSupport.gif" alt="Demonstrating Multifile input field type" width="400"/>
<img src="./assets/InteractiveInput_MultiselectToControlFlow.gif" alt="Example of using Interactive Input's multiselect to control flow" width="400"/>
<img src="./assets/InteractiveInput_GeneratedNumberCheckerFlow.gif" alt="Example of using Interactive Input to run a comparison on user input vs generated number" width="400"/>
<img src="./assets/notifier-slack-message.png" alt="Integration - Slack" width="400"/>
<img src="./assets/notifier-discord-message.png" alt="Integration - Discord" width="400"/>


## Getting Started

To get started, there are three main steps:

1. Sign up to NGROK and get your auth token if you do not already have one by [**clicking here**](https://dashboard.ngrok.com/signup)
2. Add this action `boasihq/interactive-inputs@v2` to your workflow file. See below the [various example implemenations](#example) for inspiration.
3. Use the predictable output variables from your interactive input portal to create dynamic workflows.

> Note, this action requires an ARM64 or AMD64 (x86) runner to run i.e. `ubuntu-latest`

### Sending notifications to Slack/ Discord

To send notifications to Slack/ Discord, you will need to do the  following:

1) Create your desired [Slack](#creating-a-slack-integration)/[Discord](#creating-a-discord-integration) integration token or webhook respectively.
2) Ensure you've enabled `notifier-slack-enabled` or `notifier-discord-enabled` respectively.
3) Pass the token or webhook to the action with `notifier-slack-token` or `notifier-discord-webhook`, respectively.

<details>
<summary><h4 id="creating-a-slack-integration">Creating a Slack integration</h4></summary><br>

<img src="./assets/notifier-slack-message.png" alt="Integration - Slack" width="400"/>

To create a Slack integration, follow these steps:
1. Go to the Slack API website at https://api.slack.com/apps.
2. Click on the "Create New App" button.
3. Enter a name for your app and select the workspace where you want to create the integration.
4. Click on the "Create App" button.
5. In the app's settings, navigate to the "OAuth & Permissions" section.
6. Under the "Scopes", add the following permissions:
   - `chat:write`
   - `chat:write. customise`
7. Click on the "Install App to Workspace" button.
8. Follow the instructions to install the app in your workspace.
9. Once the installation is complete, you will receive a "Bot User OAuth Access Token".
10. Copy the token and use it in your GitHub Action declaration (We recommended saving it as a secret in your GitHub repository/organisation).

> Note: If you want to send/create a notification to a specific thread, you will need to pass the thread timestamp to the action with the `notifier-slack-thread-ts` input.
>
> <img src="./assets/notifier-slack-message-to-thread.png" alt="Integration - Discord" width="400"/>

</details>

<details>
<summary><h4 id="creating-a-discord-integration">Creating a Discord integration</h4></summary><br>

<img src="./assets/notifier-discord-message.png" alt="Integration - Discord" width="400"/>

To create a Discord integration, follow these steps:
1. Go to the Discord app or web client at https://discord.com/channels/@me.
2. Right-click on the server you want to create the integration, followed by "Server Settings" and "Integrations".
3. Click on the "Webhook".
4. Click on the "New Webhook" button.
5. Select the new webhook and change the "Name" and target "Channel".
6. Press the "Copy Webhook URL" button and use it in your GitHub Action declaration (We recommended saving it as a secret in your GitHub repository/organisation).

> Note: If you want to send a notification to a specific thread, you will need to pass the thread ID to the action with the `notifier-discord-thread-id input
>
> <img src="./assets/notifier-discord-message-to-thread.png" alt="Integration - Discord" width="400"/>

</details>


## Examples

Here are various examples demonstrating how to use this action in your workflows. Note that this is not an exhaustive list of all the possible use cases. Please share your implementations with us; we will add them to this list!

It is worth noting that these example workflows can be leveraged this action in your workflow file, but to do so, you will need to replace the `NGROK_AUTHTOKEN` secret with your own, and if you wish to send notifications to Slack/ Discord, you will need to replace the `SLACK_TOKEN` and `DISCORD_WEBHOOK` secrets with your own as well as enabling the `notifier-slack-enabled` and `notifier-discord-enabled` inputs. More information on how to do this can be found in the [Getting Started](#getting-started) section.

<details>
<summary><b id="example-workflow-with-all-field-types"> •• example workflow with all field types •• </b></summary><br>

This example workflow demonstrates how the different [input field types](#input-fields-types) supported by this action can be used with their respective properties to build out many possibilities when creating a bespoke user experience when it comes the CI/CD flow for you/ your user's needs.


```yaml
name: '[Example] Interactive Inputs'

on:
  push:

jobs:
  interactive-inputs:
    timeout-minutes: 3
    runs-on: ubuntu-latest
    permissions:
      contents: write
      actions: write
    steps:
      - name: Example Interactive Inputs Step
        id: interactive-inputs
        uses: boasihq/interactive-inputs@v2
        with:
          ngrok-authtoken: ${{ secrets.NGROK_AUTHTOKEN }}
          notifier-slack-enabled: "false"
          notifier-slack-channel: "#notificaitons"
          notifier-slack-token: ${{ secrets.SLACK_TOKEN }}
          notifier-slack-thread-ts: ""
          notifier-discord-enabled: "false"
          notifier-discord-webhook: ${{ secrets.DISCORD_WEBHOOK }}
          notifier-discord-thread-id: ""
          timeout: 160
          title: 'A batch of 10 feature flags have been added to be deployed. Would you like to proceed?'
          interactive: |
            fields:
              - label: overview
                properties:
                  type: textarea
                  description: Information on what this action does
                  defaultValue: "This example is a powerful demonstration of how you can utilize the boasiHQ/interactive-inputs action to tailor the dynamic portal to your specific needs and desired output."
                  readOnly: true
              - label: custom-file
                properties:
                  display: Choose a file
                  type: file
                  description: Select the media you wish to send to the channel
                  acceptedFileTypes:
                    - image/png
                    - video/mp4
              - label: requested-files
                properties:
                  display: Upload desired files
                  type: multifile
                  required: true
                  description: Upload desired files that are to be uploaded to the runner for processing
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
                  maxNumber: 110
                  minNumber: 4
                  required: false
              - label: car
                properties:
                  display: Favourite Car
                  type: select
                  description: The name of your favourite car
                  disableAutoCopySelection: false
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
                    - Other
                  required: true
              - label: colour
                properties:
                  display: What are your favourite colours
                  type: multiselect
                  disableAutoCopySelection: true
                  choices: 
                    ["Red", "Green", "Blue", "Orange", "Purple", "Pink", "Yellow"]
                  required: true
              - label: verify
                properties:
                  display: Are you sure you want to continue?
                  defaultValue: 'false'
                  type: boolean
                  required: true'
              - label: additional-info
                properties:
                  type: textarea
                  display: More information
                  description: Any more information you would like to add

      - name: Display Outputs
        shell: bash
        run: |
          echo "Display Outputs"
          echo -e "\n==============================\n"
          echo "Detected Outputs: ${{join(steps.interactive-inputs.outputs.*, '\n')}}"
          echo -e "\n==============================\n"

      - name: List the uploaded files in the directory
        shell: bash
        run: |
          echo "Display uploaded files"
          echo -e "\n==============================\n"
          ls -la ${{ steps.interactive-inputs.outputs.requested-files }} # Use the label of the multifile/file field as the key to get the uploaded file directory
          echo -e "\n==============================\n"

```  

</details>  

<details>
<summary><b id="example-workflow-multiselect-to-control-flow"> •• multiselect to control flow •• </b></summary><br>

This is the workflow file for the ["multiselect to control flow"](./assets/InteractiveInput_MultiselectToControlFlow.gif) gif [referenced above](#see-in-action). <br>

```yaml
name: '[Example] Multiselect to invoke other flows'

on:
 workflow_dispatch:

jobs:

  interactive-inputs:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      actions: write
    steps:
      - name: Example Interactive Inputs Step
        id: interactive-inputs
        uses: boasiHQ/interactive-inputs@v2
        with:
          timeout: 300
          title: 'We need you to select your desired flow(s) to execute'
          interactive: |
            fields:
              - label: runtime-options
                properties:
                  description: Below are the available flows that can be executed
                  display: Select the flow(s) to execute
                  type: multiselect
                  choices: 
                    ["option-a", "option-b", "option-c"]
                  required: true
              - label: notes
                properties:
                  display: Optional - Additional note(s)
                  type: textarea
                  description: Additional notes on why this selection was made.
          notifier-slack-enabled: "false"
          notifier-slack-channel: "#random"
          notifier-slack-token: ${{ secrets.SLACK_TOKEN }}
          notifier-discord-enabled: "false"
          notifier-discord-webhook: ${{ secrets.DISCORD_WEBHOOK }}
          github-token: ${{ github.token }}
          ngrok-authtoken: ${{ secrets.NGROK_AUTHTOKEN }}
    outputs:
        runtime-options: ${{ steps.interactive-inputs.outputs.runtime-options }}

  option-a:
    runs-on: ubuntu-latest
    needs: ["interactive-inputs"]
    if: ${{ contains( needs.interactive-inputs.outputs.runtime-options, 'option-a' ) }}
    steps:
      - name: Run option-a
        shell: bash
        run: |
          echo "Running flow option-a "
          echo -e "\n==============================\n"
          echo "Detected Selection: ${{join(needs.interactive-inputs.outputs.runtime-options, '\n')}}"
          echo -e "\n==============================\n"

  option-b:
    runs-on: ubuntu-latest
    needs: ["interactive-inputs"]
    if: ${{ contains( needs.interactive-inputs.outputs.runtime-options, 'option-b' ) }}
    steps:
      - name: Run option-b
        shell: bash
        run: |
            echo "Running flow option-b "
            echo -e "\n==============================\n"
            echo "Detected Selection: ${{join(needs.interactive-inputs.outputs.runtime-options, '\n')}}"
            echo -e "\n==============================\n"

  option-c:
      runs-on: ubuntu-latest
      needs: ["interactive-inputs"]
      if: ${{ contains( needs.interactive-inputs.outputs.runtime-options, 'option-c' ) }}
      steps:
        - name: Run option-c
          shell: bash
          run: |
              echo "Running flow option-c "
              echo -e "\n==============================\n"
              echo "Detected Selection: ${{join(needs.interactive-inputs.outputs.runtime-options, '\n')}}"
              echo -e "\n==============================\n"
```

</details>


<details>
<summary><b id="example-workflow-generated-number-checker-flow"> •• generated number checker flow •• </b></summary><br>

This is the workflow file for the ["generated number checker flow"](./assets/InteractiveInput_GeneratedNumberCheckerFlow.gif) gif [referenced above](#see-in-action). <br>

```yaml
name: '[Example] Match generated number flow'

on:
  workflow_dispatch:

jobs:

  interactive-inputs:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      actions: write
    steps:
      - name: Generate random number
        id: get-random-number
        uses: actions/github-script@v7
        with:
          result-encoding: string
          # Script for generating random number sourced from: https://stackoverflow.com/a/7228322
          script: |
            function randomIntFromInterval(min, max) { // min and max included 
              return Math.floor(Math.random() * (max - min + 1) + min);
            }

            return randomIntFromInterval(1, 5);

      - name: Example Interactive Inputs Step
        id: interactive-inputs
        uses: boasiHQ/interactive-inputs@v2
        with:
          timeout: 300
          title: 'We need you to select your desired flow(s) to execute'
          interactive: |
            fields:
              - label: display-number
                properties:
                  description: This is a generated number
                  display: Generated Number
                  type: textarea
                  defaultValue: ${{ steps.get-random-number.outputs.result }}
                  readOnly: true
              - label: selected-number
                properties:
                  description: Match the generated number specified in the text field above
                  display: Enter the number you see above
                  type: number
                  required: true
          notifier-slack-enabled: "false"
          notifier-slack-channel: "#random"
          notifier-slack-token: ${{ secrets.SLACK_TOKEN }}
          notifier-discord-enabled: "false"
          notifier-discord-webhook: ${{ secrets.DISCORD_WEBHOOK }}
          github-token: ${{ github.token }}
          ngrok-authtoken: ${{ secrets.NGROK_AUTHTOKEN }}

      - name: Check if number matches
        id: check-number
        shell: bash
        run: |
          if [[ ${{ steps.interactive-inputs.outputs.selected-number }} -eq ${{ steps.get-random-number.outputs.result }} ]]; then
            echo "The number matches!"
          else
            echo "The number does not match."
          fi

```
</details>



### Key points

When using this action, here are a few key points to note:

- The [`multifile` input field type](#multifile-input---multifile) and the [`file` input field type](#file-input---file) will provide the path to where the uploaded files are located on the runner. You can then use this information in later stages of your workflow.
- To enable the external notifications, you will need to set the `notifier-slack-enabled` or `notifier-discord-enabled` property to `true` in the `with` object. Follow the [**Creating a Slack integration**](#creating-a-slack-integration) or [**Creating a Discord integration**](#creating-a-discord-integration) sections above for more information.
  - To send a message to a thread, you will need to set the `notifier-slack-thread-ts` or `notifier-discord-thread-id` property to the thread timestamp or thread ID, respectively.
- The portal will display fields in the order defined in the `fields` array.
- The `label` property is used to identify the input field and its corresponding output. For example, the `label` property in the `fields` array for **Continue to roll out?** is `continue-roll-out`. This means that the output will be stored in a variable called `continue-roll-out`, which can be accessed using the syntax `${{ steps.interactive-inputs.outputs.continue-roll-out }}`.
- The env `ngrok-authtoken` input is used to open the Ngrok tunnel, which is used to give access to your runner-hosted portal. It is needed to be set in the workflow file.
  - Signing up for NGROK is free and quick; it can be done [here](https://dashboard.ngrok.com/signup).
- There are various [types of input fields](#input-fields-types) that can be used, [**vist the input fields types**](#input-fields-types) in this README for more information.
- The `timeout` property sets the timeout for the interactive input. The workflow will fail if the user does not respond within the timeout period.


## Input Fields Types

The input fields shape the user interface of the interactive input. The input fields are defined in the `fields` property of the `interactive` attribute of the `with` object.

```yaml
      ...
      - name: Example Interactive Inputs Step
        id: interactive-inputs
        uses: boasihq/interactive-inputs@v2
        with:
          ...
          interactive: |
            fields:
              - label: continue-roll-out
                properties:
                  display: Continue to roll out?
                  ...
```

The `fields` property is an array of objects, each object representing a field. Each field type has its properties, some unique to the particular field type. See below the supported field types and their respective properties.

<details>
<summary><h3 id="multifile-input---multifile">Multifile Input - <code>multifile</code></h3></summary><br>


The `multifile` input field is used to allow the user to upload multiple files. It is the most commonly used to allow the 
user to upload a collection of files that can be used in the workflow.

> Note: unlike the other input fields, the `multifile` input field's output points to a direcry (the file cache), not the direct value/input provided by the user.
>
> The `acceptedFileTypes` property can be represented as a hyphenated list of strings or also an array of strings, i.e. `["image/*", "video/*"]`. [Click here](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/file#unique_file_type_specifiers) for more information on file type specifiers.

#### Example

```yaml
fields:
  - label: requested-files # Required
    properties:
      display: Upload desired files # Optional: if not specified, the title for the field will not be displayed on the portal
      type: file # Required
      required: true # Optional: If not added, will default to `false`
      description: Upload desired files that are to be uploaded to the runner for processing # Optional: If not added, "i" won't be on the portal for the field
      acceptedFileTypes: # Optional: A list of file type specifiers that the user will be able to upload (more information on file type specifiers: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/file#unique_file_type_specifiers). If not added or left empty, the user will be able to upload any file.
        - image/png # example accepted file types
        - video/mp4 # example accepted file types
```
</details>

<details>
<summary><h3 id="file-input---file">File Input - <code>file</code></h3></summary><br>


The `file` input field is used to capture a file input from the user. It is the most commonly used to allow the 
user to upload a file that can be used in the workflow.

> Note: Unlike the other input fields, the `file` input field's output points to a direcry (the file cache), not the direct value/input provided by the user.
>
> The `acceptedFileTypes` property can be represented as a hyphenated list of strings or also an array of strings, i.e. `["image/*", "video/*"]`. [Click here](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/file#unique_file_type_specifiers) for more information on file type specifiers.

#### Example

```yaml
fields:
  - label: requested-file # Required
    properties:
      display: Upload desired file # Optional: if not specified, the title for the field will not be displayed on the portal
      type: file # Required
      required: true # Optional: If not added, will default to `false`
      description: Upload desired files that are to be uploaded to the runner for processing # Optional: If not added, "i" won't be on the portal for the field
      acceptedFileTypes: [] # Optional: A list of file type specifiers that the user will be able to upload (more information: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/file#unique_file_type_specifiers). If not added or left empty, the user will be able to upload any file.

```
</details>


<details>
<summary><h3 id="text-input---text">Text Input - <code>text</code></h3></summary><br>



The text input field is used to capture text input from the user. It is the most commonly used input field type.

#### Example

```yaml
fields:
 - label: name # Required
    properties:
      display: Name # Optional: if not specified, the title for the field will not be displayed on the portal
      type: text # Required
      description: The name of the user # Optional: If not added, "i" won't be on the portal for the field
      required: true # Optional: If not added, will default to `false`
      maxLength: 20 # Optional: If not added, the user will not have a limit
      placeholder: Enter your name # Optional: If not added, the placeholder won't be displayed on the portal
      defaultValue: John Doe # Optional: If not added, the default value won't be displayed on the portal
```
</details>


<details>
<summary><h3 id="textarea-input---textarea">Textarea Input - <code>textarea</code></h3></summary><br>

The textarea input field is used to capture or display ( set `readOnly` to `true`) multi-line text input from the user. It is commonly used to capture long text input from the user.

> Note, when set to `readOnly` true, the data will still be stored in the output variable, but the user cannot change the value.

#### Example

```yaml
fields:
 - label: notes # Required
    properties:
      display: Additional note(s)  # Optional
      type: textarea # Required
      description: Additional notes on this decision  # Optional
      required: false  # Optional
      maxLength: 200  # Optional
      placeholder: Enter your notes  # Optional
      defaultValue: This is a note  # Optional
      readOnly: false # Optional: If not added, will default to `false`. If set to `true` the field will be read-only, and the user will not be able to change the value, which can be useful for displaying information to the user. 
```
</details>


<details>
<summary><h3 id="number-input---number">Number Input - <code>number</code></h3></summary><br>


The number input field is used to capture numerical input from the user.

#### Example

```yaml
fields:
 - label: cache-wipe-days # Required
    properties:
      display: Wipe cache data by (days)  # Optional
      type: number # Required
      description: The number of days to wipe cache the data for  # Optional
      required: true  # Optional
      minNumber: 0  # Optional: This is the minimum number that the user can enter
      maxNumber: 17  # Optional: This is the maximum number that the user can enter
      placeholder: Enter the number of days to wipe cache data # Optional
      defaultValue: 14  # Optional: This is the value that will be displayed on the portal and used for the output if the user enters no value
```
</details>

<details>
<summary><h3 id="boolean-input---boolean">Boolean Input - <code>boolean</code></h3></summary><br>


The boolean input field captures a boolean input from the user (`True` or `False`). It is commonly used to determine where the expected output should be `True` or `False` from the user.

#### Example

```yaml
fields:
 - label: use-interactive-inputs # Required
    properties:
      display: Should you use Interactive Inputs? # Optional
      type: boolean # Required
      description: Whether you should use Interactive Inputs in your workflows # Optional
      defaultValue: true # Optional: If not added, neither True nor False will be selected on the portal
```

</details>

<details>
<summary><h3 id="select-input---select">Select Input - <code>select</code></h3></summary><br>


The select input field captures a single selection from a list of options from the user. It is commonly used to capture when you wish to scope the user's choice for a particular set of options.

> Note, the `choices` property can be represented as a hyphenated list of strings (shown in the example below) or also an array of strings, i.e. `["US", "UK", "DE", "FR", "JP"]`.

#### Example

```yaml
fields:
 - label: country-rate-limit # Required
    properties:
      display: Which country should have a limited request rate? # Optional
      type: select # Required
      description: The country that should have requests for unregistered users rate limited # Optional
      required: false # Optional
      disableAutoCopySelection: false # Optional: If set to `true`, the user's selected choice will not be automatically copied to the clipboard.
      choices: # Required: This is the list of options the user can select. It can be generated by a previous step or a static list of options.
        - US
        - UK
        - DE
        - FR
        - JP
```

</details>


<details>
<summary><h3 id="multi-select-input---multiselect">Multi-Select Input - <code>multiselect</code></h3></summary><br>

The multi-select input field captures multiple selections from a list of user options. It is commonly used to capture when you wish to scope the user's selection for a particular set of options.

> Note, the `choices` property can be represented as a hyphenated list of strings (shown in the example below) or also an array of strings, i.e. `["US", "UK", "DE", "FR", "JP"]`.

#### Example

```yaml
fields:
 - label: countries-to-rate-limit # Required
    properties:
      display: Which countries should have a limited request rate? # Optional
      type: multiselect # Required
      description: The countries that should have requests for unregistered users rate limited # Optional
      required: false # Optional
      disableAutoCopySelection: false # Optional: If set to `true`, the user's selected choice will not be automatically copied to the clipboard.
      choices: # Required: This is the list of options the user can select. It can be generated by a previous step or a static list of options.
        - US
        - UK
        - DE
        - FR
        - JP
```
</details>


## 💻 Contributing, 🐛 Reporting Bugs & 💫 Feature Requests

We are currently developing a process to facilitate contributions. Please be patient with us! In the meantime, please create an issue if you would like to request additional features, report any unexpected behaviour, or provide any other feedback.

Soon, we will use issues to gather feedback on bugs, feature requests, and more. When testing new features or bug fixes, we will create pull requests (PRs) and keep them focused on a single feature or bug fix, allowing you to test them.

When expressing interest in a bug, enhancement, PR, or issue, please use the thumbs-up or thumbs-down emoji on the original message rather than adding duplicate comments.

### Developing Locally

If you want to contribute, fix a bug, or play around with this action locally, please follow the instructions outlined in the [**getting started** file](./gettting-started.md).


## Licence

MIT License - see [LICENSE.md](LICENSE.md) for details
