package packet

import (
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"

	"github.com/ubis/Freya/share/log"
)

// Initialized Packet
func Initialized(session *network.Session, reader *network.Reader) {
	charId := reader.ReadInt32()

	if !session.Data.Verified || !session.Data.LoggedIn || session.DataEx == nil {
		log.Errorf("User is not verified (char: %d)", charId)
		return
	}

	ctx, ok := session.DataEx.(*context)
	if !ok {
		log.Error("Unable to retrieve user context (id: %d)",
			session.Data.AccountId)
		return
	}

	// verify char id
	if (charId >> 3) != session.Data.AccountId {
		log.Errorf("User is using invalid character id (id: %d, char: %d)",
			session.Data.AccountId, charId)
		return
	}

	c := character.Character{}

	if len(session.Data.CharacterList) == 0 {
		// fetch characters
		reqList := character.ListReq{
			Account: session.Data.AccountId,
			Server:  byte(g_ServerSettings.ServerId),
		}
		resList := character.ListRes{}
		g_RPCHandler.Call(rpc.LoadCharacters, &reqList, &resList)

		session.Data.CharacterList = resList.List
	}

	// fetch character
	for _, data := range session.Data.CharacterList {
		if data.Id == charId {
			c = data
			break
		}
	}

	// check if character exists
	if c.Id != charId {
		log.Errorf("User is using invalid character id (id: %d, char: %d)",
			session.Data.AccountId, charId)
		return
	}

	// load additional character data
	req := character.DataReq{
		Server: byte(g_ServerSettings.ServerId),
		Id:     c.Id,
	}
	res := character.DataRes{}
	g_RPCHandler.Call(rpc.LoadCharacterData, req, &res)

	// serialize data
	eq, eqlen := c.Equipment.Serialize()
	inv, invlen := res.Inventory.Serialize()
	sk, sklen := res.Skills.Serialize()
	sl, sllen := res.Links.Serialize()

	pkt := network.NewWriter(INITIALIZED)
	pkt.WriteBytes(make([]byte, 57))
	pkt.WriteByte(0x00)
	pkt.WriteByte(0x14)
	pkt.WriteByte(g_ServerSettings.ChannelId)
	pkt.WriteBytes(make([]byte, 23))
	pkt.WriteByte(0xFF)
	pkt.WriteUint16(g_ServerConfig.MaxUsers)
	pkt.WriteUint32(0x8501A8C0)
	pkt.WriteUint16(0x985A)
	pkt.WriteInt32(0x01)
	pkt.WriteInt32(0x0100001F)

	pkt.WriteInt32(c.World)
	pkt.WriteInt32(0x00)
	pkt.WriteUint16(c.X)
	pkt.WriteUint16(c.Y)
	pkt.WriteUint64(c.Exp)
	pkt.WriteUint64(c.Alz)
	pkt.WriteUint64(c.WarExp)
	pkt.WriteUint32(c.Level)
	pkt.WriteInt32(0x00)

	pkt.WriteUint32(c.STR)
	pkt.WriteUint32(c.DEX)
	pkt.WriteUint32(c.INT)
	pkt.WriteUint32(c.PNT)
	pkt.WriteByte(c.SwordRank)
	pkt.WriteByte(c.MagicRank)
	pkt.WriteUint16(0x00) // padding for skillrank
	pkt.WriteUint32(0x00)
	pkt.WriteUint16(c.MaxHP)
	pkt.WriteUint16(c.CurrentHP)
	pkt.WriteUint16(c.MaxMP)
	pkt.WriteUint16(c.CurrentMP)
	pkt.WriteUint16(c.MaxSP)
	pkt.WriteUint16(c.CurrentSP)
	pkt.WriteUint16(0x00) //stats.DungeonPoints)
	pkt.WriteUint16(0x00)
	pkt.WriteInt32(0x2A30)
	pkt.WriteInt32(0x01)
	pkt.WriteUint16(0x00) //stats.SwordExp)
	pkt.WriteUint16(0x00) //stats.SwordPoint)
	pkt.WriteUint16(0x00) //stats.MagicExp)
	pkt.WriteUint16(0x00) //stats.MagicPoint)
	pkt.WriteUint16(0x00) //stats.SwordExpPoint)
	pkt.WriteUint16(0x00) //stats.MagicExpPoint)
	pkt.WriteInt32(0x00)
	pkt.WriteInt32(0x00)
	pkt.WriteInt32(0x00)  // honour pnt
	pkt.WriteUint64(0x00) // death penalty exp
	pkt.WriteUint64(0x00) // death hp
	pkt.WriteUint64(0x00) // death mp
	pkt.WriteUint16(0x00) // pk penalty // pk pna

	pkt.WriteUint32(0x8501A8C0) // chat ip
	pkt.WriteUint16(0x9858)     // chat port

	pkt.WriteUint32(0x8501A8C0) // ah ip
	pkt.WriteUint16(0x9859)     // ah port

	pkt.WriteByte(c.Nation)
	pkt.WriteInt32(0x00)
	pkt.WriteInt32(0x07) // warp code
	pkt.WriteInt32(0x07) // map code
	pkt.WriteUint32(c.Style.Get())
	pkt.WriteBytes(make([]byte, 39))

	pkt.WriteUint16(eqlen)
	pkt.WriteUint16(invlen)
	pkt.WriteUint16(sklen)
	pkt.WriteUint16(sllen)

	pkt.WriteBytes(make([]byte, 6))
	pkt.WriteUint16(0x00) // ap
	pkt.WriteUint32(0x00) // ap exp
	pkt.WriteInt16(0x00)
	pkt.WriteByte(0x00)   // blessing bead count
	pkt.WriteByte(0x00)   // active quest count
	pkt.WriteUint16(0x00) // period item count
	pkt.WriteBytes(make([]byte, 1023))

	pkt.WriteBytes(make([]byte, 128)) // quest dungeon flags
	pkt.WriteBytes(make([]byte, 128)) // mission dungeon flags

	pkt.WriteByte(0x00)              // Craft Lv 0
	pkt.WriteByte(0x00)              // Craft Lv 1
	pkt.WriteByte(0x00)              // Craft Lv 2
	pkt.WriteByte(0x00)              // Craft Lv 3
	pkt.WriteByte(0x00)              // Craft Lv 4
	pkt.WriteUint16(0x00)            // Craft Exp 0
	pkt.WriteUint16(0x00)            // Craft Exp 1
	pkt.WriteUint16(0x00)            // Craft Exp 2
	pkt.WriteUint16(0x00)            // Craft Exp 3
	pkt.WriteUint16(0x00)            // Craft Exp 4
	pkt.WriteBytes(make([]byte, 16)) // Craft Flags
	pkt.WriteUint32(0x00)            // Craft Type

	pkt.WriteInt32(0x00) // Help Window Index
	pkt.WriteBytes(make([]byte, 163))

	pkt.WriteUint32(0x00) // TotalPoints
	pkt.WriteUint32(0x00) // GeneralPoints
	pkt.WriteUint32(0x00) // QuestPoints
	pkt.WriteUint32(0x00) // DungeonPoints
	pkt.WriteUint32(0x00) // ItemPoints
	pkt.WriteUint32(0x00) // PVPPoints
	pkt.WriteUint32(0x00) // MissionWarPoints
	pkt.WriteUint32(0x00) // HuntingPoints
	pkt.WriteUint32(0x00) // CraftingPoints
	pkt.WriteUint32(0x00) // CommunityPoints
	pkt.WriteUint32(0x00) // SharedAchievments
	pkt.WriteUint32(0x00) // SpecialPoints

	pkt.WriteUint32(0x00)
	pkt.WriteUint32(0x00) // QuestsCount
	pkt.WriteUint32(0x00) // QuestFlagsCount
	pkt.WriteUint32(0x00)

	pkt.WriteByte(len(c.Name) + 1)
	pkt.WriteString(c.Name)

	pkt.WriteBytes(eq)
	pkt.WriteBytes(inv)
	pkt.WriteBytes(sk)
	pkt.WriteBytes(sl)

	session.Send(pkt)

	ctx.mutex.Lock()
	ctx.char = &c
	ctx.mutex.Unlock()
}

// Uninitialze Packet
func Uninitialze(session *network.Session, reader *network.Reader) {
	_ = reader.ReadUint16() // index
	_ = reader.ReadByte()   // map id
	_ = reader.ReadByte()   // log out

	pkt := network.NewWriter(UNINITIALZE)
	pkt.WriteByte(0) // result

	// complete - 0x00
	// fail - 0x01
	// ignored - 0x02
	// busy - 0x03
	// anti online game - 0x30

	session.Send(pkt)
}
