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
	"strings"
	//"errors"
	//"builtin"
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

	 
	var startport int	
	var prevport string 
	var port string
 	var bufOut [101](*bufio.Reader)


// To generate a random number

func randInt(min int, max int) int{

	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

// To set input and output pipes for different commands

func setPipes(i int){
	// Set up the inputpipe
  	stdinpipe, errin := cmd[i].StdinPipe()
	if errin != nil {
    log.Fatalf("failed to create pipe for STDIN: %s for %dth instance of kademlia", errin, i)
  	}
	// Set up the outputpipe	
	stdoutpipe, errout := cmd[i].StdoutPipe()
  	if errout != nil {
    log.Fatalf("failed to create pipe for STDOUT: %s for %d instance of kademlia", errout, i)
  	}
  	stdinpipes[i] = stdinpipe
  	stdoutpipes[i] = stdoutpipe
	bufOut[i] = bufio.NewReader(stdoutpipes[i])
}

// To check if a nodeID or Key is valid or not

func validID(ID string) bool{

	if len(ID) == 40{
		//fmt.Println("\nlen problem", len(ID), " ",ID)
		return true
	}
 
	return false
	
}
	
// To run a set of necessary commands on Reference Kademlia instances	
	
func runCommandsRef(i int){
 	 	
		command := "whoami"
 		//defer close the pipes
    	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
        _, err := stdinpipes[i].Write([]byte(command + "\n"))
        if err != nil {
            panic(err)
        }
		defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    //---------Starting------------------------
    
   	var result string
   
	for ;;{
   
		res,_,err := bufOut[i].ReadLine()
		if err != nil {
			//panic(err)
			return
		}
		
		defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
	//	fmt.Println(string(res))
        
		if validID(string(res)){
			//break
			result = string(res)
			break
		}
		
		/*res,_,err = bufOut[i].ReadLine()
		if err != nil {
			//panic(err)
			return
		}
		//fmt.Println(string(res))
        	res,_,err = bufOut[i].ReadLine()
		if err != nil {
			//panic(err)
			return
		}
		
	*/
	}
	
	//---------Whoami--------------------------
	//fmt.Println("whoami: ", result)
	
	
	whm := result
	whoamiArray[i] = whm	
    
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
      
	//-------------store rpc-------------------------
	 
	cm := "store "+whm+" "+whm+" "+whm
	
	//cm := "find_node "+result
	
	//fmt.Println("\n\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
			//panic(err)
			return
		}
		 //fmt.Println("Written command") 

        _,_,err = bufOut[i].ReadLine()
	//fmt.Println("AFter writting")
	if err != nil {
			//panic(err)
			return
		}
		//fmt.Println("AFter writting")
	
	storeMap[whm] = whm
	nodeIDMap[whm] = whm
	//fmt.Println("output: ", string(res))
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
	 defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
 	 
		 
	//-------------ping ipport-------------------------
    
	k := 0
	
	for j:=0;j<=i;j=j+k+1{
	
	port = strconv.Itoa(startport + j)
	//fmt.Println("pinging...", j, " from  ",i)
		
	
    cm = "ping localhost:"+port
    //fmt.Println("\n\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            //panic(err)
			return
        }
        
	//fmt.Println("Written")

        _,_,er := bufOut[i].ReadLine()
	if er != nil {
		//panic(er)
		continue
	}
	//fmt.Println("output: ", string(res))
	k =k+1
	//fmt.Printf("\nSuccessful ping from %d to %d", i, j)
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    }

}


// To test a Sample Kademlia instance using Reference Kademlia instances

