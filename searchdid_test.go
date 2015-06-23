package relaylib

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

var testSessArr [MAX_SESSION]SessionPair
var mySortArr SRArr // create a sor arry
var testCount int = 0

/*
Cmd example :
	go test
 	go test -bench . -benchtime 2s
*/

func TestMain(m *testing.M) {
	n := 5
	fmt.Println("run TestMain, Sample size=", n)
	setup(n)
	retcode := m.Run()
	os.Exit(retcode)
}

func setup(n int) {
	genMySP(n)
	if n <= 10 {
		fmt.Println("---Print Unsorted Array ----")
		fmt.Println("---------- D i d -----------------   Joined? sessionID")
		PrintSPTree()
	}

}

func TestInsertRSNode(t *testing.T) {
	fmt.Println("run TestInsertRSNode")
	mySortArr = InitSRArray(MAX_SESSION)
	for i := 0; i < testCount; i++ {
		sp := testSessArr[i]
		mySortArr.InsertRSNode(&sp)
	}
	if testCount <= 10 {
		fmt.Println("---Print Sorted Array ----")
		fmt.Println("index --------- D i d ----------------  Joined? sessionID")
		mySortArr.PrintTree() // printout current SessionPair array
	}

	if testCount != len(mySortArr) {
		t.Error("not all insertion successfully")
	}

	for j := 1; j < len(mySortArr); j++ {
		//fmt.Println("mySortArr j=", j, mySortArr[j].Did)
		if mySortArr[j].Did < mySortArr[j-1].Did { //Fail case
			//fmt.Println("pos j=", j, "val=", mySortArr[j].Did, "should >= previous ", mySortArr[j-1].Did, " but not")
			t.Error("Expect greater or equal but not ")
		}
	}
}

func TestFindUnPairSess(t *testing.T) {
	fmt.Println("run TestFindUnPairSess")
	for i := 0; i < testCount; i++ {
		sp := testSessArr[i]
		srNode, finderr := mySortArr.FindUnPairSess(sp.did)
		if finderr != nil {
			t.Error("Expect find UnPaired SessionPair but not found")
		} else if finderr == nil && srNode.SpPtr.ClientIsJoin == true {
			t.Error("Expect find ClientIsJoin false but find true")
		}
	}
}

func TestFindSess(t *testing.T) {
	fmt.Println("run TestFindSess")
	for i := 0; i < testCount; i++ {
		sp := testSessArr[i]
		srNode, finderr := mySortArr.FindSess(&sp)
		if finderr != nil {
			t.Error("Expect find SessionPair but not ")
		} else if finderr == nil && sp.SessionID != srNode.SpPtr.SessionID {
			t.Error("Expect find the same SessionID but find different")
		}
	}
}

func TestDelRSNode(t *testing.T) {
	fmt.Println("run TestDelRSNode")
	for i := 0; i < testCount; i++ {
		//mySortArr.PrintTree()
		sp := testSessArr[i]
		mySortArr.DelRSNode(&sp)

		_, finderr := mySortArr.FindSess(&sp)
		if finderr == nil {
			t.Error("Expect delete and not found but found")
		}
	}
	if len(mySortArr) > 0 {
		t.Error("Expect delete all but remain size =", len(mySortArr))
	}
}

func BenchmarkFindUnPairSess(b *testing.B) {
	if testCount > len(mySortArr) {
		for i := 0; i < testCount; i++ {
			sp := testSessArr[i]
			mySortArr.InsertRSNode(&sp)
		}
	}
	//fmt.Println("Test count", len(mySortArr))
	b.ResetTimer()
	fmt.Println("run BenchmarkFindUnPairSess")
	for i := 0; i < b.N; i++ {
		find := rand.Int() % testCount
		mySortArr.FindUnPairSess(testSessArr[find].did)
	}
}

func BenchmarkFindForLoop(b *testing.B) {
	fmt.Println("run BenchmarkFindLinear")
	for i := 0; i < b.N; i++ {
		find := rand.Int() % testCount
		findLinear(testSessArr[find].did)
	}
}

func BenchmarkSort(b *testing.B) {
	fmt.Println("run BenchmarkSort")
	var myUnSortArr SRArr = make(SRArr, 0, MAX_SESSION)

	for i := 0; i < testCount; i++ {
		sp := testSessArr[i]
		node := new(SRNode)
		node.Did = sp.did
		node.SpPtr = &sp
		myUnSortArr = append(myUnSortArr, *node)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		myUnSortArr.SortRSNode()
	}
}

func findLinear(did string) int {
	i := 0
	for i = 0; i < testCount; i++ {
		if testSessArr[i].did == did {
			break
		}
	}
	return i
}

func genMySP(n int) {
	var sp SessionPair
	var prev SessionPair
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			sp = genSession(false, i)
			prev = sp
			addMySP(sp)
		} else {
			prev.ClientIsJoin = true
			prev.SessionID = i
			addMySP(prev)
		}
	}
}

/*
genDid generate a did with 32 chars from 0-9A-Z
*/
func genDid() string {
	num1 := rand.Int63n(1099511627775)
	num2 := rand.Int63n(1099511627775)
	num3 := rand.Int63n(10000000)
	return fmt.Sprintf("%010x", num1) + fmt.Sprintf("%010x", num2) + fmt.Sprintf("%012x ", num3)
}

/*
	GenSession generate a SessionPair with random did
*/
func genSession(sessJoin bool, sessId int) SessionPair {
	ret := new(SessionPair)
	ret.did = genDid()
	ret.SessionID = sessId
	ret.ClientIsJoin = sessJoin
	return *ret
}

//TimeTrack helper function
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}

func PrintSPTree() {
	for i := 0; i < testCount; i++ {
		fmt.Println(testSessArr[i].did, testSessArr[i].ClientIsJoin, testSessArr[i].SessionID)
	}
}

func addMySP(sp SessionPair) int {
	testSessArr[testCount] = sp
	testCount++
	return testCount - 1
}
