package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import(
	"net"
	"log"
    "strings"
    "strconv"
    "container/list"
    "net/rpc"
    "fmt"
    "time"

)

const alpha = 3

const B = 160

const k = 20

const tExpire = 86400

const tRefresh = 3600

const tReplicate = 3600

const tRepublish = 86400

const KEYBYTES = 20
type KEY [KEYBYTES]byte

const VALUEBYTES = 4096
type VALUE []byte

type Bucket *list.List
var slChan, alChan, clChan chan *list.List //these are used for synchronised parallelism
var alpha1, alpha2, alpha3 chan int
var fv1, fv2, fv3 chan VALUE

// Core Kademlia type. You can put whatever state you want in this.
type Kademlia struct {
	NodeID ID
	Storage map[ID] VALUE
	Buckets [] *list.List    //bucket is an array of list
	MyContact Contact
	RefreshList [] bool
}

//this adds a contact given in the parameter, to the bucket of a kademlia node
func Update (kad *Kademlia,newContact *Contact){
	//compute the distance at first, and position
	distance := newContact.NodeID.Xor(kad.NodeID)
	pos := (159 - distance.PrefixLen())

	//if pos is 160, then it is out of bound, so return
	if(pos == -1) {
		pos=0
	}

	//setting false, so that next time it doesn't refreshed
	kad.RefreshList[pos]=false
	//we must check if the contact is already in
	//elem:=ElementInList(kad.Buckets[pos], newContact)
	contact := HasContact(kad, newContact.NodeID)
	//fmt.Println(contact)
	if  contact != nil {
	    elem :=  ElementInList(kad.Buckets[pos], contact)
	    kad.Buckets[pos].MoveToFront(elem)
	    return
	}

	//fmt.Println("Adding...")

    //if current length of the bucket is less than 20, just add it
	if kad.Buckets[pos].Len() < 20 {
		kad.Buckets[pos].PushFront(newContact) //push it in the front
		//fmt.Println("inserted in bucket ", pos," ", kad.Buckets[pos].Len())
	} else {
		//retrieve the last element of the bucket
		e := kad.Buckets[pos].Back()
		oldest := e.Value.(*Contact)

		//ping and pong oldest
		ipPort := strings.Join([]string{oldest.Host.String(), ":", strconv.FormatUint(uint64(oldest.Port),10)},"")
		client, err := rpc.DialHTTP("tcp", ipPort)
		if err != nil {
		    log.Fatal("ERR in Bueckt Add: ", err)
		}
		ping := new(Ping)
		ping.MsgID = NewRandomID()
		ping.Sender = kad.MyContact
		var pong Pong
		err = client.Call("Kademlia.Ping", ping, &pong)
		//there is error while pinging the oldest node
		if err != nil {
			//remove the old and the new one in the list
			kad.Buckets[pos].Remove(e);
    		kad.Buckets[pos].PushFront(newContact);
		} else {
			//make the old one in first position
			kad.Buckets[pos].MoveToFront(e);
		}
	}
}



// Join for a new code
func Join(kad *Kademlia, knownContact Contact) error {
    //update is performed previously in new kademlia
    Update(kad, &knownContact)
    //kCloseNodes := list.New()
    //TODO: make iterative find node such that it stores the sender in everyone's contact
    err, kCloseNodes := kad.IterativeFindNode(kad.MyContact.NodeID, true)
    kCloseNodes.Len()
    if err != nil {
      return err
    }
    return nil
}

//Refreshes the k-buckets after tRefresh time
func  Refresh (kad *Kademlia) error {
        //loop through all the buckets using i
        var i uint64
        for i = 0; i < B; i = i+1 {
                //if corresponfing alivelist is not true
                if kad.RefreshList[i] != false {
                        //selecting a random number in that list
                        //idea is simple, at first assign it to the current node id
                        //then change the particular index of that node id based on i, m, n
                        tId := kad.NodeID
                        var m uint64
                        for m = 0; m < IDBytes; m = m+1 {
                                var n uint64
                                for n = 0; n < 8; n = n+1 {
                                        if(m*8+n == i) {
                                                tId[m] = tId[m] ^ (1 <<n)
                                                break
                                        }
                                }
                        }
                        //kCloseNodes := list.New()
                        //now performing the iterative find node using that key
                        err, kCloseNodes := kad.IterativeFindNode(tId, true)
                        kCloseNodes.Len()
                        if err != nil {
                            log.Fatal("ERR in Refresh: ", err)
                        }
                }
                kad.RefreshList[i] = true //setting it true for the next iteration
        }
        return nil
}

//calls refresh every 1 hour
func CallRefresh(kad *Kademlia) {
        for ;; {
            time.Sleep(tRefresh * time.Second)
            Refresh(kad)
            fmt.Println("Refreshing..")
        }
}



