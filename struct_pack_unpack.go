package struct_pack_unpack

import (
	"errors"
	"strconv"
)

type Format struct {
	pack_element_length []int
	total_length        int
}

var pack_type_length = map[string]int{
	"c": 1,
	"b": 1,
	"B": 1,
	"?": 1,
	"h": 2,
	"H": 2,
	"i": 4,
	"I": 4,
	"l": 4,
	"L": 4,
	"q": 8,
	"Q": 8,
	"f": 4,
	"d": 8,
	"s": 1,
	"p": 1,
}

func AnalyseFmt(fmt_string string) *Format {
	total_length := 0
	fmt_length := 0
	fmt_length_every_element := make([]int, 0)
	s_q_start := 0
	s_q_end := 0
	s_q_length := 0

	for i := 0; i < len(fmt_string); i++ {
		s := string(fmt_string[i])
		value, ok := pack_type_length[s]
		if ok {
			if s != "s" && s != "p" {
				fmt_length_every_element = append(fmt_length_every_element, value)
				total_length = total_length + value
				fmt_length = fmt_length + 1
			} else {
				s_q_length, _ = strconv.Atoi(string(fmt_string[s_q_start : s_q_end+1]))
				fmt_length_every_element = append(fmt_length_every_element, s_q_length)
				total_length = total_length + s_q_length
				fmt_length = fmt_length + 1
			}
		} else {
			_, o1 := pack_type_length[string(fmt_string[i-1])]
			if o1 {
				s_q_start = i
			} else {
				_, o2 := pack_type_length[string(fmt_string[i+1])]
				if o2 {
					s_q_end = i
				}
			}
		}
	}
	return &Format{
		pack_element_length: fmt_length_every_element,
		total_length:        total_length,
	}
}

func Pack(pack_format string, bytes [][]byte) ([]byte, error) {
	var fmtInfo *Format = AnalyseFmt(pack_format)
	pack_element_length := fmtInfo.pack_element_length
	concated_bytes := make([]byte, 0)
	for i, v := range bytes {
		if len(v) != pack_element_length[i] {
			return nil, errors.New("length not compact")
		}
		concated_bytes = append(concated_bytes, v...)
	}
	return concated_bytes, nil
}

func Unpack(pack_format string, byte_array []byte) ([][]byte, error){
    var fmtInfo *Format = AnalyseFmt(pack_format)
	pack_element_length := fmtInfo.pack_element_length
    total_length := fmtInfo.total_length
    if total_length != len(byte_array){
        return nil, errors.New("length not compact")
    }
    start := 0
    unpacked_bytes := make([][]byte, 0)
    for _, v := range pack_element_length {
		current_bytes := byte_array[start:start+v]
        start = start + v
		unpacked_bytes = append(unpacked_bytes, current_bytes)
	}
    return unpacked_bytes, nil
}
