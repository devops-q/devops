# Deployment of itu-minitwit-golang

## Terraform deployment

We use terraform to initialize the infrastructure. So far the infrastructure includes:

* VM (droplet) on DigitalOcean running ubuntu with docker installed
* Static IP address (Floating IP) on DigitalOcean

### How to deploy:

#### Ensure you have the following prerequisites:

* Terraform installed
* DigitalOcean SSH key added to your account
* DigitalOcean API token

#### Steps:

1. Cd into the terraform directory:

```shell
cd deployments/terraform
```

2. Terraform initialize:

```shell
terraform init
```

3. Run Terraform plan (you will be prompted to enter your DigitalOcean API token):

```shell
terraform plan
```

4. Apply the Terraform plan:

```shell
terraform apply
```