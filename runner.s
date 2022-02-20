// We need a .s file so the Go tool does not pass -complete to go tool compile;
// that'd prevent being able to define functions with no bodies (something this
// package uses to link into private runtime methods).
