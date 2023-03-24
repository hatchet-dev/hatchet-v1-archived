## Welcome to Hatchet!

Hatchet is a continuous integration and deployment (CI/CD) solution for Terraform. It focuses on making it easier to run secure and scalable infrastructure deployment pipelines without having to build those pipelines from scratch. More specifically, it offers:

1. Remote execution of Terraform runs, such as `terraform plan` and `terraform apply`
2. Integrations with Git-based repositories, which can be configured to run pipelines automatically against pull requests and merges.
3. An extensive monitoring integration built using Open Policy Agent, an industry-standard policy as code framework.

## How does it work?

Hatchet is a self-hosted solution for Terraform management that runs on your infrastructure. After installing Hatchet, you can deploy infrastructure through either [local deployments](https://docs.hatchet.run/getting-started/modules/local-deployment) or [Github-based deployments](https://docs.hatchet.run/getting-started/modules/github-deployment).

The Hatchet architecture can be grouped into three main components:

- **Hatchet Control Plane**: This encompasses the Hatchet API server, background workers, database, and a few other services. Hatchet uses a custom build of [Temporal](https://temporal.io/) to manage the execution of Terraform runs.
- **Hatchet Worker**: The worker that executes your Terraform runs.
- **Hatchet Client**: The web interface for Hatchet. This is where you view the status of your runs and configure your pipelines.

While there are a lot of moving pieces, the [getting started](https://docs.hatchet.run/getting-started) guides will make it simple to deploy these components. After you've gotten a basic installation up and running, you can start to customize your instance by consulting the [config file references](https://docs.hatchet.run/managing-hatchet/config-file-reference).

## How is this different from other tools?

Hatchet was created due to the lack of open-source, self-hostable, and scalable alternatives to Terraform Cloud. Hatchet focuses on making it as easy as possible to manage your deployment pipelines, while also providing the flexibility to customize your deployment to your specific needs.

### Why self-hosted?

For many organizations, it's essential that sensitive data never leaves your internal infrastructure. As a result, while most solutions offer self-hosted runners, Hatchet makes it easy to self-host everything, including your control plane, credentials backend, and runners.

### Why open-source?

While there are many benefits to being open source, one strength of Hatchet is the flexibility it provides for customizing different aspects of your deployment pipelines by extending our open-source repository. For example, if you'd like to load in your cloud credentials from a custom secret storage engine, you can simply write a credential plugin for Hatchet. We aim to make everything customizable -- down to even the theme of the Hatchet dashboard.

### How are you scalable?

Hatchet is built on top of Temporal, an open-source workflow engine. Temporal is a horizontally scalable system that can run thousands of workflows in parallel. This means that you can run as many Terraform runs as you need with relatively simple configuration management.

**These three features of Hatchet -- being open-source, self-hosted, and focused on scalability -- work hand-in-hand to give you the best possible solution for infrastructure management.**
