set -xe

terraform destroy -auto-approve
terraform apply -auto-approve

echo http://$(cat terraform.tfstate | jq ".resources[1].instances[0].attributes.network_interface[0].access_config[0].nat_ip"):80/webhook | tr -d '"' | tee url.txt

gcloud container clusters get-credentials bramble --location=us-central1-a
kubectl apply -k test-k8s
