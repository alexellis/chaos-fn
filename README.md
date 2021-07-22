# chaos-fn

A function for chaos testing with OpenFaaS

## Use-cases

* Test retries on certain HTTP codes
* Test timeouts
* Test certain lengths of HTTP request bodies

## Experiment

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

