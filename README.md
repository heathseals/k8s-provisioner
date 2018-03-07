# K8S-Provisioner
This script does two things right now:
1. create a namespace
1. add a rolebinding to that namespace
    1. The rolebinding is the "admin" clusterrolebinding

This allows for quick provisioning of new users to their own dedicated namespaces.  In the future it may include additional functionality such as adding an entire ldap group at once to a namespace.

The script expects you to have a working ~/.kube/config, which you can verify by checking the output of `kubectl get nodes`.


## Install
Install glide (either `curl https://glide.sh/get | sh` or `brew install glide`), and run `glide install`

## Use
    go run provision.go -n namespace_to_create -u username

If using GKE, username will likely be in username@domain.tld format.

## Why
This script is a simple golang program that I'm using as an excuse to learn go, and work with the kuberentes go-client library.