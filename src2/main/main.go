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
    "bufio"
    "os"
    "strings"
)

import (
    "kademlia"
)

const delim = '\n'
const REFPEER = "natsu.cs.northwestern.edu:7890"
func main() {
	// By default, Go seeds its RNG with 1. This would cause every program to
	// generate the same sequence of IDs.
	//	rand.seed()
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
	//------------------------------------
	//	fmt.Println(kadem.NodeID)
	//	myNodeID := kadem.WhoAmI()
	//	fmt.Println(myNodeID)
	//------------------------------------





	//-----initializing MyContact---------
	addr, err := net.ResolveTCPAddr("tcp4", listenStr)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	host1 := addr.IP
	port1 := uint16(addr.Port)
	kadem.MyContact.NodeID =  kademlia.WhoAmI(kadem)
	kadem.MyContact.Host   = host1
	kadem.MyContact.Port   = port1
	//puttin Mycontact in the bucket 
	kadem.KademBuckets.Update_Contact(kadem.MyContact,kadem.MyContact)
	//-----registering kadem object and it's functionalities----------
	rpc.Register(kadem)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", listenStr) //telling the server is ON
	if err != nil {
		log.Fatal("Listen: ", err)
	}
	//---------------------------------------------------------------
	// Serve forever.
	go http.Serve(l, nil)
	//---------------------------------------------------------------
	// Confirm our server is up with a PING request and then exit.
	// Your code should loop forever, reading instructions from stdin and
	// printing their results to stdout. See README.txt for more details.
	client, err := rpc.DialHTTP("tcp", firstPeerStr)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
	}
	addr1, err := net.ResolveTCPAddr("tcp4", firstPeerStr)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	host2 := addr1.IP


	ping := new(kademlia.Ping)  //creating a ping object of the ping structure in the kademlia package

	//	net.LookupHost(127.0.0.1)
	//    a, err := net.LookupIP("www.google.com")
	ping.MsgID = kademlia.NewRandomID()
	ping.Sender = kadem.MyContact
	//--------------------------------------------------------------------
	var pong kademlia.Pong  //declaring a pong object of the pong structure in the kademlia package, but not initializing it
	//    calling Kademlia.Ping method passing ping as request object and pong as response object
	err = client.Call("Kademlia.Ping", ping, &pong)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	//------------------------------------------------------------------------
//	log.Printf("ping msgID: %s\n", ping.MsgID.AsString())
//	log.Printf("pong msgID: %s\n", pong.MsgID.AsString())
	fmt.Println("connecting server NodeID:",pong.Sender.NodeID.AsString())
//	fmt.Println("connecting server host:",pong.Sender.Host)
//	fmt.Println("connecting server Port:",pong.Sender.Port)
	pong.Sender.Host = host2
	kadem.KademBuckets.Update_Contact(pong.Sender,kadem.MyContact)

//--------------------------------------------------------------
/*	//---taking a string and getting the key for that string---
	a := "my name is rishabh"
	fmt.Println(a)
	h := sha1.New()
	io.WriteString(h,a)

	b:=h.Sum(nil)
	fmt.Println(b)
	c := hex.EncodeToString(b)
	f, err := kademlia.FromString(c)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	fmt.Println(f)
*/
//--------------------------------------------------------------
/*
	a := "my name is rishabh"
	f :=kademlia.StringToKey(a)
//	fmt.Println(f)

//--------------------------------------------------------------
	sreq := new(kademlia.StoreRequest)
	sreq.MsgID = kademlia.NewRandomID()
	sreq.Sender = kadem.MyContact
	sreq.Key = f
	sreq.Value = []byte(a)
	var sres kademlia.StoreResult
	fmt.Println(".....ii'm here......1")
	err = client.Call("Kademlia.Store", sreq, &sres)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	fmt.Println(".....ii'm here......2")

	//------------------------------------------------------------------------
	log.Printf("ping msgID: %s\n", sreq.MsgID.AsString())
	log.Printf("pong msgID: %s\n", sres.MsgID.AsString())

	fmt.Println(".....ii'm here......3")
	err = client.Call("Kademlia.Store", sreq, &sres)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	fmt.Println(".....ii'm here......4")

	//------------------------------------------------------------------------
	log.Printf("ping msgID: %s\n", sreq.MsgID.AsString())
	log.Printf("pong msgID: %s\n", sres.MsgID.AsString())

*/
/*
    a1 := new(kademlia.StoreRequest)  //creating a ping object of the 
    b1 := new(kademlia.StoreResult)  //creating a ping object of the 
	//---------------------------------------------------
    a2 := new(kademlia.StoreRequest)  //creating a ping object of the 
    b2 := new(kademlia.StoreResult)  //creating a ping object of the 
	//----------------------------------------------------
    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
        log.Fatal("Call: ", err)
    }
*/



//---------find value-------------------------



/*

	fvreq := new(kademlia.FindValueRequest)
	fvreq.MsgID = kademlia.NewRandomID()
	fvreq.Sender = kadem.MyContact
	fvreq.Key = f
	var fvres kademlia.FindValueResult
	fmt.Println(".....ii'm here......1")
	err = client.Call("Kademlia.FindValue", fvreq, &fvres)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	fmt.Println(".....ii'm here......2")


	fmt.Println(string(fvres.Value))

*/
//--------------------------------------------------------------------
//--------------------------------------------------------------------


	r := bufio.NewReader(os.Stdin)
//	println("enter string:")
	line, err := r.ReadString(delim)
	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
	}
