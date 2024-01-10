## Project Setup

### AWS CLI Commands (LocalStack)

```
// create bucket
aws s3 mb s3://development-bucket-ims --endpoint-url http://localhost:4566

// list all buckets
aws s3 ls --endpoint-url=http://localhost:4566 --recursive --human-readable

// view object in the bucket
aws s3 ls s3://development-bucket-ims --endpoint-url=http://localhost:4566 --recursive --human-readable

// copy everything in bucket to current working directory
aws s3 cp s3://development-bucket-ims/test/ . --recursive --endpoint-url=http://localhost:4566

// TO DELETE BUCKET WITH OBJECTS
// delete all object in bucket
aws s3 ls s3://development-bucket-ims --endpoint-url=http://localhost:4566 --recursive --human-readable
// delete empty bucket
aws s3 rb s3://development-bucket-ims --endpoint-url=http://localhost:4566
```
