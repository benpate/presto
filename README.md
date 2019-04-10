# Presto
## Magical REST interfaces for Go

Presto is a thin wrapper library that helps structure and simplify the REST interfaces you create in Go (Golang). 

Presto works side-by-side with any other routes that you create in your HTTP handler.  Using Presto, your code looks like this:

```
// Define a new service to expose online as a REST collection.
attachment := presto.NewCollection(f.AttachmentService, scope.All)

// Register HTTP methods to the service, including a list of permissions
e.GET("/attachments", attachment.List()) // public, no extra roles required
e.POST("/attachments", attachment.Post(role.InRoom)) // Must be "in room" to add a new attachment
e.GET("/attachments/:id", attachment.Get(role.InRoom)) // Must be "in room" to view an attachment
e.PUT("/attachments/:id", attachment.Put(role.InRoom, role.Owner)) // Must be "in room" and the "owner" of the attachment to update
e.DELETE("/attachments/:id", attachment.Delete(role.InRoom, role.Owner)) // Must be "in room" and the "owner" of the attachment to delete
```