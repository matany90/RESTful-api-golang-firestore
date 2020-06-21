package main

import (
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
	Value       int64
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
func getPlayers(w http.ResponseWriter, r *http.Request) {
    // define players array
    var players []FootballPlayer

    // iterate over players docs
    iter := client.Collection("players").Documents(context.Background())
    for {
            doc, err := iter.Next()
            if err == iterator.Done {
                    break
            }
            if err != nil {
                    log.Fatalf("Failed to iterate: %v", err)
            }

            // add player to players array
            data := doc.Data()
            _name, _ := data["name"].(string)
            _team, _ := data["team"].(string)
            _value, _ := data["value"].(int64)
            players = append(players, FootballPlayer{Name: _name, Team: _team, Value: _value})
    }

    // returns players json
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(players)
}

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
    router.HandleFunc("/", getPlayers).Methods("GET")
    router.HandleFunc("/", createPlayer).Methods("POST")

    // listen on port 5000
    log.Fatal(http.ListenAndServe(":5000", router))
}