package interfaces

type IAttachment interface {
	// GetName returns the name of the attachment.
	GetName() string
	// GetBody returns the body of the attachment.
	GetBody() []byte
}
