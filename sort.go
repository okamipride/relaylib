package relaylib

import (
	"errors"
	"fmt"
	"sort"
)

//SRNode is element for arrary used for search did and its SessionPair
type SRNode struct {
	keyDid  string
	valSess *SessionPair //A session Pair pointer, point to actual array
}

//SRArr is defined to implement sort interface
type SRArr []SRNode

var searchArr SRArr

// Implement sort library interface
func (arr SRArr) Len() int {
	return len(arr)
}

func (arr SRArr) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr SRArr) Less(i, j int) bool {
	return arr[i].keyDid < arr[j].keyDid
}

// InitArray init array capability to reduce slice operation cause
func InitArray(max uint64) *SRArr {
	searchArr = make(SRArr, 0, max)
	return &searchArr
}

//InsertRSNode insert a new SRNode according to given SessionPair to an sorted array
func (arr *SRArr) InsertRSNode(sp *SessionPair) error {

	did := sp.Did

	findPos := sort.Search(len(*arr), func(i int) bool {
		return (*arr)[i].keyDid >= did
	})

	if findPos < len(*arr) && (*arr)[findPos].keyDid == did {
		fmt.Println(did, " Given session has exist in sorted tree at pos =  ", findPos)
	} else {
		//fmt.Println(did, " is not Found . will be add at index=", findPos)
		node := make([]SRNode, 1, 1)
		node[0].keyDid = did
		node[0].valSess = sp
		insertNode(arr, findPos, node)
	}
	// make sure the node is inserted

	findPos = sort.Search(len(*arr), func(i int) bool {
		return (*arr)[i].keyDid >= did
	})

	if findPos < len(*arr) && (*arr)[findPos].keyDid == did {
		//node has inserted
		return nil
	} else {
		// insert failed
		errMsg := "insert " + did + " fail"
		return errors.New(errMsg)
	}
}

//DelRSNode delete a node according to provide SessionPair's did
func (arr *SRArr) DelRSNode(sp *SessionPair) error {

	did := sp.Did

	findPos := sort.Search(len(*arr), func(i int) bool {
		return (*arr)[i].keyDid >= did
	})

	if findPos < len(*arr) && (*arr)[findPos].keyDid == did {
		fmt.Println(did, " exist in sorted tree at pos =  ", findPos, " and this will be deleted")
		delNode(arr, findPos)
	} else {
		fmt.Println(did, " is not Found . no need to be deleted")
	}

	findPos = sort.Search(len(*arr), func(i int) bool {
		return (*arr)[i].keyDid >= did
	})

	if findPos < len(*arr) && (*arr)[findPos].keyDid == did {
		errMsg := "remove " + did + " faild"
		return errors.New(errMsg)
	} else {
		return nil
	}
}

// SortRSNode sort an SRArr by SRNode did from small to large
func (arr *SRArr) SortRSNode() {
	sort.Sort(arr)
}

func insertNode(arr *SRArr, idx int, node []SRNode) {
	rightLen := len(*arr) - idx + 1
	iNode := make([]SRNode, 1, rightLen)
	copy(iNode, node)
	*arr = append((*arr)[:idx], append(iNode, (*arr)[idx:]...)...)
}

func delNode(arr *SRArr, idx int) {
	copy((*arr)[idx:], (*arr)[idx+1:])
	(*arr)[len(*arr)-1] = SRNode{} // prevent memory leak, this shall be verified in the future
	*arr = (*arr)[:len(*arr)-1]
	arr.PrintTree()
}

func (arr *SRArr) PrintTree() {
	for i, c := range *arr {
		fmt.Println(i, c.keyDid)
	}
}
