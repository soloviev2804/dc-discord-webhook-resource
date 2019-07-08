# Discord Webhook Resource

A [Concourse](http://concourse-ci.org) resource that sends message discord channel by webhook.

## Getting started
Add the following [Resource Type](https://concourse-ci.org/resource-types.html) to your Concourse pipeline
```yaml
resource_types:
- name: discord-webhook-resource
  type: docker-image
  source:
    repository: v0v87/discord-webhook-resource

resources:
- name: discord-execute-webhook
  type: discord-webhook-resource
  source:
    discord:
      token: "EhCBn_2a0GsxHYJGsY8hGcRHSQV73456456dfhgdfhdfhdfgh"
      webhookId: "597473534634506580"

jobs:
  - name: send-message
    plan:
    - put: discord-execute-webhook
      params:
        messageText: "Hello from Concourse CI! ${BUILD_ID} ${BUILD_NAME} ${BUILD_JOB_NAME}"
```
