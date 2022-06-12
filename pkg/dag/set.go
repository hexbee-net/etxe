package dag

// Set is a set data structure.
type Set[T comparable] map[T]T

// Add adds an item to the set.
func (s Set[T]) Add(v T) {
	s[v] = v
}

// Delete removes an item from the set.
func (s Set[T]) Delete(v T) {
	delete(s, v)
}

// Includes returns whether a value is in the set.
func (s Set[T]) Includes(v T) bool {
	_, exists := s[v]

	return exists
}

// Intersection computes the set intersection with `other`.
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	res := make(Set[T])

	if s == nil || other == nil {
		return res
	}

	// Iterate over the smaller set for better performances.
	if len(other) < len(s) {
		s, other = other, s
	}

	for _, v := range s {
		if other.Includes(v) {
			res.Add(v)
		}
	}

	return res
}

// Difference returns a set with the elements that `s` has but `other` doesn't.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	if other == nil || len(other) == 0 {
		return s.Copy()
	}

	res := make(Set[T])

	for k, v := range s {
		if _, ok := other[k]; !ok {
			res.Add(v)
		}
	}

	return res
}

// Filter returns a set that contains the elements from the receiver
// where the given callback returns true.
func (s Set[T]) Filter(cb func(T) bool) Set[T] {
	res := make(Set[T])

	for _, v := range s {
		if cb(v) {
			res.Add(v)
		}
	}

	return res
}

// List returns the list of set elements.
func (s Set[T]) List() []T {
	if s == nil {
		return nil
	}

	res := make([]T, 0, len(s))
	for _, v := range s {
		res = append(res, v)
	}

	return res
}

// Copy returns a shallow copy of the set.
func (s Set[T]) Copy() Set[T] {
	c := make(Set[T], len(s))

	for k, v := range s {
		c[k] = v
	}

	return c
}
