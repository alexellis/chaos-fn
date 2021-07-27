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
  --set queueWorker.image=alexellis2/pro-queue-worker-demo:0.2.0
```

Deploy the function in this repository:

```bash
faas-cli deploy -f chaos-fn.yml
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

## Experiment 3 - set a timeout that expires

Setup a 502 error which may be seen when a function gives a timeout during scale from zero or its first invocation immediately after a deployment:

```bash

curl -i localhost:8080/function/chaos-fn/set --data-binary '
{	"status": 502,
	"delay": "5m",
    "body": "5m delay, then 502"}
' --header "Content-type: application/json"

```

I've picked 5 minutes `5m`, but you will need to pick a number that corresponds to your [extended timeouts](https://docs.openfaas.com/tutorials/expanded-timeouts/). Bear in mind that you may also have to extend the timeout of this function itself in the `chaos-fn.yml` file before running `faas-cli deploy` again.

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

