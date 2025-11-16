package main

import (
	"fmt"
	"sort"
)

// 1.只出现一次的数字
func singleNumber(nums []int) int {
	data := make(map[int]int)
	for _, v := range nums {
		data[v]++
	}

	for index, value := range data {
		if value == 1 {
			return index
		}
	}
	return -1
}

// 2.给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
func isValid(s string) bool {
	data := make([]string, 0, 10)
	var left, right int

	for _, v := range s {
		switch v {
		case '(':
			data = append(data, "(")
			left++
		case '[':
			data = append(data, "[")
			left++
		case '{':
			data = append(data, "{")
			left++
		case '}':
			n := len(data)
			if n > 0 && data[n-1] == "{" {
				data = data[:n-1]
			}
			right++
		case ']':
			n := len(data)
			if n > 0 && data[n-1] == "[" {
				data = data[:n-1]
			}
			right++
		case ')':
			n := len(data)
			if n > 0 && data[n-1] == "(" {
				data = data[:n-1]
			}
			right++
		}
	}

	fmt.Println(data)
	return len(data) == 0 && left == right
}

// 3.最长公共前缀,
func longestCommonPrefix(strs []string) string {
	firstStr := strs[0]
	m := len(firstStr)
	n := len(strs)
	stepFlag := true
	result := make([]byte, 0, m)

	for i := 0; i < m; i++ {
		cmpflag := true
		for j := 1; j < n; j++ {
			if i == len(strs[j]) {
				cmpflag = false
				break
			}
			fmt.Println(strs[j], i)
			stepFlag = stepFlag && (strs[j][i] == firstStr[i])

		}
		if stepFlag && cmpflag {
			result = append(result, byte(firstStr[i]))
		} else {
			break
		}
	}
	return string(result)
}

// 4.给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。这些数字按从左到右，从最高位到最低位排列。这个大整数不包含任何前导 0。
// 将大整数加 1，并返回结果的数字数组。
func plusOne(digits []int) []int {
	n := len(digits)
	data := digits[:n]

	increaseValue := 0
	if digits[n-1] < 9 {
		data[n-1] = digits[n-1] + 1
	} else if digits[n-1] == 9 {
		data[n-1] = digits[n-1] + 1

		for i := n - 1; i > 0; i-- {
			data[i] = data[i] + increaseValue
			if data[i] > 9 {
				increaseValue = data[i] / 10
				data[i] = data[i] % 10
			} else {
				increaseValue = 0
			}
		}

		data[0] = data[0] + increaseValue
		if data[0] > 9 {
			newBit := data[0] / 10
			data[0] = data[0] % 10
			data = append([]int{newBit}, data...)
		}
	}

	return data
}

// 5.删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[i-1] {
			nums = append(nums[0:i], nums[i+1:]...)
			i--
		}
	}
	return len(nums)
}

// 6.合并区间：以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。可以先对区间数组按照区间
// 的起始位置进行排序，然后使用一个切片来存储合并后的区间，遍历排序后的区间数组，将当前区间与切片中最后一个区间进行比较，
// 如果有重叠，则合并区间；如果没有重叠，则将当前区间添加到切片中。
func mergeIntervals(intervals [][]int) []int {
	data := make([]int, 0)
	for index, interval := range intervals {
		sort.Ints(interval)
		data = append(data, intervals[index]...)
	}
	sort.Ints(data)
	fmt.Println(data)

	for i := 1; i < len(data); i++ {
		if data[i] == data[i-1] {
			data = append(data[:i], data[i+1:]...)
			fmt.Println(data)
			i--
		}
	}
	fmt.Println(data, len(data))
	return data

}

//7.给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。
//你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。
//你可以按任意顺序返回答案。

func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return []int{}
}
func main() {

	data := []int{
		1, 2, 3, 2, 4, 5, 1, 3, 5,
	}
	fmt.Println(singleNumber(data))

	fmt.Println(twoSum([]int{2, 7, 11, 15, 4}, 11))

	//fmt.Println(mergeIntervals([][]int{{1, 2, 3}, {3, 4, 5}, {5, 6, 7}, {7, 8, 9}}))

	//fmt.Println(removeDuplicates([]int{1, 2, 2, 2, 3, 3, 3, 3, 4}))
	//fmt.Println(plusOne([]int{9, 9, 9}))
	//strs := []string{
	//	"abcd",
	//	"a",
	//}
	//
	//fmt.Println(longestCommonPrefix(strs))

	//	fmt.Println(isValid("([}}])"))

	//fmt.Println("Hello World")
}
