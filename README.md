# go-dev-task

## Installation

`git clone https://github.com/v1tbrah/lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ`

### Note!
* Because this is a training task:
  * This tutorial omits application configuration, so:
    * API server with phone numbers address: ":3333"
    * Exchange name for triggers: "triggers"
    * Routing key for triggers: "triggers"
  * Decorator is used instead of database. It is filled with numbers with id from 1 to 99.
  
## Getting started

* Run RMQ container `docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.11-management`
* Open first terminal. Go to project working directory. Run API server with phone numbers. For example:
   ```
   cd ~/go/src/lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ
   go run cmd/phone/main.go
   ```
* Open second terminal. Go to project working directory. Run logger exchange. For example:
   ```
   cd ~/go/src/lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ
   go run cmd/logger/main.go
   ```
* Open third terminal. Go to project working directory. Run triggers app. For example:
   ```
   cd ~/go/src/lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ
   go run cmd/app/main.go
   ```
* Let's push to triggers app:
  * Open your docker container. Open terminal in container.
  * Publish valid msg to triggers app:
    * `rabbitmqadmin publish exchange="triggers" routing_key="triggers" payload="{\"record_id\": \"1\"}"`
  * Publish invalid msg to triggers app:
    * `rabbitmqadmin publish exchange="triggers" routing_key="triggers" payload="{\"record_id\":"`
  * If everything works, you should see the corresponding messages in the terminal logger exchange.




