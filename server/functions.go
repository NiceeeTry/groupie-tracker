package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Index struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Loc struct {
	Ind []Index `json:"index"`
}

type Relation struct {
	DatesLocation map[string][]string `json:"datesLocations"`
}

type Everything struct {
	Everyone []Artist
	Location Loc
}

type ArtistInfo struct {
	Artist
	Relation
}

func GetAllArtists() ([]Artist, error) {
	var artists []Artist
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return artists, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return artists, err
	}
	err = json.Unmarshal(bytes, &artists)
	if err != nil {
		return artists, err
	}
	defer response.Body.Close()
	return artists, nil
}

func GetAllLocations() (Loc, error) {
	var location Loc
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		return location, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return location, err
	}
	err = json.Unmarshal(bytes, &location)
	if err != nil {
		return location, err
	}
	defer response.Body.Close()
	return location, nil
}

func OneArtist(id int) (Artist, error) {
	var artist Artist
	response, err := http.Get(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%d", id))
	if err != nil {
		return artist, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return artist, err
	}

	err = json.Unmarshal(bytes, &artist)
	if err != nil {
		return artist, err
	}
	return artist, nil
}

func GetLocation(id int) (Loc, error) {
	var location Loc // see the structure
	response, err := http.Get(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%d", id))
	if err != nil {
		return location, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return location, err
	}
	// fmt.Println(string(bytes))

	err = json.Unmarshal(bytes, &location)
	// fmt.Println(location.DatesLocation)
	if err != nil {
		return location, err
	}
	return location, nil
}

func Relations(id int) (Relation, error) {
	var rel Relation
	response, err := http.Get(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", id))
	if err != nil {
		return rel, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return rel, err
	}
	// fmt.Println(string(bytes))

	err = json.Unmarshal(bytes, &rel)
	// fmt.Println(rel.DatesLocation)
	if err != nil {
		return rel, err
	}
	return rel, nil
}

func Search(data Everything, searchTerm string) Everything {
	var output Everything
	var artists []Artist
	for _, result := range data.Everyone {
		if strings.Contains(strings.ToLower(result.Name), strings.ToLower(searchTerm)) || strings.Contains(strings.ToLower(result.FirstAlbum), strings.ToLower(searchTerm)) || strings.Contains(strings.ToLower(strconv.Itoa(result.CreationDate)), strings.ToLower(searchTerm)) {
			// fmt.Println(result)
			artists = append(artists, result)
		}
		// else if len(artists) == 0 {
		// 	for _, members := range result.Members {
		// 		if strings.Contains(strings.ToLower(members), strings.ToLower(searchTerm)) {
		// 			// fmt.Println("Members")
		// 			// fmt.Println(result)
		// 			artists = append(artists, members)
		// 		}
		// 	}
		// }
		for _, members := range result.Members {
			if strings.Contains(strings.ToLower(members), strings.ToLower(searchTerm)) {
				// fmt.Println("Members")
				// fmt.Println(result)
				artists = append(artists, result)
			}
		}
	}
	// if len(artists) == 0 {
	for _, result := range data.Location.Ind {
		for _, location := range result.Locations {
			if strings.Contains(strings.ToLower(location), strings.ToLower(searchTerm)) {
				// fmt.Println("location")
				// fmt.Println(location)
				art, err := OneArtist(result.ID)
				if err != nil {
					fmt.Println("error in getting artist")
				}
				artists = append(artists, art)
			}
		}
	}
	// fmt.Println(artists)
	output.Everyone = artists
	return output
	//}
}

// previous
// func Search(records []Artist, term string) []Artist {
// 	var results []Artist
// 	for _, r := range records {
// 		if strings.Contains(strings.ToLower(r.Name), strings.ToLower(term)) {
// 			results = append(results, r)
// 		}
// 	}
// 	if len(results) == 0 {
// 		return nil
// 	}
// 	return results
// }
