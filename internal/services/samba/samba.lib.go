package samba

func FilterShares(shares Shares) Shares {
	result := make(Shares)
	for k, v := range shares {
		if k != "global" && k != "printers" && k != "print$" {
			result[k] = v
		}
	}

	return result
}
