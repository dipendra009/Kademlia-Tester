package main
import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	//"strings"
	"io"
	"time"
	//"os"
	"math/rand"
	"strconv"
)

// Set your timeout
	const CommandTimeout = 5 * time.Second

	
// Set data-structures

	// var whoami array to store nodeIDs
	var whoamiArray [101]string
	
	// var store map(key, string) to store key and value
	var storeMap map[string]string

	// var ipport array to store current ipport on which kademlia is running
	var ports [101]string

	// var cmd array to array all commands
	var cmd [101]*exec.Cmd
	
	// var nodeIDs map to store nodeid and key pairs
	var nodeIDMap map[string]string
	
	// var array to store output pipe of different instances 
	var stdoutpipes [101]io.ReadCloser
	
	// var sdinpipe array to store input pipe of different instances
	var stdinpipes [101]io.WriteCloser

	// var outputfiles array to store output of different instances
	//var outputfiles [101]string
	
	// var output buffer array to read from stdoutpipe
	var bufout [101](*bufio.Reader)

	var startport = 11100
	var prevport = strconv.Itoa(startport)
	var port = strconv.Itoa(startport)
	var outfile string

func randInt(min int, max int) int{

	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func setPipes(i int){

	// Set up the inputpipe
	var err error
  	stdinpipes[i], err = cmd[i].StdinPipe()
	if err != nil {
    	log.Fatalf("failed to create pipe for STDIN: %s for %dth instance of kademlia", err, i)
  		panic(err)
  	}
  	defer stdinpipes[i].Close()
  	
	// Set up the outputpipe	
	stdoutpipes[i], err = cmd[i].StdoutPipe()
  	if err != nil {
    	log.Fatalf("failed to create pipe for STDOUT: %s for %d instance of kademlia", err, i)
  		panic(err)
  	}
  	defer stdoutpipes[i].Close()
  	
  	
  	//stdinpipes[i] = stdinpipe
  	//stdoutpipes[i] = stdoutpipe
  	 bufout[i] = bufio.NewReader(stdoutpipes[i])
    
    
    

}
	
func runCommand(command string, i int) (result string){
	// Write the input and close
	 	
	  
		

        res,_,_ := bufout[i].ReadLine()
		fmt.Printf("\nstring is %s %d %s\n",command, i ,string(res))
		fmt.Println(stdinpipes[i])
		fmt.Println( stdoutpipes[i])
		fmt.Println(bufout[i])
        //if er != nil {
          //  panic("here...er")
        //}

	 //	bufout[i] = bufio.NewReader(stdoutpipes[i])  
        _, err1 := stdinpipes[i].Write([]byte(command + "\n"))
        if err1 != nil {
            panic(err1)
        }

		defer stdinpipes[i].Close()
		
        res,_,err := bufout[i].ReadLine()
		if err != nil {
            panic("now here.....err")
        }
        fmt.Printf("\nstring is %s",string(res))
        

res,_,err = bufout[i].ReadLine()
		if err != nil {
            panic("now here.....err")
        }
        fmt.Printf("\nstring is %s",string(res))
        
        res,_,err = bufout[i].ReadLine()
		if err != nil {
            panic("now here.....err")
        }
        fmt.Printf("\nstring is %s",string(res))
        
        res,_,err = bufout[i].ReadLine()
		if err != nil {
            panic("now here.....err")
        }
        fmt.Printf("\nstring is %s",string(res))
        


    result = string(res)    
    return result
 
}
 
func main(){

	// Set Environment variable and compile mainRef and mainSample executable files
	
	// Initialize variables and data structures

	storeMap = make(map[string]string)
	nodeIDMap = make(map[string]string)

	// First loop to start instances of Kademlia and run basic commands and RPCs
	for i:= 0; i<1; i=i+1{
		
		prevport = port
		port = strconv.Itoa(startport + i)
		outfile = port+".txt"
		
		
		add1 := "localhost:"+port
		add2 := "localhost:"+prevport
		
		ports[i] = port
	
		// create new command and store in array
   		cmd[i] = exec.Command("./mainRef",add1,add2)
	
		// set input and output pipes
		setPipes(i)
		
		// start kademlia instance by starting command
	 	if err := cmd[i].Start(); err != nil {
    
			log.Fatalf("failed to start instance %dth due to :%s", i, err)
			panic(err)
			continue
  		}
  		
  		
  		
  		fmt.Printf("\nStarted %dth instance of Kademlia", i)

		

        //res1,_,err := bufout[i].ReadLine()
		//fmt.Printf("\nstring is %s\n",string(res1))
        //if err != nil {
          //  panic(err)
        //}


  		//fmt.Printf("\nStarted %dth instance of Kademlia", i)


  		
  		// run commands on i-th instance of Kademlia
		
		// whoami - store the result in whoami
		command := "whoami"
		res := runCommand(command, i)
		whoamiArray[i] = res
		fmt.Printf("\nnodeId of instance %d is %s", i, res)
	
		
	}
	
	// First loop to start instance of SampleKademlia and run basic commands and RPCs
	for i:= 100; i<101; i=i+1{
		
		prevport = port
		port = strconv.Itoa(startport + i)
		
		add1 := "localhost:"+port
		add2 := "localhost:"+prevport
		
		ports[i] = port
	
		// create new command and store in array
   		cmd[i] = exec.Command("./mainSample",add1,add2)
	
		// set input and output pipes
		setPipes(i)
		
		// start kademlia instance by starting command
	 	if err := cmd[i].Start(); err != nil {
    
			log.Fatalf("failed to start instance %dth (sample kademlia)due to :%s", i, err)
			continue
  		}
  		
  		fmt.Printf("\n\nStarted %d instance of Kademlia", i)
  		
  		// run commands on i-th instance of Kademlia
		
		// whoami - store the result in whoami
		command := "whoami"
		res := runCommand(command, i)
		whoamiArray[i] = res
		fmt.Printf("\nnodeId of instance %dth is %s", i, res)
		command = "exit"
		res = runCommand(command,i)
		
	}
	
	// Stopping all instances of Kademlia
	
	for i:=0; i<101; i++{ 
 
 		stdinpipes[i].Close()
 		stdoutpipes[i].Close()
 	
 		fmt.Printf("\nStopping %dth instance of Kademlia", i)
 		//Wait for the process to finish
  		//if err := cmd[i].Wait(); err != nil {
    	//log.Fatalf("command failed....: %s", err)
 		//}
    
		// Kill the process if it doesn't exit in time
  		defer time.AfterFunc(CommandTimeout, func() {
    		log.Printf("command timed out")
    		cmd[i].Process.Kill()
		}).Stop()
 
   
	}
	
	
}
