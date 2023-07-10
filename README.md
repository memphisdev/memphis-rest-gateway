<div align="center">
  
  ![github memphis banner](https://user-images.githubusercontent.com/70286779/229371212-8531c1e1-1a9d-4bbe-9285-b4dbb8601bfa.jpeg)
  
</div>

<div align="center">

  <h4>

**[Memphis](https://memphis.dev)** is a next-generation alternative to traditional message brokers.

  </h4>
  
  <a href="https://landscape.cncf.io/?selected=memphis"><img width="200" alt="CNCF Silver Member" src="https://github.com/cncf/artwork/raw/master/other/cncf-member/silver/white/cncf-member-silver-white.svg#gh-dark-mode-only"></a>
  
</div>

<div align="center">
  
  <img width="200" alt="CNCF Silver Member" src="https://github.com/cncf/artwork/raw/master/other/cncf-member/silver/color/cncf-member-silver-color.svg#gh-light-mode-only">
  
</div>
 
 <p align="center">
  <a href="https://memphis.dev/docs/">Docs</a> - <a href="https://twitter.com/Memphis_Dev">Twitter</a> - <a href="https://www.youtube.com/channel/UCVdMDLCSxXOqtgrBaRUHKKg">YouTube</a>
</p>

<p align="center">
<a href="https://discord.gg/WZpysvAeTf"><img src="https://img.shields.io/discord/963333392844328961?color=6557ff&label=discord" alt="Discord"></a>
<a href="https://github.com/memphisdev/memphis/issues?q=is%3Aissue+is%3Aclosed"><img src="https://img.shields.io/github/issues-closed/memphisdev/memphis?color=6557ff"></a> 
  <img src="https://img.shields.io/npm/dw/memphis-dev?color=ffc633&label=installations">
<a href="https://github.com/memphisdev/memphis/blob/master/CODE_OF_CONDUCT.md"><img src="https://img.shields.io/badge/Code%20of%20Conduct-v1.0-ff69b4.svg?color=ffc633" alt="Code Of Conduct"></a> 
<a href="https://docs.memphis.dev/memphis/release-notes/releases/v0.4.2-beta"><img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/memphisdev/memphis?color=61dfc6"></a>
<img src="https://img.shields.io/github/last-commit/memphisdev/memphis?color=61dfc6&label=last%20commit">
</p>

A simple, robust, and durable cloud-native message broker wrapped with<br>
an entire ecosystem that enables cost-effective, fast, and reliable development of modern queue-based use cases.<br><br>
Memphis enables the building of modern queue-based applications that require<br>
large volumes of streamed and enriched data, modern protocols, zero ops, rapid development,<br>
extreme cost reduction, and a significantly lower amount of dev time for data-oriented developers and data engineers.

# REST Gateway (HTTP Proxy)

## Introduction

To enable message production via HTTP calls for various use cases and ease of use, Memphis added an HTTP gateway to receive REST-based requests (=messages) and produce those messages to the required station.

Common use cases that benefit from the REST Gateway are&#x20;

* Produce events directly from a frontend
* Produce CDC events using the Debezium HTTP server
* ArgoCD webhooks
* Receive data from Fivetran/Rivery/Any ETL platform using HTTP calls

## Architecture

1. An endpoint creates an HTTP request toward the REST Gateway using **port 4444**
2. The REST gateway receives the incoming request and produces it as a message to the station

![REST gateway](https://user-images.githubusercontent.com/70286779/212469259-9f092921-63fa-4121-83cf-90f745d4b952.jpeg)


For scale requirements, the "REST gateway" component is separate from the brokers' pod and can scale out individually.

## Security Mechanisms

### JWT

Memphis REST (HTTP) gateway makes use of JWT-type identification.\
[JSON Web Tokens](https://jwt.io/) are an open, industry-standard RFC 7519 method for representing claims securely between two parties.

### API Token

Memphis REST (HTTP) gateway can be configured to use API-tokens as the default authentication method, see AUTH_METHOD below.

### AUTH_METHOD

The global setting AUTH_METHOD can be used to control the default authentication method.
AUTH_METHOD can be set to one of the values listed in the table below.

| Value      | Description                                                         |
|------------|---------------------------------------------------------------------|
| jwt        | JWT-type authentication.                                            |
| api_token  | API token authentication.                                           |
| hmac_token | HMAC (Hash-based Message Authentication Code) token authentication. |
| none       | No authentication.                                                  |

If the AUTH_METHOD configuration option is missing then authentication method defaults to JWT.


#### JWT Authentication

If AUTH_METHOD is set to 'jwt', then the configuration options in the table below can be used.

| Variable                       | Description                                            |
|--------------------------------|--------------------------------------------------------|
| JWT_SECRET                     | The secret key used to generate the JWT token.         |
| JWT_EXPIRES_IN_MINUTES         | The JWT token valid time in minutes.                   |
| REFRESH_JWT_SECRET             | The secret key used to generate the JWT refresh token. |
| REFRESH_JWT_EXPIRES_IN_MINUTES | The JWT refresh token valid time in minutes.           |

#### API Token

If AUTH_METHOD is set to api_token, then the configuration options in the table below can be used.

| Variable                       | Description                                                   |
|--------------------------------|---------------------------------------------------------------|
| API_TOKEN_HEADER               | Name of HTML-header that contains the API token.              |
| API_TOKEN                      | API token shared between the client and memphis-rest-gateway. |

Example configuration:
````
{
  "VERSION": "1.0.2",
   "AUTH_METHOD": "api_token",
   "API_TOKEN_HEADER": "X-Api-Token",
   "API_TOKEN": "EibA9ypmW7XPpn"
}
````

Example usage:
````
#!/bin/bash

now=`date`
msg="auth: api token station: api_token date: ${now}"

curl -s -XPOST 'http://localhost:4444/stations/api_station/produce/single' \
        -H 'Content-Type: application/json' \
        -H 'X-API-TOKEN: EibA9ypmW7XPpn' \
        --data-raw "{\"message\": \"$msg\"}"
````



#### HMAC Token

HMAC (Hash-based Message Authentication Code) token authentication is a method of securing and validating tokens in token-based authentication systems. 
It involves using a symmetric cryptographic algorithm and a secret key to generate and verify a signature attached to the token.

If AUTH_METHOD is set to hmac_token, then the configuration options in the table below can be used.

| Variable          | Description                                                                                                                               |
|-------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
|	HMAC_TOKEN_HEADER | Name of header of the body signature. The body signature is a MAC hex digest of the body calcuated using the TOKEN_SECRET as the hash key.|
|	HMAC_TOKEN_SECRET | A shared secret between the client and memphis-rest-gateway for calculating the body signature.                                           |
|	HMAC_TOKEN_HASH   | Hash algorithm for calculating the body signature, can be either sha256 or sha512.                                                        |

Example configuration:
````
{
  "VERSION": "1.0.2",
  "AUTH_METHOD": "hmac_token",
  "HMAC_TOKEN_HEADER": "X-Hook-Signature",
  "HMAC_TOKEN_SECRET": "super secret string",
  "HMAC_TOKEN_HASH": "sha512"
}
````

Example usage:
````
#!/usr/bin/python3 

import hashlib
import hmac
import base64
import requests

url = "http://localhost:4444/stations/hmac_station/produce/single"
msg_quoted    = "\"Quoted Message\""
msg_unquoted  = "Quoted Message"

secret  = bytes('super secret string', 'utf-8')
message = bytes(msg_quoted, 'utf-8')

hash = hmac.new(secret, message, hashlib.sha512)
sig = hash.hexdigest()

headers = {"X-Hook-Signature":  sig }

response = requests.post(url, headers=headers, json=msg_unquoted)
print("Status Code", response.status_code)
print(response.json())
````




#### None

No authentication method is used.

Example configuration:
````
{
  "VERSION": "1.0.2",
  "AUTH_METHOD": "none"
}
````

Example usage:
````
#!/bin/bash

now=`date`
msg="auth: none station: none date: ${now}"

curl -s -XPOST 'http://localhost:4444/stations/none_station/produce/single' \
        -H 'Content-Type: application/json' \
        --data-raw "{\"message\": \"$msg\"}"
````


### AUTH_METHOD per station

It is possible to set individual authentication methods per station using the 'Stations' configuration option.
Any station specific configuration overrides the global authentication methods.
All stations has to have a 'NAME' configured.

````
{
  "VERSION": "1.0.2",
  "JWT_EXPIRES_IN_MINUTES": 15,
  "REFRESH_JWT_EXPIRES_IN_MINUTES": 300,
  "Stations": [
    {
      "NAME": "none_station",
      "AUTH_METHOD": "none" 
    },
    {
      "NAME": "api_station",
      "AUTH_METHOD": "api_token",
      "API_TOKEN_HEADER": "X-Api-Token",
      "API_TOKEN": "EibA9ypmW7XPpn"
    },
    {
      "NAME": "hmac_station",
      "AUTH_METHOD": "hmac_token",
      "HMAC_TOKEN_HEADER": "X-Hook-Signature",
      "HMAC_TOKEN_SECRET": "super secret string",
      "HMAC_TOKEN_HASH": "sha512"
    },
    {
      "NAME": "jwt_station",
      "AUTH_METHOD": "jwt",
      "JWT_SECRET": "jwt_station secret",
      "REFRESH_JWT_SECRET": "jwt_station refresh secret"
    
    }
  ]
}
````

Note that if a station is configured to use 'jwt', then the path to aquire the token and to refresh a token has to include the name of the station.

Example acquiring a jwt token for a configured station:
````
#!/bin/bash

raw_token=$(curl -s -XPOST 'http://localhost:4444/auth/authenticate/jwt_station' \
                 -H 'Content-Type: application/json' \
                 --data-raw '{
                     "username": "root",
                     "password": "memphis",
                     "token_expiry_in_minutes": 1,
                     "refresh_token_expiry_in_minutes": 2
                 }')

token=$(echo $raw_token | jq -r .jwt)
refresh_token=$(echo $raw_token | jq -r .jwt_refresh_token)

now=`date`
msg="auth: jwt station: jwt date: ${now}"

curl -s -XPOST 'http://localhost:4444/stations/jwt_station/produce/single' \
        -H "Authorization: Bearer $token" \
        -H 'Content-Type: application/json' \
        --data-raw "{\"message\": \"$msg\"}"
````

Example refresing a jwt token for a configured station:
````
new_token=$(curl -s -XPOST 'http://localhost:4444/auth/refreshToken/jwt_station' \
                    -H 'Content-Type: application/json' \
                    --data-raw '{
                            "jwt_refresh_token": "$refresh_token",
                            "token_expiry_in_minutes": 1,
                            "refresh_token_expiry_in_minutes": 2
                    }' | jq -r .jwt)

now=`date`
msg="auth: jwt station: jwt date: ${now}"

curl -s -XPOST 'http://localhost:4444/stations/jwt_station/produce/single' \
        -H "Authorization: Bearer $new_token" \
        -H 'Content-Type: application/json' \
        --data-raw "{\"message\": \"$msg\"}"

````



## Sequence diagram

![Sequence diagram](https://user-images.githubusercontent.com/70286779/212469294-ebf2da3f-af30-46bc-bb42-ef860159356e.jpeg)


## Getting started

Tip: Please make sure your 'REST gateway' component is exposed either through localhost or public IP<br>
Tip: The REST gateway URL for the **sandbox** environment is:<br>
https://restgw.sandbox.memphis.dev

### Authenticate

First, you have to authenticate to get a JWT token.\
The JWT token is valid by default for 15 minutes.

#### Example:

```
curl --location --request POST 'rest_gateway:4444/auth/authenticate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "root",
    "connection_token": "memphis",
    "password": "memphis, // connect with only one of the following methods: connection token / password
    "token_expiry_in_minutes": 60,
    "refresh_token_expiry_in_minutes": 10000092
}'
```

Expected output:&#x20;

```
{"expires_in":3600000,"jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzQ3MTg0MjV9._A-fRI78fPPHL6eUFoWZjp21UYVcjXwGWiYtacYPZR8","jwt_refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNzQ3MjAzNDV9.d89acaIr4CaBp7csm-jmJv0J45YrD_slvlEOKu2rs7Q","refresh_token_expires_in":600005520000}
```

#### Parameters

`username`: Memphis application-type username\
`connection_token`: Memphis application-type connection token\
`token_expiry_in_minutes`: Initial token expiration time.\
`refresh_token_expiry_in_minutes`: When should

### Refresh Token

Before the JWT token expires, you must call the refresh token to get a new one, or after authentication failure.\
The refresh JWT is valid by default for 5 hours.

#### Example:

```
curl --location --request POST 'rest_gateway:4444/auth/refreshToken' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jwt_refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNzQ3MjA2NjB9.Furfr5EZlBlglVPSjtU4x02z_jbWhu5pIByhCRh6FU8",
    "token_expiry_in_minutes": 60,
    "refresh_token_expiry_in_minutes": 10000092
}'
```

Expected output:

```
{"expires_in":3600000,"jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzQ3MTg3NTF9.EO5ersr0kQxQNRI0XlbqzOryt-F1-MmFGXRKn2sM8Yw","jwt_refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNzQ3MjA2NzF9.E621wF_ieC-9rq4IgrsqYMPApAPS8YDgkT8R-69-Y5E","refresh_token_expires_in":600005520000}
```

### Produce a single message

Attach the JWT token to every request.\
JWT token as '`Bearer`' as a header.

#### Supported content types:

* text
* application/json
* application/x-protobuf

#### Example:

```
curl --location --request POST 'rest_gateway:4444/stations/<station_name>/produce/single' \
--header 'Authorization: Bearer eyJhbGciOiJIU**********.e30.4KOGRhUaqvm-qSHnmMwX5VrLKsvHo33u3UdJ0qYP0kI' \
--header 'Content-Type: application/json' \
--data-raw '{"message": "New Message"}'
```

#### If you don't have the option to add the authorization header, you can send the JWT via query parameters:

```
curl --location --request POST 'rest_gateway:4444/stations/<station_name>/produce/single?authorization=eyJhbGciOiJIU**********.e30.4KOGRhUaqvm-qSHnmMwX5VrLKsvHo33u3UdJ0qYP0kI' \
--header 'Content-Type: application/json' \
--data-raw '{"message": "New Message"}'
```

Expected output:

```
{"error":null,"success":true}
```

#### Error Example:

```
{"error":"Schema validation has failed: jsonschema: '' does not validate with file:///Users/user/memphisdev/memphis-rest-gateway/123#/required: missing properties: 'field1', 'field2', 'field3'","success":false}
```

### Produce a batch of messages&#x20;

Attach the JWT token to every request.\
JWT token as '`Bearer`' as a header.

#### Supported content types:

* application/json

#### Example:

```
curl --location --request POST 'rest_gateway:4444/stations/<station_name>/produce/batch' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.4KOGRhUaqvm-qSHnmMwX5VrLKsvHo33u3UdJ0qYP0kI' \
--header 'Content-Type: application/json' \
--data-raw '[
    {"message": "x"},
    {"message": "y"},
    {"message": "z"}
]'
```

Expected output:

```
{"error":null,"success":true}
```

#### Error Examples:

```
{"errors":["Schema validation has failed: jsonschema: '' does not validate with file:///Users/user/memphisdev/memphis-rest-gateway/123#/required: missing properties: 'field1'","Schema validation has failed: jsonschema: '' does not validate with file:///Users/user/memphisdev/memphis-rest-gateway/123#/required: missing properties: 'field1'"],"fail":2,"sent":1,"success":false}
```

## Support 🙋‍♂️🤝

### Ask a question ❓ about Memphis{dev} or something related to us:

We welcome you to our discord server with your questions, doubts and feedback.

<a href="https://discord.gg/WZpysvAeTf"><img src="https://amplication.com/images/discord_banner_purple.svg"/></a>

### Create a bug 🐞 report

If you see an error message or run into an issue, please [create bug report](https://github.com/memphisdev/memphis-broker/issues/new?assignees=&labels=type%3A%20bug&template=bug_report.md&title=). This effort is valued and it will help all Memphis{dev} users.


### Submit a feature 💡 request 

If you have an idea, or you think that we're missing a capability that would make development easier and more robust, please [Submit feature request](https://github.com/memphisdev/memphis-broker/issues/new?assignees=&labels=type%3A%20feature%20request).

If an issue❗with similar feature request already exists, don't forget to leave a "+1".
If you add some more information such as your thoughts and vision about the feature, your comments will be embraced warmly :)

## Contributing

Memphis{dev} is an open-source project.<br>
We are committed to a fully transparent development process and appreciate highly any contributions.<br>
Whether you are helping us fix bugs, proposing new features, improving our documentation or spreading the word - we would love to have you as part of the Memphis{dev} community.

Please refer to our [Contribution Guidelines](./CONTRIBUTING.md) and [Code of Conduct](./code_of_conduct.md).

## Contributors ✨

Thanks goes to these wonderful people ❤:<br><br>
 <a href = "https://github.com/memphisdev/memphis-broker/graphs/contributors">
   <img src = "https://contrib.rocks/image?repo=memphisdev/memphis-broker"/>
 </a>

## License 📃
Memphis is open-sourced and operates under the "Memphis Business Source License 1.0" license
Built out of Apache 2.0, the main difference between the licenses is:
"You may make use of the Licensed Work (i) only as part of your own product or service, provided it is not a message broker or a message queue product or service; and (ii) provided that you do not use, provide, distribute, or make available the Licensed Work as a Service. A “Service” is a commercial offering, product, hosted, or managed service, that allows third parties (other than your own employees and contractors acting on your behalf) to access and/or use the Licensed Work or a substantial set of the features or functionality of the Licensed Work to third parties as a software-as-a-service, platform-as-a-service, infrastructure-as-a-service or other similar services that compete with Licensor products or services."
Please check out [License](./LICENSE) to read the full text.
