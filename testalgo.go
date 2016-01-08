package main
import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	//"strings"
	"io"
	"time"
	"os"
	"math/rand"
)

// Set your timeout
	const CommandTimeout = 5 * time.Second

	var startport = 11100
	var prevport = startport
	var port = startport

// Set data-structures

// var whoami array to store nodeIDs
	var whoamiArray [101]string
	
// var store map(key, string) to store key and value
	var storeMap map[string]string

// var ipport array to store current ipport on which kademlia is running
	var ports [101]int

// var cmd array to array all commands
	var cmd [101]exec.Cmd
	
// var nodeIDs map to store nodeid and key pairs
	var nodeIDMap map[string]string
	
// var array to store output pipe of different instances 
	var stdoutpipes [101]io.ReadCloser
	
// var sdinpipe array to store input pipe of different instances
	var stdinpipes [101]io.WriteCloser

func randInt(min int, max int) int{

	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func setPipes(int i){

	// Set up the inputpipe
  	stdinpipes[i], errin := cmd[i].StdinPipe()
	if err != nil {
    log.Fatalf("failed to create pipe for STDIN: %s for %d", errin, i)
  	}

	// Set up the outputpipe	
	stdoutpipes[i], errout := cmd[i].StdoutPipe()
  	if err != nil {
    log.Fatalf("failed to create pipe for STDOUT: %s for %d", errout, i)
  	}

}
	
func runCommand(string command, int i) result string{
	// Write the input and close
	fmt.Println("Running command %s on %d", command, i)
	
	string result
  	go func() {
		string result
    	defer in.Close()
		fmt.Fprintln(stdinpipes[i], command)
       	result = stdoutpipes[i].Read()
		 
	}()
 
}


/*
  
get_contact ID
    If your buckets contain a node with the given ID,
        printf("%v %v\n", theNode.addr, theNode.port)
    If your buckers do not contain any such node, print "ERR".
  
*/

func main() {

// Set Environment variable

for i:=0; i<100; i++ {

	prevport = port
	port = startport + i
	ports[i] = port
	storeMap = make(map[string]string)
	nodeIDMap = make(map[string]string)
	
	
	// create new command and store in array
   	
	// start kademlia instance by starting command
	cmd[i] := exec.Command("./mainRef","localhost:"+port,"localhost:"+prevport)
   	
	setPipes(i)
	// run commands in the kademlia instance:
	
	// Start the process
  	if err := cmd[i].Start(); err != nil {
    
		log.Fatalf("failed to start instance %d due to :%s", i, err)
  	}
	
	// whoami - store the result in whoami
	command = "whoami"
	res := runcommand(command, i)
	whoamiArray[i] := res
	fmt.Println("nodeId of instance %d is %s", i, res)
	
	// store nodeID and value - store in the store map. for nodeID equal to whoami-10 to whoami+10
	 
	for j:= i-10;j <= i;j++{
		
		if j<0{
			j = 0
		}
		for k:=-5;k< 5; k++{
		
			key = generateKey(j, k )
			command = "store "+ whoamiArray[j]+" "+key+" "+ports[j]
			res = runcommand(command, i)
			if res != "\n"{
				fmt.Println("Error in store rpc")
			}
			else{
				storeMap[key] = ports[j]
				nodeIDMap[key] = whoamiArray[j]
			}
		}
	} 
	 
	//local find value
	for j= i-1;j<=i;j++{
		if j<0{
			j = 0
		}
		
		for k:=-5;k< 5; k++{
		
			key = generateKey(j, k)
			command = "local_find_value "+key
			res = runcommand(command, i)
			if j==i{
				if res != ports[i]{
					fmt.Println("Error local find value")
					break
				}
				continue	 
			}
			if j < i{
				if res != "ERR"{
					fmt.Println("Error in local find value")
					break
				}
				continue
			}
			else{
				break
			}
			
		}
	}
	
	  
	// find node use whoami for the previous instance of kademlia and verify with the value in nodeIDs map. for any 10 nodeid in store
		
	for j= 1;j<=5;j++{
		 
		 	k = randInt(-5, 5)
		 	nodenumber = randInt(0, i)
			key = generateKey(nodenumber, k)
			node = randInt(0,i)
			command = "find_node "+whoami[node]+" "+key
			res = runcommand(command, i)
			verifyFindNodeResult(res, node) 
			
		}
	}
		
	// findvalue for any random 10 value in store and verify with the nodeid map and the whoami map
   	for j= 1;j<=5;j++{
		 
		 	k = randInt(-5, 5)
		 	nodenumber = randInt(0, i)
			key = generateKey(nodenumber, k)
			node = randInt(0,i)
			command = "find_value "+whoami[node]+" "+key
			res = runcommand(command, i)
			if rs = veriftyFindValueResult(res, node); rs != true{
				
				verifyFindNodeResult(res, node) 
			
			}
			
		}
	}
	
	// ping checking
	
	for j= 1;j<=5;j++{
		 
		 	 
		 	 
			node = randInt(0,i)
			command = "ping "+whoami[node]
			res = runcommand(command, i)
			if res != nil{
				fmt.Println("Error in ping with nodeid")
				break
			}
			node = randInt(0,i)
			p = startport + node
			command = "ping localhost:"+string(p) 
			res = runcommand(command, i)
			if res != nil{
				fmt.Println("Error in ping with ip and port")
				break
			} 
		}
	}
	
	
	  
}

