# Presto

## Magical REST interfaces for Go

Presto is a thin wrapper library that helps structure and simplify the REST interfaces you create in [Go](https://golang.org). Its purpose is to encapsulate all of the boilerplate code that is commonly required to publish a server-side service via a REST interface.  Using Presto, your route configuration code looks like this:

```go
// Define a new service to expose online as a REST collection.
presto.NewCollection(echo.Echo, NoteFactory, "/notes").
    List().                          // Public.  No extra roles required
    Post(role.InRoom).               // Must be "in room" to add new notes.
    Get(role.InRoom).                // Must be "in room" to view existing notes.
    Put(role.InRoom, role.Owner).    // Must be "owner" to update notes.
    Patch(role.InRoom, role.Owner).  // Must be "owner" to update notes.
    Delete(role.InRoom, role.Owner). // Must be "owner" to delete notes.
    Method("action-name", customHandler, role.InRoom, role.CustomValue) // Custom POST action on this object
```

## Design Philosophy

### Clean Architecture

Presto lays the groundword to implement a REST API according to the [CLEAN architecture, first published by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).  This means decoupling business logic and databases, by injecting dependencies down through your application.  To do this in a type-safe manner, Presto requires that your services and objects fit into its interfaces, which  describe minimal behavior that each must support in order to be used by Presto. Once Presto is able to work with your business logic in an abstract way, the rest of the common code is repeated for each API endpoint you need to create.

### REST API Design Rulebook

Presto works hard to implement REST APIs according to the patterns laid out in the ["REST API Design Rulebook", by Mark Massé](https://smile.amazon.com/REST-Design-Rulebook-Mark-Masse/dp/1449310508/). This means:

* Clear route names
* Using HTTP methods (GET, PUT, POST, PATCH, DELETE) to determine the action being taken
* Using POST and URL route parameters for other API endpoints that don't fit neatly into the standard HTTP method definitions.

### Minimal Dependencies

Presto's only dependency is on the Echo router, which is a very fast, open-source router for creating HTTP servers in Go.  Our ultimate goal with this package is to remove this as a hard dependency eventually, and refactor this code to work with multiple routers in the Go ecosystem.

## Services

Presto does not replace your application business logic.  It only exposes your internal services via a REST API. Each endpoint must be linked to a corresponding
service (that matches Presto's required interface) to handle the actual loading, saving, and deleting of objects.

## Factories

The specific work of creating services and objects is pushed out to a Factory object, which provides a map of your complete domain.  The factories also manage dependencies (such as a live database connection) for each service that requires it.  Here's an example factory:

```go


```

## User Roles

It's very likely that your API requires custom authentication and authorization for each endpoint.  Since this is very custom to your application logic and environment, Presto can't automate this for you.  But, Presto does make it very easy to organize the permissions for each endpoint into a single, readable location.  Authorization requirements for each endpoint are baked into common functions called roles, and then passed in to Presto during system configuration.


```go

/// IN YOUR ROUTE CONFIGURATION

presto.NewCollection(echo.Echo, NoteFactory, "/notes").
    Post(role.InRoom) // The user must have permissions to post into the Room into which they're posting.

// IN YOUR ROLES PACKAGE

// InRoom determines if the requester has access to the Room in which this object resides.
// If so, then access to it is valid, so return a TRUE..  If not, then return a FALSE.
func InRoom(ctx echo.Context, object Object) bool {

    // Get the list of rooms that this user has access to..
    // For example, using JWT tokens in the context request headers.
    allowedRoomIDs := getRoomListFromContext(ctx)

    // Uses a type switch to retrieve the roomID from the Object interface.
    roomID, err := getRoomIDFromObject(object)

    if err != nil {
        return false
    }

    // Try to find the object.RoomID in the list of allowed rooms.
    for _, allowedRoomID := range allowedRoomIDs {
        if allowedRoomID == roomID {
            return true // If so, then you're in.
        }
    }

    // Otherwise, you are not permitted to access this object.
    return false;
}

```

## Selectors

## Boilerplate REST Endpoints

## Custom REST Endpoints