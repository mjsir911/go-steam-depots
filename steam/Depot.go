package steam

import (
	"net/url"
)

type Depot struct {
	data interface{}
}

func NewDepot(id string) (d Depot, err error) {
	return
}

func (d Depot) Name() string {
	return "1113281"
	// return d.data.ID
}
// 
// func (d Depot) Size() int64 {
// 	return d.data.Maxsize
// }

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
