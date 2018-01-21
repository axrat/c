package main

import (
	"github.com/TransAssist/goz"
	"github.com/mattn/go-shellwords"
)
import (
	_ "github.com/mattn/go-sqlite3"
)
import (
	"fmt"
	"time"
	"os"
	"os/exec"
)

var current = goz.ReturnAndExitIfError(os.Getwd()).(string) + "/"
var cr = current + "cr/"
var ok = cr + "ok"
var cmds = cr + "cmds"
var fgo = cr + "f.go"
var res = cr + "res"

const limit = 9999

func proc() {
	if !goz.Exists(cmds) {
		if !goz.Exists(fgo) {
			goz.Puts(fgo, goz.Skel())
		}
		goz.Puts(cmds, "go run f.go")
	}

	cmdslice := goz.ReadFile(cmds)
	for index, value := range cmdslice {
		//fmt.Printf("[%v]:%v\n", index, value )
		fmt.Printf("L%v:[", index)
		args, err := shellwords.Parse(value)
		for _, value := range args {
			fmt.Printf("_%v", value)
		}
		fmt.Printf("]\n")
		goz.PrintIfError(err)
		switch len(args) {
		case 1: //no args cmd
			Execute(exec.Command(args[0]))
		default: //args cmd
			Execute(exec.Command(args[0], args[1:]...))
		}
	}
}
func Execute(cmd *exec.Cmd) {
	cmd.Dir = cr
	out, err := cmd.CombinedOutput()
	goz.PrintIfError(err)
	fmt.Println(string(out))
	goz.Puts(res, string(out))
}
func before(){
	goz.PrintIfError(os.MkdirAll(cr, 0777))
	os.Chdir(cr)
}
func main() {
	t := time.NewTicker(1 * time.Second)
	count := 0
	goz.PrintIfError(os.Remove(ok))
	L:
		for {
			select {
			case <-t.C:
				if count > limit {
					fmt.Printf("ProcessLimit")
					t.Stop()
					break L
				}
				before()
				if !goz.Exists(ok) {
					fmt.Printf("#+%v\n", goz.Epoch())
					goz.Puts(ok, time.Now().Local().String())
					count++
					proc()
					fmt.Printf("#-%v@%v\n", time.Now().Format("20060102030405"), count)
				}
			}
		}
	goz.Complete()
}
