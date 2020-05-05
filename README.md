[![Go Report Card](https://goreportcard.com/badge/github.com/STreeChin/contactapi)](https://goreportcard.com/report/github.com/STreeChin/contactapi)

# Instruction

### Code structure

- Follow Domain Driven Design  and Clean Architecture  of Robert Martin.
- Follow the golang-standards' project-layout, seperate the internal service and public library. For the internal service, divide services vertically  by business domain. The service is divided into layers.
- Interfaces are provided between layers, and interface definitions are placed on the upper layer, that is, the interface is owned and defined by the user, not the implementer. Avoid affecting the user due to the modification of the implementer.

### Unit Test

- Follow TDD.
- Use GoConvey as the UT framework, use GoMock to mock the interface, and Monkey to mock the function. 
- View the coverage and coverage code details in html.

### Errors

- The principle is to wrap the error and throw it to the upper layer to deal with it.
- At this stage, the "pkg/errors" is preferred over the official "errors". Consider switching to the official package after Go 2.0 is released.

### Secrets

- Use docker secret to mange secrets and configuration by environment variables. 
- After Dockerkit becomes more available, it will switch smoothly.

# How To Run

Put the source code in any folder.

Start up the application by running `docker-compose up -d` from the deployment directory of the project.

```
D:\contactapi\deployment>docker-compose up -d
......
Successfully built e3e3297c162d
Successfully tagged contact:1.0.0
WARNING: Image for service contact was built because it did not already exist. To rebuild this image you must use `docker-compose build` or `docker-compose up --build`.
Creating redis ... done                                                             Creating mongo ... done                                                             Creating deployment_contact_1 ... done  
```

After that,  the app service, database, and cache is running. 

```
D:\contactapi\deployment>docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                      NAMES
a3dbd68e43bd        contact:1.0.0       "./contact"              47 minutes ago      Up 47 minutes       0.0.0.0:8080->8080/tcp     deployment_contact_1
c3d61a795172        mongo:4.2.6         "docker-entrypoint.s…"   47 minutes ago      Up 47 minutes       0.0.0.0:27017->27017/tcp   mongo
954c099cd831        redis:alpine3.11    "docker-entrypoint.s…"   47 minutes ago      Up 47 minutes       0.0.0.0:6379->6379/tcp     redis
```

There are 5 images.

```
D:\contactapi\deployment>docker images
REPOSITORY      TAG                 IMAGE ID            CREATED            SIZE
contact         1.0.0               e3e3297c162d        2 minutes ago      20.8MB
<none>          <none>              9b1e8104f885        2 minutes ago      726MB
mongo           4.2.6               3f3daf863757        3 days ago         388MB
redis           alpine3.11          3661c84ee9d0        3 days ago         29.8MB
golang          1.14.2-alpine3.11   dda4232b2bd5        3 days ago         370MB
```

# How To Test

In simple command line client, like Git Bash, consume the APIs provided by the server.

- Post

```
$ curl --include \
>      --request POST \
>      --header "autopilotapikey: 65263027fab7d440ba4c5f3b834fb800" \
>      --header "Content-Type: application/json" \
>      --data-binary "{
>     \"contact\": {
>         \"FirstName\": \"Slarty\",
>         \"LastName\": \"Bartfast\",
>         \"Email\": \"test@slarty.com\",
>         \"custom\": {
>             \"string--Test--Field\": \"This is a test\"
>         }
>
>   }
>
> }" \
> 'http://127.0.0.1:8080/v1/contact'

  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   272  100    65  100   207   3250  10350 --:--:-- --:--:-- --:--:-- 13600HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 04 May 2020 12:48:55 GMT
Content-Length: 65

{"contact_id":"person_AP2-8648b71c-b4c5-4b31-9318-7dad6ba302b6"}
```

- Get by Email

```
$ curl --include \
>      --header "autopilotapikey: 65263027fab7d440ba4c5f3b834fb800" \
>   'http://127.0.0.1:8080/v1/contact/test@slarty.com'

  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   594  100   594    0     0  66000      0 --:--:-- --:--:-- --:--:-- 74250HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 04 May 2020 12:51:46 GMT
Content-Length: 594

{"contact_id":"person_AP2-8648b71c-b4c5-4b31-9318-7dad6ba302b6","Email":"test@slarty.com","Twitter":"","FirstName":"Slarty","LastName":"Bartfast","Salutation":"","Company":"","NumberOfEmployees":"","Title":"","Industry":"","Phone":"","MobilePhone":"","Fax":"","Website":"","MailingStreet":"","MailingCity":"","MailingState":"","MailingPostalCode":"","MailingCountry":"","LeadSource":"","Status":"","LinkedIn":"","lists":null,"type":"","created_at":"","updated_at":"","owner_name":"","unsubscribed":false,"custom":{"Test Field":null},"_autopilot_session_id":"","_autopilot_list":"","notify":""}
```

- Get by contactID

```
$ curl --include      --header "autopilotapikey: 65263027fab7d440ba4c5f3b834fb800"   'http://127.0.0.1:8080/v1/contact/person_AP2-8648b71c-b4c5-4b31-9318-7dad6ba302b6'

  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   594  100   594    0     0   193k      0 --:--:-- --:--:-- --:--:--  193kHTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 04 May 2020 12:52:31 GMT
Content-Length: 594

{"contact_id":"person_AP2-8648b71c-b4c5-4b31-9318-7dad6ba302b6","Email":"test@slarty.com","Twitter":"","FirstName":"Slarty","LastName":"Bartfast","Salutation":"","Company":"","NumberOfEmployees":"","Title":"","Industry":"","Phone":"","MobilePhone":"","Fax":"","Website":"","MailingStreet":"","MailingCity":"","MailingState":"","MailingPostalCode":"","MailingCountry":"","LeadSource":"","Status":"","LinkedIn":"","lists":null,"type":"","created_at":"","updated_at":"","owner_name":"","unsubscribed":false,"custom":{"Test Field":null},"_autopilot_session_id":"","_autopilot_list":"","notify":""}
```

