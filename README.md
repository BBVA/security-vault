# Cloudframe Scurity Vault

This project is focused on embedding certificates into containers. 
This is achieved by providing a bridge between the requesting container and a secret's provider.
 
This program is said bridge. So far it only implements Hashicorp's Vault as secret's provider.

## Prerequisites
In case of using the lifecycle scripts, for them to work, some placeholders must be provided in the `scripts/common.sh` file with the docker registry
to use to push the built images and a rancher catalog if using rancher.

## Build
The source can be compiled to work on a host as an executable or as a container to be run in a docker host. 

### As Executable
Run following commands to run unit tests and compile the executable:

```bash
cd test
go test -v

cd ..
CGO_ENABLED=0 go build -v -a
```

### As Container
Use the provided build script in the scripts folder

```bash
cd scripts
./build.sh
```

If the script detects that is in an AWS environment, it will try to push the image to the registry.

## Run
The program can be run as a stand-alone process or as a docker container. It requires access to the Docker daemon socket
and a Vault server running.

### Stand-alone
Set the following environment variables to the desired values:
``` bash
export VAULT_SERVER=http://vault_vault-server_1:8200   # Vault server ip or dns 
export ROLE=cloudframe-dot-wtf                         # Role of the host, in vault's path for the secrets
export TOKEN_PATH=/etc/token                           # Path of the host's token
export SECRET_PATH=/tmp                                # Where to store the secrets in the requesting containers
export PERSISTENCE_PATH=/tmp/                          # Where to store leases of the secrets already given
export SOCKET=/var/run/docker.sock                     # Docker daemon token
```
Run the previously built executable.
```bash
cloudframe-security-vault
```

## Acceptance Tests
There is a suite of acceptance tests. It can be run with the following commands:
```bash
cd scripts
./run-acceptance.sh
```
Running this suite is tied to our private registry for some images so it might not be possible to do so.

## Release
The release script can do several actions:
* Make a tag to this particular version in the related git repository.
* Add the container to the rancher catalog
* Build the API documentation from the apib file, using Aglio.

Again, some of them are tied to our internal repository and might not be available.

Each of the can be invoked through the `release.sh` script and an argument.

### Create Git Tag
```bash
cd scripts
./release.sh push-tag
```

### Publish to Rancher Catalog
```bash
cd scripts
./release.sh publish-rancher-catalog
```

### Publish API Docmuentation 
```bash
cd scripts
./release.sh publish-api-docs
```

## Deploy
Likewise, there is a deploy script that will run the container version. Again, it's tied to our internal repository, so
it might not work for you.

It can be run with the following
```bash
cd scripts
./deploy.sh
```

