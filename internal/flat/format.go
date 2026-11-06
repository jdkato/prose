// Package flat defines the on-disk format for a flat averaged-perceptron
// model, plus the reader for it.
//
// # Why not gob
//
// The obvious encoding, map[string]map[string]float64, costs roughly 28 MB
// live and ~815,000 allocations to build for the shipped English model — its
// 75,447 inner maps are almost entirely header and bucket overhead around
// 1.08 MB of actual weight data.
//
// # Layout
//
// The weights are sparse (mean 3.59 classes per feature, median 2), so they
// are stored compressed-sparse-row: one offsets array indexing into parallel
// classID and weight arrays. All little-endian.
//
//	magic     [8]byte  "PROSETAG"
//	version   uint32
//	nClasses  uint32   then per class: uint16 len + bytes
//	nFeatures uint32
//	nPairs    uint32
//	blobLen   uint32   then blobLen bytes  (feature strings, concatenated,
//	                                        sorted bytewise)
//	featOff   (nFeatures+1) x uint32       (slices into the blob)
//	pairOff   (nFeatures+1) x uint32       (slices into classIDs/weights)
//	classIDs  nPairs x uint8
//	weights   nPairs x float32
//	nTagMap   uint32   then per entry: uint16 len + bytes, uint8 classID
//	nMeta     uint32   then per entry: key and value, each uint16 len + bytes
//	                                        (version 2 and later)
//
// The metadata names what the model is: its language, tag set, and the
// feature template it was trained with (see internal/perceptron). A version 1
// model has none and is read as the English Penn Treebank model.
//
// Feature strings are sliced zero-copy out of the backing bytes. The numeric
// arrays are copied once into aligned slices, since go:embed makes no
// alignment guarantee and reinterpreting unaligned bytes as []float32 is not
// safe. That copy is a few MB of memcpy — immaterial at load time.
package flat

// Magic identifies a flat perceptron model file.
const Magic = "PROSETAG"

// Version is the current format version. Bump on any layout change; the
// reader rejects anything it does not recognise.
//
// Version 2 appended the metadata section. Version 1 files are still read.
const Version = 2

// Metadata keys a model may carry.
const (
	// MetaLanguage is the BCP 47 language tag the model was trained for.
	MetaLanguage = "lang"
	// MetaTagset names the tag set the model emits: "penn" or "upos".
	MetaTagset = "tagset"
	// MetaTemplate is the feature template number, as a decimal string.
	MetaTemplate = "template"
	// MetaSource describes the training data.
	MetaSource = "source"
	// MetaLicense is the license of the training data.
	MetaLicense = "license"
)

// MaxClasses is the ceiling implied by storing class IDs as uint8. The
// shipped English model has 38.
const MaxClasses = 256
