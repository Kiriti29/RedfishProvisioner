package github.com/Kiriti29/RedfishProvisioner/kubernetes

import (
  // "github.com/redfishProvisioner/kubernetes/base"
  batchv1 "k8s.io/api/batch/v1"
)

type JobClient struct {
    j *BatchV1Client
}

func New(namespace string) *JobClient{
    clientset = base.New()
    return &JobClient{
      j: clientset.BatchV1().Jobs("metalkube")
    }
}

func (j *JobClient) CreateJob(job *batchv1.Job) bool {
    result, _ := j.Create(job)
    if result.GetObjectMeta().GetName() != ""{
      for {
          job, err := jobsClient.Get(result.GetObjectMeta().GetName(), metav1.GetOptions{})
          if err != nil {
              log.Println("Unable to fetch job")
              break
          }
          if job.Status.Failed > 0 {
            log.Println("job failed")
            break
          }
          if job.Status.Succeeded > 0 {
            log.Println("job success")
              break
          }
      }
    }
    return status
  }
}

func (j *JobClient) DeleteJob(name string, label_selector map[string]string) bool {
    j.Delete(name)
    return true
}

func (j *JobClient) GetJobs(name string, label_selector map[string]string) bool {
    return true
}

func (j *JobClient) GetJobDetails(name string) bool {
    return true
}
