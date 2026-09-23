# ssgo-frontend

# Features

Once the program server started, STUDENT, TEACHER and ADMIN user can login to system.

- STUDENT
  - prioritize their subject combinations.
  - confirm submission
  - print out reply slip for sign back by parents
- TEACHER
  - view student's application status (can be sorted by isConfirmed, classCode, classNo ...)
- ADMIN
  - includes features of TEACHER user
  - can upload students rank for subject allocation
  - can adjust subject capacity for subject allocation
  - can allocate students preference if the above 2 information is provided

### Schemas of `teacher.json` and `student.json`

```js
// src/data/teacher.json
[{
    "userAlias": String,
    "cname": String,
    "name": String
  }, ...


// src/data/student.json
[
  {
    "name": String,
    "cname": String,
    "classCode": String,
    "classNo": Number,
    "userAlias": String
  }, ...
```

## Project setup

```
npm install
```

### Compiles and hot-reloads for development

```
npm run serve
```

### Compiles and minifies for production

```
npm run build
```

### Lints and fixes files

```
npm run lint
```
