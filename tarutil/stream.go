package tarutil

import (
	"archive/tar"
	"io"
	"os"

	"shanhu.io/misc/errcode"
)

// streamFile is a file (or a zip archive) to stream into a tar stream.
type streamFile struct {
	name    string // Name to write into the tar stream.
	file    string // File to read from file system.
	zip     bool   // If to read the file as a zip file.
	content []byte // Raw content; used only when File is empty string.
	mode    int64  // File mode; used only when File is empty string.
}

func (f *streamFile) writeTo(tw *tar.Writer) error {
	if f.zip {
		return TarZipFile(tw, f.file)
	}

	if f.file != "" {
		file, err := os.Open(f.file)
		if err != nil {
			return err
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return err
		}

		mode := f.mode
		if mode == 0 {
			mode = int64(stat.Mode()) & 0777
		}

		if err := tw.WriteHeader(&tar.Header{
			Name: f.name,
			Size: stat.Size(),
			Mode: mode,
		}); err != nil {
			return err
		}

		_, err = io.Copy(tw, file)
		return err
	}

	if err := tw.WriteHeader(&tar.Header{
		Name: f.name,
		Size: int64(len(f.content)),
		Mode: f.mode,
	}); err != nil {
		return err
	}
	_, err := tw.Write(f.content)
	return err
}

// Stream is a tar stream of files (or zip files). Files are transfered in
// the order of adding.
type Stream struct {
	files []*streamFile
}

// NewStream create a new tar stream.
func NewStream() *Stream { return &Stream{} }

// NewDockerTarStream creates a stream with a docker file.
func NewDockerTarStream(dockerfile string) *Stream {
	ts := NewStream()
	ts.AddDockerfile(dockerfile)
	return ts
}

// AddDockerfile adds a DockerFile of content with mode 0600.
func (s *Stream) AddDockerfile(content string) {
	s.AddString("Dockerfile", 0600, content)
}

// AddString adds a file of name and mode into the stream,
// which content is str.
func (s *Stream) AddString(name string, mode int64, str string) {
	s.AddBytes(name, mode, []byte(str))
}

// AddBytes adds a file of name and mode into the stream,
// which content is bs.
func (s *Stream) AddBytes(name string, mode int64, bs []byte) {
	s.files = append(s.files, &streamFile{
		name:    name,
		mode:    mode,
		content: bs,
	})
}

// AddFile adds a file of name and mode into the stream,
// which content is read from file f.
func (s *Stream) addFile(name string, mode int64, f string) {
	s.files = append(s.files, &streamFile{
		name: name,
		mode: mode,
		file: f,
	})
}

// AddZipFile adds a zip file into the stream.
func (s *Stream) AddZipFile(f string) {
	s.files = append(s.files, &streamFile{
		file: f,
		zip:  true,
	})
}

type countingWriter struct {
	w io.Writer
	n int64
}

func (w *countingWriter) Write(bs []byte) (int, error) {
	n, err := w.Write(bs)
	w.n += int64(n)
	return n, err
}

// WriteTo writes the entire stream out to w.
func (s *Stream) WriteTo(w io.Writer) (int64, error) {
	cw := &countingWriter{w: w}
	tw := tar.NewWriter(cw)
	for _, f := range s.files {
		if err := f.writeTo(tw); err != nil {
			return cw.n, errcode.Annotatef(err, "write %q", f.name)
		}
	}
	err := tw.Close() // Close() might flush stuff and update cw.n
	return cw.n, err
}
