# Presto
## Magical REST interfaces for Go 

Presto is a thin wrapper library that helps structure and simplify the REST interfaces you create in Go (Golang). 

Presto works side-by-side with any other routes that you create in your HTTP handler.  Using Presto, your code looks like this:

```go
// Define a new service to expose online as a REST collection.
note := presto.NewCollection(echo.Echo, NoteService, "/notes").
    List(nil).                       // Public.  No extra roles required
    Post(role.InRoom).               // Must be "in room" to add new notes.
    Get(role.InRoom).                // Must be "in room" to view existing notes.
    Put(role.InRoom, role.Owner).    // Must be "owner" to update notes.
    Patch(role.InRoom, role.Owner).  // Must be "owner" to update notes.
    Delete(role.InRoom, role.Owner). // Must be "owner" to delete notes.


// Register HTTP methods to the service, including a list of permissions
e.GET("/notes", note.List()) // public, no extra roles required
e.POST("/notes", note.Post(role.InRoom)) // Must be "in room" to add a new note
e.GET("/notes/:id", note.Get(role.InRoom)) // Must be "in room" to view an note
e.PUT("/notes/:id", note.Put(role.InRoom, role.Owner)) // "in room" and "owner" of the note to update
e.DELETE("/notes/:id", note.Delete(role.InRoom, role.Owner)) // "in room" and "owner" of the attachment to delete
```