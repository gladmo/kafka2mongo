# kafka2mongo
Consume kafka data and write to mongodb

## how to use
```
    export MONGODB_DSN=mongodb://localhost:27017 
    export KAFKA_ADDR=localhost:9092 
    export EXCLUDE_TOPICS=__consumer_offsets
    go run main.go
```