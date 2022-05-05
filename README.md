# Match Engine

Lambda function responsible to match buy and sale orders

How to test ?


    docker-compose up


- The file test has a BUY event, to test sales with match potential must be created.
- The creation of SALES can be done through the API or using [aws cli](https://docs.aws.amazon.com/cli/latest/reference/dynamodb/put-item.html)

Here we have an example dynamoDb item em in ExpressionAttributeValues as is accepted by dynamoDb        

    {
        "operationType": {"S": "SALE"},
        "id": {"S": "2233998|cd8e90d4-32b0-40f7-a3a1-8be706474913"},
        "quantity": {"N": "200"},
        "value": {"N": "1334.99"},
        "hash": {"S": "279ba23c396ea563f33da2bcd55dc99d"},
        "operationStatus": {"S": "IN_TRADE"},
        "userId": {"N": "2233998"},
        "audit": {"M": 
            {
                "createdAt": {"S": "2022-03-30T18:22:38.625Z"},
                "updatedAt": {"S": "2022-03-30T18:22:38.625Z"},
                "updatedBy": {"S": "API"}
            }
        },
        "traceId": {"S": "cd8e90d4-32b0-40f7-a3a1-8be706474913"}
    }

If you prefer, it can be done by the API using [this example](https://github.com/joycesaquino/order-book-api#readme) according to the documentation !

If there are no records, the function is ready, just follow it on the console. I put the support logs to make testing easier :)
With the records properly created we can run the test and follow the lambda logs in the console.
    
    Run using ide in file > main_test.go

