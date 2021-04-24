package core

type ParseOption int

const (
	Second         ParseOption = 1 << iota // Seconds field, default 0
	SecondOptional                         // Optional seconds field, default 0
	Minute                                 // Minutes field, default 0
	Hour                                   // Hours field, default 0
	Dom                                    // Day of month field, default *
	Month                                  // Month field, default *
	Dow                                    // Day of week field, default *
	DowOptional                            // Optional day of week field, default *
	Descriptor                             // Allow descriptors such as @monthly, @weekly, etc.
)
