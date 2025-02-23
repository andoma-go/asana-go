package asana

import (
	"fmt"
	"time"
)

// StoryBase contains the text of a story, as used when creating a new comment
type StoryBase struct {
	// Human-readable text for the story or comment. This will
	// not include the name of the creator. Can be edited only if the story is a comment.
	//
	// Note: This is not guaranteed to be stable for a given type of story.
	// For example, text for a reassignment may not always say “assigned to …”
	// as the text for a story can both be edited and change based on the language
	// settings of the user making the request. Use the resource_subtype property to
	// discover the action that created the story.
	Text string `json:"text,omitempty"`

	// HTML formatted text for a comment. This will not include the name of the creator.
	// Can be edited only if the story is a comment.
	//
	// Note: This field is only returned if explicitly requested using the
	// opt_fields query parameter.
	HTMLText string `json:"html_text,omitempty"`

	// Whether the story should be pinned on the resource.
	// Note: This field is only present on comment and attachment stories.
	IsPinned bool `json:"is_pinned,omitempty"`
}

type Dates struct {
	DueOn   *Date      `json:"due_on,omitempty"`
	DueAt   *time.Time `json:"due_at,omitempty"`
	StartOn *Date      `json:"start_on,omitempty"`
}

type StorySubtypeFields struct {
	// Whether the text of the story has been edited after creation.
	// Note: This field is only present on comment stories.
	IsEdited bool `json:"is_edited,omitempty"`

	// Present for due_date_changed, dependency_due_date_changed
	NewDates *Dates `json:"new_dates,omitempty"`

	// Present for due_date_changed
	OldDates *Dates `json:"old_dates,omitempty"`

	// Present for name_changed
	OldName string `json:"old_name,omitempty"`
	NewName string `json:"new_name,omitempty"`

	// Present for resource_subtype_changed
	OldResourceSubtype string `json:"old_resource_subtype,omitempty"`
	NewResourceSubtype string `json:"new_resource_subtype,omitempty"`

	// Present for comment_liked, completion_liked
	Story *Story `json:"story,omitempty"`

	// Present for attachment_liked
	Attachment *Attachment `json:"attachment,omitempty"`

	// Present for assigned
	Assignee *User `json:"assignee,omitempty"`

	// Present for follower_added
	Follower *User `json:"follower,omitempty"`

	// Present for section_changed
	OldSection *Section `json:"old_section,omitempty"`
	NewSection *Section `json:"new_section,omitempty"`

	// Present for added_to_task, removed_from_task
	Task *Task `json:"task,omitempty"`

	// Present for added_to_project, removed_from_project
	Project *Project `json:"project,omitempty"`

	// Present for added_to_tag, removed_from_tag
	Tag *Tag `json:"tag,omitempty"`

	// Present for text_custom_field_changed, number_custom_field_changed, enum_custom_field_changed
	OldTextValue   string     `json:"old_text_value,omitempty"`
	NewTextValue   string     `json:"new_text_value,omitempty"`
	OldNumberValue float64    `json:"old_number_value,omitempty"`
	NewNumberValue float64    `json:"new_number_value,omitempty"`
	OldEnumValue   *EnumValue `json:"old_enum_value,omitempty"`
	NewEnumValue   *EnumValue `json:"new_enum_value,omitempty"`

	// Present for duplicate_merged, marked_duplicate, duplicate_unmerged
	DuplicateOf *Task `json:"duplicate_of,omitempty"`

	// Present for duplicated
	DuplicatedFrom *Task `json:"duplicated_from,omitempty"`

	// Present for dependency_added, dependency_removed, dependency_marked_complete, dependency_marked_incomplete,
	// dependency_due_date_changed
	Dependency *Task `json:"dependency,omitempty"`
}

// Story represents an activity associated with an object in the Asana
// system. Stories are generated by the system whenever users take actions
// such as creating or assigning tasks, or moving tasks between projects.
// Comments are also a form of user-generated story.
//
// Stories are a form of history in the system, and as such they are read-
// only. Once generated, it is not possible to modify a story.
type Story struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	StoryBase

	// Read-only. The time at which this object was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// True if the story is liked by the authorized user, false if not.
	// Note: This property only exists for stories that provide likes.
	Liked bool `json:"liked,omitempty"`

	// Read-only. Array of users who have liked this story.
	// Note: This property only exists for stories that provide likes.
	Likes []*User `json:"likes,omitempty"`

	// Read-only. The number of users who have liked this story.
	// Note: This property only exists for stories that provide likes.
	NumLikes int32 `json:"num_likes,omitempty"`

	// The user who created the story.
	CreatedBy *User `json:"created_by,omitempty"`

	// Read-only. The object this story is associated with. Currently may only
	// be a task.
	Target *Task `json:"target,omitempty"`

	// Read-only. The component of the Asana product the user used to trigger
	// the story.
	Source string `json:"source,omitempty"`

	// Read-only. The type of story. This provides fine-grained information about what
	// triggered the story’s creation. There are many story subtypes, so inspect the
	// data returned from Asana’s API to find the value for your use case.
	ResourceSubtype string `json:"resource_subtype,omitempty"`

	// A union of all possible subtype fields
	StorySubtypeFields
}

// Stories lists all stories attached to a task
func (t *Task) Stories(client *Client, opts ...*Options) ([]*Story, *NextPage, error) {
	client.trace("Listing stories for %q", t.Name)

	var result []*Story

	// Make the request
	nextPage, err := client.get(fmt.Sprintf("/tasks/%s/stories", t.ID), nil, &result, opts...)
	return result, nextPage, err
}

// CreateComment adds a comment story to a task
func (t *Task) CreateComment(client *Client, story *StoryBase) (*Story, error) {
	client.info("Creating comment for task %q", t.Name)

	result := &Story{}

	err := client.post(fmt.Sprintf("/tasks/%s/stories", t.ID), story, result)
	return result, err
}

// UpdateStory updates the story and returns the full record for the updated story.
// Only comment stories can have their text updated, and only comment stories and attachment stories can be pinned.
// Only one of text and html_text can be specified.
func (s *Story) UpdateStory(client *Client, story *StoryBase) (*Story, error) {
	client.info("Updating story %s", s.ID)

	result := &Story{}

	err := client.put(fmt.Sprintf("/stories/%s", s.ID), story, result)
	return result, err
}

func (s *Story) Delete(client *Client) error {
	client.trace("Delete story %s %s", s.ID, s.ResourceSubtype)

	err := client.delete(fmt.Sprintf("/stories/%s", s.ID))
	return err
}
