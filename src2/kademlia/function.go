package kademlia

import (
	"fmt"
	"math"
//	"net"
)

//-------1-------------------------
func Whoami(K *Kademlia) (s string){
	id := WhoAmI(K)
//	fmt.Println(id)
	s = id.AsString()
	return s
}
//-------2-------------------------
func Local_find_value(K *Kademlia, key ID) {
	s :=Local_Find_Value(K, key)
	if s!="" {
		fmt.Println(s)
	}else{
		fmt.Println("ERR: key-value not found")
	}
}
//-------3--------------------------
func Get_contact(K *Kademlia, nodeid ID) (c Contact){
	//----finding the appropriate bucket------
	dist := K.MyContact.NodeID.Xor(nodeid)	
	i := 0
	for int(dist[i])==0 && i<19 {
		i++
	}
	num := float64(dist[i])
	z :=0
	z = 160-dist.PrefixLen()
	if num>0{
		z = int(math.Log(num))
		z = z + (((19-i)*8)) // z is the k-bucket index
	}
	//--------searching fo rthe nodeid--------------------	
//	fmt.Println("get contact--searching in bucket no : ",z)
	b, e := Find_Contact(nodeid,K.KademBuckets.bucket[z])
	if b {
	 c = e.Value.(Contact)
	}else{
		fmt.Println("ERR: Contact not found")
	}
	return
}

//-----------------------------------------------------------
//----------4------------------------------------------------
//-----------iterativeStore key value------------------------
//----------5------------------------------------------------
//------------iterativeFindNode ID---------------------------
//----------6------------------------------------------------
//------------iterativeFindValue key-------------------------
//-----------------------------------------------------------


//---------7--------------------------
func Pingnode(K *Kademlia, nodeid ID) (ping Ping,pong Pong,s string){
	ping.Sender = K.MyContact
	ping.MsgID = NewRandomID()
	conta := Get_contact(K,nodeid)
//	tcpaddr :=new(net.TCPAddr)
//	tcpaddr.IP   = conta.Host
//	tcpaddr.Port = int(conta.Port)
//	addr := tcpaddr.String()
	s = fmt.Sprintf("%s:%d",conta.Host,conta.Port)
//	fmt.Println("host IP:Port string",s)
	return
/*
	//----------------------------------
	client, err := rpc.DialHTTP("tcp",s)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
	}
	err = client.Call("Kademlia.Ping", ping, &pong)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	//-----------------------------------
*/
}

//---------8-----------------------------
func Storenode(K *Kademlia, nodeid ID, key ID, value string) (sreq StoreRequest,sres StoreResult,s string){
	sreq.Sender = K.MyContact
	sreq.MsgID  = NewRandomID()
	sreq.Key    = key
	sreq.Value  = []byte(value)
	conta := Get_contact(K,nodeid)
//	tcpaddr :=new(net.TCPAddr)
//	tcpaddr.IP   = conta.Host
//	tcpaddr.Port = int(conta.Port)
//	addr := tcpaddr.String()
	s = fmt.Sprintf("%s:%d",conta.Host,conta.Port)
	return
/*
	//----------------------------------
	sreq, sres, s :=Storenode(nodeid,key,value)
	client, err := rpc.DialHTTP("tcp",s)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
	}
	err = client.Call("Kademlia.Store", sreq, &sres)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	//-----------------------------------
*/
}

//------------9-------------------------------------
func  Findnode(K *Kademlia, nodeid ID, key ID) (fnreq FindNodeRequest,fnres FindNodeResult,s string){
	fnreq.Sender = K.MyContact
	fnreq.MsgID  = NewRandomID()
	fnreq.NodeID    = key
	conta := Get_contact(K,nodeid)
//	tcpaddr :=new(net.TCPAddr)
//	tcpaddr.IP   = conta.Host
//	tcpaddr.Port = int(conta.Port)
//	addr := tcpaddr.String()
	s = fmt.Sprintf("%s:%d",conta.Host,conta.Port)
	return
/*
	//----------------------------------
	fnreq, fnres, s :=Findnode(nodeid,key)
	client, err := rpc.DialHTTP("tcp",s)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
	}
	err = client.Call("Kademlia.FindNode", fnreq, &fnres)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	//-----------------------------------
*/
}
//------------10----------------------------------------------
func Findvalue(K *Kademlia,nodeid ID, key ID) (fvreq FindValueRequest,fvres FindValueResult,s string){
	fvreq.Sender = K.MyContact
	fvreq.MsgID  = NewRandomID()
	fvreq.Key    = key
	conta :=Get_contact(K,nodeid)
//	tcpaddr :=new(net.TCPAddr)
//	tcpaddr.IP   = conta.Host
//	tcpaddr.Port = int(conta.Port)
//	addr := tcpaddr.String()
	s = fmt.Sprintf("%s:%d",conta.Host,conta.Port)
	return
/*
	//----------------------------------
	fvreq, fvres, s :=Findvalue(nodeid,key)
	client, err := rpc.DialHTTP("tcp",s)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
	}
	err = client.Call("Kademlia.FindValue", fvreq, &fvres)
	if err != nil {
		log.Fatal("Call: ", err)
	}
	//-----------------------------------
*/
}

