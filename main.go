package main

import (
	"io"
	"os"
	"fmt"
	"bufio"
	"io/ioutil"
	"archive/zip"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"crypto/sha256"
	"github.com/voxelbrain/goptions"
)

var version = "0.2.0"

// manifest data structure
type Manifest struct {
	Parts struct {
		Boot		Data	`json:"boot"`
		Fs		Data	`json:"fs"`
		Fw		Data	`json:"fw"`
		SysParam	Data	`json:"sys_params"`
		BootConfig	Data	`json:"boot_cfg"`
		Rf		Data	`json:"rf_cal_data"`
	}				`json:"parts"`
}
type Data struct {
	Addr	int	`json:"addr"`
	Size	int	`json:"size"`
	Src	string	`json:"src"`
	Hash	string	`json:"cs_sha256"`
	Fill	int	`json:"fill"`
}

// bash input options
type Options struct {
	Zipfile		string		`goptions:"-i, --input, description='Path to the firmware zip file'"`
	Output		string		`goptions:"-o, --output, description='Name of the output file'"`
	Size		int		`goptions:"-s, --size, description='Output file size in KB'"`
	Force		bool		`goptions:"-f, --force, description='Force writing to an output file that is too small'"`
	Help		goptions.Help	`goptions:"-h, --help, description='Show this help'"`
	Version		bool		`goptions:"-v, --version, description='Prints current version and exits. All other options will be ignored'"`
}

// read an input bin-file and write it to the out-buffer
func ReadFile(path string, data Data, buff *[]byte, options Options) {
	f, err  := os.Open(filepath.Join(path, data.Src))
	if err  != nil { fmt.Println("Failed to open file: ", err); os.Exit(1) }
	defer f.Close()
	b, err  := ioutil.ReadAll(f)
	if err  != nil { fmt.Println("Failed to read file: ", err); os.Exit(1) }

	// check sha256 sum
	hash    := sha256.Sum256(b)
	if hex.EncodeToString(hash[:]) != data.Hash {
		fmt.Println("SHA256 sum not matching for ", data.Src)
	}

	// cp to buffer
	if data.Addr + data.Size > len(*buff) {
		if options.Force {
			fmt.Println("Warning: Skipping file ", data.Src)
			return
		} else {
			fmt.Println("Error: Size of output file too small"); os.Exit(1)
		}
	}
	copy( (*buff)[data.Addr:data.Addr+data.Size], b)
}

// fill section of out-buffer with bytes
func ByteFill(data Data, buff *[]byte, options Options) {
	fill    := byte(data.Fill)

	if data.Addr + data.Size > len(*buff) {
		if options.Force {
			fmt.Println("Warning: Skipping fill at ", data.Addr)
			return
		} else {
			fmt.Println("Error: Size of output file too small"); os.Exit(1)
		}
	}

	for i   := data.Addr; i < data.Addr + data.Size; i++ {
		(*buff)[i] = fill
	}
}

// write output file from buffer
func WriteFile(buff *[]byte, options Options) {
	f, err  := os.Create(options.Output)
	if err  != nil { fmt.Println("Failed to create file: ", err); os.Exit(1) }
	defer f.Close()
	w       := bufio.NewWriter(f)
	_, err   = w.Write(*buff)
	if err  != nil { fmt.Println("Failed to write file: ", err); os.Exit(1) }
	err      = w.Flush()
	if err  != nil { fmt.Println("Failed to write file: ", err); os.Exit(1) }
}

// decompress a zip archive (src) to a flat folder structure (dest)
// (only suitable for achives with few files due to in-loop defer)
func Unzip(src, dest string) (err error) {
	r, err := zip.OpenReader(src)
	if err != nil { return }
	defer r.Close()

	for _, f := range r.File {
		// ignore directories
		if f.FileInfo().IsDir() { continue }

		rc, err  := f.Open()
		if err   != nil { return err }
    	        defer rc.Close()
    	        path     := filepath.Join(dest, filepath.Base(f.Name))
		out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err   != nil { return err }
		defer out.Close()
		_, err    = io.Copy(out, rc)
		if err   != nil { return err }
	}

	return nil
}

func main() {
	options := Options{
		Zipfile:  "./build/fw.zip",
		Output:   "output.bin",
		Size:     4096,
	}
	goptions.ParseAndFail(&options)

	if options.Version != false { fmt.Println(version); os.Exit(0) }
	if options.Size    <= 0     { fmt.Println("Error: Size of output file too small"); os.Exit(1) }

	// make temp dir
	tmp, err := ioutil.TempDir("", "mgos-combine-")
	if err   != nil { fmt.Println("Failed to create temporary directory: ", err); os.Exit(1) }
	defer os.RemoveAll(tmp)

	//unzip into temp dir
	err       = Unzip(options.Zipfile, tmp)
	if err   != nil { fmt.Println("Failed to extract firmware archive: ", err); os.Exit(1) }

	// read manifest
	manifest := Manifest{}
	buff     := make([]byte, options.Size * 1024)

	f, err   := os.Open(filepath.Join(tmp, "manifest.json"))
	if err   != nil { fmt.Println("Failed to open manifest: ", err); os.Exit(1) }
	defer f.Close()
	b, err   := ioutil.ReadAll(f)
	if err   != nil { fmt.Println("Failed to read manifest: ", err); os.Exit(1) }
	err       = json.Unmarshal(b, &manifest)
	if err   != nil { fmt.Println("Failed to read manifest: ", err); os.Exit(1) }

	// combine binaries to single file
	ReadFile(tmp, manifest.Parts.Boot, &buff, options)
	ReadFile(tmp, manifest.Parts.Fs, &buff, options)
	ReadFile(tmp, manifest.Parts.Fw, &buff, options)
	ReadFile(tmp, manifest.Parts.SysParam, &buff, options)
	ByteFill(manifest.Parts.BootConfig, &buff, options)
	ByteFill(manifest.Parts.Rf, &buff, options)
	WriteFile(&buff, options)
}
