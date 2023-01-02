package model

type ErrNotFound struct {

}

func (e ErrNotFound) Error() string {
	// fmt.Sprintf("not found:", )
	return "not found"
}