package load

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/jakubDoka/sterr"
	"golang.org/x/image/font"
)

// errors
var (
	ErrNotOnDisc = sterr.New("%s file not found on disc")
)

// AppData returns AppData access, the root is usually the name of your app
func AppData(root string) (Util, error) {
	ad, err := AppDataDir()
	return Util{
		Loader: OSFS{},
		Root:   path.Join(ad, root),
	}, err
}

// AppDataDir returns directory where the app data is located
func AppDataDir() (string, error) {
	var home string
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}

	if home == "" {
		return "", sterr.New("unable to optain AppData directory path")
	}

	home += "\\AppData"

	return home, nil
}

// Util brings some utility for Loader, if you specify root
// it will be prepended to path you input
type Util struct {
	Loader
	Root string
}

// OS is default os loader
var OS = Util{Loader: OSFS{}}

// Json unmarshal-s json from given path to dest
func (l Util) Json(p string, dest interface{}) error {
	bts, err := l.ReadFile(l.path(p))
	if err != nil {
		return err
	}

	return json.Unmarshal(bts, dest)
}

// LoadTTF loads TTF file into font.Face
func (l Util) LoadTTF(p string, size float64) (font.Face, error) {
	bytes, err := l.ReadFile(l.path(p))
	if err != nil {
		return nil, ErrNotOnDisc.Args("ttf").Wrap(err)
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

// LoadImage loads image from disk
func (l Util) LoadImage(p string) (*image.NRGBA, error) {
	imgFile, err := l.Open(l.path(p))
	if err != nil {
		return nil, ErrNotOnDisc.Args("image").Wrap(err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %q: %v", p, err)
	}

	var res *image.NRGBA

	switch v := img.(type) {
	case *image.NRGBA:
		res = v
	default:
		res = image.NewNRGBA(v.Bounds())
		draw.Draw(res, res.Rect, img, img.Bounds().Min, 0)
	}

	return res, nil
}

// Walk calls fs.WalkDir on Loader
func (l Util) Walk(root string, fn func(path string, d fs.DirEntry, err error) error) {
	fs.WalkDir(l, l.path(root), fn)
}

// List lists all paths in given directory. If rec == true, function will go recursively. If ext == ""
// they any extencion is valid, ext does not include dot ("go" not ".go")
func (l Util) List(root string, buff []string, rec bool, ext string) (r []string, err error) {
	r = buff
	entries, err := l.ReadDir(l.path(root))
	if err != nil {
		return
	}

	for _, e := range entries {
		p := path.Join(root, e.Name())
		if e.IsDir() {
			if rec {
				r, err = l.List(p, r, true, ext)
				if err != nil {
					return
				}
			}
		} else {
			idx := strings.LastIndex(e.Name(), ".")
			if idx == -1 {
				if ext == "" {
					r = append(r, p)
				}
			} else if ext == "" || e.Name()[idx+1:] == ext {
				r = append(r, p)
			}

		}
	}

	return
}

func (l Util) path(base string) string {
	return path.Join(l.Root, base)
}

// Loader can be a OSFS or embed.FS, so writing generic
// asset loading where it does not matter where you load from
type Loader interface {
	fs.ReadDirFS
	ReadFile(name string) ([]byte, error)
}

// OSFS uses os as source of data
type OSFS struct{}

// Open implements fs.FS interface
func (o OSFS) Open(path string) (fs.File, error) {
	return os.Open(path)
}

// ReadFile implements AssetLoader interface
func (o OSFS) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// ReadDir implements fs.ReadDirFS interface
func (o OSFS) ReadDir(path string) ([]fs.DirEntry, error) {
	return os.ReadDir(path)
}
