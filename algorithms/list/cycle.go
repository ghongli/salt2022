package list

type Node struct {
	Val  int
	post int
	Next *Node
}

// 方法一：快慢指针
// 时间复杂度：O(N)，其中 N 是链表中的节点数。
//      当链表中不存在环时，快指针将先于慢指针到达链表尾部，链表中每个节点至多被访问两次。
//      当链表中存在环时，每一轮移动后，快慢指针的距离将减小一。而初始距离为环的长度，因此至多移动 N 轮。
// 空间复杂度：O(1)。只使用了两个指针的额外空间。
func hasCycleWithFastSlow(head *Node) bool {
	if head == nil || head.Next == nil {
		return false
	}
	
	fast, slow := head.Next, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		
		if slow == fast {
			return true
		}
	}

	return false
}

// 方法二：哈希表存储
// 时间复杂度：O(N)，其中 N 是链表中的节点数。
// 空间复杂度：O(N)，其中 N 是链表中的节点数。主要为哈希表的开销，最坏情况下我们需要将每个节点插入到哈希表中一次。
func hasCycleWithMap(head *Node) bool {
	seen := map[*Node]struct{}{}
	
	for head != nil {
		if _, ok := seen[head]; ok {
			return true
		}
		
		seen[head] = struct{}{}
		head = head.Next
	}

	return false
}
