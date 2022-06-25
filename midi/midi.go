package midi

import (
	"bytes"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/FloatTech/zbputils/file"
	"github.com/pkg/errors"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

var (
	noteMap = map[string]uint8{
		"C":  60,
		"Db": 61,
		"D":  62,
		"Eb": 63,
		"E":  64,
		"F":  65,
		"Gb": 66,
		"G":  67,
		"Ab": 68,
		"A":  69,
		"Bb": 70,
		"B":  71,
	}
)

// Txt2mid 文本转midi文件,timbre是音色,filePath是midi文件路径,input是文本输入
func Txt2mid(timbre int, filePath, input string) error {
	if file.IsExist(filePath) {
		return nil
	}
	var (
		clock smf.MetricTicks
		tr    smf.Track
	)

	tr.Add(0, smf.MetaMeter(4, 4))
	tr.Add(0, smf.MetaTempo(72))
	tr.Add(0, smf.MetaInstrument("Violin"))
	tr.Add(0, midi.ProgramChange(0, uint8(timbre)))

	k := strings.ReplaceAll(input, " ", "")

	var (
		base        uint8
		level       uint8
		delay       uint32
		sleepFlag   bool
		lengthBytes = make([]byte, 0)
	)

	for i := 0; i < len(k); {
		base = 0
		level = 0
		sleepFlag = false
		lengthBytes = lengthBytes[:0]
		for {
			switch {
			case k[i] == 'R':
				sleepFlag = true
				i++
			case k[i] >= 'A' && k[i] <= 'G':
				base = noteMap[k[i:i+1]] % 12
				i++
			case k[i] == 'b':
				base--
				i++
			case k[i] == '#':
				base++
				i++
			case k[i] >= '0' && k[i] <= '9':
				level = level*10 + k[i] - '0'
				i++
			case k[i] == '<':
				i++
				for i < len(k) && (k[i] == '-' || (k[i] >= '0' && k[i] <= '9')) {
					lengthBytes = append(lengthBytes, k[i])
					i++
				}
			default:
				return errors.Errorf("无法解析第%d个位置的%c字符", i, k[i])
			}
			if i >= len(k) || (k[i] >= 'A' && k[i] <= 'G') || k[i] == 'R' {
				break
			}
		}
		length, _ := strconv.Atoi(string(lengthBytes))
		if sleepFlag {
			if length >= 0 {
				delay = clock.Ticks4th() * (1 << length)
			} else {
				delay = clock.Ticks4th() / (1 << -length)
			}
			continue
		}
		if level == 0 {
			level = 5
		}
		tr.Add(delay, midi.NoteOn(0, O(base, level), 120))
		if length >= 0 {
			tr.Add(clock.Ticks4th()*(1<<length), midi.NoteOff(0, O(base, level)))
		} else {
			tr.Add(clock.Ticks4th()/(1<<-length), midi.NoteOff(0, O(base, level)))
		}
		delay = 0
	}
	tr.Close(0)

	s := smf.New()
	s.TimeFormat = clock
	err := s.Add(tr)
	if err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = s.WriteTo(f)
	f.Close()
	return err
}

// O 原始接口,base对应C5中的C,oct对应C5中的5,base=0,oct=5,返回60
func O(base uint8, oct uint8) uint8 {
	if oct > 10 {
		oct = 10
	}

	if oct == 0 {
		return base
	}

	res := base + 12*oct
	if res > 127 {
		res -= 12
	}

	return res
}

// Name 通过音高获得音高名,n=60,返回C
func Name(n uint8) string {
	for k, v := range noteMap {
		if v%12 == n%12 {
			return k
		}
	}
	return ""
}

// Pitch 通过音高名获得音高,note=C4,返回48
func Pitch(note string) uint8 {
	k := strings.ReplaceAll(note, " ", "")
	var (
		base  uint8
		level uint8
	)
	for i := 0; i < len(k); i++ {
		switch {
		case k[i] >= 'A' && k[i] <= 'G':
			base = noteMap[k[i:i+1]] % 12
		case k[i] == 'b':
			base--
		case k[i] == '#':
			base++
		case k[i] >= '0' && k[i] <= '9':
			level = level*10 + k[i] - '0'
		}
	}
	if level == 0 {
		level = 5
	}
	return O(base, level)
}

// Mid2txt midi转txt,trackNo是音轨序号
func Mid2txt(midBytes []byte, trackNo int) (midStr string) {
	var (
		absTicksStart float64
		absTicksEnd   float64
		startNote     byte
		endNote       byte
		defaultMetric = 960.0
	)
	_ = smf.ReadTracksFrom(bytes.NewReader(midBytes), trackNo).
		Do(
			func(te smf.TrackEvent) {
				if !te.Message.IsMeta() {
					b := te.Message.Bytes()
					if te.Message.Is(midi.NoteOnMsg) && b[2] > 0 {
						absTicksStart = float64(te.AbsTicks)
						startNote = b[1]
					}
					if te.Message.Is(midi.NoteOffMsg) || (te.Message.Is(midi.NoteOnMsg) && b[2] == 0x00) {
						absTicksEnd = float64(te.AbsTicks)
						endNote = b[1]
						if startNote == endNote {
							sign := Name(b[1])
							level := b[1] / 12
							length := (absTicksEnd - absTicksStart) / defaultMetric
							midStr += sign
							if level != 5 {
								midStr += strconv.Itoa(int(level))
							}
							pow := int(math.Round(math.Log2(length)))
							if pow >= -4 && pow != 0 {
								midStr += "<" + strconv.Itoa(pow)
							}
							startNote = 0
							endNote = 0
						}
					}
					if (te.Message.Is(midi.NoteOnMsg) && b[2] > 0) && absTicksStart > absTicksEnd {
						length := (absTicksStart - absTicksEnd) / defaultMetric
						pow := int(math.Round(math.Log2(length)))
						if pow == 0 {
							midStr += "R"
						} else if pow >= -4 {
							midStr += "R<" + strconv.Itoa(pow)
						}
					}
				}
			},
		)
	return
}

// GetNumTracks 获得midi文件音轨数
func GetNumTracks(data []byte) (int, error) {
	s, err := smf.ReadFrom(bytes.NewReader(data))
	return int(s.NumTracks()), err
}