func runCommandsSampleRef(i int){
 	  	
		command := "whoami"
	 	success := 1
 		//close the pipes
    	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
        _, err := stdinpipes[i].Write([]byte(command + "\n"))
        if err != nil {
            panic(err)
        }
		defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
   
	
	 
	   //---------Starting------------------------
    
   	var result string
   
    fmt.Println("\nRunning basic commands on Sample Kademlia using Reference Kademlia instances\n" )
    for ;;{
   
		res,_,err := bufOut[i].ReadLine()
		if err != nil {
			//panic(err)
			fmt.Println("Error in starting up the sample Kademlia to be tested (may be due to whoami not displaying nodeID as string)")
			return
		}
		
		defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
	//	fmt.Println(string(res))
        
		if validID(string(res)){
			//break
			result = string(res)
			break
		}
		 
	}
	//---------Whoami--------------------------
	fmt.Println("whoami: ", result)
	
	fmt.Println("whoami command is working properly")
	
	whm := result
	whoamiArray[i] = whm	
    
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    
    //choosing a random reference kademlia
	
	 
	
	//-------------ping ipport-------------------------
    
	k := 0
	
	for j:=1;j<=i;j=j+1{
	
	port = strconv.Itoa(startport + j)
	//fmt.Println("pinging...", j, " from  ",i)
		
	
    cm := "ping localhost:"+port
    //fmt.Println("\n\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            //panic(err)
			return
        }
        
	//fmt.Println("Written")

        _,_,er := bufOut[i].ReadLine()
	if er != nil {
		//panic(er)
		continue
	}
	//fmt.Println("output: ", string(res))
	k =k+1
	//fmt.Printf("\nSuccessful ping from sample to %dth instance of reference Kademlia", j)
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    }
    
	fmt.Println("\nSelecting a random Reference Kademlia instance for testing ", )

	// Ping a random ID before running command

	for j:=1;j<=i;j=j+1{
	k = randInt(1, i-1)
	port = strconv.Itoa(startport + k)
	//fmt.Println("pinging...", k, " from  ",i)
		
	
    cm := "ping localhost:"+port
    //fmt.Println("\n\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            //panic(err)
			return
        }
        
	//fmt.Println("Written")

        _,_,er := bufOut[i].ReadLine()
	if er != nil {
		//panic(er)
		continue
	}
	//fmt.Println("output: ", string(res))
	
	fmt.Printf("Successful ping from Sample Kademlia to %dth instance of Reference Kademlia\n", k)
	break
	k =k+1
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    }
    
	
	//---------get_contact ID------------------
	
	
	
	cm := "get_contact "+whoamiArray[k]
	
	//cm := "find_node "+result+ " "+ result
	
	fmt.Printf("\nGet Contact address for %dth instance of Reference Kademlia using: ", k)
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
           //fmt.Println("continueing.......")
           // panic(err)
        }
       //fmt.Println("Written command") 

        res,_,err := bufOut[i].ReadLine()
	//fmt.Println("AFter writting")
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
//	fmt.Println("AFter writting")
	
	fmt.Println("output: ", string(res))
	if string(res)== whoamiArray[k] || strings.Contains(string(res), strconv.Itoa(startport + k)){
		fmt.Println("get_contact is working properly")
	} else{
		fmt.Println("Error in get_contact")
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
	
	 //-------------ping nodeID-------------------------
    
    cm = "ping "+whoamiArray[k]
    
	fmt.Printf("\nPing %dth instance of Reference Kademlia using: ", k)
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
         fmt.Println("Error in executing command")
            //panic(err)
       }
        
	//fmt.Println("Written")

       res,_,err = bufOut[i].ReadLine()
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
	fmt.Println("output: ", string(res))
	if string(res)== "SUCCESS"{
		fmt.Println("ping nodeID is working properly")
	}else{
		fmt.Println("Error in ping nodeID")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    //-------------ping ipport-------------------------
    
	 
	port = strconv.Itoa(startport + k)
	
	//fmt.Println("pinging...", j, " from  ",i)
		
	fmt.Printf("\nPing %dth instance of Reference Kademlia using: ", k)
	
    cm = "ping localhost:"+port
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            //panic(err)
			fmt.Println("Error in executing command")
        
			//return
        }
        
	//fmt.Println("Written")

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
		//panic(err)
		
	}
	fmt.Println("output: ", string(res))
	
	if string(res)== "SUCCESS"{
		fmt.Println("ping ipport is working properly")
	}else{
		fmt.Println("Error in ping ipport")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    
    
    
    //-------------local_find_value-------------------------
    
    cm = "local_find_value "+whm
    fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            fmt.Println("Error in executing command")
            //panic(err)
       }
        
	//fmt.Println("Written")

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
	fmt.Println("output: ", string(res))
	if string(res)== "ERR"{
		fmt.Println("local_find-value is working properly")
	}else{
		fmt.Println("Error in local_find_value")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
	
	//-------------store rpc-------------------------
	
	
	
	cm = "store "+whoamiArray[k]+" "+whoamiArray[k-1]+" "+"SUCCESS"
	
	//cm := "find_node "+result
	
	fmt.Printf("\nCalling store rpc to %dth instance of Reference Kademlia using: ", k)
	
	fmt.Println("\ncm: ", cm)
	
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
           fmt.Println("Error in executing command")
            //panic(err)
       }
       //fmt.Println("Written command") 

        res,_,err = bufOut[i].ReadLine()
	//fmt.Println("AFter writting")
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
	//fmt.Println("AFter writting")
	
	fmt.Println("output: ", string(res))
	if string(res)== "SUCCESS"{
		fmt.Println("store rpc is working properly")
		}else{
		fmt.Println("Error in store rpc")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
		
	/*
	//------------------------------------------------------------------	
	cm = "store "+whoamiArray[i]+" "+whoamiArray[i]+" "+"SUCCESSS"
	
	//cm := "find_node "+result
	setPipes(k)
	
	fmt.Printf("\nCalling store rpc to Sample Kademlia using Reference Kademlia %d using: '", k)
	
	fmt.Println("\ncm: ", cm)
	
	_, err = stdinpipes[k].Write([]byte(cm+ "\n"))
        if err != nil {
           fmt.Println("Error in executing command")
            //panic(err)
       }
       //fmt.Println("Written command") 

        res,_,err = bufOut[k].ReadLine()
	//fmt.Println("AFter writting")
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
	//fmt.Println("AFter writting")
	
	fmt.Println("output: ", string(res))
	if string(res)== "SUCCESS"{
		fmt.Println("store rpc is working properly")
		storeMap[whm] = "SUCCESSS"
		nodeIDMap[whm] = whm
	}else{
		fmt.Println("Error in store rpc")
	
		success = 0
	}
	
	defer stdinpipes[k].Close()
    	defer stdoutpipes[k].Close()
    	
	*/	
    
	//------------find_value--------------------------
	
	cm = "find_value "+whoamiArray[k]+ " "+ whoamiArray[k-1]
	
	fmt.Printf("\nCalling find_value of Sample Kademlia for %dth instance of Reference Kademlia:", k)
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            fmt.Println("Error in executing command")
            //panic(err)
       }
        

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
	fmt.Println("output: ", string(res))
	if string(res)== "SUCCESS"{
		fmt.Println("find_value is working properly")
	}else{
		fmt.Println("Error in find_value")
	
		success = 0
	}
	
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
            res,_,err = bufOut[i].ReadLine()
	if err != nil {
		panic(err)
	}
	//fmt.Println("output: ", string(res))
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
 
 	
	//--------------------------------------------------
    /*
     
	cm = "find_value "+whm+ " "+ whm
	
	fmt.Printf("Calling find_value of Sample Kademlia using %d Reference Kademlia:", k)
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[k].Write([]byte(cm+ "\n"))
        if err != nil {
            fmt.Println("Error in executing command")
            //panic(err)
       }
        

        res,_,err = bufOut[k].ReadLine()
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
	fmt.Println("output: ", string(res))
	if string(res)== storeMap[whm]{
		fmt.Println("find_value is working properly")
	}else{
		fmt.Println("Error in find_value")
	
		success = 0
	}
	
	
	defer stdinpipes[k].Close()
    	defer stdoutpipes[k].Close()
            res,_,err = bufOut[i].ReadLine()
	if err != nil {
		panic(err)
	}
	//fmt.Println("output: ", string(res))
	defer stdinpipes[k].Close()
    	defer stdoutpipes[k].Close()
 
    */	
			
					
	//----------------find_node--------------------
	cm = "find_node "+whoamiArray[k]+ " "+ whm
	
	fmt.Printf("\nUsing Sample Kademlia to call find_node of %dth instance of Reference Kademlia:",k)
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
			fmt.Println("Error in executing command")
            //panic(err)
        }
        

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
	fmt.Println("Error in executing command")
            //panic(err)
       }
	
	fmt.Println("output: ", string(res))
	if strings.Contains(string(res), nodeIDMap[whm]){
		fmt.Println("find_node is working properly")
	}else{
		fmt.Println("Error in find_node")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
 	
	 
	
	if success == 1{
    
    
    fmt.Printf("\nSuccessfully run all the basic command tests on Sample Kademlia using Reference Kademlia instances\n")
    }else{
		fmt.Println("\n Failed to execute all the basic command tests using Reference Kademlia instances")
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
 	
	fmt.Println("\nRunning iterative commands on instance of Sample Kademlia using the Reference Kademlia instances" )
    
	
	
	//----------------iterativeFindNode--------------------
	
	cm = "iterativeFindNode "+ whoamiArray[k]
	fmt.Printf("\nCalling iterativeFindNode to find closest nodes to %d instance of Reference Kademlia: ", k)
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
           //panic(err)
		   fmt.Println("Error in writing command...")
			//return
        }
        
        res,_,err = bufOut[i].ReadLine()
	if err != nil {
		fmt.Println("Can not read the iterativeFindNode command output.")
		//panic(err)
				
				//return
	}
	if err.Error() == "EOF"{
			fmt.Println("Nothing to read...no output")
			success = 0
	}else{
	fmt.Println("output: ", string(res))
	if strings.Contains(string(res), nodeIDMap[whoamiArray[k]]){
		fmt.Println("iterativeFindNode is working properly")
	}else{
		fmt.Println("Error in iterativeFind_node")
	
		success = 0
	}
	}
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()

         
	//------------iterativeFindValue--------------------------
	
	cm = "iterativeFindValue "+whoamiArray[k-1]
	fmt.Printf("\nCalling iterativeFindValue to find value on %dth or %dth instance of Reference Kademlia",k, k-1)
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            //panic(err)
			fmt.Println("Error in writing command...")
			//return
        }
        

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
		//panic(err)
		fmt.Println("Can not read the iterativeFindValue command output.")
		
			//return

	}
	if err.Error() == "EOF"{
			fmt.Println("Nothing to read...no output")
			success = 0
	}else{
	
	fmt.Println("output: ", string(res))
	if strings.Contains(string(res),"SUCCESS"){
		fmt.Println("iterativeFindValue is working properly")
	}else{
		fmt.Println("Error in iterativeFindValue")
	
		success = 0
	}
	}
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
      
    	      
    //---------iterativeStore------------------
	
	cm = "iterativeStore "+whoamiArray[k]+" "+ "SUCCESS SUCCESS"
	fmt.Println("Using iterativeStore to store data on nearest node to %dth instance of Reference Kademlia", k )
	//cm := "find_node "+result+ " "+ result
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
           //fmt.Println("continueing.......")
           // panic(err)
		   fmt.Println("Error in writing command...")
        }
       //fmt.Println("Written command") 

    defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
		    res,_,err = bufOut[i].ReadLine()
	//fmt.Println("AFter writting")
	if err != nil {
	fmt.Println("Can not read the iterativeStore command output.")
		
		
		//panic(err)
	}
