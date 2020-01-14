package main

import "errors"

// ErrAlreadyInitialized is caused by initializing the application more than once.
var ErrAlreadyInitialized = errors.New("already initialized")

// ErrLoadingConfig is caused by an error loading the configuration.
var ErrLoadingConfig = errors.New("error loading configuration")

// ErrShorteningLink is caused by an error shortening a URL.
var ErrShorteningLink = errors.New("error shortening link")
