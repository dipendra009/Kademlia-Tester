package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"net"
	"fmt"
)

// Host identification.
type Contact struct {
    NodeID ID
    Host net.IP
    Port uint16
}


// PING
type Ping struct {
    Sender Contact
    MsgID ID
}

type Pong struct {
    MsgID ID
    Sender Contact
}

func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
    // This one's a freebie.
    pong.MsgID = CopyID(ping.MsgID)
//    fmt.Println(ping.MsgID.AsString())
//    fmt.Println(ping.Sender.NodeID)
//    fmt.Println(ping.Sender.Host)
//    fmt.Println(ping.Sender.Port)

    return nil
}


// STORE
type StoreRequest struct {
    Sender Contact
    MsgID ID
    Key ID
    Value []byte
}

type StoreResult struct {
    MsgID ID
    Err error
}

func (k *Kademlia) Store(req StoreRequest, res *StoreResult) error {
//	fmt.Println("store function called")
	// TODO: Implement.
	res.MsgID = CopyID(req.MsgID)
	//update the contact in the bucket
	k.KademBuckets.Update_Contact(req.Sender,k.MyContact)
	//storing the key-value pair
	//    key := req.Key
	//    value := req.Value
	Store_keyvalue(k,req.Key, string(req.Value))
	//    STORE(key,value)

	


    return nil
}

//-----------------------------------------------------
// FIND_NODE
type FindNodeRequest struct {
    Sender Contact
    MsgID ID
    NodeID ID
}

type FoundNode struct {
    IPAddr string
    Port uint16
    NodeID ID
}

type FindNodeResult struct {
    MsgID ID
    Nodes []FoundNode
    Err error
}

func (k *Kademlia) FindNode(req FindNodeRequest, res *FindNodeResult) error {
	// TODO: Implement.
	res.MsgID = CopyID(req.MsgID)
	//update the contact in the bucket
	k.KademBuckets.Update_Contact(req.Sender,k.MyContact)

	nodes := k.KademBuckets.Local_Find_Node(req.NodeID,k.MyContact.NodeID)
	res.Nodes = nodes
//    sender :=req.Sender
//    Node :=req.NodeID
//     update(sender)
//     res.Nodes = FindNode(Node)


    return nil
}

//-------------------------------------------------------------
// FIND_VALUE
type FindValueRequest struct {
    Sender Contact
    MsgID ID
    Key ID
}
// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
    MsgID ID
    Value []byte
    Nodes []FoundNode
    Err error
}

func (k *Kademlia) FindValue(req FindValueRequest, res *FindValueResult) error {
	// TODO: Implement.
	res.MsgID = CopyID(req.MsgID)
	//update the contact in the bucket
	k.KademBuckets.Update_Contact(req.Sender,k.MyContact)

	res.Value = []byte(Local_Find_Value(k, req.Key))	
	if res.Value == nil {
	res.Nodes = k.KademBuckets.Local_Find_Node(req.Sender.NodeID,k.MyContact.NodeID)
//	res.Nodes = k.KademBuckets.Find_Nod(req.Sender,k.MyContact)
	}else{
	fmt.Println("value found :", string(res.Value) )
	}

//    sender :=req.Sender
//    key :=req.key
//     update(sender)
//     res.Value = Findvalue(key)
//     res.Nodes = FindNode(key)

    return nil
}
