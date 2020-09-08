# k8s-job-cleaner
Auto delete k8s jobs that Completed or Failed after ttl.
This tool helps you, when your kubernetes cluster can not use [TTL Controller](https://kubernetes.io/ja/docs/concepts/workloads/controllers/ttlafterfinished/#ttl%E3%82%B3%E3%83%B3%E3%83%88%E3%83%AD%E3%83%BC%E3%83%A9%E3%83%BC).
If your cluster can use `TTL Controller`, you should use it.

# Usage
```
Usage of k8s-job-cleaner:
  -kubeconfig string
    	local kubeconfig if you run from outside cluster.
  -name string
    	master name if you run from outside cluster.
  -namespace string
    	namespace that you want watch jobs.
  -ttl int
    	ttl of completed or failed job deletion minutes. (default 10)
  -url string
    	master url if you run from outside cluster.
```

# Authentication 
To execute k8s-job-cleaner, your ClusterRole attached your account need this permimssion.
```
rules:
- apiGroups:
  - batch
  resources:
  - jobs
  - jobs/status
  verbs:
  - list
  - delete
```
If you use In-Cluster service account, you need to create ServiceAccount `name: default` to target namespace.
see https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration https://github.com/kubernetes/client-go/tree/master/examples/out-of-cluster-client-configuration


# how to test in minikube
```
$ minikube start
$ kubectl apply -f example/kubernetes/test_cleaner_in_cluster_config.yaml
# check job scheduled
$ kubectl get jobs -n test-k8s-job-cleaner
# check k8s-job-cleaner delete old jobs
$ kubectl logs `kubectl get pods -n test-k8s-job-cleaner -l app=k8s-job-cleaner --no-headers -o custom-columns=":metadata.name"` -n test-k8s-job-cleaner
delete job test-job-complete
delete job test-job-fail
# check job deleted
$ kubectl get jobs -n test-k8s-job-cleaner
# clean test resouces.
$ kubectl delete -f example/kubernetes/test_cleaner_in_cluster_config.yaml
```
