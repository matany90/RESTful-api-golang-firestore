package main

import (
    _"fmt"
    "log"
    "net/http"
    "encoding/json"

    "golang.org/x/net/context"
    firebase "firebase.google.com/go"
    "cloud.google.com/go/firestore"
    "google.golang.org/api/option"
    "google.golang.org/api/iterator"
    "github.com/gorilla/mux"
)

// Define firestore client
var client *firestore.Client

// Define footballPlayer type
type FootballPlayer struct {
	Name        string
	Team        string
	Value       int
}

// Set firestore client
func setFirestoreClient() (*firestore.Client, error) {
    ctx := context.Background()
	opt := option.WithCredentialsFile("./service-account.json")
    app, err := firebase.NewApp(ctx, nil, opt)

    client, err = app.Firestore(ctx)

    return client, err
}

// Get all players
// func getPlayers(w http.ResponseWriter, r *http.Request) {
//     iter := client.Collection("players").Documents(context.Background())
//     for {
//             doc, err := iter.Next()
//             if err == iterator.Done {
//                     break
//             }
//             if err != nil {
//                     log.Fatalf("Failed to iterate: %v", err)
//             }
//             fmt.Println(doc.Data())
//     }
//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(FootballPlayer{Name: "test11"})
// }

// Create player
func createPlayer(w http.ResponseWriter, r *http.Request) {
    // build player struct
    var player FootballPlayer
    json.NewDecoder(r.Body).Decode(&player)

    // add to firestore
    _, _, err := client.Collection("players").Add(context.Background(), map[string]interface{}{
            "name": player.Name,
            "team": player.Team,
            "value": player.Value,
    })

    // Show error if needed
    if err != nil {
        log.Fatal(err)
    }

    // send res
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{
        "success": true,
    })
}

func main() {
    // set firestore client
    setFirestoreClient()

    // init router
    router := mux.NewRouter()

    // define endpoints
    // router.HandleFunc("/", getPlayers).Methods("GET")
    router.HandleFunc("/", createPlayer).Methods("POST")

    // listen on port 5000
    log.Fatal(http.ListenAndServe(":5000", router))
}