for f:=0;f<3;f++{
	
	for i:=0; i<100; i++{
		 
  
  //for random 10 nodeid in store map
  	// iterative findnode - result nodeid must be present in whoami map and one of the result must be equal to the real nodeid storing the value which can be verified using nodeids map
	
	 	
		for j= 1;j<=5;j++{
		 
			node = randInt(0,i)
			command = "iterativeFindNode "+whoami[node]
			res = runcommand(command, i)
			if rs := verifyItFindNodeResult(res, node); rs != true{
			
				fmt.Println("Error in iterative find node :%s", res)
			} 
			
		}
	
		
	// iterative findvalue - for random 10 values in the store and nodeid maps. result can be verifeid using the store map
	// iterativeFindValue key
    //printf("%v %v\n", ID, value), where ID refers to the node that finally
    returned the value. If you do not find a value, print "ERR".
	
	for j= 1;j<=5;j++{
		 
		 	k = randInt(-5, 5)
		 	nodenumber = randInt(0, i)
			key = generateKey(nodenumber, k)
			
			command = "iterativeFindValue "+key
			res = runcommand(command, i)
			if rs = veriftyFindValueResult(res, node); rs != true{
				
				fmt.Println("Error in iterative find value : %s", res) 
				break	
			}
			
		}
	
	
	// iterative store - store 10 random ipport on nodeids, and then check using findvalue using the xor distance method
	  
	//iterativeStore key value
    //Perform the iterativeStore operation and then print the ID of the node that
    //received the final STORE operation.
	
			j = randInt(-10, 10)
			key = generateKey(j, i )
			command = "iterativeStore " "+key+" "+ports[i]
			res = runcommand(command, i)
			if res != "\n"{
				fmt.Println("Error in iterativeStore : %s", res)
				break
			}
			else{
				verifyIterativeStore(res, i)
				 
			}
		}
	} 

 
 //Sample Kademlia - mainSample
 
 	i = 100
 	port = startport+100
 	cmd[100] = exec.Command("./mainSample","localhost:"+port,"localhost:"+prevport)
	
	setPipes(i)
	// run commands in the kademlia instance:
	
	// Start the process
  	if err := cmd[i].Start(); err != nil {
    
		log.Fatalf("failed to start instance %d due to :%s", i, err)
  	}
	
	
	// whoami - store the result in whoami
	command = "whoami"
	res := runcommand(command, i)
	whoamiArray[i] := res
	fmt.Println("nodeId of instance %d is %s", i, res)
	
	// store nodeID and value - store in the store map. for nodeID equal to whoami-10 to whoami+10
	 
	for j:= i-10;j <= i;j++{
		
		if j<0{
			j = 0
		}
		for k:=-5;k< 5; k++{
		
			key = generateKey(j, k )
			command = "store "+ whoamiArray[j]+" "+key+" "+ports[j]
			res = runcommand(command, i)
			if res != "\n"{
				fmt.Println("Error in store rpc")
			}
			else{
				storeMap[key] = ports[j]
				nodeIDMap[key] = whoamiArray[j]
			}
		}
	} 
	 
	//local find value
	for j= i-1;j<=i;j++{
		if j<0{
			j = 0
		}
		
		for k:=-5;k< 5; k++{
		
			key = generateKey(j, k)
			command = "local_find_value "+key
			res = runcommand(command, i)
			if j==i{
				if res != ports[i]{
					fmt.Println("Error local find value")
					break
				}
				continue	 
			}
			if j < i{
				if res != "ERR"{
					fmt.Println("Error in local find value")
					break
				}
				continue
			}
			else{
				break
			}
			
		}
	}
	
	  
	// find node use whoami for the previous instance of kademlia and verify with the value in nodeIDs map. for any 10 nodeid in store
		
	for j= 1;j<=5;j++{
		 
		 	k = randInt(-5, 5)
		 	nodenumber = randInt(0, i)
			key = generateKey(nodenumber, k)
			node = randInt(0,i)
			command = "find_node "+whoami[node]+" "+key
			res = runcommand(command, i)
			verifyFindNodeResult(res, node) 
			
		}
	}
		
	// findvalue for any random 10 value in store and verify with the nodeid map and the whoami map
   	for j= 1;j<=5;j++{
		 
		 	k = randInt(-5, 5)
		 	nodenumber = randInt(0, i)
			key = generateKey(nodenumber, k)
			node = randInt(0,i)
			command = "find_value "+whoami[node]+" "+key
			res = runcommand(command, i)
			if rs = veriftyFindValueResult(res, node); rs != true{
				
				verifyFindNodeResult(res, node) 
			
			}
			
		}
	}
	
	// ping checking
	
	for j= 1;j<=5;j++{
		 
		 	 
		 	 
			node = randInt(0,i)
			command = "ping "+whoami[node]
			res = runcommand(command, i)
			if res != nil{
				fmt.Println("Error in ping with nodeid")
				break
			}
			node = randInt(0,i)
			p = startport + node
			command = "ping localhost:"+string(p) 
			res = runcommand(command, i)
			if res != nil{
				fmt.Println("Error in ping with ip and port")
				break
			} 
		}