func NewKademlia() (kademlia *Kademlia) {
    // TODO: Assign yourself a random ID and prepare other state here.
    //initialising channels
    slChan = make (chan *list.List, 3*alpha)  //buffered channels
    alChan = make (chan *list.List, 3*alpha)
    clChan = make (chan *list.List, 3*alpha)
    alpha1 = make (chan int, 0)  //unbuffered channels to synch with function completions
    alpha2 = make (chan int, 0)
    alpha3 = make (chan int, 0)

    fv1 = make (chan VALUE, 0)  //unbuffered channels to synch with function completions
    fv2 = make (chan VALUE, 0)
    fv3 = make (chan VALUE, 0)


    //intiialing kademlia struct values
    kademlia = new (Kademlia)
    kademlia.Storage = make(map[ID]VALUE)
    kademlia.MyContact.NodeID = NewRandomID()
    kademlia.Buckets = make([](*list.List), IDBytes*8)
    kademlia.RefreshList = make([]bool, IDBytes * 8)
    for i:=0; i<B; i=i+1 {
	   kademlia.Buckets[i] = list.New()
	   kademlia.RefreshList[i] = true
    }
    //Update(kademlia, &kademlia.MyContact)
    kademlia.NodeID = kademlia.MyContact.NodeID
    go CallRefresh(kademlia)
    return
}

// Hash function
func HashKey(value []byte) (checksum [20]byte){
	//checksum = sha1.Sum(value)
	return
}

//---------------------RPC helper functions------------------------------------------------

func StoreData(kad *Kademlia, key ID, Value []byte) error {
	kad.Storage[key] = Value
	//fmt.Println("Stored: ", key, ":",Value)
	return nil
}

func FindNodeData(kad *Kademlia, NodeID ID)(Nodes []FoundNode, Err error){

	l := list.New()

	//fmt.Println("Finding nodes for ", NodeID)
	//fmt.Println(kad.Buckets)

	 for i:=0; i < IDBytes * 8 ; i = i+1{
		if kad.Buckets[i].Len() != 0{
		//fmt.Println(kad.Buckets[i].Len(), kad.Buckets[i])
		l.PushFrontList(kad.Buckets[i])
		}
	 }
	SortList(l, NodeID)

	var n []FoundNode

	if k < l.Len(){


	n = make([]FoundNode, k)
	}

	if k > l.Len(){
	n = make([]FoundNode, l.Len())
	}

	for i:=0; i < l.Len(); {
		if i < k{
			c := l.Front()
			l.Remove(c)

			con := c.Value.(*Contact)
			//fmt.Println(con.NodeID.AsString())
			if con.NodeID.AsString() != "0"{
			//fmt.Println(con.Host, con.Port, con.NodeID)
			n[i] = FoundNode{con.Host.String(), con.Port, con.NodeID}
			i = i+1
			}
		}
	}

	//fmt.Println("Found: ", n)

	Nodes = n
	Err = nil

	return


}

func FindValueData(kad *Kademlia, key ID)(Value []byte, Nodes []FoundNode, Err error){


	//fmt.Println("Searching for ",key)
	Value, ok := kad.Storage[key]
	if ok {

		//fmt.Println("ok...", Value)

		Err = nil
		Nodes = nil

	    return

	}

	// Find k closest nodes
	Nodes, Err = FindNodeData(kad, key)
	return
}


/////////-----------NS Lookup Helper Functions------/////////////////////////////
//to perform 'contains' check in a list
func ContainsInList (l *list.List, c *Contact) bool {
	//simple, iterate over the list and check if any element's node id matches the given element's node id
	//if yes, return true, else return false
	for e := l.Front(); e != nil; e = e.Next(){
		if e.Value.(*Contact).NodeID.Equals(c.NodeID){
			return true
		}
	}
	return false
}

//similar as contains in list, just returns the element here
func ElementInList (l *list.List, c *Contact) *list.Element {
	//simple, iterate over the list and check if any element's node id matches the given element's node id
	//if yes, return true, else return false
	for e := l.Front(); e != nil; e = e.Next(){
		if e.Value.(*Contact).NodeID.Equals(c.NodeID){
			return e
		}
	}
	return nil
}

//to sort a list, based on the distance from the given key in ascending order
func SortList(l *list.List, keyId ID){
    //a new list to keep the sorted elements
	sorted := list.New()
    //iterate over the given list
    for e :=l.Front();e!=nil; e=e.Next() {
        var changed bool = false
        //compute distance
		var c *Contact
        c = e.Value.(*Contact)
        dist := c.NodeID.Xor(keyId)
        //loop over the sorted array to find position for the current element
        if sorted.Len()==0 {
            sorted.PushFront(c)
            continue
        }
        for s:=sorted.Front(); s!=nil; s=s.Next(){
		   c1 := s.Value.(*Contact)
		   dist1 := c1.NodeID.Xor(keyId)
		   larger := dist1.Less(dist)
		   // if larger is true then insert the element in the current position
		   // end the loop
		   if larger {
		        changed=true
                sorted.InsertBefore(c,s)
                break
            }
        }
        //if not changed, insert the element at the back
        if !changed {
            sorted.PushBack(c)
        }
    }
    //push the sorted list into the parameter of the function
    l = list.New()
    l.PushFrontList(sorted)
    return
}

