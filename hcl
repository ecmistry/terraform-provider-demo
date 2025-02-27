hclCopyresource "gravitee_api" "example" {
  name               = "example-api"
  description        = "This is an example API"
  api_version        = "1.0"
  definition_version = "V4"
  type               = "PROXY"

  listeners {
    type = "HTTP"
    paths {
      path = "/example"
    }
    entrypoints {
      type = "http-proxy"
    }
  }

  endpoint_groups {
    name = "default-group"
    type = "http-proxy"
    endpoints {
      name                  = "default"
      type                  = "http-proxy"
      inherit_configuration = false
      configuration = {
        target = "https://api.example.com"
      }
    }
  }

  analytics {
    enabled = true
    logging {
      mode {
        entrypoint = true
        endpoint   = true
      }
      phase {
        request  = true
        response = true
      }
      content {
        headers = true
        payload = true
      }
    }
  }

  flows {
    name      = "default-flow"
    selectors = []
    enabled   = true
  }

  auto_start = true
}
gravitee_plan
Manages a Gravitee API plan.
hclCopyresource "gravitee_plan" "example" {
  api_id             = gravitee_api.example.id
  name               = "keyless-plan"
  description        = "A keyless plan for testing"
  definition_version = "V4"
  security_type      = "KEY_LESS"
  mode               = "STANDARD"
  auto_publish       = true
}
gravitee_subscription
Manages a Gravitee API subscription.
hclCopyresource "gravitee_subscription" "example" {
  api_id         = gravitee_api.example.id
  application_id = "application-id"
  plan_id        = gravitee_plan.example.id
  auto_validate  = true

  consumer_configuration {
    entrypoint_id = "webhook"
    channel       = "/channel1"
    entrypoint_configuration {
      callback_url = "https://webhook.example.com/callback"
      headers {
        name  = "X-Custom-Header"
        value = "custom-value"
      }
    }
  }

  metadata = {
    "feature" = "example"
  }
}
Data Sources
gravitee_api
Retrieves information about a specific Gravitee API.
hclCopydata "gravitee_api" "example" {
  id = "api-id"
}

output "api_name" {
  value = data.gravitee_api.example.name
}
Copy
Let's also provide some example usage patterns to demonstrate how to use the provider:

```markdown
## Example: Creating an HTTP Proxy API

This example demonstrates creating a simple HTTP proxy API in Gravitee with Terraform.

```hcl
provider "gravitee" {
  management_url = "https://apim.example.com"
  username       = "admin"
  password       = "password"
}

resource "gravitee_api" "http_proxy" {
  name               = "example-proxy"
  description        = "HTTP Proxy API Example"
  api_version        = "1.0"
  definition_version = "V4"
  type               = "PROXY"

  listeners {
    type = "HTTP"
    paths {
      path = "/example-proxy"
    }
    entrypoints {
      type = "http-proxy"
    }
  }

  endpoint_groups {
    name = "default-group"
    type = "http-proxy"
    endpoints {
      name                  = "default"
      type                  = "http-proxy"
      weight                = 1
      inherit_configuration = false
      configuration = {
        target = "https://api.example.com/echo"
      }
    }
  }

  analytics {
    enabled = true
    logging {
      mode {
        entrypoint = true
        endpoint   = true
      }
      phase {
        request  = true
        response = true
      }
      content {
        headers = true
        payload = true
      }
    }
  }

  auto_start = true
}

resource "gravitee_plan" "keyless" {
  api_id             = gravitee_api.http_proxy.id
  name               = "Keyless"
  description        = "Keyless plan for public access"
  definition_version = "V4"
  security_type      = "KEY_LESS"
  mode               = "STANDARD"
  auto_publish       = true
}
Example: Event Consumption API with Webhook
This example demonstrates creating an API for event consumption using webhooks.
hclCopyresource "gravitee_api" "event_consumption" {
  name               = "event-consumption"
  description        = "Event Consumption with Webhook"
  api_version        = "1.0"
  definition_version = "V4"
  type               = "MESSAGE"

  listeners {
    type = "SUBSCRIPTION"
    entrypoints {
      type = "webhook"
    }
  }

  endpoint_groups {
    name = "default-group"
    type = "kafka"
    endpoints {
      name                  = "default"
      type                  = "kafka"
      weight                = 1
      inherit_configuration = false
      configuration = {
        bootstrapServers = "kafka:9092"
      }
      shared_configuration_override = {
        "consumer.enabled"           = "true"
        "consumer.topics"            = "demo"
        "consumer.autoOffsetReset"   = "earliest"
      }
    }
  }

  flows {
    name      = "message-filtering"
    selectors = []
    subscribe {
      name        = "Filter messages"
      description = "Apply filter to incoming messages"
      enabled     = true
      policy      = "message-filtering"
      configuration = {
        filter = "{#jsonPath(#message.content, '$.feature') == #subscription.metadata.feature}"
      }
    }
    enabled = true
  }

  auto_start = true
}

resource "gravitee_plan" "subscription" {
  api_id             = gravitee_api.event_consumption.id
  name               = "Subscription Plan"
  description        = "Plan for webhook subscriptions"
  definition_version = "V4"
  security_type      = "subscription"
  mode               = "PUSH"
  auto_publish       = true
}

resource "gravitee_subscription" "webhook" {
  api_id         = gravitee_api.event_consumption.id
  application_id = "application-id"
  plan_id        = gravitee_plan.subscription.id
  auto_validate  = true

  consumer_configuration {
    entrypoint_id = "webhook"
    channel       = "/channel1"
    entrypoint_configuration {
      callback_url = "https://webhook.example.com/callback"
      headers {
        name  = "X-Custom-Header"
        value = "custom-value"
      }
    }
  }

  metadata = {
    "feature" = "demo"
  }
}
Copy
Finally, let's create a Makefile to simplify the build process:

```makefile
# Makefile for building the Gravitee Terraform Provider

NAMESPACE=gravitee
NAME=gravitee
VERSION=0.1.0
OS_ARCH=darwin_amd64

default: build

.PHONY: build
build:
	go build -o terraform-provider-${NAME}_v${VERSION}

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp terraform-provider-${NAME}_v${VERSION} ~/.terraform.d/plugins/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: test
test:
	go test -v ./...

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v ./...

.PHONY: clean
clean:
	rm -f terraform-provider-${NAME}_v${VERSION}
Let's also create a Go module initialization file:
goCopy// go.mod
module github.com/gravitee-io/terraform-provider-gravitee

go 1.17

require (
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.19.0
)