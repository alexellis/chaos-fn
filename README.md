# chaos-fn

A function for chaos testing with OpenFaaS

## Use-cases

* Test retries on certain HTTP codes
* Test timeouts
* Test certain lengths of HTTP request bodies

## Setup

Install openfaas:

```bash
kind cluster create

arkade install openfaas \
  --set queueWorker.image=alexellis2/pro-queue-worker-demo:0.1.1
```

## Experiment 1

Cause the API to start failing

```bash

curl -i localhost:8080/function/chaos-fn/set --data-binary '
{	"status": 500,
	"delay": "1s",
    "body": "1 second delay, then 500"}
' --header "Content-type: application/json"

```

Observe it:

```bash

curl -i localhost:8080/function/chaos-fn

```

Fix it


```bash

curl -i localhost:8080/function/chaos-fn/set --data-binary '
{"status": 200,
	"delay": "1ms",
    "body": "1ms second delay, then 200"}
' --header "Content-type: application/json"

```


```bash

curl -i localhost:8080/function/chaos-fn

```

## Experiment 2 - the PRO async queue worker

Set up a HTTP code that may be retried

```bash

curl -i localhost:8080/function/chaos-fn/set --data-binary '
{	"status": 500,
	"delay": "1s",
    "body": "1 second delay, then 500"}
' --header "Content-type: application/json"

```

View the logs:

```bash
kubectl logs deploy/queue-worker -n openfaas -f

```

Now invoke:

```bash

curl -i localhost:8080/async-function/chaos-fn -d ""

```

Observe the retrying mechanism.

Whenever you like, fix the error to allow the next retry to complete:

```bash

curl -i localhost:8080/function/chaos-fn/set --data-binary '
{"status": 200,
	"delay": "1ms",
    "body": "1ms second delay, then 200"}
' --header "Content-type: application/json"

```
