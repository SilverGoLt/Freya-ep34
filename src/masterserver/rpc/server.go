package rpc

import (
    "share/models/server"
    "share/rpc"
)

// ServerRegister RPC Call
func ServerRegister(c *rpc.Client, r *server.RegisterReq, s *server.RegisterResp) error {
    var response = server.RegisterResp{}

    switch r.Type {
    case server.LOGIN_SERVER:
        response.Registered = true
        g_ServerManager.NewServer(server.Server{r, c})
        log.Infof("Server type: LoginServer (src: %s)", c.GetEndPnt())
    case server.GAME_SERVER:
        response.Registered = true
        g_ServerManager.NewServer(server.Server{r, c})
        log.Infof("Server type: GameServer (type: %d, server: %d, channel: %d, src: %s)",
            r.ServerType, r.ServerId, r.ChannelId, c.GetEndPnt())
    default:
        log.Errorf("Unknown server type (src %s)", c.GetEndPnt())
    }

    *s = response
    return nil
}

// ServerList RPC Call
func ServerList(c *rpc.Client, r *server.ListReq, s *server.ListResp) error {
    *s = server.ListResp{g_ServerManager.GetGameServerList()}
    return nil
}