# go-hidrive

Package `go_hidrive` is a simple client SDK library for HiDrive cloud storage
(mainly provided by [Strato](https://www.strato.de/cloud-speicher/) provider) aimed to be used with Go (Golang).

Currently, the following implementation are available: `DirApi`, `FileApi` and `ShareApi`.

All methods accept url.Values as a set of request parameters.
You can also use `Parameters` objects to simplify parameters gathering required for request.

Example reading file from HiDrive:

```go
import "golang.org/x/oauth2"
import hidrive "github.com/Burmuley/go-hidrive"

func main() {
    oauth2config := oauth2.Config{
        ClientID:     "hi_drive_client_id",
        ClientSecret: "hi_drive_client_secret",
        Endpoint: oauth2.Endpoint{
            AuthURL:   hidrive.StratoHiDriveAuthURL,
            TokenURL:  hidrive.StratoHiDriveTokenURL,
            AuthStyle: 0,
        },
        Scopes: []string{"user", "rw"},
    }

    token := &oauth2.Token{
        RefreshToken: "hi_drive_oauth2_refresh_token",
    }

    client := oauth2config.Client(context.Background(), token)
    fileApi := hidrive.NewFileApi(client, StratoHiDriveAPIV21)

    rdr, err := fileApi.GetFile(context.Background(), NewParameters().SetPath("/public/test_file.txt").Values)

    if err != nil {
        fmt.Println(err)
        return
    }
	
    contents, err := io.ReadAll(rdr)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(contents)
}
```