//	fmt.Println("AFter writting")
	if err.Error() == "EOF"{
			fmt.Println("Nothing to read...no output")
			success = 0
	}else{
	
	fmt.Println("output: ", string(res))
	if validID(string(res)){
		fmt.Println("iterativeStore rpc is working properly")
		storeMap[whm] = "SUCCESS SUCCESS"
	}else{
		fmt.Println("Error in iterativeStore")
	
		success = 0
	}
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
   
	    
	 if success == 1{
    
    fmt.Printf("\nSuccessfully run all iterative commands on Sample Kademlia using reference Kademlia instances\n")
    }else{
		fmt.Println("\nFailed to execute all iterative commands using reference Kademlia instances\n")
	}
	
    
}


func runCommandsSampleSample(i int){
 	 command := "whoami"
	 success := 1
 	//close the pipes
    	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
        _, err := stdinpipes[i].Write([]byte(command + "\n"))
        if err != nil {
            panic(err)
        }
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    //---------Starting------------------------
    
   	var result string
   
    fmt.Println("\nRunning basic commands on instance of Sample Kademlia using single instance without any Reference Kademlia instances\n" )
    for ;;{
   
		res,_,err := bufOut[i].ReadLine()
		if err != nil {
			//panic(err)
			fmt.Println("Error in starting up the Sample Kademlia to be tested (may be due to whoami not displaying nodeID as string)")
			return
		}
		
		defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
	//	fmt.Println(string(res))
        
		if validID(string(res)){
			//break
			result = string(res)
			break
		}
		 
	}
	//---------Whoami--------------------------
	fmt.Println("whoami: ", result)
	
	fmt.Println("whoami command is working properly")
	
	whm := result
	whoamiArray[i] = whm	
    
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    
    //---------get_contact ID------------------
	
	cm := "get_contact "+whm
	
	//cm := "find_node "+result+ " "+ result
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
           fmt.Println("continueing.......")
           // panic(err)
        }
       //fmt.Println("Written command") 

        res,_,err := bufOut[i].ReadLine()
	//fmt.Println("AFter writting")
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
//	fmt.Println("AFter writting")
	
	fmt.Println("output: ", string(res))
	if string(res)== whm || strings.Contains(string(res), strconv.Itoa(startport + i)){
		fmt.Println("get_contact is working properly")
	} else{
		fmt.Println("Error in get_contact")
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
	
	 //-------------ping nodeID-------------------------
    
    cm = "ping "+whm
    fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
         fmt.Println("Error in executing command")
            //panic(err)
       }
        
	//fmt.Println("Written")

       res,_,err = bufOut[i].ReadLine()
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
	fmt.Println("output: ", string(res))
	if string(res)== "SUCCESS"{
		fmt.Println("ping nodeID is working properly")
	}else{
		fmt.Println("Error in ping nodeID")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    //-------------ping ipport-------------------------
    
	 
	port = strconv.Itoa(startport + i)
	//fmt.Println("pinging...", j, " from  ",i)
		
	
    cm = "ping localhost:"+port
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            //panic(err)
			fmt.Println("Error in executing command")
        
			//return
        }
        
	//fmt.Println("Written")

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
		//panic(err)
		
	}
	fmt.Println("output: ", string(res))
	
	if string(res)== "SUCCESS"{
		fmt.Println("ping ipport is working properly")
	}else{
		fmt.Println("Error in ping ipport")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
    
    
    
    //-------------local_find_value-------------------------
    
    cm = "local_find_value "+whm
    fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            fmt.Println("Error in executing command")
            //panic(err)
       }
        
	//fmt.Println("Written")

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
	fmt.Println("output: ", string(res))
	if string(res)== "ERR"{
		fmt.Println("local_find-value is working properly")
	}else{
		fmt.Println("Error in local_find_value")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
	
	//-------------store rpc-------------------------
	
	
	
	cm = "store "+whm+" "+whm+" "+"SUCCESS"
	
	//cm := "find_node "+result
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
           fmt.Println("Error in executing command")
            //panic(err)
       }
       //fmt.Println("Written command") 

        res,_,err = bufOut[i].ReadLine()
	//fmt.Println("AFter writting")
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
	//fmt.Println("AFter writting")
	
	fmt.Println("output: ", string(res))
	if string(res)== "SUCCESS"{
		fmt.Println("store rpc is working properly")
		storeMap[whm] = "SUCCESS"
		nodeIDMap[whm] = whm
	}else{
		fmt.Println("Error in store rpc")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
	//------------find_value--------------------------
	
	cm = "find_value "+whm+ " "+ whm
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            fmt.Println("Error in executing command")
            //panic(err)
       }
        

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
		fmt.Println("Error in executing command")
            //panic(err)
       }
	fmt.Println("output: ", string(res))
	if string(res)== storeMap[whm]{
		fmt.Println("find_value is working properly")
	}else{
		fmt.Println("Error in find_value")
	
		success = 0
	}
	
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
            res,_,err = bufOut[i].ReadLine()
	if err != nil {
		panic(err)
	}
	//fmt.Println("output: ", string(res))
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
 
    
     	
	//----------------find_node--------------------
	cm = "find_node "+whm+ " "+ whm
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
			fmt.Println("Error in executing command")
            //panic(err)
        }
        

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
	fmt.Println("Error in executing command")
            //panic(err)
       }
	
	fmt.Println("output: ", string(res))
	if strings.Contains(string(res), nodeIDMap[whm]){
		fmt.Println("find_node is working properly")
	}else{
		fmt.Println("Error in find_node")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
 	
	
	if success == 1{
    
    
    fmt.Printf("\nSuccessfully run all basic commands on Sample Kademlia using single kademlia instance\n")
    }else{
		fmt.Println("\n Failed to pass all basic commands test using same instance\n")
	}
	 
	fmt.Println("\nRunning iterative commands on instance of Sample Kademlia using only this instance" )
    
	
	//----------------iterativeFindNode--------------------
	cm = "iterativeFindNode "+ whm
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
           //panic(err)
			//return
        }
        

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
		//panic(err)
			//return
	}
	fmt.Println("output: ", string(res))
	if strings.Contains(string(res), nodeIDMap[whm]){
		fmt.Println("iterativeFindNode is working properly")
	}else{
		fmt.Println("Error in iterativeFind_node")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()

         
	//------------iterativeFindValue--------------------------
	
	cm = "iterativeFindValue "+whm
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
            //panic(err)
			//return
        }
        

        res,_,err = bufOut[i].ReadLine()
	if err != nil {
		//panic(err)
			//return

	}
	fmt.Println("output: ", string(res))
	if strings.Contains(string(res),storeMap[whm]){
		fmt.Println("iterativeFindValue is working properly")
	}else{
		fmt.Println("Error in iterativeFindValue")
	
		success = 0
	}
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
      
    	      
    //---------iterativeStore------------------
	
	cm = "iterativeStore "+whm+" "+ "SUCCESS"
	
	//cm := "find_node "+result+ " "+ result
	
	fmt.Println("\ncm: ", cm)
	_, err = stdinpipes[i].Write([]byte(cm+ "\n"))
        if err != nil {
           fmt.Println("continueing.......")
           // panic(err)
        }
       //fmt.Println("Written command") 

    defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
    
		    res,_,err = bufOut[i].ReadLine()
	//fmt.Println("AFter writting")
	if err != nil {
		
		//panic(err)
	}
