package steam

import (
	"strconv"
	"github.com/Philipp15b/go-steam"
	"github.com/Philipp15b/go-steam/protocol"
	"github.com/Philipp15b/go-steam/protocol/steamlang"
	"github.com/Philipp15b/go-steam/protocol/protobuf"
	"github.com/Jleagle/valve-data-format-go/vdf"
)

type PICSClient struct {
	C *steam.Client
}

type PICSProductInfoEvent struct {
	Apps map[string]AppManifest
}

func (c PICSClient) HandlePacket(p *protocol.Packet) {
	switch (p.EMsg) {
	case steamlang.EMsg_ClientPICSProductInfoResponse:
		var ret PICSProductInfoEvent
		body := new(protobuf.CMsgClientPICSProductInfoResponse)
		p.ReadProtoMsg(body)
		ret.Apps = make(map[string]AppManifest, len(body.GetApps()))
		for _, app := range body.GetApps() {
			vdf, _ := vdf.ReadBytes(app.GetBuffer())
			ret.Apps[strconv.Itoa(int(app.GetAppid()))] = AppManifest{vdf}
		}
		c.C.Emit(&ret)
	}
}
