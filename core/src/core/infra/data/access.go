package data

import (
	"context"
	"core/infra/data/uuid"
	"core/utils"
	"errors"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

//------------------------------------------------------------------
// Access Data
//------------------------------------------------------------------

const accessKind = "Access"

type Access struct {
	AToken    uuid.UID           `datastore:"atoken"`
	GuestID   uuid.UID           `datastore:"guestid"`
	NickName  string             `datastore:"nick"`
	Photo     string             `datastore:"photo"`
	FCMToken  string             `datastore:"fcm"`
	Position  appengine.GeoPoint `datastore:"position"`
	IsValid   bool               `datastore:"is_valid"`
	CreatedIn int64              `datastore:"created_in"`
	UpdatedIn int64              `datastore:"updated_in"`
}

func (a *Access) MakeAsInvalid() {
	a.IsValid = false
	a.UpdatedIn = utils.Now()
}

//------------------------------------------------------------------
// Access Data Manager
//------------------------------------------------------------------

type AccessDataMgr struct {
	ctx context.Context
}

func (mgr AccessDataMgr) NewAccess(g Guest, u User, fcmToken string, lat, lng float64) (*Access, error) {
	//
	position := appengine.GeoPoint{Lat: lat, Lng: lng}
	//
	ac := new(Access)
	ac.AToken = uuid.NewUID("session")
	ac.GuestID = g.ID
	ac.NickName = u.NickName
	ac.Photo = u.Picture
	ac.Position = position
	ac.FCMToken = fcmToken
	ac.IsValid = true
	ac.CreatedIn = utils.Now()
	ac.UpdatedIn = utils.Now()
	//
	err := mgr.Store(ac)
	if err != nil {
		return nil, err
	}
	//
	return ac, nil
}

func (mgr AccessDataMgr) Store(m *Access) error {
	return storeEntity(mgr.ctx, accessKind, m.AToken, m)
}

func (mgr AccessDataMgr) GetAccess(stoken uuid.UID) (*Access, error) {
	//result
	rslt := &Access{}
	//get entity
	err := findEntityByID(mgr.ctx, memberKind, stoken, rslt)
	if err != nil {
		return nil, ErrEntityNotFound
	}
	return rslt, nil
}

func (mgr AccessDataMgr) GetAccessByGuest(guest Guest) (*Access, error) {
	//
	dst := make([]Access, 0, 1)
	//
	var q = datastore.NewQuery(accessKind).
		Filter("is_valid=", true).
		Filter("guestid=", guest.ID).
		Limit(1)
	//
	if _, err := q.GetAll(mgr.ctx, &dst); err != nil {
		return nil, errors.New("Erro ao tentar recuperar uma lista de entidades. " +
			"Original(" + err.Error() + ").")
	}
	//
	s := len(dst)
	if s == 0 {
		return nil, nil
	}
	return &dst[0], nil
}

func (mgr AccessDataMgr) InvalidAccess(a *Access) error {
	if a != nil {
		a.MakeAsInvalid()
		return mgr.Store(a)
	}
	return ErrEntityMustBeNotNull
}

func (mgr AccessDataMgr) InvalidOldAccess(guest Guest) error {
	access, _ := mgr.GetAccessByGuest(guest)
	return mgr.InvalidAccess(access)
}

func (mgr AccessDataMgr) InvalidAllOldAccess(guest Guest) error {
	//recover user
	q := datastore.NewQuery(accessKind).
		Filter("guestid=", string(guest.ID)).
		Limit(100)
	//recuperar os acessos e invalida-los
	it := q.Run(mgr.ctx)
	for {
		//retrieve old access
		old := new(Access)
		_, err := it.Next(old)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return err
		}
		//invalid old access
		if err = mgr.InvalidAccess(old); err != nil {
			return err
		}
	}
	return nil
}

// func InvalidAllOldAccessByExpeired(ctx appengine.Context) {
// 	//
// 	//recover user
// 	expire := util.Now() - (1000 * 3600 * 48)
// 	var q = datastore.NewQuery(accessEntityName).Filter("is_valid=", true).Filter("created_in<=", expire).Limit(10)
// 	//recuperar os acessos e invalida-los
// 	it := q.Run(ctx)
// 	for {
// 		//
// 		dst := new(Access)
// 		_, err := it.Next(dst)
// 		if err != nil {
// 			break
// 		}
// 		if err == datastore.Done {
// 			break
// 		}
// 		//
// 		dst.MakeAsInvalid(ctx)
// 	}
// }

// func ListLatestAccess(ctx appengine.Context, count int) ([]Access, error) {
// 	es := make([]Access, 0, count)
// 	if err := listAllEntities(ctx, accessEntityName, "-created_in", count, &es); err != nil {
// 		return nil, errors.New("List not found. Original-> " + err.Error())
// 	}
// 	return es, nil
// }

// func ListLatestAccessActivies(ctx appengine.Context, guestId string, count int) ([]Access, error) {
// 	//
// 	dst := make([]Access, 0, count)
// 	//
// 	var q = datastore.NewQuery(accessEntityName).
// 		Filter("is_valid=", true).
// 		Order("-created_in").
// 		Limit(count)
// 	//
// 	if _, err := q.GetAll(ctx, &dst); err != nil {
// 		return []Access{}, errors.New("Erro ao tentar recuperar uma lista de entidades. " +
// 			"Original(" + err.Error() + ").")
// 	}
// 	//
// 	return dst, nil
// }

// func ListAccessByUser(ctx appengine.Context, guestId string, count int) ([]Access, error) {
// 	//
// 	dst := make([]Access, 0, count)
// 	//
// 	var q = datastore.NewQuery(accessEntityName).
// 		Filter("guestid=", guestId).
// 		Filter("is_valid=", true).
// 		Order("-created_in").
// 		Limit(count)
// 	//
// 	if _, err := q.GetAll(ctx, &dst); err != nil {
// 		return []Access{}, errors.New("Erro ao tentar recuperar uma lista de entidades. " +
// 			"Original(" + err.Error() + ").")
// 	}
// 	//
// 	return dst, nil
// }

func NewAccessDataMgr(ctx context.Context) AccessDataMgr {
	return AccessDataMgr{ctx}
}
