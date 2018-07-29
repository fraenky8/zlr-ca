# Simple REST-API  - Design Decisions
> "Document your architectural decisions and automate the setup" 

### App Structure
Tried to stick to golang standard package layout:

[https://github.com/golang-standards/project-layout]

[https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1]

Designed to match a clean architecture to enabling better testing (dependecy injections, decoupling), maintainance and extension.

### http-Framework
##### choice: 

gin [https://github.com/gin-gonic/gin]

##### why 
read several posts about other frameworks to compare but decided to stick with gin 
because i worked with it already and it is quite stable and popular

### Authorization
Using `BasicAuth` from the `gin` examples. To communicate with the api, use one of the following 
hardcoded accounts:
```
var accounts = gin.Accounts{
	"frank": "fr4nk!",
	"seb":   "th!rstY",
	"sarah": "!c3cre4M",
}
```
 example using `frank`:
 ```
 GET /icecreams/602 HTTP/1.1
 Host: localhost:8080
 Authorization: Basic ZnJhbms6ZnI0bmsh
 ```

### json-Response structure
##### choice
jsend [https://labs.omniti.com/labs/jsend]

##### why 
plain and simple, less overhead:

```
{
    "status": "[success|fail|error]",
    ["message": "some info here"]
    "data": {
        "icecreams": [
            { ... }
        ]
    }
}
```

### Database
##### choice: 
postgresql

##### why
"The World's Most Advanced Open Source Relational Database"

why relational: because I'm most experienced and family with

### Deployment
With `docker-compose` consisting of a `postgres` and an `zlrca` service. Database sets up with all data provided in 
`cmd/import/icecream.json`. Rest-Api is running default on Port `8080`.

### Improvements
- adding more tests
  - table driven tests / subtests
- implement CRUD for ingredients & sourcing values
- replace int to int64
- graceful shutdown
- moving repos in subfolder postgres because they implement postgres sql syntax and therefore cannot be reused



