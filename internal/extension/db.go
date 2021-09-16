package extension

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/rs/zerolog/log"
)

func (s *Server) updateReverseGeo(db *sql.DB, key int64) error {

	pipeline := s.Pi.Name

	// get long/lat using key
	selectSQL := fmt.Sprintf("select latitude, longitude from %s.starlink where key = $1", pipeline)

	row := db.QueryRow(selectSQL, key)
	var long float64
	var lat float64
	err := row.Scan(&lat, &long)
	if err != nil {
		log.Error().Stack().Err(err).Msg("error getting row")
		return err
	}

	longString := fmt.Sprintf("%f", long)
	latString := fmt.Sprintf("%f", lat)

	result, err := reverseGeo(longString, latString)
	if err != nil {
		log.Error().Stack().Err(err).Msg("error getting reverse value")
		return err
	}

	//update the row with the reverse geo value
	updateSQL := fmt.Sprintf("update %s.starlink set (reverselocation, lastupdated) = ($1, now()) where key = $2", pipeline)

	_, err = db.Exec(updateSQL, result, key)
	if err != nil {
		return err
	}

	return nil
}

func reverseGeo(long string, lat string) (result string, err error) {
	url := fmt.Sprintf("https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=%s&longitude=%s&localityLanguage=en", lat, long)
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return result, err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}

		log.Info().Msg(fmt.Sprintf("read %d \n", len(body)))

		obj, parseError := oj.ParseString(string(body))
		if parseError != nil {
			log.Error().Stack().Err(parseError).Msg("error in oj.ParseString")
			return result, parseError
		}

		localityX, err := jp.ParseString("$.locality")
		if err != nil {
			return result, err
		}

		var tmp []interface{}
		var locality string
		tmp = localityX.Get(obj)
		locality = tmp[0].(string)
		log.Info().Msg("locality here " + locality)

		cityX, err := jp.ParseString("$.city")
		if err != nil {
			return result, err
		}
		var city string
		tmp = cityX.Get(obj)
		city = tmp[0].(string)
		log.Info().Msg("city here " + city)

		postcodeX, err := jp.ParseString("$.postcode")
		if err != nil {
			return result, err
		}
		var postcode string
		tmp = postcodeX.Get(obj)
		postcode = tmp[0].(string)
		log.Info().Msg("postcode here " + postcode)

		//fmt.Println(string(body))
		result = fmt.Sprintf("city: %s locality: %s postcode: %s\n", city, locality, postcode)
		return result, nil
	}

	return result, fmt.Errorf("invalid http status code %d\n", resp.StatusCode)
}
