package rotate

// LogBase is common structure which all stages of log should have.
type Base struct {
	file File
}

// Source is a original file before rotated.
type Source struct {
	*Base
	owner Owner
}

// Evacuate moves original log to temporary storage and wait owner to release handle to the file.
func (source *Source) Evacuate() (*Temp, error) {
	var temp *Temp
	notifyErr := source.owner.NotifyRelease(source.file)
	if notifyErr != nil {
		return temp, notifyErr
	}
	temp = &Temp{}
	return temp, nil
}

// Temp is a temporary state of rotation.
type Temp struct {
	*Base
}

func (temp *Temp) Compress() (*Archive, error) {
	archive := &Archive{}
	return archive, nil
}

type Archive struct {
	*Base
}

func (archive *Archive) Finalize() error {
	return nil
}

// RunRotate
func RunRotate(src *Source) error {
	temp, everr := src.Evacuate()
	if everr != nil {
		return everr
	}
	archive, cmerr := temp.Compress()
	if cmerr != nil {
		return cmerr
	}
	return archive.Finalize()
}
