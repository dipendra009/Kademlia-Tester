package kademlia
//import "fmt"
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

// Core Kademlia type. You can put whatever state you want in this.



//------storing both key and value.....storing key is not necessary it is taken care by hashmap map[key]-----
//------to invalidate a key value pair------
type KeyValue struct{
	key ID
	value string
}
//-----------BucketNode structure------------------
/*
type BucketNode struct{
	IP_Add		net.IP
	UDP_Port	uint16
	Node_Id		ID
	Next		*BucketNode
}
*/
//---------------------------------------------



//----------Kademlia node structure-------------------------------------


type Kademlia struct {
	MyContact Contact
	NodeID ID
	mapval map[ID]KeyValue
	KademBuckets NodeBucket
}


func NewKademlia() *Kademlia {
    // TODO: Assign yourself a random ID and prepare other state here.
	KAD := new(Kademlia)
	KAD.NodeID = NewRandomID()
	KAD.mapval = make(map[ID]KeyValue)
	KAD.KademBuckets = NewNodeBucket(20)
    return KAD
}

//----------------------------------
func WhoAmI(K *Kademlia) ID {
	return K.NodeID
}




//------------------------------------------------------------------------------------------------
// type KeyValue struct{
//	key ID
//	value string
// }
//-------------------------
//we will be making a hashmap
// var m map[ID]KeyValue
//--------------------------
// mapval = make(map[ID]KeyValue)...it should be inside Kademlia structure
//--------------------------------------------------------------------------------------------------





//-----function to invalidate a key-value pair-----------
func Local_delete_Value(A *Kademlia, key ID) {
 delete(A.mapval,key)
}

//-----retrieving the value------------------------------
func Local_Find_Value(A *Kademlia, key ID) string {
 return ( A.mapval[key].value)
}
//------storing the value-------------------------------
func Store_keyvalue(A *Kademlia, key ID,value string) {
	A.mapval[key]=KeyValue{key,value}
}
//-----------------------------------------------------




