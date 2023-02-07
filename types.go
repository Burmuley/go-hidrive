package go_hidrive

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

// Time represents date and time information for the API.
type Time time.Time

// MarshalJSON turns Time into JSON (in Unix-time/UTC).
func (t *Time) MarshalJSON() ([]byte, error) {
	secs := time.Time(*t).Unix()
	return []byte(strconv.FormatInt(secs, 10)), nil
}

// UnmarshalJSON turns JSON into Time.
func (t *Time) UnmarshalJSON(data []byte) error {
	secs, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Time(time.Unix(secs, 0))
	return nil
}

// HiDriveObject represents HiDrive object - directory or file
type HiDriveObject struct {
	Path         string `json:"path"`
	Type         string `json:"type"`
	ID           string `json:"id"`
	ParentID     string `json:"parent_id"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	MemberCount  int64  `json:"nmembers"`
	MTime        Time   `json:"mtime"`
	CTime        Time   `json:"ctime"`
	MetaHash     string `json:"mhash"`
	MetaOnlyHash string `json:"mohash"`
	NHash        string `json:"nhash"`
	CHash        string `json:"chash"`
	Teamfolder   bool   `json:"teamfolder"`
	Readable     bool   `json:"readable"`
	Writable     bool   `json:"writable"`
	Shareable    bool   `json:"shareable"`
	MIMEType     string `json:"mime_type"`
}

func (h *HiDriveObject) UnmarshalJSON(b []byte) error {
	type HiDriveObjectAlias HiDriveObject
	defaultObject := HiDriveObjectAlias{
		Size:        -1,
		MemberCount: -1,
	}

	err := json.Unmarshal(b, &defaultObject)
	if err != nil {
		return err
	}
	name, err := url.PathUnescape(defaultObject.Name)
	if err == nil {
		defaultObject.Name = name
	}

	*h = HiDriveObject(defaultObject)
	return nil
}
