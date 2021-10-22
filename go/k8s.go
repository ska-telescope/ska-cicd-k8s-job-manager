package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func listDevices(clientset kubernetes.Clientset, p Parameters) {
	if p.help || p.config == nil {
		printListHelp()
		return
	}

	format := "%60.60s | %-30.30s | %-25.25s | %-s\n"
	fmt.Printf(format, "DEVICE", "NODE", "NAMESPACE", "POD")

	// Get all nodes with the tpm label
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{LabelSelector: "tpm==true"})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// For each node
	for _, node := range nodes.Items {
		// Get the pods
		pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", node.Name)})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// Print the pods using a managed device
		for _, pod := range pods.Items {
			podDef := fmt.Sprint(pod)

			for _, device := range p.config {
				if strings.Contains(podDef, device) {
					fmt.Printf(format, device, node.Name, pod.Namespace, pod.Name)
				}
			}
		}
	}
}

func listJobs(clientset kubernetes.Clientset, p Parameters) {
	if p.help || p.config == nil {
		printListHelp()
		return
	}

	format := "%60.60s | %-30.30s | %-25.25s | %-s\n"
	fmt.Printf(format, "DEVICE", "JOB", "NAMESPACE", "TIME_TO_LIVE (Seconds)")

	// Get all jobs
	jobs, err := clientset.BatchV1().Jobs("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// For each job
	for _, job := range jobs.Items {
		var remaining_time int64
		remaining_time_str := ""
		// If the job has finished and has a timeout after finish
		if job.Status.CompletionTime != nil && job.Spec.TTLSecondsAfterFinished != nil {
			// Then the job remaining time is (${_TTL_AFTER_FINISHED} - (NOW - ${_FINISH_TIME}))
			remaining_time = int64(*job.Spec.TTLSecondsAfterFinished) - (time.Now().Unix() - job.Status.CompletionTime.Unix())
			// Otherwise, If the job has not finished and has an active timeout
		} else if job.Status.CompletionTime == nil && job.Spec.ActiveDeadlineSeconds != nil {
			// Then the job remaining time is (${_TTL} - (NOW - ${_START_TIME}))
			remaining_time = *job.Spec.ActiveDeadlineSeconds - (time.Now().Unix() - job.Status.StartTime.Unix())
			// If none of the previous conditions apply
		} else {
			// Then there is no timeout set
			remaining_time_str = "NoTimeoutSet"
		}

		if remaining_time > 0 {
			remaining_time_str = fmt.Sprint(remaining_time)
		} else if remaining_time_str == "" {
			remaining_time_str = "Terminating..."
		}

		// Print the job using a managed device
		jobDef := fmt.Sprint(job)
		for _, device := range p.config {
			if strings.Contains(jobDef, device) {
				fmt.Printf(format, device, job.Name, job.Namespace, remaining_time_str)
			}
		}
	}
}

func createJob(clientset kubernetes.Clientset, p Parameters) {
	if p.help || p.config == nil || p.job == "" || p.namespace == "" || p.image == "" || p.device == "" || p.command == nil {
		printCreateHelp()
		return
	}

	if p.config == nil {
		fmt.Println("WARNING: No configuration file exists with the list of managed devices. This device will not appear in the list command.")
		fmt.Println("WARNING: Please see `kubectl job-manager config` for more information.")
		fmt.Println("---")
	}

	if !isElementExist(p.config, p.device) {
		fmt.Println("WARNING: The device is not in the configuration file. It will not appear in the list command.")
		fmt.Println("WARNING: Please see `kubectl job-manager config` for more information.")
		fmt.Println("---")
	}

	// Resource definition, used to set requests and limits
	resourcedef := apiv1.ResourceList{
		apiv1.ResourceName(p.device): resource.MustParse("1"),
	}

	// Job definition, used to create the job
	jobdef := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.job,
			Namespace: p.namespace,
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds:   &p.ttl,
			TTLSecondsAfterFinished: &p.ttlafterfinished,
			Parallelism:             &p.parallelism,
			Completions:             &p.completions,
			BackoffLimit:            &p.backofflimit,
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					RestartPolicy: p.restartpolicy,
					Containers: []apiv1.Container{
						{
							Name:            p.job,
							Image:           p.image,
							Command:         p.command,
							ImagePullPolicy: p.imagepullpolicy,
							Resources: apiv1.ResourceRequirements{
								Requests: resourcedef,
								Limits:   resourcedef,
							},
						},
					},
				},
			},
		},
	}

	createOptions := metav1.CreateOptions{}
	if p.dryrun {
		createOptions.DryRun = append(createOptions.DryRun, "All")
	}

	newjob, err := clientset.BatchV1().Jobs(p.namespace).Create(context.Background(), jobdef, createOptions)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if !p.dryrun {
		fmt.Printf("job.batch/%s created\n", newjob.Name)
	} else {
		printJobYaml(newjob, p)
	}
}

func deleteJob(clientset kubernetes.Clientset, p Parameters) {
	if p.help || p.job == "" || p.namespace == "" {
		PrintDeleteHelp()
		return
	}

	fg := metav1.DeletePropagationBackground
	deleteOptions := metav1.DeleteOptions{PropagationPolicy: &fg, GracePeriodSeconds: new(int64)}
	err := clientset.BatchV1().Jobs(p.namespace).Delete(context.Background(), p.job, deleteOptions)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("job.batch %s deleted\n", p.job)
}
