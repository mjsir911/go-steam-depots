package steam

import (
	"net/url"
	"github.com/Jleagle/valve-data-format-go/vdf"
	"strconv"
)

type Depot struct {
	ID string
	MaxSize int
	Manifests map[string]string
}

// func NewDepot(id string) (d Depot, err error) {
// 	return
// }

// for now use keyvalues, eventually want to use something more
// encoding/json-like, with type assertions & all
func NewDepot(kv vdf.KeyValue) (d Depot, err error) {
	d.ID = kv.Key
	child, _ := kv.GetChild("MaxSize")
	d.MaxSize, _ = strconv.Atoi(child.Value)
	child, _ = kv.GetChild("manifests")
	d.Manifests = child.GetChildrenAsMap()
	return
}

func (d Depot) Name() string {
	return d.ID
}

func (d Depot) Size() int64 {
	return int64(d.MaxSize)
}

func (d Depot) URL() url.URL {
	// http://cache20-iad1.steamcontent.com/depot/1113281
	return url.URL{
		Scheme: "http",
		Host: "cache2-iad1.steamcontent.com", // todo: dynamic
		Path: "/depot/" + d.Name(),
	}
}

func (d Depot) GetManifest(id string) (Manifest, error) {
	return DownloadManifest(ManifestUrl(d.URL(), id))
}
