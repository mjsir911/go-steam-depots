package steam

import (
	"github.com/Jleagle/valve-data-format-go/vdf"
)

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
