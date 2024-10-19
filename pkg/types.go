package kubequota

type CPUMilicore int64
type MemBytes int64
type StorageBytes int64
type ClaimsNum int16
type DivideByZeroError struct {
	message string
}
type Percentage struct {
	Parts int64
	Whole int64
}

func (m *MemBytes) ToBytes() int64 {
	return int64(*m)
}

func (s *StorageBytes) ToBytes() int64 {
	return int64(*s)
}

func (d DivideByZeroError) Error() string {
	return d.message
}

func (p *Percentage) Percentage() (float64, error) {
	if p.Whole == 0 && p.Parts != 0 {
		return 0, DivideByZeroError{message: "whole was 0 while parts was non-zero this would end up in a divide by zero error"}
	}
	if p.Parts == 0 {
		return 0, nil
	}
	//nolint:gomnd // 100 here, in terms of a percentage, is pretty self-explanitory
	return float64(p.Parts) / float64(p.Whole) * 100, nil
}
