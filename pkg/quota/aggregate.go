package quota

func (t *TotalQuota) Sum() *WorkloadQuota {
	wl := WorkloadQuota{
		Request: &ComputeQuota{},
		Limit:   &ComputeQuota{},
	}
	for _, a := range t.AggregateQuotas {
		wl.Add(a.Sum())
	}
	return &wl
}

func (a *AggregateQuota) Sum() *WorkloadQuota {
	wl := WorkloadQuota{
		Request: &ComputeQuota{},
		Limit:   &ComputeQuota{},
	}
	for _, q := range a.WorkloadQuotas {
		wl.Add(q)
	}
	return &wl
}

func (w *WorkloadQuota) Add(o *WorkloadQuota) {
	w.Request.Add(o.Request)
	w.Limit.Add(o.Limit)
}

func (r *ComputeQuota) Add(o *ComputeQuota) {
	r.CPU += o.CPU
	r.Mem += o.Mem
}
