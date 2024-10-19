package quota

func (t *NamespaceWorkloadQuota) Sum() *WorkloadQuota {
	wl := WorkloadQuota{
		Request: &ComputeQuota{},
		Limit:   &ComputeQuota{},
		StorageQuota: &StorageQuota{
			Ephemeral:      &EphemeralQuota{},
			StorageClasses: make(map[string]*StorageClassQuota),
		},
	}
	for _, a := range t.PodQuotas {
		wl.Add(a.Sum())
	}
	return &wl
}

func (a *PodQuota) Sum() *WorkloadQuota {
	wl := WorkloadQuota{
		Request: &ComputeQuota{},
		Limit:   &ComputeQuota{},
		StorageQuota: &StorageQuota{
			Ephemeral:      &EphemeralQuota{},
			StorageClasses: make(map[string]*StorageClassQuota),
		},
	}
	for _, q := range a.WorkloadQuotas {
		wl.Add(q)
	}
	return &wl
}

func (w *WorkloadQuota) Add(o *WorkloadQuota) {
	w.Request.Add(o.Request)
	w.Limit.Add(o.Limit)
	w.StorageQuota.Add(o.StorageQuota)
}

func (r *ComputeQuota) Add(o *ComputeQuota) {
	r.CPU += o.CPU
	r.Mem += o.Mem
}

func (e *EphemeralQuota) Add(o *EphemeralQuota) {
	e.Requests += o.Requests
	e.Limits += o.Limits
}

func (s *StorageClassQuota) Add(o *StorageClassQuota) {
	s.Claims += o.Claims
	s.Requests += o.Requests
}

func (s *StorageQuota) Add(o *StorageQuota) {
	s.Ephemeral.Add(o.Ephemeral)
	for scKey, scVal := range o.StorageClasses {
		sc, ok := s.StorageClasses[scKey]
		if !ok {
			sc = &StorageClassQuota{
				Name: sc.Name,
			}
			s.StorageClasses[scKey] = sc
		}
		sc.Add(scVal)
	}
}
