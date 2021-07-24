package app

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

// FileInfo는 파일 경로를 받아 ffprobe를 사용하여 데이터를 분석하고,
// width, height, fps, codec, err를 반환한다.
// ffprobe data sample
// {
//     "streams": [
//         {
//             "index": 0,
//             "codec_name": "h264",
//             "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
//             "profile": "High",
//             "codec_type": "video",
//             "codec_time_base": "1/5994",
//             "codec_tag_string": "avc1",
//             "codec_tag": "0x31637661",
//             "width": 2048,
//             "height": 1152,
//             "coded_width": 2048,
//             "coded_height": 1152,
//             "has_b_frames": 2,
//             "pix_fmt": "yuv420p",
//             "level": 62,
//             "chroma_location": "left",
//             "refs": 1,
//             "is_avc": "true",
//             "nal_length_size": "4",
//             "r_frame_rate": "2997/1",
//             "avg_frame_rate": "2997/1",
//             "time_base": "1/11988",
//             "start_pts": 0,
//             "start_time": "0.000000",
//             "duration_ts": 40,
//             "duration": "0.003337",
//             "bit_rate": "56844698",
//             "bits_per_raw_sample": "8",
//             "nb_frames": "10",
//             "disposition": {
//                 "default": 1,
//                 "dub": 0,
//                 "original": 0,
//                 "comment": 0,
//                 "lyrics": 0,
//                 "karaoke": 0,
//                 "forced": 0,
//                 "hearing_impaired": 0,
//                 "visual_impaired": 0,
//                 "clean_effects": 0,
//                 "attached_pic": 0,
//                 "timed_thumbnails": 0
//             },
//             "tags": {
//                 "handler_name": "VideoHandler",
//                 "encoder": "Lavc58.54.100 libx264"
//             }
//         }
//     ]
// }
func ImageInfo(path string) (width, height, fps float64, codec string, err error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_streams", path)
	stdout, err := cmd.Output()
	if err != nil {
		return width, height, fps, codec, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(stdout), &data)

	if err != nil {
		return width, height, fps, codec, err
	}
	if value, ok := data["streams"]; ok {
		streams := value.([]interface{})
		item := streams[0].(map[string]interface{})
		for k, v := range item {
			switch k {
			case "width":
				width = v.(float64)
			case "height":
				height = v.(float64)
			case "r_frame_rate":
				slice := strings.Split(v.(string), "/")
				r1, _ := strconv.ParseFloat(slice[0], 64)
				r2, _ := strconv.ParseFloat(slice[1], 64)
				fps = r1 / r2
			case "codec_name":
				codec = v.(string)
			}
		}
	}
	return width, height, fps, codec, err
}
