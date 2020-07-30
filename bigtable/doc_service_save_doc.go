package bigtable

import (
	"cloud.google.com/go/bigtable"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mattwelke/manydocs/utils"
)

const (
	IDPropName = "_id"
)

func (service DocService) SaveDoc(newDoc map[string]interface{}, newDocQueryPrefixes []string) (string, error) {
	newDocID := utils.NewID()
	newDoc[IDPropName] = newDocID

	// JSON encode doc
	newDocJSON, err := json.Marshal(newDoc)
	if err != nil {
		return "", fmt.Errorf("could not JSON encode doc for save doc operation: %v", err)
	}

	insertMut := bigtable.NewMutation()
	insertMut.Set("value", "value", bigtable.Now(), newDocJSON)

	if err := service.docsByDocsIDTable.Apply(context.Background(), newDocID, insertMut); err != nil {
		return "", fmt.Errorf("could not apply Bigtable insert mutation for docs by doc ID table for save doc operation: %v", err)
	}

	addedDocRefs := make([]addedDocRef, 0)

	// Save ref to doc in docs by doc ID table for "delete doc" operation later
	addedDocRefs = append(addedDocRefs, addedDocRef{
		docID:   newDocID,
		refType: addedDocRefTypeByDocID,
		rowKey:  newDocID,
	})

	for _, prefix := range newDocQueryPrefixes {
		finalRowKey := fmt.Sprintf("%s%s", prefix, utils.NewID())

		insertMut := bigtable.NewMutation()
		insertMut.Set("value", "value", bigtable.Now(), newDocJSON)

		if err := service.docsByQueryPrefixTable.Apply(context.Background(), finalRowKey, insertMut); err != nil {
			return "", fmt.Errorf("could not apply Bigtable insert mutation for docs by query prefix table for save doc operation: %v", err)
		}

		// Save ref to doc in docs by query prefix table for "delete doc" operation later
		addedDocRefs = append(addedDocRefs, addedDocRef{
			docID:   newDocID,
			refType: addedDocRefTypeByQueryPrefix,
			rowKey:  finalRowKey,
		})
	}

	for _, ref := range addedDocRefs {
		// row key is the doc ID with a UUID appended to it, so we can get all doc refs later by doc ID
		finalDocInsertPrimaryKey := fmt.Sprintf("%s%s", ref.docID, utils.NewID())

		insertMut := bigtable.NewMutation()
		insertMut.Set("data", "ref_type", bigtable.Now(), []byte(ref.refType))
		insertMut.Set("data", "row_key", bigtable.Now(), []byte(ref.rowKey))

		if err := service.addedDocRefsTable.Apply(context.Background(), finalDocInsertPrimaryKey, insertMut); err != nil {
			return "", fmt.Errorf("could not apply Bigtable insert mutation for added doc refs table for save doc operation: %v", err)
		}
	}

	return newDocID, nil
}
