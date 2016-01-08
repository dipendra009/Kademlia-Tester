package main
import (
    "bufio"
    "fmt"
    "os/exec"
)
func main() {
    // What we want to calculate
    calcs := make([]string, 1)
    calcs[0] = "whoami"//"3*3"//
    // To store the results
    results := make([]string, 1)
    //---------------------------------------------
    //cmd := exec.Command("/usr/bin/bc")
    cmd := exec.Command("./mainRef","localhost:1234","localhost:1234")
    in, err := cmd.StdinPipe()
    if err != nil {
        panic(err)
    }
    defer in.Close()
    out, err := cmd.StdoutPipe()
    if err != nil {
        panic(err)
    }
    defer out.Close()
    // We want to read line by line
    bufOut := bufio.NewReader(out)
    //bufin := bufio.NewWriter(in)
    //----------------------------------------------
    // Start the process
    if err = cmd.Start(); err != nil{
        panic(err)
    }
    // Write the operations to the process
    // for _, calc := range calcs {
        _, err1 := in.Write([]byte(calcs[0] + "\n"))
        if err != nil {
            panic(err1)
        }
//    }
defer in.Close()
var id string
    // Read the results from the process
    for i := 0; i < 2; i++ {
        result,_,err := bufOut.ReadLine()
        if err != nil {
            panic(err)
        }
        results[0] = string(result)
        id = results[0]
		fmt.Println(results[0])
    }
    // See what was calculated
//    for _, result := range results {
//        fmt.Println(result)
//    }
    // We want to read line by line
//    bufOut = bufio.NewReader(out)
    //--------------------------------------------
    // Write the operations to the process
        _, err = in.Write([]byte("get_contact "+id+ "\n"))
        if err != nil {
            panic(err)
        }
    // Read the results from the process
 //   for i := 0; i < len(results); i++ {
        result, _, err := bufOut.ReadLine()
        if err != nil {
            panic(err)
        }
        results[0] = string(result)
   // }
    // See what was calculated
    for _, result := range results {
        fmt.Println(result)
    }
}