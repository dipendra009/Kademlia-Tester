package kademlia
import (
	"math"
//	"fmt"
	"container/list"
)
type NodeBucket struct{
	bucket [160] *list.List
	K int
	count int
}
//------------------------------------------------------
func NewNodeBucket( bucketsize int) (nodebucket NodeBucket){
	nodebucket.K = bucketsize
	nodebucket.count = 0
	for i:=0 ; i<160 ; i++ {
		nodebucket.bucket[i] = list.New()
	}
	return nodebucket
}
//----------------------------------------------------
//-----------find the element based on the node-id---------
func Find_Contact(nodeid ID, bucket *list.List) (b bool, e *list.Element){
//	fmt.Println("find contact is called")
	
	e = bucket.Front()
//	fmt.Println("bucket front: ",e.Value.(Contact).NodeID)
	if e==nil{
//		fmt.Println("bucket conact nil")
	}
	b = false
	for ; e!=nil ; e=e.Next(){
	//	fmt.Println("bucket conact not nil")
		if (e.Value.(Contact).NodeID == nodeid){
		//	fmt.Println("id matched")
			b = true
			return b,e
		}
	}
	return b,e
}
//-------------------------------------------------------------------------------
//-------------------------------Update_Contact()--------------------------------
//-------------------------------------------------------------------------------
func (nodebucket *NodeBucket) Update_Contact(sender Contact,mycontact Contact) {
//	fmt.Println("Update_Contact function called")
	//take xor to find the distance
	distance := mycontact.NodeID.Xor(sender.NodeID)
//	fmt.Println(distance)
	//finding the appropriate bucket
	//math.Log(dist) //64bit....but we need to have 160 bit log
	i := 0
	for int(distance[i])==0 && i<19 {
		i++
	}
	num := float64(distance[i])
	z :=0
	if num>0{
		z = int(math.Log(num))
		z = z + (((19-i)*8)) // z is the k-bucket index
	}
//	z = 159-distance.PrefixLen()
	//-----------------------------------------------
//	fmt.Println("update contact--find contact called with bucket - ",z)	
	b, e := Find_Contact(sender.NodeID,nodebucket.bucket[z])
	if b {
		nodebucket.bucket[z].MoveToBack(e)
	}else{
//		fmt.Println("bucket count",nodebucket.bucket[z].Len())
//		fmt.Println("max limit of bucket",nodebucket.K)
		if nodebucket.bucket[z].Len() < nodebucket.K {
			c := sender
			nodebucket.bucket[z].PushBack(c)
			nodebucket.count++
		} else {
			//ping nodebucket.bucket[z].head
			//func (l *List) Front() *Element
			//x  := nodebucket.bucket[z].Front()
			//id := x.Value.(*Contact).NodeID
			//if(ping_respond_time<T){
				//func (l *List) MoveToBack(e *Element)
				//nodebucket.bucket[z].MoveToBack(x)
			//} else {
				//  func (l *List) Remove(e *Element) interface{}
				//nodebucket.bucket[z].PushBack(sender)
			//}

		}		

	}
}
//--------------------------Update_Contact() Ends-----------------------------------------------------
//--------------------------------------------------------------------------------------------------
//--------------------------------------------------------------------------------------------------










//-------------------------------------------------------------------------------
//-------------------------------Find_Node()--------------------------------
//-------------------------------------------------------------------------------
func (nodebucket *NodeBucket) Local_Find_Node(findid ID,myid ID) (nodes []FoundNode){
//	fmt.Println("Local_Find_Node function called")
	//take xor to find the distance
	distance := myid.Xor(findid)
	//finding the appropriate bucket
	//math.Log(dist) //64bit....but we need to have 160 bit log
	i := 0
	for int(distance[i])==0 && i<19 {
		i++
	}
	num := float64(distance[i])
	z :=0
	if num>0{
		z = int(math.Log(num))
		z = z + (((19-i)*8)) // z is the k-bucket index
	}
//	z = 159-distance.PrefixLen()

	//------------------------------------------------------------------
//	fmt.Println("find contact called with bucket - ",z)	
	b, e := Find_Contact(findid,nodebucket.bucket[z])
	if b {
		nodes = make([]FoundNode,1)
		// return the contact
//		fmt.Println("contact found...filleing nodes[0]")
		nodes[0].IPAddr = e.Value.(Contact).Host.String()
		nodes[0].Port   = e.Value.(Contact).Port
		nodes[0].NodeID = e.Value.(Contact).NodeID
//		fmt.Println("contact found...filled nodes[0]")

	}else{
//		fmt.Println("bucket count",nodebucket.bucket[z].Len())
//		fmt.Println("max limit of bucket",nodebucket.K)
		if nodebucket.bucket[z].Len() < nodebucket.K {
		limit:=0
		if nodebucket.K < nodebucket.count {
		nodes = make([]FoundNode,nodebucket.K)
		limit = nodebucket.K
		}else {
		nodes = make([]FoundNode,nodebucket.count)
		limit = nodebucket.count
		}
		//return k closest
			i:=0
			t:=z
			for t>=0 && i<limit {		
				e =  nodebucket.bucket[t].Front()
				for ; e!=nil && i<nodebucket.K; e=e.Next(){
					nodes[i].IPAddr = e.Value.(Contact).Host.String()
					nodes[i].Port   = e.Value.(Contact).Port
					nodes[i].NodeID = e.Value.(Contact).NodeID
					i++
				}//for
			t--
			}
			if i<limit {
			t=z
			for t<160 && i<limit {		
				e =  nodebucket.bucket[t].Front()
				for ; e!=nil && i<nodebucket.K; e=e.Next(){
					nodes[i].IPAddr = e.Value.(Contact).Host.String()
					nodes[i].Port   = e.Value.(Contact).Port
					nodes[i].NodeID = e.Value.(Contact).NodeID
					i++
				}//for
			t++
			}//for

			}
		} else {
		// return the whole bucket
			nodes = make([]FoundNode,nodebucket.K)
			e =  nodebucket.bucket[z].Front()
			i:=0
			for ; e!=nil ; e=e.Next(){
				nodes[i].IPAddr = e.Value.(Contact).Host.String()
				nodes[i].Port   = e.Value.(Contact).Port
				nodes[i].NodeID = e.Value.(Contact).NodeID
				i++
			}//for
		}//else



		
	}//else
	return nodes
}

//--------------------------Find_Node() Ends-----------------------------------------------------
//--------------------------------------------------------------------------------------------------
//--------------------------------------------------------------------------------------------------
