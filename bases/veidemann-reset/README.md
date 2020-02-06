# Reset veidemann

## Example

Reset dev environment:
```bash
kustomize build | kubectl apply --context=minikube --validate=true -f -
```

## What is cleaned up?

Tries to remove all files in (cannot remove open files):

```bash
/warcs
/validwarcs
/invalidwarcs
/delivery
/backup/oos
``` 

Deletes (empty; `r.db(database).table(table).delete()`) the following database tables :

```bash
veidemann:
    "crawl_host_group"
    "crawl_log"
    "crawled_content"
    "events"
    "executions"
    "extracted_text"
    "job_executions"
    "locks"
    "page_log"
    "storage_ref"
    "uri_queue"

report:
    "invalid_warcs"
    "valid_warcs"
```

## Clean up after veidemann-reset

The _TTLAfterFinished_ feature gate needs to be enabled for the job to be delete automatically, e.g. (minikube):

    minikube start --feature-gates=TTLAfterFinished=true

To configure minikube to enable feature gate:
   
    minikube config set feature-gates TTLAfterFinished=true

Without the feature gate the job must be removed manually:

    kubectl delete job -n veidemann veidemann-reset
    
The secret created must be removed manually:

    kubectl delete secret -n veidemann veidemann-reset-XXXXX
