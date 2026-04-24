go work init ./authservice ./productservice ./apigateway ./orderservice ./notificationservice ./shared

docker run -d --name rabbitmq -p 5672:5672 rabbitmq


ecommerce poject with microservices in go

api gateway central entrypoint with auth 
product 
order
auth
notification service

order service gets product service data via rest api interservice call
connected order with notification service via rabbitmq 


no configs but will be added soon 
db wll be updated to psql with saga and dtransactions

servies will be made stateless 

