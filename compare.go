package main

// Compare checks two file slices against each other and returnes two slices in which the first one represents
// Files which are in the first one but not in the second, and the second one represents files which are in the second,
// but not in the first.
func Compare(new []File, old []File) (newandchanged []File, deleted []File) {
	// If new files are not in the old scan, they are new or changed.
	for _, vnew := range new {
		if vnew.ifFileIsIn(old) == false {
			newandchanged = append(newandchanged, vnew)
		}
	}

	// If an old file is not in the new scan, it must have been deleted or changed.
	for _, vold := range old {
		if vold.ifFileIsIn(new) == false {
			deleted = append(deleted, vold)
		}
	}
	return
}

func (f File) ifFileIsIn(list []File) bool {
	for _, v := range list {
		if f == v {
			return true
		}
	}
	return false
}
