package data

import (
	"core/infra/data/uuid"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

//Constants
const (
	rowCountsMax = 100
)

func newEntityKey(ctx context.Context, entityKind string, uid uuid.UID) *datastore.Key {
	return datastore.NewKey(ctx, entityKind, string(uid), 0, nil)
}

func newParentKey(ctx context.Context, entityKind string, uid uuid.UID, parent *datastore.Key) *datastore.Key {
	return datastore.NewKey(ctx, entityKind, string(uid), 0, parent)
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

func storeEntityWithParent(ctx context.Context, rootKind string, rootUID uuid.UID, root interface{},
	childKind string, childUID uuid.UID, child interface{}) error {
	//
	rootKey := newEntityKey(ctx, rootKind, rootUID)
	childkey := newParentKey(ctx, childKind, childUID, rootKey)
	//
	return datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		//
		_, err := datastore.Put(ctx, rootKey, root)
		if err != nil {
			return err
		}
		_, err = datastore.Put(ctx, childkey, child)
		if err != nil {
			return err
		}
		return nil
	}, nil)
}

//remove datas from datastore
func deleteEntity(ctx context.Context, entityKind string, uid uuid.UID) error {
	key := newEntityKey(ctx, entityKind, uid)
	if err := datastore.Delete(ctx, key); err != nil {
		return ErrEntityNotDeleted.Original(err)
	}
	return nil
}

func findEntityByID(ctx context.Context, entityKind string, uid uuid.UID, ref interface{}) error {
	//key
	key := newEntityKey(ctx, entityKind, uid)
	//get entity
	err := datastore.Get(ctx, key, ref)
	if err != nil {
		return ErrEntityNotFound.Original(err)
	}
	return nil
}

func findOneEntityByFilters(ctx context.Context, entityKind string, filters map[string]interface{}, ref interface{}) error {
	//query initial
	q := datastore.NewQuery(entityKind)
	//filters
	if filters != nil {
		for k, v := range filters {
			q = q.Filter(k, v)
		}
	}
	//limit
	q = q.Limit(1)
	//get entity
	it := q.Run(ctx)
	for {
		_, err := it.Next(ref)
		if err != nil {
			return ErrEntityNotFound.Original(err)
		}
		break
	}
	//success
	return nil
}

//list entities with filters (optional)
func listEntitiesWithLimit(ctx context.Context, entityKind string, filters map[string]string, sort string, limit int, refs interface{}) error {
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
	if limit > 0 {
		q = q.Limit(limit)
	} else {
		q = q.Limit(rowCountsMax)
	}

	if _, err := q.GetAll(ctx, refs); err != nil {
		return ErrEntitiesNotListed.Original(err)
	}
	//result with success
	return nil
}

//list entities with filters (optional)
func listEntities(ctx context.Context, entityKind string, filters map[string]string, sort string, refs interface{}) error {
	return listEntitiesWithLimit(ctx, entityKind, filters, sort, 0, refs)
}
