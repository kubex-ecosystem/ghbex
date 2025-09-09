package common

import "github.com/kubex-ecosystem/ghbex/internal/defs/interfaces"

type Attachment struct {
	Name string `yaml:"name" json:"name"`
	Body []byte `yaml:"body" json:"body"`
}

func NewAttachmentType(name string, body []byte) *Attachment {
	return &Attachment{
		Name: name,
		Body: body,
	}
}

func NewAttachment(name string, body []byte) interfaces.IAttachment {
	return NewAttachmentType(name, body)
}
func (a *Attachment) GetName() string { return a.Name }
func (a *Attachment) GetBody() []byte { return a.Body }
