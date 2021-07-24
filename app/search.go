package app

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// 이미지 데이터 자료구조
type ImageData struct {
	Path       string  `json:"path"`
	FrameIn    int     `json:"framein"`
	FrameOut   int     `json:"frameout"`
	FrameRange int     `json:"framerange"`
	Seqs       []int   `json:"seqs"`
	Pad        int     `json:"pad"`
	Ext        string  `json:"ext"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
	Fps        float64 `json:"fps"`
	Codec      string  `json:"codec"`
}

// Search는 디렉토리 경로를 받아 경로 안에 image 또는 mov 데이터를 검색하여,
// 이미지 데이터 자료구조 리스트와 에러를 반환한다.
func Search(srcDir string) ([]ImageData, error) {
	_, err := os.Stat(srcDir)
	if err != nil {
		return nil, err
	}
	dataMap := make(map[string]ImageData)
	data := make([]ImageData, 0) // data list
	err = filepath.WalkDir(srcDir,
		func(path string, info fs.DirEntry, err error) error {
			switch {
			case err != nil:
				return nil
			case info.IsDir():
				return nil
			case strings.HasPrefix(info.Name(), "."):
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case ".mov", ".mp4":
				item := ImageData{
					Path: path,
					Ext:  ext,
				}
				dataMap[path] = item
			case ".jpg", ".jpeg", ".png", ".dpx", ".exr":
				pathname, seq, pad := seqInfo(path)
				if pathname == path { //single frame images(hasn't seqs number)
					item := ImageData{
						Path: path,
						Ext:  ext,
					}
					dataMap[path] = item
				} else if _, has := dataMap[pathname]; has { //path already exists
					item := dataMap[pathname]
					if seq < item.FrameIn { //find first frame
						item.FrameIn = seq
					} else if seq > item.FrameOut { //find last frame
						item.FrameOut = seq
					}
					item.Seqs = append(item.Seqs, seq) //make seq array
					dataMap[pathname] = item
				} else {
					width, height, fps, codec, err := ImageInfo(path)
					if err != nil {
						fmt.Println(err)
					}
					item := ImageData{
						Path:     pathname,
						FrameIn:  seq,
						FrameOut: seq,
						Seqs:     []int{seq},
						Pad:      pad,
						Ext:      ext,
						Width:    width,
						Height:   height,
						Fps:      fps,
						Codec:    codec,
					}
					dataMap[pathname] = item
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	//연속된 시퀀스를 가지지 않은 데이터를 찾아 분리하기 위해 다시 for문을 사용.
	//같은 파일이지만 연속되지 않은 시퀀스면 둘로 나눈다.
	//ss_0010.%04d.jpg 1001-1003, ss_0010.%04d.jpg 1006-1010
	for _, val := range dataMap {
		seqs := val.Seqs
		if seqs == nil {
			data = append(data, val)
		} else {
			for i := val.FrameIn; i < (val.FrameIn + len(seqs)); i++ {
				if !contains(seqs, i) {
					continue
				} else if !contains(seqs, i-1) {
					val.FrameIn = i
				} else if contains(seqs, i+1) {
					continue
				} else {
					val.FrameOut = i
					val.FrameRange = val.FrameOut - val.FrameIn + 1
					data = append(data, val)
				}
			}
		}
	}
	return data, nil
}

// seqInfo는 이미지 경로를 받아 시퀀스 문자열을 찾아 '%04d'의 형태로 변환하고
// 변환한 경로, 시퀀스 숫자, 시퀀스 자리수를 반환한다.
func seqInfo(path string) (output string, seq int, pad int) {
	re, _ := regexp.Compile("(.+[_.])([0-9]+)(.[a-zA-Z]+$)")
	results := re.FindStringSubmatch(path)
	if results == nil {
		return path, seq, pad
	}
	seqStr := results[2]
	seq, _ = strconv.Atoi(seqStr)
	output = results[1] + "%0" + strconv.Itoa(len(seqStr)) + "d" + results[3]
	pad = len(seqStr)
	return output, seq, pad
}

// contains는 array와 int를 받아, array에 해당 int가 있는지 찾아 bool로 반환한다.
func contains(nums []int, num int) bool {
	for _, v := range nums {
		if v == num {
			return true
		}
	}
	return false
}
