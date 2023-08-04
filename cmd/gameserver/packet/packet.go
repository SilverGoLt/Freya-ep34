package packet

import (
	"sync"

	"github.com/ubis/Freya/cmd/gameserver/def"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
)

var g_ServerConfig = def.ServerConfig
var g_ServerSettings = def.ServerSettings
var g_PacketHandler = def.PacketHandler
var g_RPCHandler = def.RPCHandler
var g_NetworkManager = def.NetworkManager

type context struct {
	char  *character.Character
	mutex sync.RWMutex
}

// Registers network packets
func RegisterPackets() {
	log.Info("Registering packets...")

	var pk = g_PacketHandler
	pk.Register(GETMYCHARTR, "GetMyChartr", GetMyChartr)
	pk.Register(NEWMYCHARTR, "NewMyChartr", NewMyChartr)
	pk.Register(DELMYCHARTR, "DelMyChartr", DelMyChartr)
	pk.Register(CONNECT2SVR, "Connect2Svr", Connect2Svr)
	pk.Register(VERIFYLINKS, "VerifyLinks", VerifyLinks)
	pk.Register(INITIALIZED, "Initialized", Initialized)
	pk.Register(UNINITIALZE, "Uninitialze", Uninitialze)
	pk.Register(GETSVRTIME, "GetSvrTime", GetSvrTime)
	pk.Register(MESSAGEEVNT, "MessageEvnt", MessageEvnt)
	pk.Register(NEWUSERLIST, "NewUserList", nil)
	pk.Register(DELUSERLIST, "DelUserList", nil)
	pk.Register(NFY_MESSAGEEVNT, "NotifyMessageEvnt", nil)
	pk.Register(SYSTEMMESSG, "SystemMessg", nil)
	pk.Register(CHARGEINFO, "ChargeInfo", ChargeInfo)
	pk.Register(SERVERENV, "ServerEnv", ServerEnv)
	pk.Register(CHECK_USR_PDATA, "CheckUserPrivacyData", CheckUserPrivacyData)
	pk.Register(BACK_TO_CHAR_LOBBY, "BackToCharLobby", BackToCharLobby)
	pk.Register(SUBPW_SET, "SubPasswordSet", SubPasswordSet)
	pk.Register(SUBPW_CHECK_REQ, "SubPasswordCheckRequest", SubPasswordCheckRequest)
	pk.Register(SUBPW_CHECK, "SubPasswordCheck", SubPasswordCheck)
	pk.Register(SUBPW_FIND_REQ, "SubPasswordFindRequest", SubPasswordFindRequest)
	pk.Register(SUBPW_FIND, "SubPasswordFind", SubPasswordFind)
	pk.Register(SUBPW_DEL_REQ, "SubPasswordDelRequest", SubPasswordDelRequest)
	pk.Register(SUBPW_DEL, "SubPasswordDel", SubPasswordDel)
	pk.Register(SUBPW_CHG_QA_REQ,
		"SubPasswordChangeQARequest", SubPasswordChangeQARequest)
	pk.Register(SUBPW_CHG_QA, "SubPasswordChangeQA", SubPasswordChangeQA)
	pk.Register(SET_CHAR_SLOT_ORDER, "SetCharacterSlotOrder", SetCharacterSlotOrder)
	pk.Register(CHANNEL_LIST, "ChannelList", ChannelList)
	pk.Register(CHANNEL_CHANGE, "ChannelChange", ChannelChange)
	pk.Register(CHAR_DEL_CHK_SUBPW,
		"CharacterDeleteCheckSubPassword", CharacterDeleteCheckSubPassword)
}
