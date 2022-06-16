# Filter Server

## Overview

This is the filter service, messages can be sent to this service over a REST API, these messages will be validated and stored in the database, any approvals will also be stored as well as the rejected messages. 
All images will need seperate approval, once each approval is approved the message will be revalidated. if any image is rejected, the whole message is rejected.

## configuration

There is a small yml config file stored on the root of the application, this config file can be used to tell the application to start up using different databases, you can chose Bbolt (https://github.com/etcd-io/bbolt), or mongoDB. You can set the URL for the language server that is required for checking banned words abd you can set the host and port the http service will listen on. if you change the port, you will need to update the dockerfile


## To Run

To run this service, the Language Service must be running as well.

## Using Docker
build the image using `docker build -t filter .`
run a container using `docker run -it --rm -p 8080:8080 filter`

## Usage

There are 6 REST apis, 1 API is used to send in the messages, the others are used to get stored messages, approvals and rejected messages. 1 is for approving, 1 is for rejecting. below is the apis is greater detail.


**POST** `/api/validate`

This requires a JSON body of

```
{
    "id": "1e969744-1e55-42a0-84c6-80d4fea2f1fd",
    "body": "# Simple Message\n\n\n\nA simple message"
}
```

The id and body is required. if the id and body are not present the message will be rejected.

**GET** `/api/messages`

This returns a list of all messages in the system.

```
{
    "messages": [
        {
            "id": "16",
            "body": "# Simple Message\n\n\n\n![The Eiffel Tower](https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg/240px-Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg)",
            "actions": [
                {
                    "id": "1e969744-1e55-42a0-84c6-80d4fea2f1fd",
                    "status": "pending",
                    "reason": "image [https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg/240px-Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg] requires approval"
                }
            ],
            "status": "awaiting approval",
            "reasons": "message contains image that require approval"
        },
        {
            "id": "8",
            "body": "# Simple Message\n\n\n\n[google](https://www.google.com)",
            "status": "rejected",
            "reasons": "message body contains external links"
        },
        {
            "id": "9",
            "body": "# Simple Message\n\n\n\n[google](https://www.google.com)",
            "status": "rejected",
            "reasons": "message body contains external links"
        }
    ],
    "updated": "2022-06-16T07:37:27.2521303Z"
}
```

**GET** `/api/rejected`

This returns a list of all rejected messages in the system.

```
{
    "rejected": [
        {
            "id": "10",
            "body": "# Simple Message\n\n\n\n[google](https://www.google.com)",
            "status": "rejected",
            "reasons": "message body contains external links"
        }
    ],
    "updated": "2022-06-16T07:52:30.9034105Z"
}
```
**GET** `/api/approvals`

This returns a list of all approvals in the system.

```
{
    "approvals": [
        {
            "id": "1e969744-1e55-42a0-84c6-80d4fea2f1fd",
            "status": "pending",
            "messageId": "16",
            "reason": "image [https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg/240px-Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg] requires approval"
        },
        {
            "id": "89f3a7e7-ae11-42bb-9405-90d29debbf29",
            "status": "pending",
            "messageId": "15",
            "reason": "image [https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg/240px-Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg] requires approval"
        }
    ],
    "updated": "2022-06-16T07:52:37.5681565Z"
}
```

The id of these approval messages is used to approve or reject the image

**POST** `/api/approvals/:id/approve`

Using the id provided, this will approve the image. If there are multiple images in the message, all images must be approved before the message is r-eevaluated

**POST** `/api/approvals/:id/reject`

Using the id provided, this will reject the image. Messages with a rejected image will be updated and stored in the rejected store. If there are multiple images in the message and one is rejected, the whole message is rejected. 

## design decisons and changes

I've used bbolt because its a little embedding key,value store, which is fast and not memory based. Ive also added a MongoDB driver to show that its possible to have other databases attached.

ideally the banned words would be cached here rather than fetching from the server all the time. The `updated` field on the banned words list would be used to acheive this.

Orignally I felt that using a Message Queue like SQS, rabbitMQ or RedPanda would be best for this. Then you could have a number of consumers here that just processed the messages in the queue. I went with a REST api for simplistic sake, however its very easy to added a message queue listener and keep the REST api as well.

I was also torn on offloading the validating to another goroutine and ending the endpoint sooner. this would again allow multiple messages to be handled at the same time. Using this idea it would be best to offload the approvals to another service.

My logger of choice is Logrus, its a bit old now but it still works perfectly. I would recommend moving to Zap or Zero however.