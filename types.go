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

// Object represents HiDrive object - directory or file
type Object struct {
	Path         string         `json:"path"`
	Type         string         `json:"type"`
	ID           string         `json:"id"`
	ParentID     string         `json:"parent_id"`
	Name         string         `json:"name"`
	Size         int64          `json:"size"`
	MemberCount  int64          `json:"nmembers"`
	Members      []*Object      `json:"members"`
	MTime        Time           `json:"mtime"`
	CTime        Time           `json:"ctime"`
	MetaHash     string         `json:"mhash"`
	MetaOnlyHash string         `json:"mohash"`
	NHash        string         `json:"nhash"`
	CHash        string         `json:"chash"`
	Teamfolder   bool           `json:"teamfolder"`
	Readable     bool           `json:"readable"`
	Writable     bool           `json:"writable"`
	Shareable    bool           `json:"shareable"`
	MIMEType     string         `json:"mime_type"`
	RShare       []*ShareObject `json:"rshare"`
}

func (h *Object) UnmarshalJSON(b []byte) error {
	type HiDriveObjectAlias Object
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

	*h = Object(defaultObject)
	return nil
}

// ShareObject represents HiDrive Share object
type ShareObject struct {
	ID           string `json:"id"`
	Path         string `json:"path"`
	Status       string `json:"status"`
	FileType     string `json:"file_type"`
	Count        int    `json:"count"`
	Created      Time   `json:"created"`
	HasPassword  bool   `json:"has_password"`
	Encrypted    bool   `json:"is_encrypted"`
	LastModified Time   `json:"last_modified"`
	MaxCount     int    `json:"maxcount"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	PID          string `json:"pid"`
	Readable     bool   `json:"readable"`
	Remaining    int    `json:"remaining"`
	ShareType    string `json:"share_type"`
	Size         int    `json:"size"`
	TTL          int    `json:"ttl"`
	URI          string `json:"uri"`
	ValidUntil   Time   `json:"valid_until"`
	ViewMode     string `json:"viewmode"`
	Writable     bool   `json:"writable"`
}

func (s *ShareObject) UnmarshalJSON(b []byte) error {
	type HiDriveShareObjectAlias ShareObject
	defaultObject := HiDriveShareObjectAlias{
		Size:      -1,
		TTL:       -1,
		MaxCount:  -1,
		Count:     -1,
		Remaining: -1,
	}

	err := json.Unmarshal(b, &defaultObject)
	if err != nil {
		return err
	}

	*s = ShareObject(defaultObject)
	return nil
}

type ShareInviteStatus struct {
	To      string `json:"to"`
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type ShareInviteResponse struct {
	Done   []ShareInviteStatus `json:"done"`
	Failed []ShareInviteStatus `json:"failed"`
}
