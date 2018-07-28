## Design Decisions
"Document your architectural decisions" 

#### App Structure
tried to stick to golang standard package layout 

[https://github.com/golang-standards/project-layout]

[https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1]

#### http-Framework
##### choice: 

gin [https://github.com/gin-gonic/gin]

##### why 
read several posts about other frameworks to compare but decided to stick with gin 
because i worked with it already and it is quite stable and popular

#### Database
##### choice: 
postgresql

##### why
because of reasons

#### Improvements
- adding more tests
- abstract / encapsulate database/storage
- implement CRUD for ingredients & sourcing values
- replace int to int64



