package main

import (
	"bufio"
	"fmt"
	"os"

	batchv1 "k8s.io/api/batch/v1"
)

// Checks if an array contains a given element
func isElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Loads the managed devices from the configuration file
func loadManagedDevices() []string {
	configLocation := findConfigLocation()
	if configLocation != "" {
		devices, err := readLines(configLocation)
		if err != nil {
			return nil
		}
		return devices
	}
	return nil
}

// Loads the managed devices from the configuration file
func findConfigLocation() string {
	// Checks for the file at "$HOME/.kube/job-manager.conf" first
	if _, err := os.Stat(os.Getenv("HOME") + "/.kube/job-manager.conf"); err != nil {
		// Checks for the file at "/etc/kubernetes/job-manager.conf" after
		if _, err := os.Stat("/etc/kubernetes/job-manager.conf"); err != nil {
			return ""
		}
		return "/etc/kubernetes/job-manager.conf"
	}
	return os.Getenv("HOME") + "/.kube/job-manager.conf"
}

func printJobYaml(job *batchv1.Job, p Parameters) {
	var yaml string = ""
	yaml += "apiVersion: batch/v1\n"
	yaml += "kind: Job\n"
	yaml += "metadata:\n"
	yaml += "  name: " + job.Name + "\n"
	yaml += "  namespace: " + job.Namespace + "\n"
	yaml += "spec:\n"
	yaml += "  activeDeadlineSeconds: " + fmt.Sprint(*job.Spec.ActiveDeadlineSeconds) + "\n"
	yaml += "  ttlSecondsAfterFinished: " + fmt.Sprint(*job.Spec.TTLSecondsAfterFinished) + "\n"
	yaml += "  parallelism: " + fmt.Sprint(*job.Spec.Parallelism) + "\n"
	yaml += "  completions: " + fmt.Sprint(*job.Spec.Completions) + "\n"
	yaml += "  template:\n"
	yaml += "    spec:\n"
	yaml += "      containers:\n"
	yaml += "      - name: " + job.Name + "\n"
	yaml += "        image: " + job.Spec.Template.Spec.Containers[0].Image + "\n"
	yaml += "        command: " + parseCommand(job.Spec.Template.Spec.Containers[0].Command) + "\n"
	yaml += "        resources:\n"
	yaml += "          requests:\n"
	yaml += "            " + p.device + ": 1\n"
	yaml += "          limits:\n"
	yaml += "            " + p.device + ": 1\n"
	yaml += "        imagePullPolicy: " + fmt.Sprint(job.Spec.Template.Spec.Containers[0].ImagePullPolicy) + "\n"
	yaml += "      restartPolicy: " + fmt.Sprint(job.Spec.Template.Spec.RestartPolicy) + "\n"
	yaml += "  backoffLimit: " + fmt.Sprint(*job.Spec.BackoffLimit) + "\n"
	fmt.Print(yaml)
}

func parseCommand(command []string) string {
	var ret string = "["
	var del string = ""
	for _, word := range command {
		ret += del + "\"" + word + "\""
		del = ", "
	}
	return ret + "]"
}
