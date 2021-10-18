# SKA CICD K8s Job Manager Plugin
SKA Job Manager plugin for specialized hardward in Kubernetes.

This repository hosts the job-manager plugin for kubectl, which allows to list managed devices and jobs using them.

## Instalation
To install the plugin, the [kubectl-job_manager](kubectl-job_manager) executable must be placed on the system's `$PATH`.

Once done, the plugin can be used through `kubectl`:
- `kubectl job-manager`

## Configuration
The plugin only searches for devices currently in use on nodes labelled with `tpm=true`. So be sure to label every node with specialized hardware devices.

The plugin also requires a list of managed devices to be available in order to list only the appropriate jobs and devices.

To configure this, write a file where each line contains the fully qualified name (as used in Kubernetes jobs) of each of the managed devices. Finally, place this file at `$HOME/.kube/job-manager.conf` or `/etc/kubernetes/job-manager.conf` for the plugin to find.

If both locations contain a configuration file, the `$HOME/.kube/job-manager.conf` file takes priority so each user can customize it.

## Available Commands
The plugin has five available commands:

```
$ kubectl job-manager
Kubectl job manager plugin for Kubernetes cluster.

Available Commands:
  list        List information about jobs and devices.
  create      Create and run a job with associated devices.
  delete      Delete an existing job.
  version     Displays the plugin version.
  config      Displays currently loaded configuration file with the list of managed devices.
```

### `config`
The `config` command allows the user to validate the current configuration file in use and its contents (example uses the account `ubuntu` with a ``):

```
$ kubectl job-manager config
Loaded configuration file: /home/ubuntu/.kube/job-manager.conf
---
dummy/dummyDev
---

Configuration file location priority:
  - /home/ubuntu/.kube/job-manager.conf
  - /etc/kubernetes/job-manager.conf
```
This example shows the user `ubuntu` running the `config` command, where [this](samples/job-manager.conf) configuration file was correctly loaded from `$HOME/.kube/job-manager.conf`.

## `version`
The `version` command simply shows the version of the plugin running on the system:

```
$ kubectl job-manager version
0.1.0
```

## `create`

```
$ kubectl job-manager create
Create and run a job with associated devices.

Examples:
  # Run a job using a device on the default namespace.
  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c '["echo", "Hello!"]'

  # Print the request of the job using a device on the default namespace.
  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c '["echo", "Hello!"]' --dry-run

  # Run a job using a device on the default namespace for debugging purposes.
  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c '["echo", "Hello!"]' --debug

  # Run a job that should be deleted after 10 minutes.
  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c '["echo", "Hello!"]' --ttl 600

  # Run a job that is allowed to backoff 3 times.
  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c '["echo", "Hello!"]' --backoff-limit 3

  # Run a job that must complete 2 times.
  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c '["echo", "Hello!"]' --completions 2

Options:
  -j, --job <jobname>             Set the job name.
  -n, --namespace <namespace>     Set the namespace to create the job on.
  -i, --image <image>             Set the image to be run on the job.
  -c, --command <command>         Set the command to be executed on the image.
  -d, --device <device>           Set a device resource request and limit for the job.
      --ttl <seconds>             Time in seconds before the job is terminated. Defaults to 7200.
      --completions <n>           Set the number of times the job must run. Defaults to 1.
      --backoff-limit <n>         Set the number of times the job can backoff. Defaults to 1.
      --dry-run                   Do not send the create request to the server and print it to stdout.
      --debug <seconds>           Sets the job to debug mode, keeping it alive for the duration in seconds and preventing restarts.

Usage:
  kubectl job-manager create -j <jobname> -n <namespace> -i <image> -c <commmand> -d <device> [--ttl <n>] [--completions <n>] [--backoff-limit <n>] [--dry-run] [--debug]
```

Example:
```
$ kubectl job-manager create -j dummyjob -n default -i busybox -c '["sleep", "600"]' -d 'dummy/dummyDev'
job.batch/dummyjob created
```

## `list`
The `list` command can be used for two purposes:
- list all devices currently in use, through the subcommand `list devices`.
- list all jobs currently using a device, through the subcommand `list jobs`.

```
$ kubectl job-manager list
List information about jobs and devices.

IMPORTANT:
  A configuration file with the list of managed devices must be present to list jobs and devices.
  See the `kubectl job-manager config` command for more information.

Examples:
  # List all devices currently in use.
  kubectl job-manager list devices

  # List all jobs currently using a device.
  kubectl job-manager list jobs

Available Commands:
  devices     List all devices currently in use.
  jobs        List all jobs currently using a device.

Usage:
  kubectl job-manager list <command>
```

### `list devices`
The `list devices` command allows to see the list of all managed devices currently in use by the cluster, including the node, namespace and pod. For a device to appear in the list, it must exist in the plugin configuration file.

```
$ kubectl job-manager create -j dummyjob -n default -i busybox -c '["sleep", "600"]' -d dummy/dummyDev
job.batch/dummyjob created
$ kubectl job-manager list devices
                    DEVICE | NODE           | NAMESPACE         | POD
            dummy/dummyDev | minikube       | default           | dummyjob-fnp5f
```

### `list jobs`
The `list jobs` command allows to see the list of all jobs that are using one of the managed devices, including the namespace and the current time to live in seconds for the job to be automatically terminated. For a device to appear in the list, it must exist in the plugin configuration file.
