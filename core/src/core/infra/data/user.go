package data

import (
	"bytes"
	"context"
	"core/infra/data/uuid"
	"core/utils"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//-------------------------------------------------------------------
// User Data
//-------------------------------------------------------------------

const userKind = "User"

type User struct {
	ID         uuid.UID `datastore:"mid"`
	Nick       string   `datastore:"nick"`
	Email      string   `datastore:"email"`
	Photo      string   `datastore:"photo"`
	PassKey    string   `datastore:"passkey"`
	PassToken  string   `datastore:"passtoken"`
	IsVerified bool     `datastore:"is_verified"`
	IsBlocked  bool     `datastore:"is_blocked"`
	CreatedIn  int64    `datastore:"created_in,noindex"`
	UpdatedIn  int64    `datastore:"updated_in,noindex"`
}

func (u *User) SetPassKey(password string) {
	//calculated
	buf := bytes.NewBufferString(u.Email)
	buf.WriteString(password)
	hasher := md5.New()
	hasher.Write(buf.Bytes())
	data := hasher.Sum(nil)
	enc := base64.StdEncoding.EncodeToString(data)
	u.PassKey = enc
}

func (u *User) CheckPassword(password string) bool {
	//calculated
	buf0 := bytes.NewBufferString(u.Email)
	buf0.WriteString(password)
	hasher0 := md5.New()
	hasher0.Write(buf0.Bytes())
	data0 := hasher0.Sum(nil)
	enc0 := base64.StdEncoding.EncodeToString(data0)
	//expected
	enc1 := u.PassKey
	//result
	return enc0 == enc1
}

func (u *User) CheckPassToken(ptoken string) bool {
	//compare pass tokens
	if u.PassToken != ptoken {
		return false
	}
	//decoding - exp
	expected, err := base64.StdEncoding.DecodeString(u.PassToken)
	if err != nil {
		return false
	}
	//extract time
	exp := string(expected)
	values := strings.Split(exp, ":")
	if len(values) != 2 {
		return false
	}
	i64, err := strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		return false
	}
	crr := utils.Now()
	//
	return crr < i64
}

func (u *User) ResetPassToken() string {
	//calculate data
	r := rand.New(rand.NewSource(99))
	t0 := r.Uint32()
	t1 := time.Now().Add(3 * time.Minute).UnixNano()
	//data
	c := fmt.Sprintf("%d:%d", t0, t1)
	//pass token
	u.PassToken = base64.StdEncoding.EncodeToString([]byte(c))
	//result
	return u.PassToken
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

func (m UserDataMgr) NewUser(id uuid.UID, nick, email, password string) *User {
	//create instance
	u := new(User)
	u.ID = id
	u.Nick = nick
	u.Email = email
	u.PassKey = password
	u.IsBlocked = false
	u.CreatedIn = utils.Now()
	u.UpdatedIn = utils.Now()
	//result
	return u
}

func (m UserDataMgr) Get(uid uuid.UID) (*User, error) {
	gu := new(User)
	if err := findEntityByID(m.ctx, userKind, uid, gu); err != nil {
		return nil, ErrEntityNotFound.Original(err)
	}
	return gu, nil
}

func (m UserDataMgr) GetByNick(nick string) (*User, error) {
	//
	filter := make(map[string]interface{})
	filter["nick="] = nick
	//
	gu := new(User)
	if err := findOneEntityByFilters(m.ctx, userKind, filter, gu); err != nil {
		return nil, err
	}
	return gu, nil
}

func (m UserDataMgr) GetByEmail(email string) (*User, error) {
	//
	filter := make(map[string]interface{})
	filter["email="] = email
	//
	gu := new(User)
	if err := findOneEntityByFilters(m.ctx, userKind, filter, gu); err != nil {
		return nil, err
	}
	return gu, nil
}

func (m UserDataMgr) Store(u *User) error {
	return storeEntity(m.ctx, userKind, u.ID, u)
}

func (m UserDataMgr) Block(uid uuid.UID) error {
	gu, err := m.Get(uid)
	if err != nil {
		return err
	}
	gu.IsBlocked = true
	gu.UpdatedIn = utils.Now()
	return m.Store(gu)
}

func NewUserDataMgr(ctx context.Context) UserDataMgr {
	return UserDataMgr{ctx}
}
