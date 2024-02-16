package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	NEW_ACCOUNT_THRESHOLD = time.Minute * 20
	MENTIONS_THRESHOLD    = 5
	FETCH_DELAY           = time.Second * 5
)

var (
	API_TOKEN = os.Getenv("API_TOKEN")
	HOSTNAME  = os.Getenv("HOSTNAME")
	DEBUG     = os.Getenv("DEBUG") == "true"
)

var (
	lastNoteId string = makeAidx(time.Now())
)

type MiUser struct {
	Id             string  `json:"id"`
	Username       string  `json:"username"`
	Host           string  `json:"host"`
	AvatarBlurhash *string `json:"avatarBlurhash"`
}

type MiNote struct {
	Id         string   `json:"id"`
	User       MiUser   `json:"user"`
	Visibility string   `json:"visibility"`
	Mentions   []string `json:"mentions"`
	Text       string   `json:"text"`
}

type NotesDeleteApiRequest struct { // POST /api/notes/delete
	NoteId string `json:"noteId"`
}

type AdminSuspendUserApiRequest struct { // POST /api/admin/suspend-user
	UserId string `json:"userId"`
}

type NotesGlobalTimelineApiRequest struct { // POST /api/notes/global-timeline
	WithRenotes bool   `json:"withRenotes"`
	Limit       int    `json:"limit"`
	SinceId     string `json:"sinceId"`
}

func main() {
	for {
		var notes []MiNote

		err := requestApi("notes/global-timeline", &NotesGlobalTimelineApiRequest{
			WithRenotes: false,
			Limit:       100,
			SinceId:     lastNoteId,
		}, &notes)
		if err != nil {
			fmt.Println(err)

			time.Sleep(FETCH_DELAY)

			continue
		}

		if len(notes) > 0 {
			lastNoteId = notes[len(notes)-1].Id

			debugLog("Updated last note id: %s\n", lastNoteId)
		}

		debugLog("Fetched %d notes\n", len(notes))

		for _, note := range notes {
			if isTargetNote(note) {
				debugLog("Found target note: %s\n", note.Id)

				err := requestApi("notes/delete", &NotesDeleteApiRequest{
					NoteId: note.Id,
				}, nil)
				if err != nil {
					fmt.Println(err)
				}

				err = requestApi("admin/suspend-user", &AdminSuspendUserApiRequest{
					UserId: note.User.Id,
				}, nil)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		time.Sleep(FETCH_DELAY)
	}
}

func requestApi(endpoint string, body any, target any) error {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}
	bodyMap := make(map[string]interface{})
	err = json.Unmarshal(bodyJson, &bodyMap)
	if err != nil {
		return err
	}
	bodyMap["i"] = API_TOKEN
	bodyJson, err = json.Marshal(bodyMap)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/api/%s", HOSTNAME, endpoint)

	bodyReader := bytes.NewReader(bodyJson)
	resp, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if target != nil {
		err = json.NewDecoder(resp.Body).Decode(target)
		if err != nil {
			return err
		}
	}

	debugLog("Requested to %s and got response: %s\n", url, resp.Status)

	return nil
}

func parseAidx(id string) (time.Time, error) {
	const TIME2000 = 946684800000
	const TIME_LENGTH = 8
	timeInt, err := strconv.ParseInt(id[:TIME_LENGTH], 36, 64)
	if err != nil {
		return time.Time{}, err
	}
	timeInt += TIME2000
	return time.Unix(0, timeInt*int64(time.Millisecond)), nil
}

func makeAidx(t time.Time) string {
	const TIME2000 = 946684800000
	timeInt := t.UnixNano() / int64(time.Millisecond)
	timeInt -= TIME2000
	return strconv.FormatInt(timeInt, 36)
}

func isTargetNote(note MiNote) bool {
	userIsNew := assumeNoError(parseAidx(note.User.Id)).(time.Time).After(time.Now().Add(-NEW_ACCOUNT_THRESHOLD))
	manyMentions := len(note.Mentions) >= MENTIONS_THRESHOLD
	isPublic := note.Visibility == "public"
	noAvatar := note.User.AvatarBlurhash == nil

	debugLog("%d", len(note.Mentions))

	debugLog("Note %s: userIsNew=%t, manyMentions=%t, isPublic=%t, noAvatar=%t\n", note.Id, userIsNew, manyMentions, isPublic, noAvatar)

	return userIsNew && manyMentions && isPublic && noAvatar
}

func assumeNoError(ret interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return ret
}

func debugLog(format string, args ...interface{}) {
	if DEBUG {
		fmt.Printf(format, args...)
	}
}
