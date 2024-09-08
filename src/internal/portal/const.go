package portal

const (
	// InputFieldLabelUriVariableId holds the identifer used for the input label in the URI
	InputFieldLabelUriVariableId = "inputFieldVariableId"

	// ErrKeyInvalidInputFieldId is returned when the input field label cannot be found for
	// a targetted request
	ErrKeyInvalidInputFieldId = "InvalidInputFieldId"

	// ErrNoFilesProvidedWithUploadRequest is returned when no files are detected in the
	// form data of a request
	ErrNoFilesProvidedWithUploadRequest = "NoFilesProvidedWithUploadRequest"

	// ErrKeyNoInputFieldCacheDirFound is returned when no cache directory is found for a given input field label
	ErrKeyNoInputFieldCacheDirFound = "NoInputFieldCacheDirFound"

	// ErrKeyUnableToReadCacheDir is returned when the cache directory cannot be read
	ErrKeyUnableToReadCacheDir = "UnableToReadCacheDir"

	// ErrKeyUnableToRemoveCacheDirContents is returned when the cache directory contents cannot be removed
	ErrKeyUnableToRemoveCacheDirContents = "UnableToRemoveCacheDirContents"
)
