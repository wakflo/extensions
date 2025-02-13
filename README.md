# Wakflo Automation Framework

## Overview

The Wakflo Automation Framework is a comprehensive library for building, managing, and executing workflows with seamless integrations. It is designed to enable developers to efficiently connect and automate tasks across various platforms and services. With support for numerous integrations and connectors, the framework provides the foundation for scalable, reliable, and robust automation solutions.

This project includes:
- **Workflow engine implementation** for executing workflows and managing their life cycles.
- **Integration connectors** for interacting with third-party platforms and APIs.
- Tools for **logging**, **error handling**, and **state management** to enable developer-friendly automation designs.

## Key Features

- **Rich Integrations Library**: Supports connectors for popular platforms like Google Drive, Trello, Notion, Shopify, Zoom, and more.
- **Error Handling and Resilience**: In-built mechanisms for managing recoverable errors and marking steps as "cancelled" upon failures.
- **Step-by-Step Execution**: Modular architecture for tracking and managing execution states of each step in a workflow.
- **Scalable Context Management**: Provides execution metadata and secrets management for workflows.
- **Extensible Framework**: Built with a pluggable, extensible architecture to allow easy addition of new connectors or extensions.

## Pre-requisites

- **Go version**: Requires Go SDK 1.19+.
- **Dependencies**:
    - [`ent`](https://entgo.io): Database entity modeling.
    - [`zerolog`](https://github.com/rs/zerolog): Logging support.
    - [`github.com/wakflo/go-sdk`](https://github.com/wakflo/go-sdk): Wakflo SDK for core capabilities.
- **Wakflo Extensions Package**: Includes integration connectors (`github.com/wakflo/extensions`).
- **Database**: Ensure the required database configurations are set up if workflow state persistence is needed.

## Getting Started

### Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/wakflo/wakflo.git
    cd wakflo
    ```

2. Install dependencies using `go mod`:
    ```bash
    go mod tidy
    ```

3. Add the `wakflo/extensions` package for seamless integration support:
    ```bash
    go get github.com/wakflo/extensions
    ```

4. Build the project:
    ```bash
    go build ./...
    ```

### Example Usage

Here's a simple example to start a workflow using the `FlowRunner`:

```go
package main

import (
    "context"
    "log"

    "github.com/wakflo/wakflo/internal/data/ent"
    "github.com/wakflo/wakflo/internal/ant"
)

func main() {
    ctx := context.Background()

    // Initialize required components, such as the repository, logger, etc.
    repo := ent.NewRepository()
    logger := ant.NewLogger()
    
    // Create a new FlowRunner instance
    flowRun := &ent.WorkflowRun{} // Replace with actual workflow run data.
    flowVersion := &ent.WorkflowVersion{} // Define your workflow version.
    flowRunner := ant.NewFlowRunner(flowRun, flowVersion, repo, logger, nil, nil)

    // Run the workflow
    err := flowRunner.RunFlow(ctx)
    if err != nil {
        log.Fatalf("Workflow execution failed: %v", err)
    }

    log.Println("Workflow executed successfully!")
}
```

### Step Monitoring & Error Handling

The framework automatically tracks the status of each step and updates their statuses. Use the `MarkCancelledSteps` method to handle dependent workflows when a failure occurs. Example:

```go
err := flowRunner.MarkCancelledSteps(ctx, failedStep)
if err != nil {
    logger.Error(err, "Unable to mark steps as cancelled")
}
```

### Adding an Integration

To add a new connector or integration:
1. Create a new integration within the `wakflo/extensions/internal/integrations` folder.
2. Register the integration in `RegisterIntegrations()` within the `wakflo/extensions` package.
3. Implement your integration logic using the SDK interface provided within `github.com/wakflo/go-sdk`.

#### Example: Adding a Custom CRM Integration

```go
package crm

import (
    "github.com/wakflo/sdk"
)

var Integration = &sdk.Registration{
    Name:    "CustomCRM",
    Version: "1.0.0",
    Actions: []sdk.Action{
        {
            ID:          "fetchContacts",
            Description: "Fetch contacts from Custom CRM",
            Execute: func(ctx context.Context, input sdk.Input, meta sdk.Metadata) (sdk.Output, error) {
                // Add your action execution logic here
                return sdk.Output{
                    "contacts": [...] // Example output
                }, nil
            },
        },
    },
}
```

Then, add it to `RegisterIntegrations`:

```go
func RegisterIntegrations() map[string]sdk.RegistrationMap {
    plugins := []*sdk.Registration{
        crm.Integration, // CustomCRM Integration
    }
    ...
}
```

## Supported Connectors

- **Google Drive**
- **Trello**
- **Shopify**
- **Zoom**
- **Notion**
- **Asana**
- **Monday**
- **OpenAI**
- ...and many more! (See the connectors in `github.com/wakflo/extensions` for the full list)

## Contributing

We welcome contributions to improve and expand the framework! To get started:
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Submit a pull request with a clear description of the changes.

### Contributing Guidelines
- Ensure your changes are thoroughly tested.
- Use Go conventions for coding style.
- Update the documentation for any new features or changes.

## License

This project is licensed under the **Apache License 2.0**. See the [LICENSE](LICENSE) file for details.

---

This README provides a high-level overview of the Wakflo Automation Framework. For more detailed documentation, refer to individual packages or code comments within the library. Feel free to reach out with contributions or questions!