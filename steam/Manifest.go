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
		fmt.Println(len);
		fmt.Println(buf);
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
	d.Files = make([]File, len(d.GetMappings()))
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