for f:=0;f<3;f++{
	
	
		 
  
  //for random 10 nodeid in store map
  	// iterative findnode - result nodeid must be present in whoami map and one of the result must be equal to the real nodeid storing the value which can be verified using nodeids map
	
	 	
		for j= 1;j<=5;j++{
		 
			node = randInt(0,i)
			command = "iterativeFindNode "+whoami[node]
			res = runcommand(command, i)
			if rs := verifyItFindNodeResult(res, node); rs != true{
			
				fmt.Println("Error in iterative find node :%s", res)
			} 
			
		}
	
		
	// iterative findvalue - for random 10 values in the store and nodeid maps. result can be verifeid using the store map
	// iterativeFindValue key
    //printf("%v %v\n", ID, value), where ID refers to the node that finally
    returned the value. If you do not find a value, print "ERR".
	
	for j= 1;j<=5;j++{
		 
		 	k = randInt(-5, 5)
		 	nodenumber = randInt(0, i)
			key = generateKey(nodenumber, k)
			
			command = "iterativeFindValue "+key
			res = runcommand(command, i)
			if rs = veriftyFindValueResult(res, node); rs != true{
				
				fmt.Println("Error in iterative find value : %s", res) 
				break	
			}
			
		}
	
	
	// iterative store - store 10 random ipport on nodeids, and then check using findvalue using the xor distance method
	  
	//iterativeStore key value 
    
    //Perform the iterativeStore operation and then print the ID of the node that
    //received the final STORE operation.
	
			j = randInt(-10, 10)
			key = generateKey(j, i )
			command = "iterativeStore " "+key+" "+ports[i]
			res = runcommand(command, i)
			if res != "\n"{
				fmt.Println("Error in iterativeStore : %s", res)
				break
			}
			else{
				verifyIterativeStore(res, i)
				 
			}		
	} 

 // Stopping all processes
 
 for i=0;i<101;i++{ 
 
 	// Wait for the process to finish
  	if err = cmd[i].Wait(); err != nil {
    log.Fatalf("command failed....: %s", err)
 }
    
// Kill the process if it doesn't exit in time
  	defer time.AfterFunc(CommandTimeout, func() {
    log.Printf("command timed out")
    cmd[i].Process.Kill()
}).Stop()
 
   
}