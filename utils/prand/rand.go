// 随机数生成工具, 更多随机数生成方法，可对照标准库“math/rand”
package prand

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

var (
	letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52 英文字符
	symbols    = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"                   // 32 特殊字符
	digits     = "0123456789"                                           // 10 数字字符
	characters = letters + digits + symbols                             // 94

	// bufferChan is the buffer for random bytes,
	// every item storing 4 bytes.
	bufferChan = make(chan []byte, bufferChanSize)
)

const (
	// Buffer size for uint32 random number.
	bufferChanSize = 10000
)

func init() {
	go asyncProducingRandomBufferBytesLoop()
}

// asyncProducingRandomBufferBytes is a named goroutine, which uses an asynchronous goroutine
// to produce the random bytes, and a buffer chan to store the random bytes.
// So it has high performance to generate random numbers.
func asyncProducingRandomBufferBytesLoop() {
	var step int
	for {
		buffer := make([]byte, 1024)
		if n, err := rand.Read(buffer); err != nil {
			panic("error reading random buffer from system")
		} else {
			// The random buffer from system is very expensive,
			// so fully reuse the random buffer by changing
			// the step with a different number can
			// improve the performance a lot.
			// for _, step = range []int{4, 5, 6, 7} {
			for _, step = range []int{4} {
				for i := 0; i <= n-4; i += step {
					bufferChan <- buffer[i : i+4]
				}
			}
		}
	}
}

// Intn 方法返回大于等0且不大于max的随机整数，即：[0, max]
//
// Note that:
// 1. The `max` can only be greater than 0, or else it returns `max` directly;
// 2. The result is greater than or equal to 0, but less than `max`;
// 3. The result number is 32bit and less than math.MaxUint32.
func Intn(max int) int {
	if max <= 0 {
		return max
	}
	n := int(binary.LittleEndian.Uint32(<-bufferChan)) % max
	if (max > 0 && n < 0) || (max < 0 && n > 0) {
		return -n
	}
	return n
}

// Bytes 方法用于返回指定长度的二进制[]byte数据
func Bytes(n int) []byte {
	if n <= 0 {
		return nil
	}
	i := 0
	b := make([]byte, n)
	for {
		copy(b[i:], <-bufferChan)
		i += 4
		if i >= n {
			break
		}
	}
	return b
}

// Int方法返回min到max之间的随机整数，支持负数，包含边界，即：[min, max]
func Int(min, max int) int {
	if min >= max {
		return min
	}
	if min >= 0 {
		return Intn(max-min+1) + min
	}
	// As `Intn` dose not support negative number,
	// so we should first shift the value to right,
	// then call `Intn` to produce the random number,
	// and finally shift the result back to left.
	return Intn(max+(0-min)+1) - (0 - min)
}

// 方法用于返回指定长度的数字、字符，第二个参数symbols用于指定是否返回的随机字符串中包含特殊字符，默认false
func String(n int, symbols ...bool) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = Bytes(n)
	)
	for i := range b {
		if len(symbols) > 0 && symbols[0] {
			b[i] = characters[numberBytes[i]%94]
		} else {
			b[i] = characters[numberBytes[i]%62]
		}
	}
	return string(b)
}

// Duration 返回一个随机 time.Duration 类型值，即: [min, max].
func Duration(min, max time.Duration) time.Duration {
	multiple := int64(1)
	if min != 0 {
		for min%10 == 0 {
			multiple *= 10
			min /= 10
			max /= 10
		}
	}
	n := int64(Int(int(min), int(max)))
	return time.Duration(n * multiple)
}

// StrRand 方法是一个比较高级的方法，用于从给定的字符列表中选择指定长度的随机字符串返回，
// 并且支持unicode字符，例如中文。例如，Str("中文123abc", 3)将可能会返回1a文的随机字符串。
func StrRand(s string, n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b     = make([]rune, n)
		runes = []rune(s)
	)
	if len(runes) <= 255 {
		numberBytes := Bytes(n)
		for i := range b {
			b[i] = runes[int(numberBytes[i])%len(runes)]
		}
	} else {
		for i := range b {
			b[i] = runes[Intn(len(runes))]
		}
	}
	return string(b)
}

// Digits 方法用于返回指定长度的随机数字字符串
func Digits(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = Bytes(n)
	)
	for i := range b {
		b[i] = digits[numberBytes[i]%10]
	}
	return string(b)
}

// Letters 方法用于返回指定长度的随机英文字符串.
func Letters(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = Bytes(n)
	)
	for i := range b {
		b[i] = letters[numberBytes[i]%52]
	}
	return string(b)
}

// Symbols 方法用于返回指定长度的随机特殊字符串.
func Symbols(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = Bytes(n)
	)
	for i := range b {
		b[i] = symbols[numberBytes[i]%32]
	}
	return string(b)
}

// Perm 产生一个整数切片，范围[0,n]
func Perm(n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}

// Meet 用于指定一个数num和总数total，往往 num<=total，并随机计算是否满足num/total的概率。
// 例如，Meet(1, 100)将会随机计算是否满足百分之一的概率。
func Meet(num, total int) bool {
	return Intn(total) < num
}

// MeetProb 用于给定一个概率浮点数prob，往往 prob<=1.0，并随机计算是否满足该概率。
// 例如，MeetProb(0.005)将会随机计算是否满足千分之五的概率。
func MeetProb(prob float32) bool {
	return Intn(1e7) < int(prob*1e7)
}
