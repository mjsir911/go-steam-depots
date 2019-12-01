package steam

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"
	"io"
	"github.com/golang/protobuf/proto"
	"net/url"
	"net/http"
	"github.com/mjsir911/szip"
	"errors"
)

type Chunk struct {
	ContentManifestPayload_FileMapping_ChunkData
}
func (c Chunk) Read(b []byte) (n int, err error) {
	// download
	// decrypt
	// decompress
	return
}
func (c Chunk) Name() string {
	return hex.EncodeToString(c.GetSha())
}
func (c Chunk) Size() int64 {
	return int64(c.GetCbOriginal())
}


type File struct {
	ContentManifestPayload_FileMapping
	Chunks []Chunk
}
func (f File) Read(b []byte) (n int, err error) {
	n = 0
	for _, chunk := range f.Chunks {
		var m int
		m, err = chunk.Read(b[n:])
		n += m
		if err != nil {
			return
		}
	}
	return
}


func (f File) Name() string {
	return f.GetFilename()
}

func (f File) Size() int64 {
	return int64(f.GetSize())
}

/* A manifest is sort of like a directory */
type Manifest struct {
	ContentManifestPayload
	ContentManifestMetadata
	ContentManifestSignature
	Files []File
}
func (i Manifest) Name() string {
	return strconv.FormatUint(i.GetGidManifest(), 10)
}
func (i Manifest) Size() int64 {
	return int64(i.GetCbDiskOriginal())
}
func (i Manifest) ModTime() time.Time {
	return time.Unix(int64(i.GetCreationTime()), 0)
}

func NewManifest(r io.Reader) (d Manifest, err error) {
	for {
		var typ, len uint32
		if err = binary.Read(r, binary.LittleEndian, &typ); err != nil {
			break
		}
		if err = binary.Read(r, binary.LittleEndian, &len); err != nil {
			break
		}
		buf := make([]byte, len)
		if _, err = io.ReadFull(r, buf); err != nil {
			break
		}
		switch typ {
		case 0x71F617D0:
			if err = proto.Unmarshal(buf, &d.ContentManifestPayload); err != nil {
				return
			}
		case 0x1F4812BE:
			if err = proto.Unmarshal(buf, &d.ContentManifestMetadata); err != nil {
				return
			}
		case 0x1B81B817:
			if err = proto.Unmarshal(buf, &d.ContentManifestSignature); err != nil {
				return
			}
		default:
			err = errors.New("unknown type: " + string(typ))
			return
		}
	}
	if err != nil && err != io.EOF {
		return
	}
	d.Files = make([]File, len(d.GetMappings()))
	for i, file := range d.GetMappings() {
		d.Files[i].ContentManifestPayload_FileMapping = *file

		// if len(file.Chunks) > 1 {
			// fmt.Println(file)
			// for _, chunk := range file.Chunks {
				//fmt.Println(chunk)
		// 		fmt.Println(hex.EncodeToString(chunk.GetSha()))
		// 	}
		 // }
	}
	//fmt.Println(signature)
	return
}

func DownloadManifest(url url.URL) (d Manifest, err error) {
	r, err := http.Get(url.String())
	if err != nil {
		return
	}
	unzipper, err := szip.NewReader(r.Body)
	if err != nil {
		return
	}
	_, err = unzipper.Next()
	if err != nil {
		return
	}
	d, err = NewManifest(unzipper)
	if err != nil {
		return
	}

	return
}

func ManifestUrl(depot url.URL, id string) url.URL {
	// http://cache20-iad1.steamcontent.com/depot/1113281/manifest/8153978177929726511/5
	depot.Path += "/manifest/" + id + "/5"
	return depot
}
