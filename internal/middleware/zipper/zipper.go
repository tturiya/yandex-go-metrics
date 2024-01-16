package zipper

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/tturiya/iter5/internal/util"
)

type GZipResponseWritterWrapper struct {
	http.ResponseWriter
}

func (w *GZipResponseWritterWrapper) Write(b []byte) (int, error) {
	var buff bytes.Buffer
	gzw := gzip.NewWriter(&buff)
	n, err := gzw.Write(b)
	if err != nil {
		return n, err
	}
	written, err := w.ResponseWriter.Write(buff.Bytes())
	return written, err
}

func GZipper(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var (
			enc               = r.Header["Content-Encoding"]
			clientEnc         = r.Header["Accept-Encoding"]
			targEnc           = "gzip"
			contentIsEncoded  = util.SliceContainsE[string](enc, targEnc)
			clientAcceptsGzip = false
		)
		if len(clientEnc) > 0 {
			clientAcceptsGzip = strings.Contains(clientEnc[0], targEnc)
		}
		if contentIsEncoded {
			fz, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Fatalln("Failed decompressing")
			}
			defer fz.Close()

			b, err := io.ReadAll(fz)
			if err != nil {
				log.Fatalln("Failed decompressing")
			}
			r.Body = io.NopCloser(bytes.NewReader(b))

		}

		if clientAcceptsGzip {
			fmt.Println("got here")
			fmt.Println("client accepts encodings: ", clientEnc[0])
			gzrw := &GZipResponseWritterWrapper{w}
			next.ServeHTTP(gzrw, r)
			return
		} else {
			fmt.Println("client accepts encodings: ", clientEnc[0])
			fmt.Println("client accepts encodings: ", clientAcceptsGzip)
			fmt.Println(clientEnc[0] == targEnc)
			next.ServeHTTP(w, r)
			return
		}
	}

	return http.HandlerFunc(fn)
}
