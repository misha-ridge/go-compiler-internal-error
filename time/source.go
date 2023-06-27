package time

// Source is a function generating a time.Time value
type Source func() Time

// ConstantSource always returns the same constant predefined time.Time that was given to it
// Constant implementation of Source
func ConstantSource(t Time) Source {
	return func() Time { return t }
}
