package steam

import (
	"github.com/Jleagle/valve-data-format-go/vdf"
	"strings"
)

type OS string
const (
	Windows OS = "windows"
	Linux      = "linux"
	Macos      = "macos"
)
func GetOS(s string) OS {
	return OS(s)
}

func fillOSList(s string) (oslist map[OS]struct{}) {
	oslist = make(map[OS]struct{})
	for _, s := range strings.Split(s, ",") {
		oslist[OS(s)] = struct{}{}
	}
	return
}

type AppManifest struct {
	kv vdf.KeyValue
}

func (am AppManifest) ID() string {
	if id, found := am.kv.GetChild("appid"); found == true {
		return id.Value
	}
	return ""
}
func (am AppManifest) GetDepots() []Depot {
	depots, found := am.kv.GetChild("depots")
	if found == false {
		return nil
	}

	depotKeys := depots.GetChildrenAsMap()
	ret := make([]Depot, 0, len(depotKeys))

	for k := range depots.GetChildrenAsMap() {
		if depot, found := depots.GetChild(k); found == true {
			ret = append(ret, Depot{depot})
		} else {
		// should never happen
		}
	}
	return ret
}

func (am AppManifest) GetOSList() map[OS]struct{} {
	common, _ := am.kv.GetChild("common")
	oslist_kv, _ := common.GetChild("oslist")
	return fillOSList(oslist_kv.Value)
}
