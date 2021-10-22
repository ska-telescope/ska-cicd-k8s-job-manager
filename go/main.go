package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Parameters struct {
	action           string
	subaction        string
	job              string
	namespace        string
	image            string
	command          []string
	device           string
	ttl              int64
	ttlafterfinished int32
	parallelism      int32
	completions      int32
	backofflimit     int32
	restartpolicy    v1.RestartPolicy
	imagepullpolicy  v1.PullPolicy
	dryrun           bool
	config           []string
	help             bool
}

var params = Parameters{
	action:           "",
	subaction:        "",
	job:              "",
	namespace:        "",
	image:            "",
	command:          nil,
	device:           "",
	ttl:              7200,
	ttlafterfinished: 0,
	parallelism:      1,
	completions:      1,
	backofflimit:     1,
	restartpolicy:    v1.RestartPolicyOnFailure,
	imagepullpolicy:  v1.PullIfNotPresent,
	dryrun:           false,
	config:           loadManagedDevices(),
	help:             false,
}

func main() {
	version := "0.1.0"

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Parse arguments for paramaters
	parseArguments()

	// Chose the appropriate function to call
	switch params.action {
	case "create":
		createJob(*clientset, params)
	case "list":
		{
			switch params.subaction {
			case "devices":
				listDevices(*clientset, params)
			case "jobs":
				listJobs(*clientset, params)
			default:
				printListHelp()
			}
		}
	case "delete":
		deleteJob(*clientset, params)
	case "version":
		fmt.Println(version)
	case "config":
		printConfig()
	case "", "help", "--help", "-h":
		printHelp()
	default:
		{
			printHelp()
			os.Exit(1)
		}
	}
}

func parseArguments() {
	i := 1
	params.action = os.Args[i]
	i++

	if params.action == "list" && len(os.Args) > 2 {
		params.subaction = os.Args[i]
		i++
	}

	for i < len(os.Args) {
		switch os.Args[i] {
		case "-j", "--job":
			{
				i++
				params.job = os.Args[i]
			}
		case "-n", "--namespace":
			{
				i++
				params.namespace = os.Args[i]
			}
		case "-i", "--image":
			{
				i++
				params.image = os.Args[i]
			}
		case "-c", "--command":
			{
				i++
				params.command = strings.Split(os.Args[i], " ")
			}
		case "-d", "--device":
			{
				i++
				params.device = os.Args[i]
			}
		case "--ttl":
			{
				i++
				val, err := strconv.ParseInt(os.Args[i], 10, 64)
				if err != nil {
					panic(err)
				}
				params.ttl = val
			}
		case "--completions":
			{
				i++
				val, err := strconv.ParseInt(os.Args[i], 10, 32)
				if err != nil {
					panic(err)
				}
				params.completions = int32(val)
			}
		case "--backoff-limit":
			{
				i++
				val, err := strconv.ParseInt(os.Args[i], 10, 32)
				if err != nil {
					panic(err)
				}
				params.backofflimit = int32(val)
			}
		case "--dry-run":
			{
				params.dryrun = true
			}
		case "--debug":
			{
				i++
				val, err := strconv.ParseInt(os.Args[i], 10, 64)
				if err != nil {
					panic(err)
				}
				params.ttlafterfinished = int32(val)
				params.restartpolicy = v1.RestartPolicyNever
			}
		case "help", "--help", "-h":
			{
				params.help = true
			}
		default:
			fmt.Printf("Unknown option %s.\n", os.Args[i])
			os.Exit(2)
		}

		i++
	}
}
