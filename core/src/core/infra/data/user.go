package data

import (
	"bytes"
	"context"
	"core/infra/data/uuid"
	"core/utils"
	"crypto/md5"
	"encoding/base64"
)

//-------------------------------------------------------------------
// User Data
//-------------------------------------------------------------------

const userKind = "User"

type User struct {
	ID        uuid.UID `datastore:"guestid"`
	NickName  string   `datastore:"nick"`
	Email     string   `datastore:"email"`
	Name      string   `datastore:"name"`
	PassKey   string   `datastore:"passkey"`
	FirstName string   `datastore:"first_name"`
	LastName  string   `datastore:"last_name"`
	Gender    string   `datastore:"gender"`
	Picture   string   `datastore:"picture"`
	Avatar    string   `datastore:"avatar"`
	Age       int64    `datastore:"age"`
	Verified  bool     `datastore:"verified"`
	IsBlocked bool     `datastore:"is_blocked"`
	CreatedIn int64    `datastore:"created_in,noindex"`
	UpdatedIn int64    `datastore:"updated_in,noindex"`
}

func (u *User) CheckPassword(password string) bool {
	//expected
	buf0 := bytes.NewBufferString(u.Email)
	buf0.WriteString(password)
	hasher0 := md5.New()
	hasher0.Write(buf0.Bytes())
	data0 := hasher0.Sum(nil)
	enc0 := base64.StdEncoding.EncodeToString(data0)
	//calculated
	buf1 := bytes.NewBufferString(u.Email)
	buf1.WriteString(password)
	hasher1 := md5.New()
	hasher1.Write(buf1.Bytes())
	data1 := hasher1.Sum(nil)
	enc1 := base64.StdEncoding.EncodeToString(data1)
	//result
	return enc0 == enc1
}

//-------------------------------------------------------------------
// Methods - implementations of Entity
//-------------------------------------------------------------------

// func (p User) Id() string {
//     return p.GuestID
// }

// func (p User) EName() string {
//     return userEntityName
// }

// func (p *User) Save(ctx appengine.Context) error{
//     return save(ctx, p)
// }

// //-------------------------------------------------------------------
// // Operations
// //-------------------------------------------------------------------

// func GetUser(ctx appengine.Context, id string) *User {
//     e := new(User)
//     if err := findEntity(ctx, userEntityName, id, e); err != nil {
//         return nil
//     }
//     return e
// }

// func GetUserByFBId(ctx appengine.Context, fbid string) *User {
//     //
//     var q = datastore.NewQuery(userEntityName).Filter("fbid=", fbid).Limit(1)
//     //
//     dst := make([]User, 0, 1)
//     if _, err := q.GetAll(ctx, &dst); err != nil {
//         return nil
//     }
//     //
//     if len(dst) > 0 {
//         return &dst[0]
//     } else {
//         return nil
//     }
// }

// func UserCount(ctx appengine.Context) int {
//     var q = datastore.NewQuery(userEntityName)
//     if count, err := q.Count(ctx); err == nil {
//         return count
//     }
//     return 0
// }

//------------------------------------------------------------------
// Guest Data Manager
//------------------------------------------------------------------

type UserDataMgr struct {
	ctx context.Context
}

func (m UserDataMgr) NewGuest() *User {
	//create instance
	gu := new(User)
	gu.ID = uuid.NewUID("guest")
	gu.IsBlocked = false
	gu.CreatedIn = utils.Now()
	gu.UpdatedIn = utils.Now()
	//result
	return gu
}

func (m UserDataMgr) GetUser(uid uuid.UID) *User {
	gu := new(User)
	if err := findEntityByID(m.ctx, guestKind, uid, gu); err != nil {
		return nil
	}
	return gu
}

func (m UserDataMgr) GetUserByNick(nick string) (*User, error) {
	//
	filter := make(map[string]interface{})
	filter["nick="] = nick
	//
	gu := new(User)
	if err := findOneEntityByFilters(m.ctx, guestKind, filter, gu); err != nil {
		return nil, err
	}
	return gu, nil
}

func (m UserDataMgr) GetUserByEmail(email string) (*User, error) {
	//
	filter := make(map[string]interface{})
	filter["email="] = email
	//
	gu := new(User)
	if err := findOneEntityByFilters(m.ctx, guestKind, filter, gu); err != nil {
		return nil, err
	}
	return gu, nil
}

func (m UserDataMgr) Store(u *User) error {
	return storeEntity(m.ctx, userKind, u.ID, u)
}

func (m UserDataMgr) Block(uid uuid.UID) error {
	gu := m.GetUser(uid)
	gu.IsBlocked = true
	gu.UpdatedIn = utils.Now()
	return m.Store(gu)
}

func NewUserDataMgr(ctx context.Context) UserDataMgr {
	return UserDataMgr{ctx}
}
