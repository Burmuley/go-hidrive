package go_hidrive

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Parameters struct {
	url.Values
}

func NewParameters() *Parameters {
	return &Parameters{url.Values{}}
}

/*
SetPath adds "path" parameter to the query - path to a filesystem object

Can be used in the following methods:
  - [Dir.Get]
  - [Dir.Create]
  - [Dir.Delete]
  - [File.Get]
  - [File.Delete]
  - [Share.Get]
  - [Share.Create]
  - [Sharelink.Create]
  - [Meta.Get]
  - [Meta.Update]
*/
func (p *Parameters) SetPath(path string) *Parameters {
	p.Set("path", path)
	return p
}

/*
SetPid adds "pid" parameter to the query.

The public id is a path and encoding independent representation
of a specific filesystem object. Also returned and referred to as id in data related responses.

Can be used in the following methods:
  - [Dir.Get]
  - [Dir.Create]
  - [Dir.Delete]
  - [File.Get]
  - [File.Delete]
  - [Share.Get]
  - [Share.Create]
  - [Sharelink.Create]
  - [Meta.Get]
  - [Meta.Update]
*/
func (p *Parameters) SetPid(pid string) *Parameters {
	p.Set("pid", pid)
	return p
}

/*
SetMembers - adds "members" parameter to the query - list of directory content types to be included in the members part of the response

Valid values are:
  - all     - include all contents
  - none    - do not return any members
  - dir     - include sub-directories   (not in combination with none or all)
  - file    - include files             (not in combination with none or all)
  - symlink - include symlinks          (not in combination with none or all)

Can be used in the following methods:
  - [Dir.Get]
*/
func (p *Parameters) SetMembers(members []string) *Parameters {
	p.Set("members", strings.Join(members, ","))
	return p
}

/*
SetLimit - adds "limit" parameter to the query - limit the number of directory entries returned, starting from an optional offset.

Both <limit> and <offset> need to be nonnegative integer values.

The returned amount of entries may be less than requested.
To get all directory entries it is always recommended to check the nmembers field and issue another request with an <offset> updated accordingly.

A value of none or 0 for <limit> signifies to return as many entries as is feasible. This also works when combined with an offset.

Can be used in the following methods:
  - [Dir.Get]
*/
func (p *Parameters) SetLimit(limit uint, offset uint) *Parameters {
	p.Set("limit", fmt.Sprintf("%d,%d", offset, limit))
	return p
}

/*
SetFields - adds "fields" parameter to the query - list of value types that will be included in the response.

The performance of the call might be influenced by the amount of information requested.
Therefore, it is recommended to use a "need to know" approach instead of "get all".

Can be used in the following methods:
  - [Dir.Get]
  - [Share.Get]
  - [Sharelink.Get]
  - [Meta.Get]

Valid values for [Dir.Get]:
  - category                - string    - object category (audio, image, etc.)
  - chash (*)               - string    - recursive hashvalue for the directory
  - ctime                   - timestamp - ctime of the object
  - has_dirs                - bool      - contains subdirs?
  - id                      - string    - path id (pid) of the directory
  - members                 - array     - include information on dir contents
  - members.category        - string    - object category (audio, image, etc.)
  - members.chash (*)       - string    - recursive hashvalue for a contained directory
  - members.ctime           - timestamp - ctime of contained objects
  - members.has_dirs        - bool      - does a contained dir contain subdirs?
  - members.id              - string    - path id (pid) of contained object
  - members.image.exif      - object    - selected exif data of contained images
  - members.image.height    - int       - height of contained images
  - members.image.width     - int       - width of contained images
  - members.mhash (*)       - string    - meta hash of contained objects
  - members.mime_type       - string    - MIME type of contained files
  - members.mohash (*)      - string    - meta only hash of contained objects
  - members.mtime           - timestamp - mtime of contained objects
  - members.name            - string    - URL-Encoded name of contained objects
  - members.nmembers (*)    - int       - number of members of a contained directory
  - members.nhash (*)       - string    - name hash of contained objects
  - members.path            - string    - URL-Encoded path of contained objects
  - members.parent_id       - string    - path id (pid) of the members parent
  - members.parent.id       - string    - path id (pid) of the members parent
  - members.parent.writable - bool      - write-permission of the members parent
  - members.readable        - bool      - read-permission for contained objects
  - members.rshare          - array     - sharing information (details below)
  - members.size (*)        - int       - recursive size of a contained directory
  - members.type            - string    - dir/file/symlink (see param "members")
  - members.writable        - bool      - write-permission for contained objects
  - members.shareable       - bool      - share-permission for the contained objects
  - members.teamfolder      - bool      - indicates whether the contained object is a teamfolder or not
  - mhash (*)               - string    - meta hash of the object
  - mohash (*)              - string    - meta only hash of the object
  - mtime                   - timestamp - mtime of the directory
  - name                    - string    - URL-Encoded name of the directory
  - nhash (*)               - string    - name hash of the object
  - nmembers                - int       - number of members in the directory
  - parent_id               - string    - path id (pid) of the parent directory
  - parent.id               - string    - path id (pid) of the parent directory
  - parent.writable         - bool      - write-permission for the parent directory
  - path                    - string    - URL-Encoded path of the directory
  - readable                - bool      - read-permission for the directory
  - rshare                  - object    - sharing information (details below)
  - size (*)                - int       - recursive size of the directory
  - type                    - string    - dir/file/symlink (see param "members")
  - writable                - bool      - write-permission for the directory
  - shareable               - bool      - share-permission for the directory
  - teamfolder              - bool      - indicates whether the directory is a teamfolder or not

Valid values for [Share.Get]:
  - count           - int       - the number of successfully completed downloads
  - created         - int       - UNIX timestamp
  - file_type       - string    - 'dir'
  - has_password    - bool
  - is_encrypted    - bool      - is the given share encrypted
  - id              - string    - the unique share id
  - last_modified   - int       - UNIX timestamp
  - maxcount (*)    - int       - maximum number of share-tokens
  - name            - int       - name of the shared directory
  - password (*)    - string    - optional password for the share
  - path            - string    - path of the shared directory
  - pid             - string    - path id of the shared directory
  - readable        - bool
  - remaining       - int       - number of remaining available share tokens
  - share_type      - string    - 'sharedir'
  - size            - int       - size of the shared directory
  - status          - string    - 'valid', 'invalid' or 'expired'
  - ttl             - int       - time-to-live, in seconds, possibly negative
  - uri             - string    - url of the shared directory
  - valid_until     - int       - UNIX timestamp
  - viewmode        - string    - single letter. influences the share folder display
  - writable        - bool
*/
func (p *Parameters) SetFields(fields []string) *Parameters {
	p.Set("fields", strings.Join(fields, ","))
	return p
}

