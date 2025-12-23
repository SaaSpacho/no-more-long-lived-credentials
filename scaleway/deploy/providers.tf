terraform {
  required_version = "~> 1.0"
  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = "~> 2.0"
    }
  }

  backend "s3" {
    bucket       = "no-more-long-lived-credentials-terrafo-statebucket-sylakfqewpdn"
    key          = "scaleway/terraform.tfstate"
    region       = "eu-central-1"
    use_lockfile = true
    encrypt      = true
  }
}

provider "scaleway" {
  zone   = "fr-par-1"
  region = "fr-par"
}
