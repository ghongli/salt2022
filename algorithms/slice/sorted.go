// Package slice 合并两个有序数组
package slice

import (
	"sort"
)

// 方法一：直接合并后排序
func mergeWithDirect(nums1 []int, m int, nums2 []int, _ int) {
	copy(nums1[m:], nums2)
	sort.Ints(nums1)
}

// 方法二：
// func mergeWithList(nums1 []int, m int, nums2 []int, n int) {
// 	sorted := make([]int, 0, m+n)
// 	l1, l2 := 0, 0
// 	for l1 < m || l2 < n {
// 		if l1 == m {
// 			sorted = append(sorted, nums2[l2:]...)
// 		}
// 		if l2 == n {
// 			sorted = append(sorted, nums1[l1:]...)
// 		}
// 		if nums1[l1] < nums2[l2] {
// 			sorted = append(sorted, nums1[l1])
// 			l1++
// 		} else {
// 			sorted = append(sorted, nums2[l2])
// 			l2++
// 		}
// 	}
// 	copy(nums1, sorted)
// }

// 方法三：逆向双指针
func mergeWithReversed(nums1 []int, m int, nums2 []int, n int) {
	for l1, l2, tail := m-1, n-1, m+n-1; l1 >= 0 || l2 >= 0; tail-- {
		var cur int
		if l1 == -1 {
			cur = nums2[l2]
			l2--
		} else if l2 == -1 {
			cur = nums1[l1]
			l1--
		} else if nums1[l1] > nums2[l2] {
			cur = nums1[l1]
			l1--
		} else {
			cur = nums2[l2]
			l2--
		}
		nums1[tail] = cur
	}
}