/*
SetSortBy - adds "sort" parameter to the request - determines the order of the members in the result.

They can be sorted by name, category, mtime, type, or size. The default sort order is ascending.
Prefix the sort criterion with a dash '-' for descending order.
The first criteria in the comma-separated list take precedence over the others.

Names are sorted case-insensitive according to the locale determined by the sort_lang parameter.
Numbers are compared by their numerical value.

The size of a directory is the recursive sum of file sizes of all files it contains.
The size of a directory in a snapshot is sorted as 0 and not reported.

With the value "none" the output is unsorted.

Can be used in the following methods:
  - [Dir.Get]
*/
func (p *Parameters) SetSortBy(sortBy string) *Parameters {
	p.Set("sort", sortBy)
	return p
}

/*
SetSortLang - adds "sort_lang" parameter to the request - Determines the locale used for sorting.

Currently allowed values are `de_DE`, `en_US` and `sv_SE`.

Can be used int the following methods:
  - [Dir.Get]
*/
func (p *Parameters) SetSortLang(lang string) *Parameters {
	p.Set("sort_lang", lang)
	return p
}

/*
SetOnExist - adds "on_exist" parameter to the request - Optional parameter to determine the API behavior
in case of a conflict with an existing filesystem object.

Valid values are:
  - "autoname"  - find another name if the destination already exists

Can be used in the following methods:
  - [File.Upload]
*/
func (p *Parameters) SetOnExist(onExists string) *Parameters {
	p.Set("on_exist", onExists)
	return p
}

/*
SetMTime - adds "mtime" parameter to the request - the modification time (mtime) of the file system target
to be set after the operation.

Can be used in the following methods:
  - [Dir.Create]
  - [File.Upload]
  - [Meta.Update]
*/
func (p *Parameters) SetMTime(t time.Time) *Parameters {
	p.Set("mtime", fmt.Sprint(t.Unix()))
	return p
}

/*
SetParentMTime - adds "parent_mtime" parameter to the query - the modification time (mtime) of the file system
target's parent folder to be set after the operation.

Can be used in the following methods:
  - [Dir.Create]
  - [Dir.Delete]
  - [File.Delete]
  - [File.Upload]
*/
func (p *Parameters) SetParentMTime(t time.Time) *Parameters {
	p.Set("parent_mtime", fmt.Sprint(t.Unix()))
	return p
}

/*
SetRecursive - adds "recursive" parameter to the request - if `true`, the call will also delete non-empty directories
and their contents recursively without throwing a 409 Conflict error.

Can be used in the following methods:
  - [Dir.Delete]
*/
func (p *Parameters) SetRecursive(recursive bool) *Parameters {
	p.Set("on_exist", fmt.Sprint(recursive))
	return p
}

