package main

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

type JobCleaner struct {
	api batchv1.JobInterface
	ttl time.Duration
}

func NewJobCleaner(api batchv1.JobInterface, ttl time.Duration) *JobCleaner {
	return &JobCleaner{
		api: api,
		ttl: ttl,
	}
}

func (j *JobCleaner) Watch(ctx context.Context) error {
	for {
		if err := j.Delete(ctx); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (j *JobCleaner) Delete(ctx context.Context) error {
	jobs, err := j.api.List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, job := range jobs.Items {
		if j.jobCompletedOrFailed(job) {
			fmt.Printf("delete job %s\n", job.Name)
			if err := j.api.Delete(ctx, job.Name, metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
	}
	return nil

}
func (j *JobCleaner) jobCompletedOrFailed(job v1.Job) bool {
	now := time.Now()
	for _, c := range job.Status.Conditions {
		if c.LastProbeTime.Add(j.ttl).Before(now) && c.Type == v1.JobComplete || c.Type == v1.JobFailed {
			return true
		}
	}
	return false
}
