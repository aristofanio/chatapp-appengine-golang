package data

import (
	"context"
	"core/data/uuid"
	"core/err"

	"google.golang.org/appengine/datastore"
)

//Constants
const (
	rowCountsMax = 100
)

//Errors
var (
	ErrEntityNotFound    = err.NewErr(1001, "Entidade não encontrada")
	ErrEntityNotDeleted  = err.NewErr(1002, "Entidade não pode ser excluída")
	ErrEntitiesNotListed = err.NewErr(1003, "Entidades não podem ser listadas")
)

func newEntityKey(ctx context.Context, entityKind string, uid uuid.UID) *datastore.Key {
	return datastore.NewKey(ctx, entityKind, string(uid), 0, nil)
}

//store datas from datastore
func storeEntity(ctx context.Context, entityKind string, uid uuid.UID, e interface{}) error {
	key := newEntityKey(ctx, entityKind, uid)
	_, err := datastore.Put(ctx, key, e)
	if err != nil {
		return err
	}
	return nil
}

//remove datas from datastore
func delete(ctx context.Context, entityKind string, uid uuid.UID) error {
	key := newEntityKey(ctx, entityKind, uid)
	if err := datastore.Delete(ctx, key); err != nil {
		return ErrEntityNotDeleted.Original(err)
	}
	return nil
}

func findEntity(ctx context.Context, entityKind string, uid uuid.UID, ref interface{}) error {
	//key
	key := newEntityKey(ctx, entityKind, uid)
	//get entity
	err := datastore.Get(ctx, key, ref)
	if err != nil {
		return ErrEntityNotFound.Original(err)
	}
	return nil
}

//list entities with filters (optional)
func listEntities(ctx context.Context, entityKind string, filters map[string]interface{}, sort string, refs interface{}) error {
	//query initial
	q := datastore.NewQuery(entityKind)
	//filters
	if filters != nil {
		for k, v := range filters {
			q = q.Filter(k, v)
		}
	}
	//sort query
	if sort != "" {
		q = q.Order(sort)
	}
	//limit
	q = q.Limit(rowCountsMax)
	if _, err := q.GetAll(ctx, refs); err != nil {
		return ErrEntitiesNotListed.Original(err)
	}
	//result with success
	return nil
}
