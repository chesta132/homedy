package samba

func FilterShares(shares Shares) Shares {
	delete(shares, "global")
	delete(shares, "printers")
	delete(shares, "print$")
	return shares
}
