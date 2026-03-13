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

func isPathExist(shares Shares, share Share) bool {
	for _, _share := range shares {
		if _share.Path == share.Path {
			return true
		}
	}
	return false
}
