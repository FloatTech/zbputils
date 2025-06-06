package chat

// findRepeatedPattern 查找字符串中连续重复次数超过 n 的子串
func findRepeatedPattern(s string, n int) (string, int) {
	runes := []rune(s) // 将字符串转换为 rune 切片，支持 Unicode 字符
	length := len(runes)

	// 遍历所有可能的子串长度，从 1 到字符串长度的一半
	for size := 1; size <= length/2; size++ {
		// 遍历字符串的每个起始位置
		for i := 0; i <= length-size*n; i++ {
			pattern := runes[i : i+size] // 当前子串模式
			count := 1                   // 初始化重复次数

			// 检查后续是否连续出现相同的子串
			for j := i + size; j+size <= length; j += size {
				if string(runes[j:j+size]) == string(pattern) {
					count++
					if count >= n {
						return string(pattern), count
					}
				} else {
					break
				}
			}
		}
	}
	return "", 0 // 未找到满足条件的重复子串
}
