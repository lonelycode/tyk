# Generated by: tyk-ci/wf-gen
# Generated on: Tue 11 May 04:10:51 UTC 2021

# Generation commands:
# ./pr.zsh -title copy proto -repos tyk -branch releng/python -base releng/python -p
# m4 -E -DxREPO=tyk


data "terraform_remote_state" "integration" {
  backend = "remote"

  config = {
    organization = "Tyk"
    workspaces = {
      name = "base-prod"
    }
  }
}

output "tyk" {
  value = data.terraform_remote_state.integration.outputs.tyk
  description = "ECR creds for tyk repo"
}

output "region" {
  value = data.terraform_remote_state.integration.outputs.region
  description = "Region in which the env is running"
}
