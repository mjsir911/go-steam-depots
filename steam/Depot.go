package steam

import (
	"net/url"
	"github.com/Jleagle/valve-data-format-go/vdf"
	"strconv"
)

// for now use keyvalues, eventually want to use something more
// encoding/json-like, with type assertions & all
type Depot struct {
	kv vdf.KeyValue
}

func (d Depot) Name() string {
	return d.kv.Key
}

func (d Depot) GetManifestMap() map[string]string {
	manifests_kv, _ := d.kv.GetChild("manifests")
	return manifests_kv.GetChildrenAsMap()
}

func (d Depot) Size() int64 {
	size, _ := d.kv.GetChild("MaxSize")
	size2, _ := strconv.Atoi(size.Value)
	return int64(size2)
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

func (d Depot) GetOSList() map[OS]struct{} {
	config, _ := d.kv.GetChild("config")
	oslist_kv, _ := config.GetChild("oslist")
	return fillOSList(oslist_kv.Value)
}
func (d Depot) GetManifestByLabel(label string) (Manifest, error) {
	return d.GetManifest(d.GetManifestMap()[label])
}
