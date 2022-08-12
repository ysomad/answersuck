package mime

import "errors"

type Type string

const (
	typeImageJPEG = Type("image/jpeg")
	typeImagePNG  = Type("image/png")
	typeImageWEBP = Type("image/webp")
	typeAudioMP4  = Type("audio/mp4")
	typeAudioAAC  = Type("audio/aac")
	typeAudioMPEG = Type("audio/mpeg")
)

var types = map[Type]struct{}{
	typeImageJPEG: struct{}{},
	typeImagePNG:  struct{}{},
	typeImageWEBP: struct{}{},
	typeAudioMP4:  struct{}{},
	typeAudioAAC:  struct{}{},
	typeAudioMPEG: struct{}{},
}

var imageTypes = map[Type]struct{}{
	typeImageJPEG: struct{}{},
	typeImagePNG:  struct{}{},
	typeImageWEBP: struct{}{},
}

var (
	errInvalid      = errors.New("invalid mime type")
	errInvalidImage = errors.New("provided mime type is not valid image mime type")
)

func NewType(s string) (Type, error) {
	mt := Type(s)
	_, ok := types[mt]
	if !ok {
		return "", errInvalid
	}
	return mt, nil
}

func NewImageType(s string) (Type, error) {
	mt := Type(s)
	_, ok := imageTypes[mt]
	if !ok {
		return "", errInvalidImage
	}
	return mt, nil
}

func (t Type) IsImage() bool {
	_, ok := imageTypes[t]
	return ok
}