//	fmt.Println(line)
	token := strings.Fields(line)
	for token[0]!= "exit" || token[0]== "" {
	switch token[0]{
	case "whoami":
//		fmt.Println("this is 1")
		a := kademlia.Whoami(kadem)
		fmt.Println(a)

	case "local_find_value":
//		fmt.Println("this is 2")
		key,err := kademlia.FromString(token[1])
//		fmt.Println(id)	
		if err == nil {
			kademlia.Local_find_value(kadem,key)
		}
	case "get_contact":
//		fmt.Println("this is 3")
		id,err := kademlia.FromString(token[1])
		fmt.Println(id)	
		var c kademlia.Contact
		if err == nil {
			c = kademlia.Get_contact(kadem,id)
				fmt.Println("Node id: ",c.NodeID)
				fmt.Println("IP Addr: ",c.Host)
				fmt.Println("Port   : ",c.Port)
		}
		
	case "iterativeStore":
//		fmt.Println("this is 4")
	case "iterativeFindNode":
//		fmt.Println("this is 5")
	case "iterativeFindValue":
//		fmt.Println("this is 6")
	case "ping":
//		fmt.Println("this is 7")
		id,err := kademlia.FromString(token[1])
		if err == nil {
			ping,pong,s := kademlia.Pingnode(kadem,id)
			client, err := rpc.DialHTTP("tcp",s) //REFPEER
//			s = s + ":"
			if err != nil {
				log.Fatal("DialHTTP: ", err)
			}
			err = client.Call("Kademlia.Ping", ping, &pong)
			if err != nil {
				log.Fatal("Call: ", err)
			}
			log.Printf("ping msgID: %s\n", ping.MsgID.AsString())
			log.Printf("pong msgID: %s\n", pong.MsgID.AsString())
		}

	case "store":
//		fmt.Println("this is 8")

		id,err := kademlia.FromString(token[1])
		if err != nil {
		}
		key,err := kademlia.FromString(token[2])
		if err != nil {
		}

//		val:=token[3]
		val := ""
		for t:=3;t<len(token);t++ {
		val = val +string(token[t])+" "
		} 

		
		sreq, sres, s :=kademlia.Storenode(kadem,id,key,val)
//			s = s + ":"
		client, err := rpc.DialHTTP("tcp",s)//REFPEER) //s
		if err != nil {
			log.Fatal("DialHTTP: ", err)
		}
		err = client.Call("Kademlia.Store", sreq, &sres)
		if err != nil {
			log.Fatal("Call: ", err)
		}

//		log.Printf("sreq msgID: %s\n", sreq.MsgID.AsString())
//		log.Printf("sres msgID: %s\n", sres.MsgID.AsString())




	case "find_node":
//		fmt.Println("this is 9")
		id,err := kademlia.FromString(token[1])
		if err != nil {
		}
		key,err := kademlia.FromString(token[2])
		if err != nil {
		}
		fnreq, fnres, s :=kademlia.Findnode(kadem,id,key)
//			s = s + ":"
		client, err := rpc.DialHTTP("tcp",s)//REFPEER) //s
		if err != nil {
			log.Fatal("DialHTTP: ", err)
		}
		err = client.Call("Kademlia.FindNode", fnreq, &fnres)
		if err != nil {
			log.Fatal("Call: ", err)
		}
		log.Printf("fnreq msgID: %s\n", fnreq.MsgID.AsString())
		log.Printf("fnres msgID: %s\n", fnres.MsgID.AsString())
		var temp_contact kademlia.Contact
		l := len(fnres.Nodes)
		for j:=0 ; j < l ;j++ {
		if fnres.Nodes[j].NodeID.AsString() != "0"{
		fmt.Println("NodeID : ",fnres.Nodes[j].IPAddr)
		fmt.Println("Port   : ",fnres.Nodes[j].Port)
		fmt.Println("Host   : ",fnres.Nodes[j].NodeID.AsString())
		
		peer3 := fmt.Sprintf("%s:%d",fnres.Nodes[j].IPAddr,fnres.Nodes[j].Port)
	addr3, err := net.ResolveTCPAddr("tcp4", peer3)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	

		temp_contact.NodeID = fnres.Nodes[j].NodeID
		temp_contact.Host = addr3.IP  //net.IP(fnres.Nodes[j].IPAddr)
		temp_contact.Port = fnres.Nodes[j].Port
		kadem.KademBuckets.Update_Contact(temp_contact,kadem.MyContact)
		}
		}//for 


	case "find_value":
//		fmt.Println("this is 10")
		id,err := kademlia.FromString(token[1])
		if err != nil {
		}
		key,err := kademlia.FromString(token[2])
		if err != nil {
		}
		fvreq, fvres, s := kademlia.Findvalue(kadem,id,key)
//			s = s + ":"
		client, err := rpc.DialHTTP("tcp",s)//REFPEER) //s
		if err != nil {
			log.Fatal("DialHTTP: ", err)
		}
		err = client.Call("Kademlia.FindValue", fvreq, &fvres)
		if err != nil {
			log.Fatal("Call: ", err)
		}
//		log.Printf("fvreq msgID: %s\n", fvreq.MsgID.AsString())
//		log.Printf("fvres msgID: %s\n", fvres.MsgID.AsString())
		
		fmt.Println(string(fvres.Value))
	case "exit":
		break
	default :
		fmt.Println("enter the correct input")
	}//switch
	line, err := r.ReadString(delim)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	token = strings.Fields(line)
	}//for

}//main