//	fmt.Println("AFter writting")
	
	fmt.Println("output: ", string(res))
	if validID(string(res)){
		fmt.Println("iterativeStore rpc is working properly")
		storeMap[whm] = "SUCCESS SUCCESS"
	}else{
		fmt.Println("Error in iterativeStore")
	
		success = 0
	}
	
	
	defer stdinpipes[i].Close()
    	defer stdoutpipes[i].Close()
     
	 if success == 1{
    
    fmt.Printf("\nSuccessfully run all iterative commands on Sample Kademlia using single instance\n")
    }else{
		fmt.Println("\nFailed to execute all iterative commands using same instance of Sample Kademlia")
	}
	

     
    
}


func main(){

	// Set Environment variable and compile mainRef and mainSample executable files	
	// Initialize variables and data structures
	startport = randInt(2222, 6000)
	storeMap = make(map[string]string)
	nodeIDMap = make(map[string]string)
	
	prevport = strconv.Itoa(startport)
	port = strconv.Itoa(startport)
	
	//var ref, sample string
	//var totalNum int
	
	fmt.Println("\nTesting kademliaSample against multiple instances of kademliaRef")
	//fmt.Println("Enter name of Reference Kademlia and Sample Kademlia executable files:")
	//fmt.Scanf("%s %s", &ref, &sample)
	
	//fmt.Println("Enter number of reference instances for testing: ")
	//fmt.Scanf("%d", &totalNum)
	
	// First loop to start instance of SampleKademlia and run basic commands and RPCs
	for i:= 0; i<1; i=i+1{
		
		prevport = port
		port = strconv.Itoa(startport + i)
		add1 := "localhost:"+port
		
 		ports[i] = port
	
		// create new command and store in array
   		cmd[i] = exec.Command("./kademliaSample",add1,add1)
	
		// set input and output pipes
		setPipes(i)
		
		// start kademlia instance by starting command
	 	if err := cmd[i].Start(); err != nil {
			log.Fatalf("failed to start instance %dth (sample kademlia)due to :%s", i, err)
			continue
		}
   		// run commands on i-th instance of Kademlia
	 	
 		runCommandsSampleSample(i)
 	}
	
	fmt.Println("\nStarting multiple instances of reference Kademlia for further testing")
	  
	// First loop to start instances of Reference Kademlia and run some commands and RPCs
	for i:= 1; i<50; i=i+1{
		
		
		prevport = port
		port = strconv.Itoa(startport + i)
		add1 := "localhost:"+port
		//add2 := "localhost:"+prevport
		ports[i] = port
		//fmt.Println("add: ", add1, add2)
		// create new command and store in array
   		cmd[i] = exec.Command("./kademliaRef",add1,add1) //"localhost:1234","localhost:1234")
		// set input and output pipes
		setPipes(i)
		// start kademlia instance by starting command
	 	if err := cmd[i].Start(); err != nil {
			log.Fatalf("failed to start instance %dth due to :%s", i, err)
			continue
  		}
  		
  		//fmt.Printf("Started %dth instance of Reference Kademlia...\n", i)
  		
		//go DeferClose(i)
		
		
  		// run commands on i-th instance of Kademlia
		
		runCommandsRef(i)
		 

	}
	 
			
	fmt.Println("\nTesting the sample Kademlia using the existing Reference Kademlia instances")
				
	// First loop to start instance of SampleKademlia and run basic commands and RPCs
	for i:= 50; i<51; i=i+1{
		
		prevport = port
		port = strconv.Itoa(startport + i)
		
		add1 := "localhost:"+port
//		add2 := "localhost:"+prevport
		
		ports[i] = port
	
		// create new command and store in array
   		cmd[i] = exec.Command("./kademliaSample",add1,add1)
	
		// set input and output pipes
		setPipes(i)
		
		// start kademlia instance by starting command
	 	if err := cmd[i].Start(); err != nil {
			log.Fatalf("failed to start instance %dth (sample kademlia)due to :%s", i, err)
			continue
		}
  		//fmt.Printf("\n\nStarted %d instance of Kademlia", i)
  		// run commands on i-th instance of Kademlia
	 	
		 
		runCommandsSampleRef(i)
	}
	
	fmt.Println("Stopping all the running instances of Kademlia")
	
	// Stopping all instances of Kademlia
	for i:=0; i<51; i++{
 		stdinpipes[i].Close()
 		stdoutpipes[i].Close()
 		//fmt.Printf("\nStopping %dth instance of Kademlia", i)
 		//Wait for the process to finish
  		//if err := cmd[i].Wait(); err != nil {
    	//        log.Fatalf(string(i), ":command failed....: %s", err)
 		//}
		// Kill the process if it doesn't exit in time
 		//fmt.Println("After waiting")
		defer time.AfterFunc(CommandTimeout, func() {
   		fmt.Println("After timeout")
		log.Printf("command timed out")
   		cmd[i].Process.Kill()
		}).Stop()
	}

	fmt.Println("\n\nThank you for using this Kademlia Tester\n")
}


 



