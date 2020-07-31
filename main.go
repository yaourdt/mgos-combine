package main

import (
	"os"
	"fmt"
	"bufio"
	"io/ioutil"
	"encoding/hex"
	"encoding/json"
	"crypto/sha256"
	"github.com/voxelbrain/goptions"
)

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
	Manifest	string		`goptions:"-m, --manifest, description='Path to the JSON manifest file'"`
	Output		string		`goptions:"-o, --output, description='Name of the output file'"`
	Size		int		`goptions:"-s, --size, description='Output file size in KB'"`
	Force		bool		`goptions:"-f, --force, description='Force writing to an output file that is too small'"`
	Help		goptions.Help	`goptions:"-h, --help, description='Show this help'"`
	Version		bool		`goptions:"-v, --version, description='Prints current version and exits. All other options will be ignored'"`
}

// read an input bin-file and write it to the out-buffer
func ReadFile(data Data, buff *[]byte, options Options) {
	f, err  := os.Open(data.Src)
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

func main() {
	options := Options{
		Manifest: "manifest.json",
		Output:   "output.bin",
		Size:     4096,
	}
	goptions.ParseAndFail(&options)

	if options.Version != false { fmt.Println("0.1.0"); os.Exit(0) }
	if options.Size    <= 0     { fmt.Println("Error: Size of output file too small"); os.Exit(1) }

	manifest := Manifest{}
	buff     := make([]byte, options.Size * 1024)

	f, err   := os.Open(options.Manifest)
	if err   != nil { fmt.Println("Failed to open file: ", err); os.Exit(1) }
	defer f.Close()
	b, err   := ioutil.ReadAll(f)
	if err   != nil { fmt.Println("Failed to read file: ", err); os.Exit(1) }
	err       = json.Unmarshal(b, &manifest)
	if err   != nil { fmt.Println("Failed to read file: ", err); os.Exit(1) }

	ReadFile(manifest.Parts.Boot,       &buff, options)
	ReadFile(manifest.Parts.Fs,         &buff, options)
	ReadFile(manifest.Parts.Fw,         &buff, options)
	ReadFile(manifest.Parts.SysParam,   &buff, options)
	ByteFill(manifest.Parts.BootConfig, &buff, options)
	ByteFill(manifest.Parts.Rf,         &buff, options)
	WriteFile(&buff, options)
}
