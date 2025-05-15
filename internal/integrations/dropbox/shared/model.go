package shared

type FileMove struct {
	AllowOwnershipTransfer bool   `json:"allow_ownership_transfer"`
	Autorename             bool   `json:"autorename"`
	FromPath               string `json:"from_path"`
	ToPath                 string `json:"to_path"`
}

type FileMoveMetadata struct {
	Metadata struct {
		Tag            string `json:".tag"`
		ClientModified string `json:"client_modified"`
		ContentHash    string `json:"content_hash"`
		FileLockInfo   struct {
			Created        string `json:"created"`
			IsLockholder   bool   `json:"is_lockholder"`
			LockholderName string `json:"lockholder_name"`
		} `json:"file_lock_info"`
		HasExplicitSharedMembers bool   `json:"has_explicit_shared_members"`
		ID                       string `json:"id"`
		IsDownloadable           bool   `json:"is_downloadable"`
		Name                     string `json:"name"`
		PathDisplay              string `json:"path_display"`
		PathLower                string `json:"path_lower"`
		PropertyGroups           []struct {
			Fields []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"fields"`
			TemplateID string `json:"template_id"`
		} `json:"property_groups"`
		Rev            string `json:"rev"`
		ServerModified string `json:"server_modified"`
		SharingInfo    struct {
			ModifiedBy           string `json:"modified_by"`
			ParentSharedFolderID string `json:"parent_shared_folder_id"`
			ReadOnly             bool   `json:"read_only"`
		} `json:"sharing_info"`
		Size int64 `json:"size"`
	} `json:"metadata"`
}
