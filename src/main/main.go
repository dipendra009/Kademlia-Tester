   package main



import (

    "flag"

    "fmt"

    "log"

    "math/rand"

    "net"

    "net/http"

    "net/rpc"

    "time"

    "container/list"

  //  "io"

    "os"

   "bufio"

   "strings"

)



import (

    "kademlia"

)





func main() {

    // By default, Go seeds its RNG with 1. This would cause every program to

    // generate the same sequence of IDs.

    rand.Seed(time.Now().UnixNano())



    // Get the bind and connect connection strings from command-line arguments.

    flag.Parse()

    args := flag.Args()

    if len(args) != 2 {

        log.Fatal("Must be invoked with exactly two arguments!\n")

    }

    listenStr := args[0]

    firstPeerStr := args[1]



    fmt.Printf("kademlia starting up!\n")

    kadem := kademlia.NewKademlia()





	tcpAddr, err := net.ResolveTCPAddr("tcp4", args[0])



	kadem.MyContact.Host = tcpAddr.IP

	kadem.MyContact.Port = uint16(tcpAddr.Port)



    rpc.Register(kadem)

    rpc.HandleHTTP()

    l, err := net.Listen("tcp", listenStr)

    if err != nil {

        log.Fatal("Listen: ", err)

    }



    // Serve forever.

    go http.Serve(l, nil)



    // Confirm our server is up with a PING request and then exit.

    // Your code should loop forever, reading instructions from stdin and

    // printing their results to stdout. See README.txt for more details.

    client, err := rpc.DialHTTP("tcp", firstPeerStr)

    if err != nil {

        log.Fatal("DialHTTP: ", err)

    }

    ping := new(kademlia.Ping)

    ping.MsgID = kademlia.NewRandomID()

    var pong kademlia.Pong

    err = client.Call("Kademlia.Ping", ping, &pong)

    if err != nil {

        log.Fatal("Call: ", err)

    }



    //log.Printf("ping msgID: %s\n", ping.MsgID.AsString())

   // log.Printf("pong msgID: %s\n", pong.MsgID.AsString())







	tcpAdd, err := net.ResolveTCPAddr("tcp4", args[1])



	host := tcpAdd.IP

	port := uint16(tcpAdd.Port)

	//fmt.Println(pong.Sender)



	nodeID := pong.Sender.NodeID



	cont := kademlia.Contact{nodeID, host, port}

	kademlia.Join(kadem, cont)





	for true {

		var command, arg1, arg2, arg3 string

		fmt.Scanf("%s", &command,"\n")

		//fmt.Print("command:", command)



		switch command{



		case "whoami":

		//Print your node ID

		kademlia.WhoAmI(kadem)



		case "local_find_value":

		//If your node has data for the given key, print it.

		//If your node does not have data for the given key, you should print "ERR".

		fmt.Scanf("%s", &arg1)

		key := arg1

		kademlia.LocalFindValue(kadem, key)



		case "get_contact":

		//If your buckets contain a node with the given ID, printf("%v %v\n",theNode.addr,theNode.port)

		//If your buckers do not contain any such node, print "ERR".

		fmt.Scanf("%s", &arg1)

		ID := arg1

		//fmt.Println(ID)

		kademlia.GetContact(kadem,ID)



		case "iterativeStore":

		//Perform the iterativeStore operation and then print the ID of the node that received the final STORE operation.

		fmt.Scanf("%s %s", &arg1, &arg2)

		key := arg1

		keyId, err := kademlia.FromString(key)

		if err != nil {

		}

		value := arg2

		//fmt.Println(key,value)

		storereq := kademlia.StoreRequest{}

		var kclosenodes *list.List

		storereq.Sender = kadem.MyContact

		storereq.MsgID = kademlia.NewRandomID()

		storereq.Key = keyId

		storereq.Value = []byte(value)

		kadem.IterativeStore(storereq, kclosenodes)







		case "iterativeFindNode":

		//Print a list of ≤ k closest nodes and print their IDs. You should collect the IDs in a slice and print that.

		fmt.Scanf("%s", &arg1)

		ID := arg1

		//fmt.Println(ID)

		ID2, err := kademlia.FromString(ID)

		if err != nil {

		}

		//kadem.IterativeFindNode(ID2, false)



		case "iterativeFindValue":

		//printf("%v %v\n"  ID, value), where ID refers to the node that finally returned the value. If you do not find a value, print "ERR".

		fmt.Scanf("%s", &arg1)

		key := arg1

		//fmt.Println(key)

		keyId, err := kademlia.FromString(key)

		if err != nil {

		}
               
                findvalueresult := kademlia.FindValueResult{}

		kadem.IterativeFindValue(keyId, &findvalueresult)



		case "ping":

		fmt.Scanf("%s", &arg1)

		nodeID := arg1

		//fmt.Println(nodeID)

		if strings.ContainsAny(nodeID, ":"){



			client, err := rpc.DialHTTP("tcp", arg1)

    		if err != nil {

        		log.Fatal("DialHTTP: ", err)

    		}

    		ping := new(kademlia.Ping)

    		ping.MsgID = kademlia.NewRandomID()

    		var pong kademlia.Pong

    		err = client.Call("Kademlia.Ping", ping, &pong)

    		if err != nil {

        	log.Fatal("Call: ", err)

    		fmt.Println("Err")

    		continue

    		}



			tcpAdd, err := net.ResolveTCPAddr("tcp4", arg1)



			host := tcpAdd.IP

			port := uint16(tcpAdd.Port)

			//fmt.Println(pong.Sender)



			nodeID := pong.Sender.NodeID



			cont := kademlia.Contact{nodeID, host, port}

			kademlia.Update(kadem, &cont)





 			fmt.Println("Ping success")

			continue

 		}



		isPing := kademlia.CallPings(kadem,nodeID)

		if isPing != true{

			fmt.Println("Error in Ping")

		}



		case "store":

		//Perform a store and print a blank line.

		fmt.Scanf("%s %s", &arg1, &arg2)

		nodeID := arg1

		key := arg2

		r := bufio.NewReader(os.Stdin)

		arg3, err = r.ReadString(10)

		//fmt.Fscanln(r, &arg3)

		//arg3 = bufio.ReadString('\n')

		value := arg3

		//fmt.Println("ID: ", nodeID, "key: ", key, "value", value)

		isStored := kademlia.CallStore(kadem,nodeID, key, value)

		if isStored != true{

			fmt.Println("Error in Store")

		}





		case "find_node":

		//Perform a find_node and print its results as for iterativeFindNode.

		fmt.Scanf("%s %s", &arg1, &arg2)

		nodeID := arg1

		key := arg2

		//fmt.Println(nodeID, key)

		isFound := kademlia.CallFindNode(kadem, nodeID, key)

		if isFound != true{

			fmt.Println("Error in finding")

		}





		case "find_value":

		//Perform a find_value. If it returns nodes, print them as for find_node. If it returns a value, print the value as in iterativeFindValue.



		fmt.Scanf("%s %s", &arg1, &arg2)

		nodeID := arg1

		key := arg2

		//fmt.Println(nodeID, key)

		isFound := kademlia.CallFindValue(kadem,nodeID, key)

		if isFound != true{

			fmt.Println("UnSuccessful ")

		}
		break;
		
		
		case "exit":
			break;



		}

		if command == "exit"{
			break;
		}
		

	}

}





/*



NUID = a8aa31dfd6ae65ae2dae759bd64bf8f3836aa12f



How to check:



./bin/main localhost:8080 natsu.cs.northwestern.edu:7890



whoami



gives YOURID



ping NUID





find_node NUID YOURID



store NUID KEY data



find_value NUID KEY



now do iterative check............





















*/


