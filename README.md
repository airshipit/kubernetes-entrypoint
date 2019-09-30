# Kubernetes Entrypoint

[![Build Status](https://api.travis-ci.org/stackanetes/kubernetes-entrypoint.svg?branch=master "Build Status")](https://travis-ci.org/stackanetes/kubernetes-entrypoint)
[![Container Repository on Quay](https://quay.io/repository/stackanetes/kubernetes-entrypoint/status "Container Repository on Quay")](https://quay.io/repository/stackanetes/kubernetes-entrypoint)
[![Go Report Card](https://goreportcard.com/badge/stackanetes/kubernetes-entrypoint "Go Report Card")](https://goreportcard.com/report/stackanetes/kubernetes-entrypoint)


============

Kubernetes-entrypoint enables complex deployments on top of Kubernetes.

## Overview

Kubernetes-entrypoint is meant to be used as a container entrypoint, which means it has to bundled in the container.
Before launching the desired application, the entrypoint verifies and waits for all specified dependencies to be met.

Kubernetes-entrypoint queries the Kubernetes API directly, and each container is self-aware of its dependencies and their states.
Therefore, no centralized orchestration layer is required to manage deployments, and scenarios (such as failure recovery or pod migration) become easy.

## Usage

Kubernetes-entrypoint reads the dependencies out of environment variables passed into a container.
There is only one required environment variable `COMMAND` which specifies a command (arguments delimited by whitespace) which has to be executed when all dependencies are resolved:

`COMMAND="sleep inf"`

## Latest features

Extending functionality of kubernetes-entrypoint by adding an ability to specify dependencies in different namespaces. The new format for writing dependencies is `namespace:name`, with the exception of pod dependencies (which use JSON).
In order to ensure backward compatibility, if `namespace` is omitted, it is assumed that dependencies are running in the same namespace as kubernetes-entrypoint, just like in previous versions.
This feature is not implemented for Container, Config or Socket dependency because the different namespace is irrelevant for those cases.

For instance:
`
DEPENDENCY_SERVICE=mysql:mariadb,keystone-api
`

The new entrypoint will resolve mariadb in the mysql namespace and keystone-api in the same namespace as kubernetes-entrypoint was deployed in.

## Supported types of dependencies

All dependencies are passed as environment variables with the format `DEPENDENCY_<NAME>`, delimited by a colon.
For dependencies to be effective please use [readiness probes](http://kubernetes.io/docs/user-guide/production-pods/#liveness-and-readiness-probes-aka-health-checks) for all containers.

### Service
Checks whether given kubernetes service has at least one endpoint.
Example:

`DEPENDENCY_SERVICE=mariadb,keystone-api`

### Container
Within a pod composed of multiple containers, kubernetes-entrypoint waits for the containers specified by their names to start.
This dependency requires a `POD_NAME` environment variable which can be easily passed through the [downward api](http://kubernetes.io/docs/user-guide/downward-api/).
Example:

`DEPENDENCY_CONTAINER=nova-libvirt,virtlogd`

### Daemonset
Checks if a specified daemonset is already running on the same host
This dependency requires a `POD_NAME` environment variable which can be easily passed through the [downward api](http://kubernetes.io/docs/user-guide/downward-api/).
The `POD_NAME` variable is mandatory and is used to resolve dependencies.
Example:

`DEPENDENCY_DAEMONSET=openvswitch-agent`

A simple example of how to use downward API to get `POD_NAME` can be found [here](https://raw.githubusercontent.com/kubernetes/kubernetes.github.io/master/docs/user-guide/downward-api/dapi-pod.yaml).

### Job
Checks if a given job or set of jobs with matching name and/or labels succeeded at least once.
In order to use labels, `DEPENDENCY_JOBS_JSON` must be used.
DEPENDENCY_JOBS is supported as well for backward compatibility.
Examples:

`DEPENDENCY_JOBS_JSON='[{"namespace": "foo", "name": "nova-init"}, {"labels": {"initializes": "neutron"}}]'`
`DEPENDENCY_JOBS=nova-init,neutron-init'`

### Config
This dependency performs a container level templating of configuration files. It can template an ip address `{{ .IP }}` and hostname `{{ .HOSTNAME }}`.
Templated config has to be stored in an arbitrary directory `/configmaps/<name_of_file>/<name_of_file>`.
This dependency requires an `INTERFACE_NAME` environment variable to know which interface to use to obtain the ip address.
Example:

`DEPENDENCY_CONFIG=/etc/nova/nova.conf`

Kubernetes-entrypoint will look for the configuration file `/configmaps/nova.conf/nova.conf`, template the `{{ .IP }} and {{ .HOSTNAME }}` tags, and then save the file as `/etc/nova/nova.conf`.

### Socket
Checks whether a given file exists and that the container has rights to read it.
Example:

`DEPENDENCY_SOCKET=/var/run/openvswitch/ovs.socket`

### Pod
Checks if at least one pod matching the specified labels is already running, by default anywhere in the cluster, or use `"requireSameNode": true` to require a pod on the same node.
Labels are specified using JSON, as seen in the example below.
This dependency requires a `POD_NAME` env which can be easily passed through the [downward api](http://kubernetes.io/docs/user-guide/downward-api/).
The `POD_NAME` variable is mandatory and is used to resolve dependencies.
Example:

`DEPENDENCY_POD_JSON='[{"namespace": "foo", "labels": {"k1": "v1", "k2": "v2"}}, {"labels": {"k1": "v1", "k2": "v2"}, "requireSameNode": true}]'`

## Image

Build process for image is triggered after each commit.
Can be found [here](https://quay.io/repository/stackanetes/kubernetes-entrypoint?tab=tags), and pulled by executing:
`docker pull quay.io/stackanetes/kubernetes-entrypoint:v0.1.0`

## Examples

[OpenStack Helm](https://github.com/openstack/openstack-helm-infra/blob/master/helm-toolkit/templates/snippets/_kubernetes_entrypoint_init_container.tpl) uses kubernetes-entrypoint to manage dependencies when deploying OpenStack on Kubernetes.
