package bigtable

const (
	addedDocRefTypeByDocID       = "addedDocRefTypeByDocID"
	addedDocRefTypeByQueryPrefix = "addedDocRefTypeByQueryPrefix"
)

// addedDocRef is a reference to a previously saved document. It includes the
// appropriate Bigtable table can be used and the row key of the Bigtable row
// so that the right row within that table can be deleted.
type addedDocRef struct {
	docID   string
	refType string
	rowKey  string
}
