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

func Search(data Everything, searchTerm string) (Everything, error) {
	var output Everything
	ids := make(map[int]int)
	var artists []Artist
	for _, result := range data.Everyone {
		if strings.Contains(strings.ToLower(result.Name), strings.ToLower(searchTerm)) || strings.Contains(strings.ToLower(result.FirstAlbum), strings.ToLower(searchTerm)) || strings.Contains(strings.ToLower(strconv.Itoa(result.CreationDate)), strings.ToLower(searchTerm)) {
			if _, ok := ids[result.ID]; ok {
				continue
			} else {
				artists = append(artists, result)
				ids[result.ID] += 1
			}
		}

		for _, members := range result.Members {
			if strings.Contains(strings.ToLower(members), strings.ToLower(searchTerm)) {
				if _, ok := ids[result.ID]; ok {
					continue
				} else {
					artists = append(artists, result)
					ids[result.ID] += 1
				}
			}
		}

	}
	for _, result := range data.Location.Ind {
		for _, location := range result.Locations {
			if strings.Contains(strings.ToLower(location), strings.ToLower(searchTerm)) {
				if _, ok := ids[result.ID]; ok {
					continue
				} else {
					art, err := OneArtist(result.ID)
					if err != nil {
						fmt.Println("error in getting artist")
						return output, err
					}
					artists = append(artists, art)
					ids[result.ID] += 1
				}
			}
		}
	}
	// fmt.Println(artists)
	output.Everyone = artists
	return output, nil
	//}
}

func CreationDate(data Everything, range1 string, range2 string) (Everything, error) {
	var output Everything
	year1, err := strconv.Atoi(range1)
	if err != nil {
		fmt.Println("error in years")
		return output, err
	}
	year2, err := strconv.Atoi(range2)
	if err != nil {
		fmt.Println("error in years")
		return output, err
	}
	var artists []Artist
	for _, result := range data.Everyone {
		if year1 <= result.CreationDate && year2 >= result.CreationDate {
			artists = append(artists, result)
		}
	}

	// fmt.Println(artists)
	output.Everyone = artists
	return output, nil
	//}
}

func FirstAlbumDates(data Everything, range1 string, range2 string) (Everything, error) {
	var output Everything
	year1, err := strconv.Atoi(range1)
	if err != nil {
		fmt.Println("error in years")
		return output, err
	}
	year2, err := strconv.Atoi(range2)
	if err != nil {
		fmt.Println("error in years")
		return output, err
	}
	var artists []Artist
	for _, result := range data.Everyone {
		data := strings.Split(result.FirstAlbum, "-")
		dataNumber, err := strconv.Atoi(data[len(data)-1])
		if err != nil {
			fmt.Println("error in years")
			return output, err
		}
		if year1 <= dataNumber && year2 >= dataNumber {
			artists = append(artists, result)
		}
	}

	// fmt.Println(artists)
	output.Everyone = artists
	return output, nil
	//}
}

func MembersNumber(data Everything, arr []string) (Everything, error) {
	var output Everything
	var memberNum []int
	for _, val := range arr {
		n, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("error in years")
			return output, err
		}
		memberNum = append(memberNum, n)
	}
	var artists []Artist
	for _, result := range data.Everyone {
		for _, val := range memberNum {
			if len(result.Members) == val {
				artists = append(artists, result)
			}
		}
	}

	// fmt.Println(artists)
	output.Everyone = artists
	return output, nil
	//}
}

func LocationSearch(data Everything, search string) (Everything, error) {
	var output Everything
	var artists []Artist
	ids := make(map[int]int)
	for _, result := range data.Location.Ind {
		for _, location := range result.Locations {
			if strings.Contains(strings.ToLower(location), strings.ToLower(search)) {
				if _, ok := ids[result.ID]; ok {
					continue
				} else {
					art, err := OneArtist(result.ID)
					if err != nil {
						fmt.Println("error in getting artist")
						return output, err
					}
					artists = append(artists, art)
					ids[result.ID] += 1
				}
			}
		}
	}
	output.Everyone = artists
	return output, nil
	//}
}
