# Conventions

## Resolver

Maps a GraphQL query to a service method. Validates the input values and types,
but no business logic.

## Command and query

The only place in the application that may call the repository.

## Naming conventions

These names are supposed to help people understand the codebase by recognizing
names from other projects, but are not necessary to follow.

### Commands

Action + Subject
: e.g. `SendStudentGrades`

### Queries

Action + Quantity + Subject + Filter (optional)
: e.g. `FindOneStudentById`, `FindManyStudents`

### Events

TODO

### Entities

TODO

### GraphQL queries and mutations

Who + Action + Subject + `s` if plural
: e.g. `PrincipalListStudentsGrades`

### Controllers

TODO

- Input:
Action + Subject + Input
: e.g. `ListStudentsGradesInput`
