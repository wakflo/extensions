// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package googledrive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"slices"

	fastshot "github.com/opus-domini/fast-shot"
	"google.golang.org/api/drive/v3"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://oauth2.googleapis.com/token"
	sharedAuth = autoform.NewOAuthField("https://accounts.google.com/o/oauth2/auth", &tokenURL, []string{
		"https://www.googleapis.com/auth/drive",
	}).Build()
)

// File The metadata for a file. Some resource methods (such as
// `files.update`) require a `fileId`. Use the `files.list` method to
// retrieve the ID for a file.
type File struct {
	// CreatedTime: The time at which the file was created (RFC 3339
	// date-time).
	CreatedTime string `json:"createdTime,omitempty"`

	// Description: A short description of the file.
	Description string `json:"description,omitempty"`

	// DriveID: Output only. ID of the shared drive the file resides in.
	// Only populated for items in shared drives.
	DriveID string `json:"driveId,omitempty"`

	// ExplicitlyTrashed: Output only. Whether the file has been explicitly
	// trashed, as opposed to recursively trashed from a parent folder.
	ExplicitlyTrashed bool `json:"explicitlyTrashed,omitempty"`

	// ExportLinks: Output only. Links for exporting Docs Editors files to
	// specific formats.
	ExportLinks map[string]string `json:"exportLinks,omitempty"`

	// FileExtension: Output only. The final component of
	// `fullFileExtension`. This is only available for files with binary
	// content in Google Drive.
	FileExtension string `json:"fileExtension,omitempty"`

	// FullFileExtension: Output only. The full file extension extracted
	// from the `name` field. May contain multiple concatenated extensions,
	// such as "tar.gz". This is only available for files with binary
	// content in Google Drive. This is automatically updated when the
	// `name` field changes, however it is not cleared if the new name does
	// not contain a valid extension.
	FullFileExtension string `json:"fullFileExtension,omitempty"`

	// HasAugmentedPermissions: Output only. Whether there are permissions
	// directly on this file. This field is only populated for items in
	// shared drives.
	HasAugmentedPermissions bool `json:"hasAugmentedPermissions,omitempty"`

	// HasThumbnail: Output only. Whether this file has a thumbnail. This
	// does not indicate whether the requesting app has access to the
	// thumbnail. To check access, look for the presence of the
	// thumbnailLink field.
	HasThumbnail bool `json:"hasThumbnail,omitempty"`

	// ID: The ID of the file.
	ID string `json:"id,omitempty"`

	// Kind: Output only. Identifies what kind of resource this is. Value:
	// the fixed string "drive#file".
	Kind string `json:"kind,omitempty"`

	// Md5Checksum: Output only. The MD5 checksum for the content of the
	// file. This is only applicable to files with binary content in Google
	// Drive.
	Md5Checksum string `json:"md5Checksum,omitempty"`

	// MimeType: The MIME type of the file. Google Drive attempts to
	// automatically detect an appropriate value from uploaded content, if
	// no value is provided. The value cannot be changed unless a new
	// revision is uploaded. If a file is created with a Google Doc MIME
	// type, the uploaded content is imported, if possible. The supported
	// import formats are published in the About resource.
	MimeType string `json:"mimeType,omitempty"`

	// ModifiedByMe: Output only. Whether the file has been modified by this
	// user.
	ModifiedByMe bool `json:"modifiedByMe,omitempty"`

	// ModifiedByMeTime: The last time the file was modified by the user
	// (RFC 3339 date-time).
	ModifiedByMeTime string `json:"modifiedByMeTime,omitempty"`

	// ModifiedTime: he last time the file was modified by anyone (RFC 3339
	// date-time). Note that setting modifiedTime will also update
	// modifiedByMeTime for the user.
	ModifiedTime string `json:"modifiedTime,omitempty"`

	// Name: The name of the file. This is not necessarily unique within a
	// folder. Note that for immutable items such as the top level folders
	// of shared drives, My Drive root folder, and Application Data folder
	// the name is constant.
	Name string `json:"name,omitempty"`

	// OriginalFilename: The original filename of the uploaded content if
	// available, or else the original value of the `name` field. This is
	// only available for files with binary content in Google Drive.
	OriginalFilename string `json:"originalFilename,omitempty"`

	// Size: Output only. Size in bytes of blobs and first party editor
	// files. Won't be populated for files that have no size, like shortcuts
	// and folders.
	Size int64 `json:"size,omitempty,string"`

	// Spaces: Output only. The list of spaces which contain the file. The
	// currently supported values are 'drive', 'appDataFolder' and 'photos'.
	Spaces []string `json:"spaces,omitempty"`

	// Starred: Whether the user has starred the file.
	Starred bool `json:"starred,omitempty"`

	// TeamDriveID: Deprecated: Output only. Use `driveId` instead.
	TeamDriveID string `json:"teamDriveId,omitempty"`

	// Trashed: Whether the file has been trashed, either explicitly or from
	// a trashed parent folder. Only the owner may trash a file, and other
	// users cannot see files in the owner's trash.
	Trashed bool `json:"trashed,omitempty"`

	// TrashedTime: The time that the item was trashed (RFC 3339 date-time).
	// Only populated for items in shared drives.
	TrashedTime string `json:"trashedTime,omitempty"`

	// Version: Output only. A monotonically increasing version number for
	// the file. This reflects every change made to the file on the server,
	// even those not visible to the user.
	Version int64 `json:"version,omitempty,string"`

	// ViewedByMe: Output only. Whether the file has been viewed by this
	// user.
	ViewedByMe bool `json:"viewedByMe,omitempty"`

	// ViewedByMeTime: The last time the file was viewed by the user (RFC
	// 3339 date-time).
	ViewedByMeTime string `json:"viewedByMeTime,omitempty"`

	// ViewersCanCopyContent: Deprecated: Use `copyRequiresWriterPermission`
	// instead.
	ViewersCanCopyContent bool `json:"viewersCanCopyContent,omitempty"`

	// WebContentLink: Output only. A link for downloading the content of
	// the file in a browser. This is only available for files with binary
	// content in Google Drive.
	WebContentLink string `json:"webContentLink,omitempty"`

	// WebViewLink: Output only. A link for opening the file in a relevant
	// Google editor or viewer in a browser.
	WebViewLink string  `json:"webViewLink,omitempty"`
	FileData    *string `json:"fileData"`
}

