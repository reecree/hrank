package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var COUNT = 5

type Node struct {
	Parent     *Node
	Left       *Node
	Right      *Node
	Value      int32
	Height     int
	ChildCount int
}

type AVL struct {
	Root *Node
	size int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (n *Node) fixChildCountAndHeight() {
	n.ChildCount = 0
	if n.Left != nil {
		n.ChildCount = n.Left.ChildCount + 1
	}
	if n.Right != nil {
		n.ChildCount += n.Right.ChildCount + 1
	}
	lHeight, rHeight := 0, 0
	if n.Left != nil {
		lHeight = n.Left.Height
	}
	if n.Right != nil {
		rHeight = n.Right.Height
	}
	n.Height = max(rHeight, lHeight) + 1
}

func print2DUtil(root *Node, space int) {
	// if times > 15 {
	// 	return
	// }
	// Base case
	if root == nil {
		return
	}

	// Increase distance between levels
	space += COUNT

	// Process right child first
	print2DUtil(root.Right, space)

	// Print current node after space
	// count
	//fmt.Println()
	for i := COUNT; i < space; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("%d %d %d\n", root.Value, root.ChildCount, root.BalanceFactor())

	// Process left child
	print2DUtil(root.Left, space)
}

// Wrapper over print2DUtil()
func (avl *AVL) print2D() {
	// Pass initial space count as 0
	print2DUtil(avl.Root, 0)
	fmt.Println(avl.size)
}

func (avl *AVL) GetMedian() float32 {
	size := avl.size
	if size == 0 {
		return -1
	}
	root := avl.Root

	if size%2 == 0 {
		b := root.getNode(size/2 - 1).Value
		a := root.getNode(size / 2).Value
		return 0.5 * float32(a+b)
	}
	return float32(root.getNode(size / 2).Value)
}

func (n *Node) getNode(index int) *Node {
	current := n
	for {
		leftSubtreeSize := current.getLeftSubtreeSize()

		if index == leftSubtreeSize {
			return current
		}

		if index > leftSubtreeSize {
			index -= (leftSubtreeSize + 1)
			current = current.Right
		} else {
			current = current.Left
		}
	}
}

func (n *Node) getLeftSubtreeSize() int {
	if n == nil {
		return 0
	}
	tmp := n.ChildCount

	if n.Right != nil {
		tmp -= (n.Right.ChildCount + 1)
	}

	return tmp
}

func (n *Node) BalanceFactor() int {
	var rHeight, lHeight = 0, 0
	if n.Right != nil {
		rHeight = n.Right.Height
	}
	if n.Left != nil {
		lHeight = n.Left.Height
	}
	return rHeight - lHeight
}

func (avl *AVL) Delete(key int32) {
	avl.size--
	avl.Root = deleteNode(avl.Root, key)
}

func deleteNode(root *Node, val int32) *Node {
	if root == nil {
		return root
	}

	if val < root.Value {
		root.Left = deleteNode(root.Left, val)
	} else if val > root.Value {
		root.Right = deleteNode(root.Right, val)
	} else {
		// node with only one child or no child
		if root.Left != nil && root.Right != nil {
			temp := root.findMin()

			root.Value = temp.Value

			// Delete the inorder successor
			root.Right = deleteNode(root.Right, temp.Value)
		} else {
			var temp *Node
			if root.Left != nil {
				temp = root.Left
			} else if root.Right != nil {
				temp = root.Right
			} else {
				return nil
			}

			*root = *temp
		}
	}

	root.fixChildCountAndHeight()
	b := root.BalanceFactor()

	// Left Left Case
	if b > 1 {
		bRight := root.Right.BalanceFactor()
		if bRight < 0 {
			return rightLeftRotate(root)
		}
		return leftRotate(root)
	}

	// Right Right Case
	if b < -1 {
		bLeft := root.Left.BalanceFactor()
		if bLeft > 0 {
			return leftRightRotate(root)
		}
		return rightRotate(root)
	}

	return root
}

func (n *Node) findMin() *Node {
	cur := n.Right
	for {
		if cur.Left == nil {
			return cur
		}
		cur = cur.Left
	}
}

func (avl *AVL) Insert(val int32) {
	avl.size++
	avl.Root = insertR(avl.Root, val)
}

func insertR(root *Node, val int32) *Node {
	if root == nil {
		root = &Node{Value: val, Height: 1}
		return root
	}

	walkLeft, walkRight := false, false
	if val == root.Value {
		b := root.BalanceFactor()
		if b > 0 {
			walkLeft = true
		} else if b < 0 {
			walkRight = true
		} else {
			if rand.Intn(2) == 0 {
				walkLeft = true
			} else {
				walkRight = true
			}
		}
	}

	if val < root.Value || walkLeft {
		root.Left = insertR(root.Left, val)
		if root.BalanceFactor() < -1 {
			if root.Left.BalanceFactor() < 0 {
				root = rightRotate(root)
			} else {
				root = leftRightRotate(root)
			}
		}
	} else if val > root.Value || walkRight {
		root.Right = insertR(root.Right, val)
		if root.BalanceFactor() > 1 {
			if root.BalanceFactor() > 0 {
				root = leftRotate(root)
			} else {
				root = rightLeftRotate(root)
			}
		}
	}

	root.fixChildCountAndHeight()
	return root
}

func leftRotate(root *Node) *Node {
	node := root.Right
	root.Right = node.Left
	node.Left = root
	root.fixChildCountAndHeight()
	node.fixChildCountAndHeight()
	return node
}

func rightRotate(root *Node) *Node {
	node := root.Left
	root.Left = node.Right
	node.Right = root
	root.fixChildCountAndHeight()
	node.fixChildCountAndHeight()
	return node
}

func rightLeftRotate(X *Node) *Node {
	X.Right = rightRotate(X.Right)
	X = leftRotate(X)
	return X
}

func leftRightRotate(X *Node) *Node {
	X.Left = leftRotate(X.Left)
	X = rightRotate(X)
	return X
}

// Complete the activityNotifications function below.
func activityNotifications(expenditure []int32, d int32) int32 {
	normalD := int(d)
	avl := AVL{}
	var notifications int32
	for i, exp := range expenditure {
		if i >= normalD {
			if int32(avl.GetMedian()*2) <= exp {
				notifications++
			}
			avl.Insert(exp)
			avl.Delete(expenditure[i-normalD])
		} else {
			avl.Insert(exp)
		}
	}
	return notifications
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	nd := strings.Split(readLine(reader), " ")

	nTemp, err := strconv.ParseInt(nd[0], 10, 64)
	checkError(err)
	n := int32(nTemp)

	dTemp, err := strconv.ParseInt(nd[1], 10, 64)
	checkError(err)
	d := int32(dTemp)

	expenditureTemp := strings.Split(readLine(reader), " ")

	var expenditure []int32

	for i := 0; i < int(n); i++ {
		expenditureItemTemp, err := strconv.ParseInt(expenditureTemp[i], 10, 64)
		checkError(err)
		expenditureItem := int32(expenditureItemTemp)
		expenditure = append(expenditure, expenditureItem)
	}

	result := activityNotifications(expenditure, d)

	fmt.Fprintf(writer, "%d\n", result)

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
