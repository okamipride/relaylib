package relaylib

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

//SRNode store a did and a pointer to a SessionPair which has the same did
type SRNode struct {
	Did   string
	SpPtr *SessionPair //A session Pair pointer, point to actual array
}

//SRArr is a slice of SRNode(s)。 it is used for store SRNode and thus could perform variety of operation
type SRArr []SRNode
//Array used for searching only
var searchArr SRArr

// Implement sort library interface
func (arr SRArr) Len() int {
	return len(arr)
}

func (arr SRArr) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr SRArr) Less(i, j int) bool {
	var ret bool = false
	if arr[i].Did < arr[j].Did {
		ret = true
	} else if arr[i].Did == arr[j].Did {
		if arr[i].SpPtr.ClientIsJoin == false {
			ret = true
		}
	}
	return ret
}

// InitArray  returns a SRArr with max capability
// Use InitArray to get a pointer of SRArr to perform operation
func InitSRArray(max uint64) SRArr {
	searchArr = make(SRArr, 0, max)
	return searchArr
}

//InsertRSNode insert a new SRNode into SRArr, according to a given SessionPair
func (arr *SRArr) InsertRSNode(sp *SessionPair) error {

	did := sp.did

	findPos := sort.Search(len(*arr), func(i int) bool {
		return (*arr)[i].Did >= did
	})
	node := make([]SRNode, 1, 1)
	node[0].Did = did
	node[0].SpPtr = sp

	insertNode(arr, findPos, node)
	if findPos < len(*arr) {
		for ; findPos < len(*arr); findPos++ {
			if (*arr)[findPos].SpPtr.SessionID == sp.SessionID {
				break
			}
		}
		if findPos < len(*arr) {
			return nil
		} else {
			erMsg := "Add SesssionPair id=" + string(sp.did) + " Fail"
			return errors.New(erMsg)
		}
	} else {
		erMsg := "Add SesssionPair id=" + string(sp.did) + " Fail"
		return errors.New(erMsg)
	}

}

//DelRSNode delete a node according to given SessionPair's did
func (arr *SRArr) DelRSNode(sp *SessionPair) error {

	did := sp.did
	//fmt.Println("DelRSNode did=", did, "sid=", sp.SessionID)
	findPos := sort.Search(len(*arr), func(i int) bool {
		return (*arr)[i].Did >= did
	})

	if findPos < len(*arr) {
		for ; findPos < len(*arr); findPos++ {
			if (*arr)[findPos].SpPtr.SessionID == sp.SessionID {
				break
			}
		}
		if findPos < len(*arr) {
			delNode(arr, findPos)
		} else {
			return errors.New("no session found")
		}
	} else {
		return errors.New("no session found")
	}

	// Check if has deleted in the SRArr
	findPos = sort.Search(len(*arr), func(i int) bool {
		return (*arr)[i].Did >= did
	})

	// Double Check
	_, findErr := arr.FindSess(sp)
	if findErr != nil {
		errMsg := "remove session fail sessionid=" + string(sp.SessionID) + "and session exist"
		return errors.New(errMsg)
	}
	return nil
}

// SortRSNode sort a SRArr
func (arr *SRArr) SortRSNode() {
	sort.Sort(arr)
}

// FindSess find a SessionPair in
func (arr *SRArr) FindSess(sp *SessionPair) (SRNode, error) {
	did := sp.did
	findPos := sort.Search(len(*arr), func(i int) bool {
		return (*arr)[i].Did >= did
	})

	if findPos < len(*arr) {
		for ; findPos < len(*arr); findPos++ {
			if (*arr)[findPos].SpPtr.SessionID == sp.SessionID {
				break
			}
		}
		if findPos < len(*arr) {
			return (*arr)[findPos], nil
		} else {
			return SRNode{}, errors.New("no session found")
		}
	} else {
		return SRNode{}, errors.New("no session found")
	}

}

//FindUnusedSessByDid search for the did in SRArr and return a SessionPair pointer with the same did and the position of the session in the origin array
func (arr *SRArr) FindUnPairSess(did string) (SRNode, error) {

	// find the small index that match did
	findPos := sort.Search(len(*arr), func(i int) bool {
		return (*arr)[i].Did >= did
	})

	//find the first ClientIsJoin == false session after small did
	if findPos < len(*arr) {
		for findPos < len(*arr) && (*arr)[findPos].SpPtr.ClientIsJoin != false {
			findPos = findPos + 1
		}
		return (*arr)[findPos], nil
	} else {
		return SRNode{}, errors.New("not found")
	}
}

func (arr *SRArr) IsDidExist(did string) bool {

	findPos := sort.Search(len(*arr), func(i int) bool {
		return (*arr)[i].Did >= did
	})

	//find the first ClientIsJoin == false session after small did
	if findPos < len(*arr) {
		return true
	}
	return false
}

func insertNode(arr *SRArr, idx int, node []SRNode) {
	rightLen := len(*arr) - idx + 1
	iNode := make([]SRNode, 1, rightLen)
	copy(iNode, node)
	*arr = append((*arr)[:idx], append(iNode, (*arr)[idx:]...)...)
	//arr.PrintTree()
}

func delNode(arr *SRArr, idx int) {
	copy((*arr)[idx:], (*arr)[idx+1:])
	(*arr)[len(*arr)-1] = SRNode{} // prevent memory leak, this shall be verified in the future
	*arr = (*arr)[:len(*arr)-1]
	//arr.PrintTree()
}

func (arr *SRArr) PrintTree() {
	for i, c := range *arr {
		fmt.Println(i, c.Did, c.SpPtr.ClientIsJoin, c.SpPtr.SessionID)
	}
}

//TimeTrack helper function
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}
