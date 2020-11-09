package rxgo_test

import (
	"fmt"
	"testing"
	"time"
	"github.com/Hide-on-bush2/rxgo"
)

type observer struct {
	name string
}

func (o observer) OnNext(x interface{}) {
	fmt.Println(o.name, "observed value ", x)
}

func (o observer) OnError(e error) {
	fmt.Println(o.name, "Error ", e)
}

func (o observer) OnCompleted() {
	fmt.Println(o.name, "Down ")
}

func TestMain(t *testing.T) {

	// test Subscribe on any
	ob := rxgo.Just(10, 20, 30).Map(dd)
	ob1 := ob.Map(dd).SubscribeOn(rxgo.ThreadingIO).Debug(true).Map(dd)
	ob1.Subscribe(func(x int) {
		fmt.Println("Just", x)
	})

	ob = rxgo.Just(0, 12, 7, 34, 2).Filter(func(x int) bool {
		return x < 10
	}).SubscribeOn(rxgo.ThreadingIO)
	ob.Subscribe(
		func(x int) {
			fmt.Println("Filter", x)
		})
}

func dd(x int) int { return 2 * x }

func TestObserver(t *testing.T) {
	var s rxgo.Observer = observer{"test observer"}
	rxgo.Just(1, 2, 3).Subscribe(s)
}

func TestTreading(t *testing.T) {
	flow := rxgo.Just(10, 20, 30).Map(func(x int) int {
		return x + 1
	})
	/* 	.FlatMap(func(x int) *rxgo.Observable {
		return rxgo.Just(x+1, x+2)
	}).SubscribeOn(rxgo.ThreadingIO) */

	go flow.Subscribe(observer{"test flatMap"})
	//time.Sleep(time.Nanosecond * 1000)
	flow.Subscribe(observer{"test flatMap again"})
	time.Sleep(time.Microsecond * 1000)
}

func TestDebounce(t *testing.T){
	rxgo.Just("Faker","Theshy","Doinb","Clearlove","Showmaker","Nuguri","Uzi","dasima").Map(func(x string)string{
		if x != "Faker"{
			time.Sleep(1*time.Millisecond)
		}
		return x
	}).Debounce(2*time.Millisecond).Subscribe(func(x string){
		if x != "dasima" {
			t.Errorf("Debounce Fail")
		}
	})
}

func TestDistinct(t *testing.T){
	var all = map[string]bool{}
	rxgo.Just("Faker","Theshy","Doinb","Clearlove","Showmaker","Nuguri","Uzi","dasima", "Faker").Distinct().Subscribe(func(x string){
		if _, ok := all[x]; !ok{
			all[x] = true
			return
		}
		t.Errorf("Distinct Fail")
	})
}

func TestElementAt(t *testing.T){
	rxgo.Just("Faker","Theshy","Doinb","Clearlove","Showmaker","Nuguri","Uzi","dasima").ElementAt("Faker").Subscribe(func(x string){
		if x != "Faker" {
			t.Errorf("ElementAt Fail")
		}
	})
}

func TestFirst(t *testing.T){
	rxgo.Just("Faker","Theshy","Doinb","Clearlove","Showmaker","Nuguri","Uzi","dasima").First(func(x interface{})bool{return true}).Subscribe(func(x string){
		if x != "Faker" {
			t.Errorf("First Fail")
		}
	})
}

func TestIgnoreElements(t *testing.T){
	rxgo.Just("Faker","Theshy","Doinb","Clearlove","Showmaker","Nuguri","Uzi","dasima").IgnoreElements().Subscribe(func(x string){
		t.Errorf("IgnoreElements Fail")
	})
}

func TestLast(t *testing.T){
	rxgo.Just("Faker","Theshy","Doinb","Clearlove","Showmaker","Nuguri","Uzi","dasima").Last().Subscribe(func(x string){
		if x != "dasima"{
			t.Errorf("Last Fail")
		}	
	})
}

func TestSkip(t *testing.T){
	var skiparr = []string{"Showmaker","Nuguri","Uzi","dasima"}
	var count = 0
	rxgo.Just("Faker","Theshy","Doinb","Clearlove","Showmaker","Nuguri","Uzi","dasima").Skip(4).Subscribe(func(x string){
		if x != skiparr[count]{
			t.Errorf("Skip Fail")
		}
		count++
	})
}

func TestSkiplast(t *testing.T){
	var skiparr = []string{"Faker","Theshy","Doinb"}
	var count = 0
	rxgo.Just("Faker","Theshy","Doinb","Clearlove","Showmaker","Nuguri","Uzi","dasima").Skiplast(5).Subscribe(func(x string){
		if x != skiparr[count]{
			t.Errorf("Skiplast Fail")
		}
		count++
	})
}

func TestTake(t *testing.T){
	var takearr = []string{"Faker","Theshy","Doinb"}
	var count = 0
	rxgo.Just("Faker","Theshy","Doinb","Clearlove","Showmaker","Nuguri","Uzi","dasima").Take(3).Subscribe(func(x string){
		if x != takearr[count]{
			t.Errorf("Take Fail")
		}
		count++
	})
}

func TestTakelast(t *testing.T){
	var takelastarr = []string{"Nuguri","Uzi","dasima"}
	var count = 0
	rxgo.Just("Faker","Theshy","Doinb","Clearlove","Showmaker","Nuguri","Uzi","dasima").Takelast(3).Subscribe(func(x string){
		if x != takelastarr[count]{
			t.Errorf("Takelast Fail")
		}
		count++
	})
}