/*
SetDir - adds "dir" parameter to the request - the path to the filesystem object (directory) to be used as the target
for the upload operation.

The shortest possible path is "/" and it will always refer to the topmost directory accessible by the authenticated user.
For a regular HiDrive user this is the HiDrive "root". If used with a share access_token it will be the shared directory.

This value must not contain path elements "." or ".." and must not end with a slash "/".

Note: if used in combination with a dir_id, this value is not allowed to start with "/" either.

Note: this is always a parent directory and must not contain the intended filename.
Use the SetName method to specify the file name.

Can be used in the following methods:
  - [File.Upload]
*/
func (p *Parameters) SetDir(dir string) *Parameters {
	p.Set("dir", dir)
	return p
}

/*
SetDirId - adds "dir_id" parameter to the request - the pulic id (pid) of the target filesystem object.
(Or, if used in combination with dir, its parent directory.)

Note: a pid is not persistent upon changes (rename/move) to a filesystem object.
So after this operation, the dir_id may no longer be valid.
However, the current value will be part of the returned information (as parent_id) after a successful request.

Can be used in the following methods:
  - [File.Upload]
*/
func (p *Parameters) SetDirId(id string) *Parameters {
	p.Set("dir_id", id)
	return p
}

/*
SetName - adds "name" parameter to the request - the intended filename.

The name parameter is mandatory for binary uploads. It is forbidden for multipart/formdata uploads, where the name has
to be specified as "filename" parameter within the content-disposition header.

Can be used in the following methods:
  - [File.Upload]
*/
func (p *Parameters) SetName(name string) *Parameters {
	p.Set("name", name)
	return p
}

/*
SetFilePath - parses the path provided and uses the last part as file name (field "name"),
the rest of the path is defined in the "dir" parameter.

Use this method to simplify settings upload path with a single string instead of calling two methods to set "dir" and
"name separately.

Can be used in the following methods:
  - [File.Upload]
*/
func (p *Parameters) SetFilePath(path string) *Parameters {
	elems := strings.Split(path, "/")
	fName := elems[len(elems)-1]
	dir := strings.Join(elems[:len(elems)-1], "/")
	p.SetDir(dir)
	p.SetName(fName)
	return p
}

/*
SetMaxCount - adds "maxcount" parameter to the request - number of share tokens that can be issued for this share.

When not provided, the value will be set to unlimited, if the user's tariff supports it, otherwise to the maximum
value permissible.

Can be used in the following methods:
  - [Share.Create]
  - [Sharelink.Create]
  - [Sharelink.Update]
*/
func (p *Parameters) SetMaxCount(count int) *Parameters {
	p.Set("maxcount", fmt.Sprint(count))
	return p
}

/*
SetPassword - adds "password" parameter to the request - optional protection for the share.

Consider this recommended, especially the closer the share is set to the root directory.
This parameter must be omitted for encrypted shares which require salt, share_access_key, pw_sharekey.

Can be used in the following methods:
  - [Share.Create]
  - [Sharelink.Create]
  - [Sharelink.Update]
*/
func (p *Parameters) SetPassword(password string) *Parameters {
	p.Set("password", password)
	return p
}

/*
SetWritable - adds "writable" parameter to the request - This option can be set to allow write access to the shared filesystem object.

Note: This includes deletion and modification of existing content.

Can be used in the following methods:
  - [Share.Create]
*/
func (p *Parameters) SetWritable(writable bool) *Parameters {
	p.Set("writable", fmt.Sprint(writable))
	return p
}

/*
SetTTL - adds "ttl" parameter to the request - share expiry.

A positive number defining seconds from now. Not specifying a value sets ttl to the tariff maximum.

Can be used in the following methods:
  - [Share.Create]
  - [Sharelink.Create]
  - [Sharelink.Update]
*/
func (p *Parameters) SetTTL(ttl uint) *Parameters {
	p.Set("ttl", fmt.Sprint(ttl))
	return p
}

/*
SetSalt - adds "salt" parameter to the request - Random salt value generated by the hdcrypt library for encrypted shares.

If this parameter is present, the share is created as 'encrypted' and share_access_key as well as pw_sharekey must
also be present. The password parameter must be omitted because encrypted shares rely on a challenge-response
authentication that only requires knowledge of the share_access_key.

Note: this attribute cannot be removed from a share.

Can be used in the following methods:
  - [Share.Create]
*/
func (p *Parameters) SetSalt(salt string) *Parameters {
	p.Set("salt", salt)
	return p
}

/*
SetShareAccessKey - adds "share_access_key" parameter to the request - Authentication key provided by the `hdcrypt`
library for encrypted shares. Requires `password` to be absent and salt and `pw_sharekey` to be present.

Can be used in the following methods:
  - [Share.Create]
*/
func (p *Parameters) SetShareAccessKey(key string) *Parameters {
	p.Set("share_access_key", key)
	return p
}

