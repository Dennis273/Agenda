package entity

import (
	"encoding/json"
	"fmt"
	"os"
)

type Meeting struct {
	Title        string
	Holder       string
	Participator []string
	StartTime    int64
	EndTime      int64
}

const meetingFilePath string = "./meeting.json"

func CreateMeeting(meeting Meeting) (bool, string) {
	allMeetings := readMeetingsFromFile()
	allUsers := readUsersFromFile()
	currentUser := readCurrentUserFromFile()
	if currentUser == "" {
		return false, "You should login first"
	}

	if len(meeting.Participator) == 0 {
		return false, "A meeting should at least has a participator"
	}
	meeting.Holder = currentUser
	var tempMeeting = meeting
	tempMeeting.Participator = make([]string, 0)
	// check meeting title
	for _, temp := range allMeetings {
		if temp.Title == meeting.Title {
			return false, "meeting exit"
		}
	}

	// check user
	for i, user1 := range meeting.Participator {
		var flag = true
		for _, user2 := range allUsers {
			if user1 == user2.Username {
				flag = false
				break
			}
		}
		if flag {
			return false, fmt.Sprintf("%s is not exit", user1)
		}
		if user1 == meeting.Holder {
			return false, fmt.Sprint("You could not attend this meeting as participator")
		}
		for _, user2 := range meeting.Participator[i+1:] {
			if user1 == user2 {
				return false, fmt.Sprint("You could not repeatly attend a user to a same meeting")

			}
		}
	}

	for _, temp := range allMeetings {
		if (meeting.StartTime <= temp.StartTime && meeting.EndTime > temp.StartTime) || (meeting.StartTime < temp.EndTime && meeting.EndTime >= temp.EndTime) {
			if temp.Holder == meeting.Holder {
				return false, "You are busy at that time"
			}
			for _, user1 := range meeting.Participator {
				if user1 == meeting.Holder {
					return false, fmt.Sprintf("%s is busy at that time", user1)
				}
				for _, user2 := range temp.Participator {
					if meeting.Holder == user2 {
						return false, "You are busy at that time"
					}
					if user1 == user2 {
						return false, fmt.Sprintf("%s is busy at that time", user1)
					}
				}
			}
		}
	}
	allMeetings = append(allMeetings, meeting)
	writeMeetingsIntoFile(allMeetings)
	Info("Create meeting " + meeting.Title + " success")
	return true, ""
}

func ModifyMeeting(title string, adds []string, removes []string) (bool, string) {
	allMeetings := readMeetingsFromFile()
	allUsers := readUsersFromFile()
	usernames := make(map[string]string, len(allUsers))
	currentUser := readCurrentUserFromFile()
	if currentUser == "" {
		panic("You should login first")
	}
	for _, temp := range allUsers {
		usernames[temp.Username] = temp.Username
	}
	for _, temp := range adds {
		if _, ok := usernames[temp]; !ok {
			panic(fmt.Sprintf("User %s is not exit", temp))
		}
	}
	for _, temp := range removes {
		if _, ok := usernames[temp]; !ok {
			panic(fmt.Sprintf("User %s is not exit", temp))
		}
	}
	var meeting *Meeting
	var meetingIndex int
	for i, temp := range allMeetings {
		if title == temp.Title {
			if temp.Holder == currentUser {
				meetingIndex = i
				meeting = &allMeetings[meetingIndex]
			} else {
				panic(fmt.Sprintf("You are not the holder of meeting %s", title))
			}
		}
	}
	if meeting == nil {
		panic(fmt.Sprintf("Meeting %s is not exit", title))
	}

	for _, temp := range adds {
		if success, err := AddMemberToMeeting(meeting, temp, allMeetings); !success {
			panic(err)
		}
	}
	for _, temp := range adds {
		for _, p := range meeting.Participator {
			if temp == p {
				panic(fmt.Sprintf("You can not repeatly add %s to the meeting", temp))
			}
		}
		meeting.Participator = append(meeting.Participator, temp)
	}
	for _, temp := range removes {
		if success, err := RemoveMemberFromMeeting(meeting, temp, allMeetings); !success {
			panic(err)
		}
	}
	for _, temp := range removes {
		var index = -1
		for i, p := range meeting.Participator {
			if temp == p {
				index = i
				break
			}
		}
		if index == -1 {
			panic(fmt.Sprintf("You can not repeatly remove %s from the meeting", temp))
		}
		meeting.Participator = append(meeting.Participator[:index], meeting.Participator[index+1:]...)
	}
	if len(meeting.Participator) == 0 {
		allMeetings = append(allMeetings[:meetingIndex], allMeetings[meetingIndex+1:]...)
	}
	writeMeetingsIntoFile(allMeetings)
	Info("ModifyMeeting " + title + " success")
	return true, ""
}