//it will return k closest nodes currently residing from the bucket
func GetKCloseNodes (kad *Kademlia, keyId ID) (*list.List) {
      //first job is to find k-closest nodes for the ID keyId
      //find the distance with the current node and use the prefix len to determine the index
      //fmt.Println(kad.NodeID.AsString())
      kCloseNodes := list.New()
      distance := keyId.Xor(kad.NodeID)
      index := (159 - distance.PrefixLen())

      //out of bound
      if index==-1 {
        index=0
      }

      //retrieve the relative bucket using the index
      buck := kad.Buckets[index]
      //a list to store all the nodes known

      //iterate over this bucket to store k-contacts in the bucket
      for e := buck.Front(); e != nil && kCloseNodes.Len() < k; e = e.Next(){
      	 c:= e.Value.(*Contact)
      	 kCloseNodes.PushFront(c)
      }
      //if less than 20 contacts in the closest bucket
      if kCloseNodes.Len() < k {
        //initialising values
        adder := 1
        counter := 0
        j := index
      	for ; counter < 160; counter++ {
      		// j is always negated until hit 0
      		j = j+adder
      		if j==160 {
      		    j=index-1 //then j is increased until reached 160
      		    adder=-1
      		}
      		if j==-1 {
      		    break
      		}
                //retrieve current bucket
      		curBuck := kad.Buckets[j]
      		// if current bucket has some contacts
      		if curBuck.Len() > 0 {
      			//add contacts until shortlist is full
      			for e := curBuck.Front(); e != nil && kCloseNodes.Len() < 20; e = e.Next() {
      				c := e.Value.(*Contact)
      				if !ContainsInList(kCloseNodes,c) {   //if shortList doesn't contains the contact, add it
      					kCloseNodes.PushFront(c)
      				}
      			}
      		}
      		if kCloseNodes.Len()>=20 {
      		    break
      		}
      	}
      }
      return kCloseNodes
}


