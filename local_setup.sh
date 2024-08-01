
# Step 1 - install Terraform

wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install terraform
terraform --version


# Step 2 - install Go

ARCH=arm64  # TODO: change this to your architecture
wget https://go.dev/dl/go1.22.5.linux-$ARCH.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.5.linux-$ARCH.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile
go version

# Step 3 - Configure local TF provider dev override

echo "
provider_installation {
    dev_overrides {
        \"hashicorp.com/edu/theta\" = \"/home/$USER/go/bin\"
    }

    direct {}
}" >> ~/.terraformrc

# Step 4 - Set up the repo

git clone https://github.com/igorperic17/thetaform.git
cd thetaform
echo "
email    = "your_theta_account_email"       # TODO: change these to your credentials
password = "your_theta_account_password"
hf_token = "your_hf_token"
" >> ./local.tfvars

# Step 5 - Build the provider
go mod tidy
go install .

# Step 6 - Provision resources on Theta EdgeCloud!
terraform apply --var-file=local.tfvars

# Edit 4_notebook.tf to make edits to the deployment if you wish

