package presto

import "github.com/benpate/data/journal"

type testPerson struct {
	PersonID        string `bson:"personId"`
	Name            string `bson:"name"`
	Email           string `bson:"email"`
	Age             int    `bson:"age"`
	journal.Journal `bson:"journal"`
}

func (testPerson *testPerson) ID() string {
	return testPerson.PersonID
}
