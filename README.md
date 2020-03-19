# kafka2mongo
Consume kafka data and write to mongodb

## how to use
```
    MONGODB_DSN=mongodb://localhost:27017 \
    KAFKA_ADDR=localhost:9092 \
    EXCLUDE_TOPICS=__consumer_offsets
```