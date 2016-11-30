package redmine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type timeEntriesResult struct {
	TimeEntries []TimeEntry `json:"time_entries"`
}

type timeEntryResult struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

type timeEntryRequest struct {
	TimeEntry struct {
		IssueId    int     `json:"issue_id,omitempty"`
		ProjectId  int     `json:"project_id,omitempty"`
		SpentOn    string  `json:"created_on,omitempty"`
		Hours      float32 `json:"hours"`
		ActivityId int     `json:"activity_id,omitempty"`
		Comments   string  `json:"comments,omitempty"`
	} `json:"time_entry"`
}

type TimeEntry struct {
	Id        int    `json:"id"`
	Project   IdName `json:"project"`
	Issue     Id     `json:"issue"`
	User      IdName `json:"user"`
	Activity  IdName `json:"activity"`
	Hours     float32
	SpentOn   string `json:"created_on"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	Comments  string `json:"comments"`
}

func (t *TimeEntry) Request() *timeEntryRequest {
	r := &timeEntryRequest{}
	r.TimeEntry.IssueId = t.Issue.Id
	r.TimeEntry.ProjectId = t.Project.Id
	r.TimeEntry.SpentOn = t.SpentOn
	r.TimeEntry.Hours = t.Hours
	r.TimeEntry.ActivityId = t.Activity.Id
	r.TimeEntry.Comments = t.Comments
	return r
}

func (c *client) TimeEntries(projectId int) ([]TimeEntry, error) {
	res, err := c.Get(c.endpoint + "/projects/" + strconv.Itoa(projectId) + "/time_entries.json?key=" + c.apikey + c.getPaginationClause())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntriesResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return r.TimeEntries, nil
}

func (c *client) TimeEntry(id int) (*TimeEntry, error) {
	res, err := c.Get(c.endpoint + "/time_entries/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntryResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return &r.TimeEntry, nil
}

func (c *client) CreateTimeEntry(timeEntry TimeEntry) (*TimeEntry, error) {
	s, err := json.Marshal(timeEntry.Request())
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.endpoint+"/time_entries.json?key="+c.apikey, strings.NewReader(string(s)))

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntryResult
	if res.StatusCode != 201 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return &r.TimeEntry, nil
}

func (c *client) UpdateTimeEntry(timeEntry TimeEntry) error {
	s, err := json.Marshal(timeEntry.Request())
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", c.endpoint+"/time_entries/"+strconv.Itoa(timeEntry.Id)+".json?key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		decoder := json.NewDecoder(res.Body)
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	if err != nil {
		return err
	}
	return err
}

func (c *client) DeleteTimeEntry(id int) error {
	req, err := http.NewRequest("DELETE", c.endpoint+"/time_entries/"+strconv.Itoa(id)+".json?key="+c.apikey, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}
