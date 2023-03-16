## Welcome to Hatchet!

Hatchet is a continuous integration and deployment (CI/CD) solution for Terraform. It focuses on making it easier to run secure and scalable infrastructure deployment pipelines without having to build those pipelines from scratch. More specifically, it offers:

1. Remote execution of Terraform runs, such as `terraform plan` and `terraform apply`
2. Integrations with Git-based repositories, which can be configured to run pipelines automatically against pull requests and merges.
3. An extensive monitoring integration built using Open Policy Agent, an industry-standard policy as code framework.

## How is this different from other tools?

Hatchet was created due to the lack of open-source and self-hostable alternatives to Terraform Cloud. While Terraform Enterprise is offered as a self-hosted alternative to Terraform Cloud, Hatchet provides a self-hosted setup from Day 1 without requiring an expensive enterprise license.

### Why self-hosted?

For many organizations, it's essential that sensitive data never leaves your internal infrastructure. As a result, while most solutions offer self-hosted runners, Hatchet makes it easy to self-host everything, including your control plane, credentials backend, and runners.

### Why open-source?

While there are many benefits to being open source, one strength of Hatchet is the flexibility it provides for customizing different aspects of your deployment pipelines by extending our open-source repository. For example, if you'd like to load in your cloud credentials from a custom secret storage engine, you can simply write a credential plugin for Hatchet. We aim to make everything customizable -- down to even the theme of the Hatchet dashboard.

**These two features of Hatchet -- being open-source and self-hosted -- work hand-in-hand to give you the best possible solution for infrastructure management.**
