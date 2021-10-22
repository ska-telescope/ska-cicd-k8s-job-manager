package main

import "fmt"

func printHelp() {
	fmt.Println("Kubectl job manager plugin for Kubernetes cluster.")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  list        List information about jobs and devices.")
	fmt.Println("  create      Create and run a job with associated devices.")
	fmt.Println("  delete      Delete an existing job.")
	fmt.Println("  version     Displays the plugin version.")
	fmt.Println("  config      Displays currently loaded configuration file with the list of managed devices.")
}

func PrintDeleteHelp() {
	fmt.Println("Delete an existing job.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Delete the job named my-job.")
	fmt.Println("  kubectl job-manager delete -j my-job")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -j, --job <jobname>     The job name to delete.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  kubectl job-manager delete -j <jobname>")
}

func printCreateHelp() {
	fmt.Println("Create and run a job with associated devices.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Run a job using a device on the default namespace.")
	fmt.Println("  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c \"echo Hello!\"")
	fmt.Println()
	fmt.Println("  # Print the request of the job using a device on the default namespace.")
	fmt.Println("  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c \"echo Hello!\" --dry-run")
	fmt.Println()
	fmt.Println("  # Run a job using a device on the default namespace for debugging purposes.")
	fmt.Println("  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c \"echo Hello!\" --debug")
	fmt.Println()
	fmt.Println("  # Run a job that should be deleted after 10 minutes.")
	fmt.Println("  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c \"echo Hello!\" --ttl 600")
	fmt.Println()
	fmt.Println("  # Run a job that is allowed to backoff 3 times.")
	fmt.Println("  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c \"echo Hello!\" --backoff-limit 3")
	fmt.Println()
	fmt.Println("  # Run a job that must complete 2 times.")
	fmt.Println("  kubectl job-manager create -j my-job -n default -d xilinx.com/fpga -i busybox -c \"echo Hello!\" --completions 2")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -j, --job <jobname>             Set the job name.")
	fmt.Println("  -n, --namespace <namespace>     Set the namespace to create the job on.")
	fmt.Println("  -i, --image <image>             Set the image to be run on the job.")
	fmt.Println("  -c, --command <command>         Set the command to be executed on the image.")
	fmt.Println("  -d, --device <device>           Set a device resource request and limit for the job.")
	fmt.Println("      --ttl <seconds>             Time in seconds before the job is terminated. Defaults to 7200.")
	fmt.Println("      --completions <n>           Set the number of times the job must run. Defaults to 1.")
	fmt.Println("      --backoff-limit <n>         Set the number of times the job can backoff. Defaults to 1.")
	fmt.Println("      --dry-run                   Do not send the create request to the server and print it to stdout.")
	fmt.Println("      --debug <seconds>           Sets the job to debug mode, keeping it alive for the duration in seconds and preventing restarts.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  kubectl job-manager create -j <jobname> -n <namespace> -i <image> -c <commmand> -d <device> [--ttl <n>] [--completions <n>] [--backoff-limit <n>] [--dry-run] [--debug]")
}

func printListHelp() {
	fmt.Println("List information about jobs and devices.")
	fmt.Println()
	fmt.Println("IMPORTANT:")
	fmt.Println("  A configuration file with the list of managed devices must be present to list jobs and devices.")
	fmt.Println("  See the `kubectl job-manager config` command for more information.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # List all devices currently in use.")
	fmt.Println("  kubectl job-manager list devices")
	fmt.Println()
	fmt.Println("  # List all jobs currently using a device.")
	fmt.Println("  kubectl job-manager list jobs")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  devices     List all devices currently in use.")
	fmt.Println("  jobs        List all jobs currently using a device.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  kubectl job-manager list <command>")
}

func printConfig() {
	config := loadManagedDevices()
	if config == nil {
		fmt.Println("No configuration file found.")
	} else {
		fmt.Println("Loaded configuration: " + findConfigLocation())
		fmt.Println("---")
		for _, device := range config {
			fmt.Println(device)
		}
		fmt.Println("---")
	}
	fmt.Println()
	fmt.Println("Configuration file location priority:")
	fmt.Println("  - ${HOME}/.kube/job-manager.conf")
	fmt.Println("  - /etc/kubernetes/job-manager.conf")
}
