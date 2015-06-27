package main

import(
		"os"
		"fmt"
		"reflect"
		"strings"
		"strconv"
		"net/http"
		"io/ioutil"
		"encoding/json"
		"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
		"github.com/gorilla/mux"
		
)

type imageLinks struct {
	Link string 				`bson:"link"`
	Upvote int 					`bson:"upvote"`
	Downvote int 				`bson:"downvote"`
}

type Result struct {
		Photos struct {
		Page int 				`json: "page"`
		Pages int 				`json: "pages"`
		PerPage int 			`json: "perpage"`
		Total string 			`json: "total"`
		Photo []struct {
				Id string 		`json: "id"`
				Owner string 	`json: "owner"`
				Secret string 	`json: "secret"`
				Server string 	`json: "server"`
				Farm int 		`json: "farm"`
				Title string 	`json: "title"`
				IsPublic int 	`json: "ispublic"`
				IsFriend int 	`json: "isfriend"`
				IsFamily int 	`json: "isfamily`
				} 				`json: "photo"`
		} 						`json: "photos"`
	Stat string 				`json: "stat"`
}

func main() {
	
	router := mux.NewRouter().StrictSlash(true);
	router.HandleFunc("/Index", Index)
	router.HandleFunc("/UpVote",UpVoteRoute)
	router.HandleFunc("/DownVote",DownVoteRoute)
	//log.Fatal(http.ListenAndServe(":8080", router))
}

func UpVoteRoute(w http.ResponseWriter, r *http.Request){

  link := r.URL.Query().Get("imagelink")
  w.Header().Set("Content-Type", "application/json");
  session, err := mgo.Dial("mongodb://krithghosh:mercury29@ds061631.mongolab.com:61631/flickrimagedb")
	if err != nil {
		fmt.Printf("%s", err)
        os.Exit(1)
	}
  defer session.Close()
 
  session.SetMode(mgo.Monotonic, true)
  c := session.DB("flickrimagedb").C("image_links_votes")

  err = c.Update(bson.M{"link": link}, bson.M{"$inc": bson.M{"upvote": 1}})
  if err != nil {
    fmt.Printf("Can't update document %v\n", err)
    os.Exit(1)
  }
}

func DownVoteRoute(w http.ResponseWriter, r *http.Request){

  link := r.URL.Query().Get("imagelink")
  w.Header().Set("Content-Type", "application/json");
  session, err := mgo.Dial("mongodb://krithghosh:mercury29@ds061631.mongolab.com:61631/flickrimagedb")
	if err != nil {
		fmt.Printf("%s", err)
        os.Exit(1)
	}
  defer session.Close()
 
  session.SetMode(mgo.Monotonic, true)
  c := session.DB("flickrimagedb").C("image_links_votes")

  err = c.Update(bson.M{"link": link}, bson.M{"$inc": bson.M{"downvote": 1}})
  if err != nil {
    fmt.Printf("Can't update document %v\n", err)
    os.Exit(1)
  }
}

func Index(w http.ResponseWriter, r *http.Request) {

  w.Header().Set("Content-Type", "application/json");
	session, err := mgo.Dial("mongodb://krithghosh:mercury29@ds061631.mongolab.com:61631/flickrimagedb")
	if err != nil {
		fmt.Printf("%s", err)
        os.Exit(1)
	}
 
	defer session.Close()
 
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("flickrimagedb").C("image_links_votes")

     checkResult := &imageLinks{}
     // Create a slice to begin with
     myType := reflect.TypeOf(checkResult)
     slice := reflect.MakeSlice(reflect.SliceOf(myType), 10, 10)
     // Create a pointer to a slice value and set it to the slice
     x := reflect.New(slice.Type())
     x.Elem().Set(slice)
     err = c.Find(bson.M{}).All(x.Interface())
     if err != nil {
        response, err := json.Marshal(x.Interface())
        if err != nil{
            fmt.Printf("%s", err)
            os.Exit(2)
         }
        fmt.Fprintf(w, string(response))
     } else {
        url := "https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=99d139e54d5c9fb2e7b766901a6e7420&text=cute+puppies&per_page=12&format=json&nojsoncallback=1"
        res, err := http.Get(url)
         if err != nil{
            fmt.Printf("%s", err)
            os.Exit(3)
         }
       	body, err := ioutil.ReadAll(res.Body)
        if err != nil{
            fmt.Printf("%s", err)
            os.Exit(4)
         }

         jsonData := &Result{}
         err = json.Unmarshal(body, jsonData)

         for value := range jsonData.Photos.Photo{
         		s1 :=  []string{"https://farm", ".staticflickr.com/"}
         		s2 :=  []string{strings.Join(s1, strconv.Itoa(jsonData.Photos.Photo[value].Farm)), "/"}
         		s3 :=  []string{strings.Join(s2, jsonData.Photos.Photo[value].Server), "_"}
         		s4 :=  []string{strings.Join(s3, jsonData.Photos.Photo[value].Id), ".jpg"}
         		s := strings.Join(s4, jsonData.Photos.Photo[value].Secret)
              	singleReuslt := imageLinks{}
              	err = c.Find(bson.M{"link": s}).One(&singleReuslt)
              	if err != nil {
                    err = c.Insert(&imageLinks{Link: s, Upvote: 0, Downvote: 0})
                    if err != nil {
    				fmt.Printf("%s", err)
            		os.Exit(5)
    				}
            	}
        }

         allResult := &imageLinks{}
         // Create a slice to begin with
    	   myType := reflect.TypeOf(allResult)
    	   slice := reflect.MakeSlice(reflect.SliceOf(myType), 10, 10)
    	   // Create a pointer to a slice value and set it to the slice
    	   x := reflect.New(slice.Type())
    	   x.Elem().Set(slice)
         err = c.Find(bson.M{}).All(x.Interface())
         response, err := json.Marshal(x.Interface())
         if err != nil{
            fmt.Printf("%s", err)
            os.Exit(6)
         }
         fmt.Fprintf(w, string(response))
      }
}