/*
SetPwShareKey - adds "pw_sharekey" parameter to the request - Password protected Share Key provided by the `hdcrypt`
library for encrypted shares. Requires `password` to be absent and `salt` and `share_access_key` to be present.

Can be used in the following methods:
  - [Share.Create]
*/
func (p *Parameters) SetPwShareKey(key string) *Parameters {
	p.Set("pw_sharekey", key)
	return p
}

/*
SetId - adds "id" parameter to the request - a share id as returned by [Share.Get] or [Share.Create].

Can be used in the following methods:
  - [Share.Create]
  - [Sharelink.Get]
  - [Sharelink.Create]
  - [Sharelink.Update]
  - [Sharelink.Delete]
*/
func (p *Parameters) SetId(id string) *Parameters {
	p.Set("id", id)
	return p
}

/*
SetRecipient - adds "recipient" parameter to the request - A RFC822-compliant, UTF-8 encoded e-mail address.

The parameter can be specified multiple times to send an invitation to more than one recipient at once.

Note: If the address is preceded by a string (e.g. "Bob Test" <bob@example.com>), the specified string is used as
salutation in the generated mail without modification. It is recommended to specify names as "Firstname Lastname"
instead of "Lastname, Firstname".

Can be used in the following methods:
  - [Share.Invite]
*/
func (p *Parameters) SetRecipient(recip string) *Parameters {
	p.Set("recipient", recip)
	return p
}

/*
SetMsg - adds "msg" parameter to the request - A UTF-8 encoded message text that will be included in the e-mail.

Can be used in the following methods:
  - [Share.Invite]
*/
func (p *Parameters) SetMsg(msg string) *Parameters {
	p.Set("msg", msg)
	return p
}

/*
SetSrc - adds "src" parameter to the request - the path to the filesystem object to be used as the source for
the operation.

Example: /users/example/something

It is not allowed to use only "/" as a source, since it is the parent for all possible targets.
Note: if used in combination with a src_id, this value is not allowed to start with "/" either.

Can be used in the following methods:
  - [File.Copy]
  - [File.Move]
*/
func (p *Parameters) SetSrc(src string) *Parameters {
	p.Set("src", src)
	return p
}

/*
SetSrcId - adds "src_id" parameter to the request - the public id (`pid`) of the source filesystem object
(or when used in combination with src, its parent directory.)

Example: b1489258310.123

Note: a pid is not persistent upon changes (rename/move) to a filesystem object.
So after this operation, the src_id may no longer be valid. However, the current value will be part of the returned
information (as id) after a successful request.

Can be used in the following methods:
  - [File.Copy]
  - [File.Move]
*/
func (p *Parameters) SetSrcId(srcId string) *Parameters {
	p.Set("src_id", srcId)
	return p
}

/*
SetDst - adds "dst" parameter to the request - the path to the filesystem object to be used as the destination for
the operation.

Example: /users/example/archive/something

Note: if used in combination with a dst_id, this value is not allowed to start with "/".

Can be used in the following methods:
  - [File.Copy]
  - [File.Move]
*/
func (p *Parameters) SetDst(dst string) *Parameters {
	p.Set("dst", dst)
	return p
}

/*
SetDstId - adds "dst_id" parameter to the request - if provided, this must always be the pid of a parent directory
of the dst.

Example: b1489258310.123

Can be used in the following methods:
  - [File.Copy]
  - [File.Move]
*/
func (p *Parameters) SetDstId(dstId string) *Parameters {
	p.Set("dst_id", dstId)
	return p
}

/*
SetSrcParentMTime - adds "src_parent_mtime" parameter to the request - the modification time (mtime) of the source
parent folder to be set after the operation.

Can be used in the following methods:
  - [File.Move]
*/
func (p *Parameters) SetSrcParentMTime(t time.Time) *Parameters {
	p.Set("src_parent_mtime", fmt.Sprint(t.Unix()))
	return p
}

/*
SetDstParentMTime - adds "dst_parent_mtime" parameter to the request - the modification time (mtime) of the destination
parent folder to be set after the operation.

Can be used in the following methods:
  - [File.Copy]
  - [File.Move]
*/
func (p *Parameters) SetDstParentMTime(t time.Time) *Parameters {
	p.Set("dst_parent_mtime", fmt.Sprint(t.Unix()))
	return p
}

/*
SetPreserveMTime - adds "preserve_mtime" parameter to the request - the modification time (mtime) of the target will be
copied from the source after the operation.

Can be used in the following methods:
  - [File.Copy]
*/
func (p *Parameters) SetPreserveMTime(pmTime bool) *Parameters {
	p.Set("preserve_mtime", fmt.Sprint(pmTime))
	return p
}
