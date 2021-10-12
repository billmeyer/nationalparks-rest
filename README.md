# National Parks - REST Api

This project represents a web services (written in Go) that can be used to make RESTful web services to lookup National Park data stored in a MySQL database.

It utilizes a MySQL database that can be found at [https://github.com/billmeyer/nationalparks-mysql]() and is a dependency.

## Setup

1. Follow the instructions at [https://github.com/billmeyer/nationalparks-mysql]() for installing and running the MySQL database that this service uses.  The mysql instance can be run as either a Docker container or deployed to Kubernetes.

    > NOTE: Make note of the IP Address the MySQL Server instance.  It will be used below.

2. Next, clone the repository

    ```bash
    $ git clone https://github.com/billmeyer/nationalparks-rest.git
    ```

3. Change your working directory into the project directory

    ```bash
    $ cd nationalparks-rest
    ```

4. Edit `config.sh` and modify the entries to conform to your specific deployment of MySQL and this service.  For example, when running the service locally, the `BACKEND_URL` might be `http://localhost:8080`.  When running in a Docker container or on Kubernetes, it would be reachable via a different IP address/port combination.

## Run

> NOTE: This service was developed using Go version 1.17.1.  Be sure a recent version has been downloaded from [https://golang.org/dl/]() and that the `go` executable is the system path.

1. Source your `config.sh` to be sure your environment variables are set properly:

   ```bash
   $ source config.sh
   ```

2. From the `nationalparks-rest` directory, run a local instance of the server using:

   ```bash
   $ go run cmd/main/server.go 
{"message":"National Parks REST API Service","severity":"info","timestamp":"2021-10-12T16:24:27.375633-05:00"}
{"message":"Using MySQL instance at 192.168.3.230:3306","severity":"info","timestamp":"2021-10-12T16:24:27.375871-05:00"}
{"message":"Server started at 0.0.0.0:8080","severity":"info","timestamp":"2021-10-12T16:24:27.376107-05:00"}
   ```
   
   The startup output should reflect the MySQL instance you have configured and the service is listening on interface 0.0.0.0, port 8080.

3. Confirm the Health Check endpoint is reachable by running the `health-check.sh` script:

   ```bash
   $ ./health-check.sh
   HTTP/1.1 200 OK
   Access-Control-Expose-Headers: Server-Timing
   Content-Type: application/json
   Server-Timing: traceparent;desc="00-423c205380e2894a99814a02b8547dd1-c032884227869afb-01"
   Vary: Origin
   Date: Tue, 12 Oct 2021 21:26:29 GMT
   Content-Length: 24
   
   "API is up and running"
   ```

   If you see the `"API is up and running"` result, the service is reachable.

4. Confirm the service can communicate with the MySQL instance by running the `get-a-park.sh` script:

   ```bash
   $ ./get-a-park.sh
   HTTP/1.1 200 OK
   Access-Control-Expose-Headers: Server-Timing
   Content-Type: application/json
   Server-Timing: traceparent;desc="00-3b1610a82d56f9834177ed79a689cbc3-3b6e084c58c396e8-01"
   Vary: Origin
   Date: Tue, 12 Oct 2021 21:28:10 GMT
   Content-Length: 250
   {"id":1,"location_num":"ADAM","location_name":"Adams National Historical Park","address":"135 Adams Street","city":"Quincy","state":"MA","zip_code":2169,"phone_num":"(617) 770-1175","fax_num":"(617) 472-7562","latitude":42.2564,"longitude":-71.0112}
   ```

### Review traces

1. Both of the above tests will invoke the REST api which will, in turn, produce traces that are sent to Splunk Observability. Confirm these traces are arriving in Splunk Observability by visiting [https://app.us1.signalfx.com/#/apm/troubleshooting]().

2. Click on the **Services** dropdown and add a service named `nationalparks-rest`:

   ![Service Filter](images/service-filter.png)

   Dismiss the filter and you should see a single view of the `nationalparks-rest` service and it's dependency on the `mysql` service:

   ![National Parks Service](images/national-parks-apm.png)

3. Take a closer look at the two purple blips on the **Service Requests & Errors** chart:

   ![Service Requests & Errors](images/service-requests-errors.png)

4. Hover over the first blip and click the peak on the chart:

   ![Health Check Trace](images/health-check-trace.png)

   This should reveal the Health Check test we ran first in the prior section.  You'll notice the trace ID (`423c205380e2894a99814a02b8547dd1` in this example) matches the trace ID that was returned from the service call as a `Server-Timing` HTTP Header: 

   ```
   Server-Timing: traceparent;desc="00-423c205380e2894a99814a02b8547dd1-c032884227869afb-01"
   ```

5. Likewise, clicking on the peak of the second blip:

   ![National Parks Trace](images/national-parks-trace.png)

   reveals the trace for the second test of running the `get-a-park.sh` script.

## Evaluate
