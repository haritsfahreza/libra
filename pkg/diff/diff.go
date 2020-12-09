package diff

//ChangeType represents the objects change type
type ChangeType string

const (
	//New represents new object change type
	New ChangeType = "new"

	//Removed represents removed object change type
	Removed ChangeType = "removed"

	//Changed represents changed object change type
	Changed ChangeType = "changed"
)

//Diff represents the different between two objects
type Diff struct {
	ChangeType ChangeType  `json:"change_type"`
	ObjectType string      `json:"object_type"`
	ObjectID   string      `json:"object_id"`
	Field      string      `json:"field,omitempty"`
	Old        interface{} `json:"old,omitempty"`
	New        interface{} `json:"new,omitempty"`
}
