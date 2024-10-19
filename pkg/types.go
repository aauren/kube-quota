package kubequota

type CPUMilicore int64
type MemBytes int64
type StorageBytes int64
type ClaimsNum int16

func (m *MemBytes) ToBytes() int64 {
	return int64(*m)
}

func (s *StorageBytes) ToBytes() int64 {
	return int64(*s)
}
