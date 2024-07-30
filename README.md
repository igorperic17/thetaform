# Thetaform

![Header Image](img/header.png)

Thetaform is a Terraform provider for Theta EdgeCloud. It lets you manage your Theta Network infrastructure as code using Terraform configuration files, bringing standard DevOps enterprise best practices to the Theta distributed compute network.

## Features

- **Infrastructure as Code**: Define and provision your Theta EdgeCloud resources using Terraform configuration files.
- **Automation**: Use Terraform’s automation capabilities to manage, update, and destroy your Theta Network infrastructure.
- **Consistency**: Ensure consistent deployment and management of your Theta infrastructure across different environments.

## Getting Started

### Prerequisites

- **Terraform**: You need Terraform installed on your local machine. [Download Terraform](https://www.terraform.io/downloads).
- **GitHub Secrets**: Add your Thetaform email and password to GitHub secrets for deployment actions.

### Setup

1. **Clone the Repository**

   ```bash
   git clone https://github.com/igorperic17/thetaform.git
   cd thetaform

2. **Configure Terraform**

TODO: add link to local provider dev setup

3. **Configure Theta account credential**

TODO: edit your .tfvars file

4. **Deploy**

Run "terraform apply" and you should see the link to your new IPython Notebook at the end of the successful deployment in your terminal.

Prerequisites
Go: Install Go (version 1.16 or later).
Terraform: Ensure Terraform is installed and available in your PATH.
Setup
Clone the Repository

bash
Copy code
git clone https://github.com/your-username/thetaform.git
cd thetaform
Install Dependencies

Ensure you have Go modules enabled. Run the following command to download necessary dependencies:

bash
Copy code
go mod tidy
Build the Provider

Use the following command to build the provider:

bash
Copy code
go build -o terraform-provider-theta
Test the Provider

Run your tests to ensure everything is working correctly:

bash
Copy code
go test ./...
Link the Provider

You can test your provider locally by creating a symbolic link to it in your Terraform plugin directory:

bash
Copy code
mkdir -p ~/.terraform.d/plugins/registry.example.com/edu/theta/0.1.0/linux_amd64
ln -s $(pwd)/terraform-provider-theta ~/.terraform.d/plugins/registry.example.com/edu/theta/0.1.0/linux_amd64/terraform-provider-theta
Run Terraform

Initialize Terraform to use your local provider:

bash
Copy code
terraform init
Documentation
For more details, visit the Thetaform official webpage.

Contributing
We welcome contributions to Thetaform! If you have any improvements or bug fixes, please submit a pull request or open an issue.

License
Thetaform is licensed under the MIT License.

Contact
For any questions or support, please contact us at support@example.com.