type listFile struct {
	drive.File
	FileData string `json:"fileData"`
}

type ListFileResponse struct {
	Files            []listFile `json:"files"`
	IncompleteSearch bool       `json:"incompleteSearch"`
	Kind             string     `json:"kind"`
}

var googleType = []string{
	"application/vnd.google-apps.document",
	"application/vnd.google-apps.spreadsheet",
	"application/vnd.google-apps.presentation",
}

func handleFileContent(ctx *sdk.RunContext, files []*drive.File, driveService *drive.Service) ([]File, error) {
	outputs := make([]File, len(files))

	for i, file := range files {
		buf := new(bytes.Buffer)
		if slices.Contains(googleType, file.MimeType) {
			rsp, err := driveService.Files.Export(file.Id, file.MimeType).Download()
			if err != nil {
				return nil, err
			}
			defer rsp.Body.Close()
			_, err = buf.ReadFrom(rsp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			rsp, err := driveService.Files.Get(file.Id).Download()
			if err != nil {
				return nil, err
			}

			defer rsp.Body.Close()
			_, err = buf.ReadFrom(rsp.Body)
			if err != nil {
				return nil, err
			}
		}

		ext, err := mime.ExtensionsByType(file.MimeType)
		if err != nil {
			return nil, err
		}

		fileName := fmt.Sprintf("%s.%s", file.Name, ext)
		if file.Name != "" {
			fileName = file.Name
			// fileName = strings.Replace()
		}

		fileData, err := ctx.Files.PutWorkflow(ctx.Metadata, fileName, buf)
		if err != nil {
			return nil, err
		}

		out := File{
			ID:                file.Id,
			MimeType:          file.MimeType,
			Kind:              file.Kind,
			Name:              file.Name,
			Version:           file.Version,
			Description:       file.Description,
			CreatedTime:       file.CreatedTime,
			DriveID:           file.DriveId,
			Trashed:           file.Trashed,
			FileExtension:     file.FileExtension,
			FullFileExtension: file.FullFileExtension,
			Size:              file.Size,
			OriginalFilename:  file.OriginalFilename,
			WebViewLink:       file.WebViewLink,
			FileData:          fileData,
		}
		outputs[i] = out
	}

	return outputs, nil
}

func downloadFile(ctx *sdk.RunContext, driveService *drive.Service, fileID string, fileName *string) (*string, error) {
	file, err := driveService.Files.Get(fileID).Do()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	var rsp *http.Response
	if slices.Contains(googleType, file.MimeType) {
		rsp, err = driveService.Files.Export(file.Id, file.MimeType).Download()
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()
		_, err = buf.ReadFrom(rsp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		rsp, err = driveService.Files.Get(file.Id).Download()
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()
		_, err = buf.ReadFrom(rsp.Body)
		if err != nil {
			return nil, err
		}
	}

	ext, err := mime.ExtensionsByType(file.MimeType)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("%s.%s", file.Name, ext)
	if fileName != nil {
		name = fmt.Sprintf("%s.%s", *fileName, ext)
	}

	return ctx.Files.PutWorkflow(ctx.Metadata, name, buf)
}

func getParentFoldersInput() *sdkcore.AutoFormSchema {
	getParentFolders := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		client := fastshot.NewClient("https://www.googleapis.com/drive/v3").
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/files").Query().
			AddParams(map[string]string{
				"q": "mimeType='application/vnd.google-apps.folder' and trashed = false",
			}).Send()
		if err != nil {
			return nil, err
		}

		if rsp.IsError() {
			return nil, errors.New(rsp.StatusText())
		}

		bytes, err := io.ReadAll(rsp.RawBody())
		if err != nil {
			return nil, err
		}

		var body ListFileResponse
		err = json.Unmarshal(bytes, &body)
		if err != nil {
			return nil, err
		}

		return body.Files, nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Parent Folder").
		SetDescription("select parent folder").
		SetDynamicOptions(&getParentFolders).
		SetDependsOn([]string{"connection"}).
		SetRequired(false).Build()
}

func getFoldersInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getParentFolders := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		client := fastshot.NewClient("https://www.googleapis.com/drive/v3").
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/files").Query().
			AddParams(map[string]string{
				"q": "mimeType='application/vnd.google-apps.folder' and trashed = false",
			}).Send()
		if err != nil {
			return nil, err
		}

		if rsp.IsError() {
			return nil, errors.New(rsp.StatusText())
		}

		bytes, err := io.ReadAll(rsp.RawBody())
		if err != nil {
			return nil, err
		}

		var body ListFileResponse
		err = json.Unmarshal(bytes, &body)
		if err != nil {
			return nil, err
		}

		return body.Files, nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getParentFolders).
		SetRequired(required).Build()
}

var includeTeamFieldInput = autoform.NewBooleanField().
	SetDisplayName("Include Team Drives").
	SetDescription("Determines if folders from Team Drives should be included in the results.").
	SetDefaultValue(false).
	Build()
