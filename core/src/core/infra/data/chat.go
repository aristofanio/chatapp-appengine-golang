package data

import (
	"core/infra/data/uuid"
	"core/utils"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
)

//-------------------------------------------------------------------
// Chat Data
//-------------------------------------------------------------------

const chatKind = "Chat"

type ChatMember struct {
	Member   uuid.UID           `datastore:"member"`
	Nick     string             `datastore:"nick"`
	Photo    string             `datastore:"photo"`
	FCMToken string             `datastore:"fcm"`      //(by latest access)
	Position appengine.GeoPoint `datastore:"position"` //(by latest access)
}

type Chat struct {
	ID          uuid.UID   `datastore:"id"`
	From        ChatMember `datastore:"from"`
	To          ChatMember `datastore:"to"`
	IsBlocked   bool       `datastore:"is_blocked"`
	IsActive    bool       `datastore:"is_active"`
	IsVisible   bool       `datastore:"is_visible"`
	Offset      int64      `datastore:"offset"`
	LastContact int64      `datastore:"last_contact"`
	CreatedIn   int64      `datastore:"created_in,noindex" json:"create_in"`
	UpdatedIn   int64      `datastore:"updated_in,noindex" json:"update_in"`
}

//-------------------------------------------------------------------
// Chat Manager Data
//-------------------------------------------------------------------

type ChatDataMgr struct {
	ctx context.Context
}

func (m ChatDataMgr) NewChat(inviteID uuid.UID, to, from Member) (*Chat, error) {
	//members
	chatTo := ChatMember{
		Member: to.ID,
		Nick:   to.Nick,
		Photo:  to.Photo,
	}
	chatFrom := ChatMember{
		Member: from.ID,
		Nick:   from.Name,
		Photo:  from.Photo,
	}
	//chats
	chat := new(Chat)
	chat.ID = inviteID
	chat.To = chatTo
	chat.From = chatFrom
	chat.IsBlocked = false
	chat.IsActive = true
	chat.IsVisible = true //??
	chat.Offset = 0
	chat.LastContact = 0
	chat.CreatedIn = utils.Now()
	chat.UpdatedIn = utils.Now()
	//
	err := storeEntity(m.ctx, chatKind, inviteID, chat)
	if err != nil {
		return nil, err
	}
	//
	return chat, nil
}

// func (m ChatDataMgr) GetChat(ctx appengine.Context, id string) *Chat {
// 	e := new(Chat)
// 	if err := findEntity(ctx, chatEntityName, id, e); err != nil {
// 		return nil
// 	}
// 	return e
// }

// func GetChatByPair(ctx appengine.Context, fromId, toId string) *Chat {
// 	//
// 	dst := make([]Chat, 0, 1)
// 	//
// 	var q = datastore.NewQuery(chatEntityName).
// 		Filter("is_blocked=", false).
// 		Filter("is_active=", true).
// 		Filter("is_visible=", true).
// 		Filter("from.guestid=", fromId).
// 		Filter("to.guestid=", toId).
// 		Limit(1)
// 	//
// 	if _, err := q.GetAll(ctx, &dst); err != nil {
// 		return nil
// 	}
// 	//
// 	if len(dst) == 0 {
// 		return nil
// 	} else {
// 		return &dst[0]
// 	}
// }

// func InactivateChat(ctx appengine.Context, fromId, toId string) error {
// 	//desativar todos os chats 'ativos' com now() - last_contact > (600*1000*1000) --> 10m
// 	//tornar os chats correspondentes como invisível
// 	//e desativar todos os chats 'ativos'
// 	dst := make([]Chat, 0, 1000)
// 	//
// 	var q = datastore.NewQuery(chatEntityName).
// 		Filter("to.guestid=", toId).
// 		Filter("from.guestid=", fromId).
// 		Filter("is_active=", true).
// 		Limit(1000)
// 	//
// 	if _, err := q.GetAll(ctx, dst); err != nil {
// 		return errors.New("Erro ao tentar recuperar uma lista de entidades. " +
// 			"Original(" + err.Error() + ").")
// 	}
// 	//
// 	s := len(dst)
// 	for i := 0; i < s; i++ {
// 		//current
// 		r := dst[i]
// 		r.IsActive = false
// 		r.UpdatedIn = util.Now()
// 		r.Save(ctx)
// 	}
// 	//
// 	return nil
// }

