/*
Copyright © 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/awslabs/clencli/cobra/model"
	"github.com/awslabs/clencli/helper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// GetModelFromFlags fills the parameters onto the Unsplash Random Photo Parameters struct
func GetModelFromFlags(cmd *cobra.Command) model.UnsplashRandomPhotoParameters {
	var params model.UnsplashRandomPhotoParameters

	params.Query, _ = cmd.Flags().GetString("query")
	params.Collections, _ = cmd.Flags().GetString("collections")
	params.Featured, _ = cmd.Flags().GetBool("featured")
	params.Username, _ = cmd.Flags().GetString("username")
	params.Orientation, _ = cmd.Flags().GetString("orientation")
	params.Filter, _ = cmd.Flags().GetString("filter")
	params.Size, _ = cmd.Flags().GetString("size")

	return params
}

func buildURL(params model.UnsplashRandomPhotoParameters, cred model.Credential) string {
	clientID := cred.AccessKey
	url := fmt.Sprintf("https://api.unsplash.com/photos/random?client_id=%s", clientID)

	if len(params.Collections) > 0 {
		url += fmt.Sprintf("&collections=%s", params.Collections)
	}

	if len(params.Query) > 0 {
		url += fmt.Sprintf("&query=%s", params.Query)
	}

	url += fmt.Sprintf("&featured=%t", params.Featured)

	if len(params.Username) > 0 {
		url += fmt.Sprintf("&username=%s", params.Username)
	}

	if len(params.Orientation) > 0 {
		url += fmt.Sprintf("&orientation=%s", params.Orientation)
	}

	if len(params.Filter) > 0 {
		url += fmt.Sprintf("&filter=%s", params.Filter)
	}

	return url
}

// DownloadPhoto downloads a photo and saves into downloads/unsplash/ folder
// It creates the downloads/ folder if it doesn't exists
func DownloadPhoto(params model.UnsplashRandomPhotoParameters, cred model.Credential, photoSizes []string) error {
	response, err := RequestRandomPhoto(params, cred)
	if err != nil {
		return err
	}

	dumpUnsplashRandomPhotoResponse(response)

	dirPath, err := helper.CreateDirectoryNamedPath("downloads/unsplash/" + params.Query)
	if err != nil {
		return err
	}

	if params.Size == "all" {
		for _, pSize := range photoSizes {
			if pSize != "all" {
				params.Size = pSize
				err = helper.DownloadFile(getPhotoURLBySize(params, response), dirPath, response.ID+"-"+pSize+".jpeg")
				if err != nil {
					return err
				}
			}
		}
	} else {
		err = helper.DownloadFile(getPhotoURLBySize(params, response), dirPath, response.ID+".jpeg")
	}

	return err
}

// DownloadPhotoByID TODO ...
func DownloadPhotoByID(r model.UnsplashGetPhotoResponse, size string) error {
	dirPath, err := helper.CreateDirectoryNamedPath("downloads/unsplash/" + r.ID)
	if err != nil {
		return err
	}

	if size == "raw" || size == "all" {
		err = helper.DownloadFile(r.Urls.Raw, dirPath, r.ID+"-raw.jpeg")
		if err != nil {
			return fmt.Errorf("unable to download photo\n%v", err)
		}
	}

	if size == "full" || size == "all" {
		err = helper.DownloadFile(r.Urls.Full, dirPath, r.ID+"-full.jpeg")
		if err != nil {
			return fmt.Errorf("unable to download photo\n%v", err)
		}
	}

	if size == "regular" || size == "all" {
		err = helper.DownloadFile(r.Urls.Regular, dirPath, r.ID+"-regular.jpeg")
		if err != nil {
			return fmt.Errorf("unable to download photo\n%v", err)
		}
	}

	if size == "small" || size == "all" {
		err = helper.DownloadFile(r.Urls.Small, dirPath, r.ID+"-small.jpeg")
		if err != nil {
			return fmt.Errorf("unable to download photo\n%v", err)
		}
	}

	if size == "thumb" || size == "all" {
		err = helper.DownloadFile(r.Urls.Thumb, dirPath, r.ID+"-thumb.jpeg")
		if err != nil {
			return fmt.Errorf("unable to download photo\n%v", err)
		}
	}

	return err

}

// GetPhotoURLBySize return the photo URL based on the given size
func getPhotoURLBySize(p model.UnsplashRandomPhotoParameters, r model.UnsplashRandomPhotoResponse) string {
	switch p.Size {
	case "thumb":
		return r.Urls.Thumb
	case "small":
		return r.Urls.Small
	case "regular":
		return r.Urls.Regular
	case "full":
		return r.Urls.Full
	case "raw":
		return r.Urls.Raw
	default:
		return r.Urls.Small
	}
}

// RequestRandomPhoto retrieves a single random photo, given optional filters.
func RequestRandomPhoto(params model.UnsplashRandomPhotoParameters, cred model.Credential) (model.UnsplashRandomPhotoResponse, error) {
	var response model.UnsplashRandomPhotoResponse
	url := buildURL(params, cred)

	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return response, fmt.Errorf("unexpected error while performing GET on Unsplash API \n%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return response, fmt.Errorf("unexpected error while reading Unsplash response \n%v", err)
		}

		json.Unmarshal(bodyBytes, &response)
	}

	return response, err
}

// GetPhoto TODO..
func GetPhoto(id string, cred model.Credential) (model.UnsplashGetPhotoResponse, error) {
	var response model.UnsplashGetPhotoResponse

	clientID := cred.AccessKey
	url := fmt.Sprintf("https://api.unsplash.com/photos/%s?client_id=%s", id, clientID)

	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return response, fmt.Errorf("unexpected error while performing GET on Unsplash API \n%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return response, fmt.Errorf("unexpected error while reading Unsplash response \n%v", err)
		}

		json.Unmarshal(bodyBytes, &response)
	}

	return response, err
}

// SaveGetPhotoResult TODO ...
func SaveGetPhotoResult(r model.UnsplashGetPhotoResponse) {
	d, err := yaml.Marshal(r)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	helper.WriteFile("unsplash.yaml", d)
}

func dumpUnsplashRandomPhotoResponse(r model.UnsplashRandomPhotoResponse) {
	d, err := yaml.Marshal(r)
	if err != nil {
		logrus.Fatalf("error: %v", err)
	}
	helper.WriteFile("unsplash.yaml", d)
}