//---------------------------------------------------------------------
// Has Contact looks in bucket and returns the contact for a given ID
func HasContact(kad *Kademlia, NodeID ID) (contact *Contact){

	distance := kad.NodeID.Xor(NodeID)
	index := (159 - distance.PrefixLen())
	//fmt.Println("index: ", index,kad.Buckets)
	//fmt.Println("index: ", index)
	if index==-1 {
	    index=0
	}
	if kad.Buckets[index] == nil || kad.Buckets[index].Len() == 0 {
		return nil
	}
	l := kad.Buckets[index]
	//fmt.Println("l ", l)

	if l == nil{

		return nil
	}

	//if ContainsInList(l, x) == false{
		//  return nil
	//}

	//simple, iterate over the list and check if any element's node id matches the given element's node id
	//if yes, return true, else return false
	for e := l.Front(); e != nil; e = e.Next(){
		if e.Value.(*Contact).NodeID.Equals(NodeID){
			return e.Value.(*Contact)
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------

// Commands

// WhoAmI command

func WhoAmI (kad *Kademlia) {
	//fmt.Println(kad.MyContact.NodeID)
	fmt.Println(kad.MyContact.NodeID.AsString())
	return
}

 // Local Find Value Command

// Local Find Value for a given key
// If your node has data for the given key, print it.
// If your node does not have data for the given key, you should print "ERR"

func LocalFindValue(kad *Kademlia, key string){
	if key, err := FromString(key);  err == nil {
		if value, ok := kad.Storage[key]; ok {
			fmt.Println(string(value))
			return
		}
		fmt.Println("ERR")
		return
	}
	fmt.Println("ERR")
}

// Get Contact Command

// Get Contact for a given key
// If your buckets contain a
// node with the given ID, printf("%v %f \n", theNode.addr, theNode.port)
// If your node does not have data for the given key, you should print "ERR".

func GetContact (kad *Kademlia, ID string) {
	if ID, err := FromString(ID); err == nil{
		if value := HasContact (kad, ID); value != nil{
			fmt.Printf("%v %d \n", value.Host, value.Port)
			return
		}
		print("ERR\n")
		return

	}
}

// Iterative Store command

// Perform the iterativeStore operation and then print the ID of the node that
//********************************Need Implementation *************************
// recieved the final STORE operation

func (kad *Kademlia) IterativeStore(req StoreRequest, kCloseNodes *list.List) error {
      //at first we have to call iterative findnode function with the key/id value
      //and provide kCloseNodes list to store the k-closest nodes found
      kCloseNodes = list.New()
      err, ka := kad.IterativeFindNode(req.Key, true)
      kCloseNodes = ka

      //it will return k closest ids to the value, or an error
      //if no error, we have to call the Store function from the rpc.go file
      if err == nil {
          //no contacts in the bucket, that means this node has to store the value
          if kCloseNodes==nil || kCloseNodes.Len() <=0 {
              StoreData(kad, req.Key, req.Value)
              return err
          }
          //else store the value to each of the k-nodes?
          //FIX ME: am i correct?
          var tmp string
          for e := kCloseNodes.Front(); e != nil; e = e.Next() {
          	    // store the tmp node id, so that we can print the node id of the last store recipient node
              	c := e.Value.(*Contact)
              	if e.Next() == nil {
                   tmp = c.NodeID.AsString()
                }
              	//retreve the ip:port string from the contact
                ipPort := strings.Join([]string{c.Host.String(), ":", strconv.FormatUint(uint64(c.Port),10)},"")

                strReq:= StoreRequest{kad.MyContact, NewRandomID(), req.Key, req.Value}
                var strRep StoreResult
                //establish the tcp connection
                client, err1 := rpc.DialHTTP("tcp", ipPort)
                if err1 != nil {
                    //log.Fatal("Call: ", err1)
                    continue
                }
                //make the store call for the contact
                err2 := client.Call("Kademlia.Store", strReq, &strRep)
                if err2 != nil {
                    //log.Fatal("Call: ", err1)
                    continue
                }
          }
          //print last node id
          fmt.Println("store ", tmp)
      }
      //handle error
      return err
}

// Iterative Find Node command

// Print a list of <= k closest nodes and print their IDs.
// You should collect the IDs in a slice and print that

// ***********************************Need Implementation************************

func (kad *Kademlia) IterativeFindNode(keyId ID, join bool) (error, *list.List) {
    //first job is to find k-closest nodes for the ID keyId
    //find the distance with the current node and use the prefix len to determine the index
    kCloseNodes := list.New()
    shortList := list.New()
    shortList = GetKCloseNodes(kad, keyId)

    //if shortList is still 0, then no closest node is found! so return..
    if shortList.Len() == 0 {
       fmt.Println("ERR: No nodes found")
       return nil, kCloseNodes
    }

    // a variable to keep track of the closest node found till now
    //sort the list and intialize the front element to closest node
    SortList(shortList, keyId)
    closestNode := shortList.Front().Value.(*Contact)

    //we have to iteratively perform find_node rpc now
    //for this, at first alpha contacts must be chosen from the k-shortlist contacts
    //then to each fo these alpha contacts find_node rpc must be sent
    //at each step, we should keep k-closest (to the given key) of the all known nodes in shortList
    //we also have to:
    // 1. keep track of list of already queried and currently querying nodes
    // 2. drop inactive nodes
    // 3. keep track of the current closest node
    //the iteration will end when there will be no new closest contacts found in a parallel query, or
    // zero contacts are returned from a query
    activeList := list.New()
    currentList := list.New()
    //pushing dummy values
    cnt := new(Contact)
    activeList.PushFront(cnt)
    currentList.PushFront(cnt)

    for ;; {
        nextList := list.New()
        //alpha contacts are selected each time
        i:=0
        for e:=shortList.Front(); e != nil && i<alpha; e = e.Next() {
                	c := e.Value.(*Contact)
                	if !ContainsInList (activeList, c) && !ContainsInList (currentList, c) {
               			nextList.PushFront(c)
               			i=i+1
               		}
        }

        //hasmore will be false if no new contacts found in the previous loop
        if nextList.Len() == 0 {
            break
        }

        //set the channel values
        slChan <-shortList
        alChan <-activeList
        clChan <-currentList

        //perform strict parallelism here
        i=0
        for e:=nextList.Front(); e != nil; e = e.Next() {
        	fnRes := FindNodeResult{}
        	sender := kad.MyContact
        	fnReq := FindNodeRequest {sender, NewRandomID(), keyId}
            go ParallelFindNode(e, i, fnReq, &fnRes, kad)
            i = i+1
        }

        //wait for all the functions to finish
        if i==1 {
            i = <-alpha1
        } else if i==2 {
            i = <-alpha1
            i = <-alpha2
        } else {
            i = <-alpha1
            i = <-alpha2
            i = <-alpha3
        }
        activeList = <-alChan
        currentList = <-clChan
        shortList = <-slChan

        //sort the list according to distance to key
        SortList(shortList, keyId)

        //trim the list to closest k nodes
        if shortList.Len() > k {
           tmpList := list.New()
           i=0
           for e:= shortList.Front(); e != nil && i < k; e= e.Next() {
        	  tmpList.PushBack(e.Value.(*Contact))
        	  i++
           }
           shortList = tmpList
        }


        if shortList.Len()==0 {
            break
        }

        //check whether the closest node has changed from previous iteration
        if closestNode.NodeID != shortList.Front().Value.(*Contact).NodeID {
            closestNode = shortList.Front().Value.(*Contact)
        } else {       // no new closest node found in the last iteration, so the iteration loop should stop here
           oldList:= shortList
           for e:=oldList.Front(); e!=nil; e = e.Next() {
             c := e.Value.(*Contact)
             if !ContainsInList(activeList, c) && !ContainsInList(currentList, c) {
                fnRes := FindNodeResult{}
        	    sender := kad.MyContact
        	    fnReq := FindNodeRequest {sender, NewRandomID(), keyId}

        	    if c.Port != 0 {
			        peer := fmt.Sprintf("%s:%d",c.Host, c.Port)
                	client, err := rpc.DialHTTP("tcp", peer)
                	if err != nil {
                	    log.Fatal("Call: ", err)
                	    continue
                	}
                	err1 := client.Call("Kademlia.FindNode", fnReq, fnRes)
              		if err1 != nil {
                		log.Fatal("Call: ", err1)
                		continue
               		}
               		Update(kad, c)
        	     }

        	   //process the newly found nodes
        	    for i := 0; i < len(fnRes.Nodes); i++ {
            		c := Contact {fnRes.Nodes[i].NodeID, net.ParseIP(fnRes.Nodes[i].IPAddr), fnRes.Nodes[i].Port}
    				//check
    				if !ContainsInList(shortList,&c) {
    				   shortList.PushBack(c)
    				}
    			}

    			//sort the list according to distance to key
    			SortList(shortList, keyId)

                //trim the list to closest k nodes
                if shortList.Len() > k {
                   tmpList := list.New()
                   i=0
                   for e:= shortList.Front(); e != nil && i < k; e= e.Next() {
                	  tmpList.PushBack(e.Value.(*Contact))
                	  i++
                   }
                   shortList = tmpList
                }
             }
           }
           break
        }
    }
    //copying the value of shortList to kCloseNodes
    kCloseNodes=shortList
    //print ids
    if !join {
        for e:=kCloseNodes.Front();e != nil;e = e.Next() {
            c := e.Value.(*Contact)
           fmt.Println(c.NodeID.AsString())
        }
    }
	//return true
    return nil, kCloseNodes
}


//parallel function for find node
func ParallelFindNode(e *list.Element, i int, fnReq FindNodeRequest, fnRes *FindNodeResult, kad *Kademlia) {
    c := e.Value.(*Contact)
    currentList := <-clChan
    activeList := <-alChan
    shortList := <-slChan
    if c.Port != 0 {
        currentList.PushFront(c)
        peer := fmt.Sprintf("%s:%d",c.Host, c.Port)
        client, err := rpc.DialHTTP("tcp", peer)
        if err != nil {
            log.Fatal("Call in pf: ", err)
            shortList.Remove(e)
            currentList.Remove(e)
            alChan <-activeList
            clChan <-currentList
            slChan <-shortList
            CallReturn(i)
            return
        }
        err1 := client.Call("Kademlia.FindNode", &fnReq, fnRes)
        if err1 != nil {
            log.Fatal("Call in pf: ", err1)
            shortList.Remove(e)
            currentList.Remove(e)
            alChan <-activeList
            clChan <-currentList
            slChan <-shortList
            CallReturn(i)
            return
        }
        //update the bucket of the sender
        Update(kad, c)
        activeList.PushFront(c)
        currentList.Remove(e)
    }
    //process the newly found nodes
    for j := 0; j < len(fnRes.Nodes); j++ {
        c := Contact {fnRes.Nodes[j].NodeID, net.ParseIP(fnRes.Nodes[j].IPAddr), fnRes.Nodes[j].Port}
        //check
        if !ContainsInList(shortList,&c) {
            shortList.PushBack(&c)
        }
    }
    slChan <-shortList
    alChan <-activeList
    clChan <-currentList
    CallReturn(i)
}

func CallReturn(i int) {
    //writing to the appropriate channel for function finish signaling
    if i==0 {
        alpha1 <-1
    } else if i==1 {
        alpha2 <-1
    }else if i==2 {
       alpha3 <-1
    }
}


// Iterative Find Value Command

// Printf("%v %v\n", ID, value) where ID refers to the node
// that finally returned the value.
// If you do not find a value, print "ERR".

///////////////**Another implementation of find value********/////////////////////////////////
func (kad *Kademlia) IterativeFindValue(keyId ID, result *FindValueResult) (error){
     //at first checking whether the "kad" node contains the value
    value := kad.Storage[keyId]
    //if values is not present
    if value ==  nil {
        //first job is to find k-closest nodes for the ID keyId
        //find the distance with the current node and use the prefix len to determine the index
        shortList := list.New()
        shortList = GetKCloseNodes(kad, keyId)

        //if shortList is still 0, then no closest node is found! so return..
        if shortList.Len() == 0 {
               result.Value = nil   /***********************/
               result.Nodes = nil  /**********************/
               fmt.Println("ERR")
               return nil
        }

        // a variable to keep track of the closest node found till now
        //sort the list and intialize the front element to closest node
        SortList(shortList, keyId)
        closestNode := shortList.Front().Value.(*Contact)

        //we have to iteratively perform find_node rpc now
        //for this, at first alpha contacts must be chosen from the k-shortlist contacts
        //then to each fo these alpha contacts find_node rpc must be sent
        //at each step, we should keep k-closest (to the given key) of the all known nodes in shortList
        //we also have to:
        // 1. keep track of list of already queried and currently querying nodes
        // 2. drop inactive nodes
        // 3. keep track of the current closest node
        //the iteration will end when there will be no new closest contacts found in a parallel query, or
        // zero contacts are returned from a query
        activeList := list.New()
        currentList := list.New()
        //pushing dummy values
        cnt := new(Contact)
        activeList.PushFront(cnt)
        currentList.PushFront(cnt)

        done := false
        var prevContact *Contact = nil

        for ;; {
            nextList := list.New()
            //alpha contacts are selected each time
            i := 0
            for e:=shortList.Front(); e != nil && i<alpha; e = e.Next() {
                    	c := e.Value.(*Contact)
                    	if !ContainsInList(activeList, c) && !ContainsInList(currentList, c) {
                   			nextList.PushFront(c)
                   			i=i+1
                   		}
            }

            //hasmore will be false if no new contacts found in the previous loop
            if nextList.Len() == 0 {
                break
            }
///////////////////////Make change here, to perform find value, if value is found, do the other things //////////////////////
//////////////////////Other things: store the value to the closest node, that hasn't stored it ///////////////////////

            //set the channel values
            slChan <-shortList
            alChan <-activeList
            clChan <-currentList

            //perform strict parallelism here
            i=0
            for e:=nextList.Front(); e != nil; e = e.Next() {
                sender := kad.MyContact
                //create the request, result
                fvRequest := FindValueRequest {sender, NewRandomID(), keyId}
                var fvResult FindValueResult
                //setup the connection and other things
                go ParallelFindValue(e, i, fvRequest, &fvResult, kad)
                i=i+1
            }

            //wait for all the functions to finish
            var p, q, r VALUE
            if i==1 {
                p = <-fv1
            } else if i==2 {
                p = <-fv1
                q = <-fv2
            } else {
                p = <-fv1
                q = <-fv2
                r = <-fv3
            }

            //if one of p, q, r is 1, then a result has been found and stop
            if p!=nil || q!=nil || r!=nil {
                var pre *Contact = nil
                var c *Contact = nil
                var val VALUE
                if p!=nil {
                   pre=prevContact
                   val=p
                   c=nextList.Front().Value.(*Contact)
                } else if  q!=nil {
                   pre=nextList.Front().Value.(*Contact)
                   val=q
                   c=nextList.Front().Next().Value.(*Contact)
                }else{
                   pre=nextList.Front().Next().Value.(*Contact)
                   val=r
                   c=nextList.Front().Next().Next().Value.(*Contact)
                }
                if c!=nil {
                    fmt.Println(c.NodeID.AsString(), " ", string(val))
                    // Fill the result
                    result.Value = val   //set the appropriate value
                    result.Nodes = make([]FoundNode, 1, 1)    //make the found node array with 1 element
                    result.Nodes[0].NodeID = c.NodeID
                }
                if pre!=nil {
                    ipPort := strings.Join([]string{pre.Host.String(), ":", strconv.FormatUint(uint64(pre.Port),10)},"")
                    strReq := StoreRequest{kad.MyContact, NewRandomID(), keyId, val}
                    var strRes StoreResult
                    client, err := rpc.DialHTTP("tcp", ipPort)
                    if err != nil {
                        log.Fatal("Call ", err)
                    } else {
                        err = client.Call("Kademlia.Store", strReq, &strRes)
                        if err != nil {
                            log.Fatal("Call: ", err)
                        }
                    }
                }
                done=true
            }

            activeList = <-alChan
            currentList = <-clChan
            shortList = <-slChan

            //value is found, so break
            if done==true {
                return nil
            }
            //tracking the previous contact
            prevContact = nextList.Back().Value.(*Contact)

            //sort the list according to distance to key
             SortList(shortList, keyId)

            //trim the list to closest k nodes
            if shortList.Len() > k {
               tmpList := list.New()
               i=0
               for e:= shortList.Front(); e != nil && i < k; e= e.Next() {
            	  tmpList.PushBack(e.Value.(*Contact))
            	  i++
               }
               shortList = tmpList
            }

            if shortList.Len()==0 {
                break
            }

            //check whether the closest node has changed from previous iteration
            if closestNode.NodeID != shortList.Front().Value.(*Contact).NodeID {
                closestNode = shortList.Front().Value.(*Contact)
            } else {       // no new closest node found in the last iteration, so the iteration loop should stop here
                oldList:=shortList
                for e:=oldList.Front(); e!=nil; e = e.Next() {
                    c := oldList.Front().Value.(*Contact)
                    if !ContainsInList(activeList, c) && !ContainsInList(currentList, c) {
                        c := e.Value.(*Contact)
                        sender := kad.MyContact
                        //create the request, result
                        fvRequest := FindValueRequest {sender, NewRandomID(), keyId}
                        var fvResult FindValueResult
                        //setup the connection and other things
                        if c.Port != 0 {
                            ipPort := strings.Join([]string{c.Host.String(), ":", strconv.FormatUint(uint64(c.Port),10)},"")
                            client, err := rpc.DialHTTP("tcp", ipPort)
                            if err != nil {
                                log.Fatal("Err: ", err)
                                continue
                            }
                            err = client.Call("Kademlia.FindValue", fvRequest, &fvResult)
                            if err != nil {
                                log.Fatal("Err: ", err)
                                continue
                            }

                            Update(kad, c)

                            //if no error found, check if any value is recieved
                            //if any vale is recieved, we must break the iteration and return the value
                            if fvResult.Value != nil {
                                // Fill the result
                                result.Value = fvResult.Value   //set the appropriate value
                                result.Nodes = make([]FoundNode, 1, 1)    //make the found node array with 1 element
                                result.Nodes[0].NodeID = c.NodeID   //set the node id

                                //print the value and id
                                fmt.Println(c.NodeID.AsString()," ", string(fvResult.Value))

                                //When an iterativeFindValues succeeds, the initiator then has to store the key, value pair at the
                                //closest node seen which did not have the value
                                //as according to our implementation, kCloseNodes will be returned sorted, based on the distance
                                //So, this closest node, that did not have the value would be the previous node of the list
                                //If this is the first node of the list, we dont have to do anything!
                                var pre *Contact
                                if e == oldList.Front() {   //not sure
                                    //nothing to do
                                    pre = prevContact
                                } else {
                                    pre = e.Prev().Value.(*Contact)
                                }

                                if pre != nil {
                                    //store the value to previous contact
                                    ipPort := strings.Join([]string{pre.Host.String(), ":", strconv.FormatUint(uint64(pre.Port),10)},"")
                                    strReq := StoreRequest{kad.MyContact, NewRandomID(), keyId, fvResult.Value}
                                    var strRes StoreResult
                                    client, err = rpc.DialHTTP("tcp", ipPort)
                                    if err != nil {
                                        log.Fatal("Call ", err)
                                    }
                                    err = client.Call("Kademlia.Store", strReq, &strRes)
                                    if err != nil {
                                        log.Fatal("Call: ", err)
                                    }
                                }
                                break
                            }

                            //trim the list to closest k nodes
                            if shortList.Len() > k {
                               tmpList := list.New()
                               i=0
                               for e:= shortList.Front(); e != nil && i < k; e= e.Next() {
                            	  tmpList.PushBack(e.Value.(*Contact))
                            	  i++
                               }
                               shortList = tmpList
                            }

                        }
                    }

                }
                break
            }
        }
    } else {   //else value is found and in this node, so return the value
		result.Value = value    //set the appropriate value
		result.Nodes = make([]FoundNode, 1, 1)    //make the found node array with 1 element
		result.Nodes[0].NodeID = kad.MyContact.NodeID    //set the node id

		//print the value and id
		fmt.Println(kad.MyContact.NodeID.AsString(), " ", string(value))
		return nil
    }
    fmt.Println("ERR")
    return nil
}


//parallel function for find node
func ParallelFindValue(e *list.Element, i int, fvRequest FindValueRequest, fvResult *FindValueResult, kad *Kademlia) {
    c := e.Value.(*Contact)
    currentList := <-clChan
    activeList := <-alChan
    shortList := <-slChan

    if c.Port != 0 {
        currentList.PushFront(c)
        ipPort := strings.Join([]string{c.Host.String(), ":", strconv.FormatUint(uint64(c.Port),10)},"")
        client, err := rpc.DialHTTP("tcp", ipPort)
        if err != nil {
            shortList.Remove(e)
            currentList.Remove(e)
            alChan <-activeList
            clChan <-currentList
            slChan <-shortList
            CallReturnValue(i, nil)
            return
        }
        err = client.Call("Kademlia.FindValue", fvRequest, fvResult)
        if err != nil {
            //log.Fatal("Err: ", err)
            shortList.Remove(e)
            currentList.Remove(e)
            alChan <-activeList
            clChan <-currentList
            slChan <-shortList
            CallReturnValue(i, nil)
            return
        }

        //update the bucket of the sender
        Update(kad, c)
        activeList.PushFront(c)
        currentList.Remove(e)
    }
    //if no error found, check if any value is recieved
    //if any vale is recieved, we must break the iteration and return the value
    var res VALUE
    if fvResult.Value != nil {
        res = fvResult.Value
    } else {
        //process the newly found nodes
        for i := 0; i < len(fvResult.Nodes); i++ {
            c := Contact {fvResult.Nodes[i].NodeID, net.ParseIP(fvResult.Nodes[i].IPAddr), fvResult.Nodes[i].Port}
            //check
            if !ContainsInList(shortList, &c) {
    		    shortList.PushBack(&c)
            }
        }
    }
    slChan <-shortList
    alChan <-activeList
    clChan <-currentList
    CallReturnValue(i, res)
}

func CallReturnValue(i int, val VALUE) {
    //writing to the appropriate channel for function finish signaling
    if i==0 {
        fv1 <-val
    } else if i==1 {
        fv2 <-val
    }else if i==2 {
       fv3 <-val
    }
}

// Ping a nodeID
// Perform a ping rpc

func CallPings(kad *Kademlia, nodeID string) (result bool) {

  	if nodeID, err := FromString(nodeID); err == nil {

 	contact := HasContact(kad,nodeID)

 	if nodeID == kad.MyContact.NodeID{

  		contact = &(kad.MyContact)
  	}

	if contact == nil{
	   result = false
	   return
	}
	//fmt.Println (nodeID)
	ping := new (Ping)
	var pong Pong
	ping.MsgID = NewRandomID()
	peer := fmt.Sprintf("%s:%d",contact.Host, contact.Port)
	//fmt.Println(peer)
	client, err := rpc.DialHTTP("tcp", peer)
	if err != nil{
		return false
	}

	err = client.Call("Kademlia.Ping",ping, &pong)

	  if err != nil{
		 result=false
	     return
	  }


	if ping.MsgID.Equals(pong.MsgID){
		fmt.Println("Ping success")

		//fmt.Println(pong.Sender.NodeID.AsString(), pong.Sender.Host)
		Update(kad, &pong.Sender)
	return true
	}


	}
    result=false
	return
}


// Store (key, value) in a nodeID
// Use Store rpc


func CallStore (kad *Kademlia, nodeID, key, value string) bool {

	if nodeID, err := FromString(nodeID); err == nil{

		if key, err := FromString(key); err == nil{

		contact := HasContact(kad, nodeID)

		if contact == nil{
		  return false
		}

		storereq := new (StoreRequest)
		var storeres StoreResult
		storereq.Sender = kad.MyContact
		storereq.MsgID = NewRandomID()
		storereq.Key = key
		storereq.Value = []byte(value)

	 	peer := fmt.Sprintf("%s:%d",contact.Host, contact.Port)
		fmt.Println("Storing ", value, "in ", peer)
		client, err := rpc.DialHTTP("tcp", peer)
		if err != nil{
			return false
		}

		err = client.Call("Kademlia.Store",storereq, &storeres)

		if err != nil{
			fmt.Println("Error in client calling...")
			return false
		}
		//fmt.Println("Successful Store")
		fmt.Println("")
		Update(kad, contact)
				return true

		}

		return false
		}

    return false

}


// Find Node for a given key using rpc to nodeID


func CallFindNode(kad *Kademlia, nodeID, key string)bool {

	if nodeID, err := FromString(nodeID); err == nil{

		if key, err := FromString(key); err == nil{

		contact := HasContact(kad, nodeID)

		if contact == nil{
		return false
		}

		findnodereq := new (FindNodeRequest)
		var findnoderes FindNodeResult
		findnodereq.MsgID = NewRandomID()
		findnodereq.NodeID = key
		findnodereq.Sender = kad.MyContact

	 	peer := fmt.Sprintf("%s:%d",contact.Host, contact.Port)
		fmt.Println("Finding ", key, "in ", peer)

		client, err := rpc.DialHTTP("tcp", peer)
		if err != nil{
			return false
		}

		err = client.Call("Kademlia.FindNode",&findnodereq, &findnoderes)

		if err != nil{

			return false
		}

		fmt.Println("Successful RPC")
		Update(kad, contact)


		 for i:=0; i<len(findnoderes.Nodes);i=i+1{

			node := findnoderes.Nodes[i]
			if node.NodeID.AsString() != "0" {

			peer := fmt.Sprintf("%s:%d", node.IPAddr, node.Port)

			tcpAdd, err1 := net.ResolveTCPAddr("tcp4", peer)

			if err1 != nil{
				continue
			}

			host := tcpAdd.IP
			port := uint16(tcpAdd.Port)
			//fmt.Println(node.NodeID, host, port)

				contact := Contact{node.NodeID, host, port}
				Update(kad, &contact)
				fmt.Println(node.IPAddr, node.Port, node.NodeID.AsString())
			}
		}

		return true
		}
		}
	 	fmt.Println("Err")
		return false


}
// Find Value
// Perform a find_value. If it returns nodes, print them as for find_node.
// If it returns a value, print the value as in iterativeFindValue

func CallFindValue(kad *Kademlia, nodeID, key string) bool{
	if nodeID, err := FromString(nodeID); err == nil{

		if key, err := FromString(key); err == nil{


			contact := HasContact(kad, nodeID)

		if contact == nil{
		return false
		}

		findvaluereq := new(FindValueRequest)
		var findvalueres FindValueResult
		findvaluereq.MsgID = NewRandomID()
		findvaluereq.Key = key
		findvaluereq.Sender = kad.MyContact

	 	peer := fmt.Sprintf("%s:%d",contact.Host, contact.Port)
		fmt.Println("Finding ", key, "in ", peer)
		client, err := rpc.DialHTTP("tcp", peer)

		if err != nil{
			return false

		}

		err = client.Call("Kademlia.FindValue",findvaluereq, &findvalueres)

		if err != nil{

			return false
		}

		fmt.Println("Successful RPC")

		Update(kad, contact)
		if findvalueres.Value != nil{
				fmt.Println (string(findvalueres.Value))
				return true
			}

			//fmt.Println (findvalueres.Nodes)



			for i:=0; i<len(findvalueres.Nodes);i=i+1{

			node := findvalueres.Nodes[i]
			if node.NodeID.AsString() != "0" {

			peer := fmt.Sprintf("%s:%d", node.IPAddr, node.Port)

			tcpAdd, err1 := net.ResolveTCPAddr("tcp4", peer)

			if err1 != nil{
				continue
			}

			host := tcpAdd.IP
			port := uint16(tcpAdd.Port)
			fmt.Println(node.NodeID, host, port)

			contact := Contact{node.NodeID, host, port}
			Update(kad, &contact)
			fmt.Println(node.IPAddr, node.Port, node.NodeID.AsString())
			}
			}



			return true


		}
		}
	 	fmt.Println("Err")

		return false


}
//

//----------------------------------------------------------------------------------




