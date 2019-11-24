package steam

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"
	"io"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type DepotChunk struct {
	ContentManifestPayload_FileMapping_ChunkData
}
func (c DepotChunk) Read(b []byte) (n int, err error) {
	// download
	// decrypt
	// decompress
	return
}
func (c DepotChunk) Name() string {
	return hex.EncodeToString(c.GetSha())
}
func (c DepotChunk) Size() int64 {
	return int64(c.GetCbOriginal())
}


type DepotFile struct {
	ContentManifestPayload_FileMapping
	Chunks []DepotChunk
}
func (f DepotFile) Read(b []byte) (n int, err error) {
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


func (f DepotFile) Name() string {
	return f.GetFilename()
}

func (f DepotFile) Size() int64 {
	return int64(f.GetSize())
}

/* A depot is sort of like a directory */
type Depot struct {
	ContentManifestPayload
	ContentManifestMetadata
	ContentManifestSignature
	Files []DepotFile
}
func (i Depot) Name() string {
	return strconv.FormatUint(i.GetGidManifest(), 10)
}
func (i Depot) Size() int64 {
	return int64(i.GetCbDiskOriginal())
}
func (i Depot) ModTime() time.Time {
	return time.Unix(int64(i.GetCreationTime()), 0)
}
func NewDepot(r io.Reader) (d Depot, err error) {
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
			fmt.Println("unknown type:", typ)
		}
	}
	if err != nil && err != io.EOF {
		return
	}
	d.Files = make([]DepotFile, len(d.GetMappings()))
	for i, file := range d.GetMappings() {
		d.Files[i].ContentManifestPayload_FileMapping = *file

		// if len(file.Chunks) > 1 {
			// fmt.Println(file)
			for _, chunk := range file.Chunks {
				fmt.Println(chunk)
		// 		fmt.Println(hex.EncodeToString(chunk.GetSha()))
		// 	}
		 }
	}
	//fmt.Println(signature)
	return
}