// func BlockChat(ctx appengine.Context, fromId, toId string) error {
// 	//bloquear todos os chats 'não bloqueados' cuja correspondência exista
// 	//e desativar todos os chats 'ativos'
// 	//tornar os chats correspondentes como invisível
// 	//e desativar todos os chats 'ativos'
// 	dst := make([]Chat, 0, 1000)
// 	//
// 	var q = datastore.NewQuery(chatEntityName).
// 		Filter("to.guestid=", toId).
// 		Filter("from.guestid=", fromId).
// 		Filter("is_blocked=", false).
// 		Limit(1000)
// 	//
// 	if _, err := q.GetAll(ctx, &dst); err != nil {
// 		return errors.New("Erro ao tentar recuperar uma lista de entidades. " +
// 			"Original(" + err.Error() + ").")
// 	}
// 	//
// 	s := len(dst)
// 	for i := 0; i < s; i++ {
// 		//current
// 		r := dst[i]
// 		r.IsVisible = false
// 		r.IsActive = false
// 		r.IsBlocked = true
// 		r.UpdatedIn = util.Now()
// 		r.Save(ctx)
// 	}
// 	//
// 	return nil
// }

// func MakeInvisibleChat(ctx appengine.Context, fromId, toId string) error {
// 	//tornar os chats correspondentes como invisível
// 	//e desativar todos os chats 'ativos'
// 	dst := make([]Chat, 0, 1000)
// 	//
// 	var q = datastore.NewQuery(chatEntityName).
// 		Filter("to.guestid=", toId).
// 		Filter("from.guestid=", fromId).
// 		Filter("is_visible=", true).
// 		Limit(1000)
// 	//
// 	if _, err := q.GetAll(ctx, &dst); err != nil {
// 		return errors.New("Erro ao tentar recuperar uma lista de entidades. " +
// 			"Original(" + err.Error() + ").")
// 	}
// 	//
// 	s := len(dst)
// 	for i := 0; i < s; i++ {
// 		//current
// 		r := dst[i]
// 		r.IsVisible = false
// 		r.IsActive = false
// 		r.IsBlocked = true
// 		r.UpdatedIn = util.Now()
// 		r.Save(ctx)
// 	}
// 	//
// 	return nil
// }

// func ListChats(ctx appengine.Context, guestId string) ([]Chat, error) {
// 	//listar todos os chats cujo toId sejam igual a guestId
// 	//que nao estejam invisiveis e que nao estejam blockeados
// 	dst := make([]Chat, 0, 1000)
// 	//
// 	var q = datastore.NewQuery(chatEntityName).
// 		Filter("to.guestid=", guestId).
// 		Filter("is_visible=", true).
// 		Filter("is_blocked=", false).
// 		Limit(1000)
// 	//
// 	if _, err := q.GetAll(ctx, &dst); err != nil {
// 		return nil, errors.New("Erro ao tentar recuperar uma lista de entidades. " +
// 			"Original(" + err.Error() + ").")
// 	}
// 	//
// 	return dst, nil
// }

// func ListChatsFrom(ctx appengine.Context, guestId string) ([]Chat, error) {
// 	//listar todos os chats cujo toId sejam igual a guestId
// 	//que nao estejam invisiveis e que nao estejam blockeados
// 	dst := make([]Chat, 0, 1000)
// 	//
// 	var q = datastore.NewQuery(chatEntityName).
// 		Filter("from.guestid=", guestId).
// 		Filter("is_visible=", true).
// 		Filter("is_blocked=", false).
// 		Limit(1000)
// 	//
// 	if _, err := q.GetAll(ctx, &dst); err != nil {
// 		return nil, errors.New("Erro ao tentar recuperar uma lista de entidades. " +
// 			"Original(" + err.Error() + ").")
// 	}
// 	//
// 	return dst, nil
// }

// func ChatCount(ctx appengine.Context) int {
// 	var q = datastore.NewQuery(chatEntityName)
// 	if count, err := q.Count(ctx); err == nil {
// 		return count
// 	}
// 	return 0
// }

func NewChatDataMgr(ctx context.Context) ChatDataMgr {
	return ChatDataMgr{ctx}
}
