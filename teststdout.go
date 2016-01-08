package main
import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	//"strings"
	//"io"
	"time"
	"os"
)
// Set your timeout
const CommandTimeout = 5 * time.Second
func main() {
	cmd := exec.Command("./main","localhost:1235","localhost:1235")//, ">", "outputpipe.txt")
	// Set up the input
	in, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("failed to create pipe for STDIN: %s", err)
	}
	// open output file
	fo, err := os.Create("output1.txt")
	if err != nil { panic(err) }
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	
	// Capture the output
	var b bytes.Buffer
	cmd.Stdout, cmd.Stderr = &b, &b
	
	// Start the process
	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to start command: %s", err)
	}
	
	
	// Write the input and close
	go func() {
		defer in.Close()
		fmt.Fprintln(in, "whoami")
		fmt.Fprintln(in, "exit") 
	}()
	// Wait for the process to finish
	if err := cmd.Wait(); err != nil {
		log.Fatalf("command failed....: %s", err)
	}
	// Print out the output
	fmt.Printf("Output:\n%s", b.String())
	c :=b.String()
	d :=[]byte(c)
	if _, err := fo.Write(d); err != nil {
		panic(err)
	}
	
	go func() {
		defer in.Close()
		fmt.Fprintln(in, "whoami")
		fmt.Fprintln(in, "exit") 
	}()
	// Wait for the process to finish
	//if err = cmd.Wait(); err != nil {
	//	log.Fatalf("command failed....: %s", err)
	//}
	// Print out the output
	fmt.Printf("Output:\n%s", b.String())
	c =b.String()
	d =[]byte(c)
	if _, err = fo.Write(d); err != nil {
		panic(err)
	}
	
	
	// Kill the process if it doesn't exit in time
	defer time.AfterFunc(CommandTimeout, func() {
		log.Printf("command timed out")
		cmd.Process.Kill()
	}).Stop()
}