func AddMemberToMeeting(meeting *Meeting, user string, meetings []Meeting) (bool, string) {
	if meeting.Holder == user {
		return false, "You are the holder of this meeting"
	}
	if isParticipator(user, meeting.Participator) {
		return false, user + " was a participator of this meeting"
	}
	for _, temp := range meetings {
		if user == temp.Holder || isParticipator(user, temp.Participator) && temp.Title != meeting.Title {
			if (temp.StartTime <= meeting.StartTime && temp.EndTime > meeting.StartTime) || (temp.EndTime >= meeting.EndTime && temp.StartTime < meeting.EndTime) {
				return false, user + " is busy at that time"
			}
		}
	}
	return true, ""
}

func RemoveMemberFromMeeting(meeting *Meeting, user string, meetings []Meeting) (bool, string) {
	if meeting.Holder == user {
		return false, "You are the holder of this meeting"
	}
	if !isParticipator(user, meeting.Participator) {
		return false, user + " was not a participator of this meeting"
	}
	return true, ""
}

func QueryMeeting(startTime, endTime int64) []Meeting {
	var result []Meeting
	currentUser := readCurrentUserFromFile()
	if currentUser == "" {
		panic("You should login first")
	}
	allMeetings := readMeetingsFromFile()
	for _, meeting := range allMeetings {
		if currentUser == meeting.Holder || isParticipator(currentUser, meeting.Participator) {
			if (meeting.StartTime <= startTime && meeting.EndTime >= startTime) || (meeting.EndTime >= endTime && meeting.EndTime <= endTime) {
				result = append(result, meeting)
			}
		}
	}
	Info("Query meetings success")
	return result
}

func isParticipator(username string, participartors []string) bool {
	for _, temp := range participartors {
		if username == temp {
			return true
		}
	}
	return false
}

func CancelMeeting(title string) (bool, string) {
	return true, ""
}

func QuitMeeting(title string) (bool, string) {
	return true, ""
}

func ClearMeeting() (bool, string) {
	return true, ""
}

func readMeetingsFromFile() []Meeting {
	meetings := make([]Meeting, 0)
	file, err := os.Open(meetingFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fInfo, err := os.Stat(meetingFilePath)
	if err != nil {
		panic(err)
	}
	buffer := make([]byte, fInfo.Size())
	_, err = file.Read(buffer)
	if fInfo.Size() == 0 {
		return meetings
	}
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buffer, &meetings)
	if err != nil {
		panic(err)
	}
	return meetings
}

func writeMeetingsIntoFile(meetings []Meeting) {
	data, err := json.Marshal(meetings)
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(meetingFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
}

func init() {
	// check if meeting.json exist
	if _, err := os.Stat(meetingFilePath); os.IsNotExist(err) {
		_, err := os.Create(meetingFilePath)
		if err != nil {
			panic(err)
		}
		file, err := os.OpenFile(meetingFilePath, os.O_WRONLY|os.O_TRUNC, 0777)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = file.Write([]byte("[]"))
		if err != nil {
			panic(err)
		}
	}
}
