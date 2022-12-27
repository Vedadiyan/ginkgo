package ginkgobase

type GinkgoError string

const (
	INVALID_PROTOMESSAGE_TYPE GinkgoError = "type is not protomessage"
)

func (ginkgoError GinkgoError) Error() string {
	return string(ginkgoError)
}
